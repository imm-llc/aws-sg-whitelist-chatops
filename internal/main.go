package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/google/uuid"
)

// Init intializes the AWS session and the service clients
func Init() error {
	awsSess, err = session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	if err != nil {
		log.Println("Unable to instantiate AWS session")
		return ErrAWSSession
	}
	log.Println("Instantiated AWS Session")

	sg = ec2.New(awsSess)
	dynamo = dynamodb.New(awsSess)

	return nil
}

// AddEntryToSG adds an ingress entry to an AWS SG
// AWS docs: https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.AuthorizeSecurityGroupIngress
func AddEntryToSG(ip string, username string, environment string, resource string) error {
	i := ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(securityGroupIDMap[environment]), // grab environment's security group ID from map
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(resourcePortMap[resource]), // grab resource's port number from map
				ToPort:     aws.Int64(resourcePortMap[resource]),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(fmt.Sprintf("%s/32", ip)),
						Description: aws.String(username),
					},
				},
			},
		},
	}

	result, err := sg.AuthorizeSecurityGroupIngress(&i)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return err
	}

	log.Println(result)

	err = addEntryToDynamo(ip, username, securityGroupIDMap[environment], resourcePortMap[resource], resourcePortMap[resource], environment)

	if err != nil {
		return ErrDynamoDB
	}

	return nil
}

// CheckExistingEntry ensures the IP does not already exist in the security group
func CheckExistingEntry(ip string, environment string) error {

	i := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(securityGroupIDMap[environment]),
		},
	}

	result, err := sg.DescribeSecurityGroups(i)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return err
	}

	// Iterate through existing entries and check if the IP is already allowed
	for _, p := range result.SecurityGroups[0].IpPermissions {
		for _, e := range p.IpRanges {
			if *e.CidrIp == fmt.Sprintf("%s/32", ip) {
				return ErrEntryExists
			}
		}
	}

	return nil

}

// RemoveSecuritGroupRule removes an ingress rule from a security group
func RemoveSecuritGroupRule(ip string, sgID string, fromPort int64, toPort int64, slackUsername string) error {

	i := &ec2.RevokeSecurityGroupIngressInput{
		GroupId: aws.String(sgID),
		DryRun:  aws.Bool(false),
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(fromPort),
				ToPort:     aws.Int64(toPort),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(fmt.Sprintf("%s/32", ip)),
						Description: aws.String(slackUsername),
					},
				},
			},
		},
	}

	_, err := sg.RevokeSecurityGroupIngress(i)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return err
	}

	log.Printf("Removed entry for IP '%s'", ip)

	return nil
}

func addEntryToDynamo(ip string, username string, sgID string, fromPort int64, toPort int64, environment string) error {

	item := DynamoWhitelistEntry{
		IP:          ip,
		Username:    username,
		Expiration:  time.Now().Unix() + int64(60*60*2), // Now + (60 seconds * 60 minutes * 2 hours)
		FromPort:    fromPort,
		ToPort:      toPort,
		Protocol:    "tcp",
		Environment: environment,
		Group:       sgID,
		ID:          strings.Replace(uuid.New().String(), "-", "", -1),
	}

	ddbItem, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		log.Println("Error creating DynamoDB item", err.Error())
		return err
	}

	i := &dynamodb.PutItemInput{
		Item:      ddbItem,
		TableName: aws.String(os.Getenv("DDB_TABLE_NAME")),
	}

	_, err = dynamo.PutItem(i)

	if err != nil {
		log.Println("Error adding item to DynamoDB", err.Error())
		return err
	}

	return nil
}

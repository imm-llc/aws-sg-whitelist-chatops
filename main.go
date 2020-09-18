package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	n := time.Now().Unix()
	twoHours := int64(60 * 60 * 2)
	fmt.Println(n + twoHours)

	os.Exit(0)

	awsSess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-2")})

	if err != nil {
		log.Println("Unable to instantiate AWS session")
		//return err
	}
	log.Println("Instantiated AWS Session")

	sg := ec2.New(awsSess)

	i := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String("sg-023d321b2df6e1b6d"),
		},
	}

	result, err := sg.DescribeSecurityGroups(i)

	if err != nil {
		log.Println(err)
	}

	ip := "50.232.79.90"
	for _, p := range result.SecurityGroups[0].IpPermissions {
		for _, e := range p.IpRanges {
			if *e.CidrIp == fmt.Sprintf("%s/32", ip) {
				log.Println("Already allowed")
			}
		}
	}

	inp := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String("sg-023d321b2df6e1b6d"),
		DryRun:  aws.Bool(false),
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(5432), // Hardcoding for now, eventually this should be a variable for extense-ability
				ToPort:     aws.Int64(5432),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(fmt.Sprintf("%s/32", ip)),
						Description: aws.String("Testing"),
					},
				},
			},
		},
	}

	res, err := sg.AuthorizeSecurityGroupIngress(inp)

	if err != nil {
		log.Println(err.Error()) // TODO check for AWS error
	}

	log.Println(res)

	rev := &ec2.RevokeSecurityGroupIngressInput{
		GroupId: aws.String("sg-023d321b2df6e1b6d"),
		DryRun:  aws.Bool(false),
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(5432), // Hardcoding for now, eventually this should be a variable for extense-ability
				ToPort:     aws.Int64(5432),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(fmt.Sprintf("%s/32", ip)),
						Description: aws.String("Testing"),
					},
				},
			},
		},
	}

	r, err := sg.RevokeSecurityGroupIngress(rev)

	if err != nil {
		log.Println(err)
	}

	log.Println(r)

}

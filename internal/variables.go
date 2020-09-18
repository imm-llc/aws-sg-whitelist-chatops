package internal

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	awsSess *session.Session
	sg      *ec2.EC2
	dynamo  *dynamodb.DynamoDB
	err     error

	// TODO make this a map[string]map[string]string{} for environment: resource: SG ID

	securityGroupIDMap = map[string]string{
		"staging":    os.Getenv("STAGING_SECURITY_GROUP_ID"),
		"production": os.Getenv("PRODUCTION_SECURITY_GROUP_ID"),
	}

	resourcePortMap = map[string]int64{
		"rds": 5432,
	}
	/*
		environmentResourcePortMap = map[string]map[string]int64{
			"staging": {
				"rds": 5432,
			},
			"production": {
				"rds": 5432,
			}
		}
	*/
)

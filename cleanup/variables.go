package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	awsSess *session.Session
	sg      *ec2.EC2
	//dynamo  *dynamodb.DynamoDB
	err error

	ip          string
	username    string
	fromPort    int64
	toPort      int64
	protocol    string
	environment string
	groupID     string
)

package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/imm-llc/aws-sg-whitelist-chatops/internal"
)

func lambdaHandler(ctx context.Context, e events.DynamoDBEvent) {
	internal.Init()

	for _, record := range e.Records {
		log.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
		if record.EventName != "REMOVE" {
			log.Printf("Not actioning '%s' event type, dumping record", record.EventName)
			// Print the entire event
			log.Println(e)

		} else {
			// Print new values for attributes of type String
			for name, value := range record.Change.OldImage {
				switch name {
				case "IP":
					ip = value.String()
				case "Username":
					username = value.String()
				case "FromPort":
					fromPort, _ = value.Integer()
				case "ToPort":
					toPort, _ = value.Integer()
				case "Protocol":
					protocol = value.String()
				case "Group":
					groupID = value.String()
				}
			}
			log.Println(record.Change.OldImage)
			if err := internal.RemoveSecuritGroupRule(ip, groupID, fromPort, toPort, username); err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	lambda.Start(lambdaHandler)
}

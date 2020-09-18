package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func whitelistAlreadyExists() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       "IP address already allowed",
		StatusCode: 200,
	}
}

func slackVerifyFailed() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       "Unable to verify Slack token",
		StatusCode: 200,
	}
}

func unknownEnvironment(environment string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Enviroment '%s' not recognized", environment),
		StatusCode: 200,
	}
}

func unknownAppname(app string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("App '%s' not recognized", app),
		StatusCode: 200,
	}
}

func badRequest(environment string, resource string, ip string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Expected '</command> environment resource ip', received '</command> '%s' '%s' '%s'", environment, resource, ip),
		StatusCode: 200,
	}
}

func badIP() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       "IP Address must be in x.x.x.x format",
		StatusCode: 200,
	}
}

func unauthorized() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       "You are not authorized to use this tool",
		StatusCode: 200,
	}
}

func internalError() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       "Something went wrong trying to add whitelist",
		StatusCode: 200,
	}
}

func addedIPSuccessfully() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       "Successfully added your IP address. This will expire in 2 hours.",
		StatusCode: 200,
	}
}

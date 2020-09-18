package main

import (
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/imm-llc/aws-sg-whitelist-chatops/internal"
)

func lambdaHanlder(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	parsedRequest, err := url.ParseQuery(req.Body)

	if err != nil {
		log.Println("Error parsing request body", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Unable to parse request"}, nil
	}

	values := slackRequest{
		parsedRequest.Get("token"),
		parsedRequest.Get("team_id"),
		parsedRequest.Get("team_domain"),
		parsedRequest.Get("channel_id"),
		parsedRequest.Get("channel_name"),
		parsedRequest.Get("user_id"),
		parsedRequest.Get("user_name"),
		parsedRequest.Get("command"),
		parsedRequest.Get("text"),
		parsedRequest.Get("response_url"),
		parsedRequest.Get("trigger_id"),
	}

	if values.token != os.Getenv("SLACK_VERIFICATION_TOKEN") {
		return slackVerifyFailed(), nil
	}

	// Environment, resource, IP
	splitSlackText := strings.Split(values.text, " ")

	// Make sure we have valid values
	if splitSlackText[0] == "" {
		return badRequest(splitSlackText[0], splitSlackText[1], splitSlackText[2]), nil
	} else if splitSlackText[1] == "" {
		return badRequest(splitSlackText[0], splitSlackText[1], splitSlackText[2]), nil
	} else if splitSlackText[2] == "" {
		return badRequest(splitSlackText[0], splitSlackText[1], splitSlackText[2]), nil
	}

	environment := splitSlackText[0]
	resource := splitSlackText[1]
	ip := splitSlackText[2]

	if !checkIPRegex(ip) {
		return badIP(), nil
	}

	// Ensure a valid environment was provided
	// If environment is not production, the first check will return true
	if environment != "production" && environment != "staging" {
		return unknownEnvironment(environment), nil
	}

	// Ensure user's Slack ID is allowed to add whitelist
	validUserIDs := strings.Split(allowedSlackUserIDS[environment], ",")
	for _, uid := range validUserIDs {
		if uid == values.userID {
			log.Printf("Adding entry for Slack user %s %s", values.userID, values.userName)
			// If the user is allowed to add entries, try to add their IP
			return wrapAddWhitelist(values.userName, ip, environment, resource)
		}
	}

	// If user didn't match in the above loop, they're not authorized
	return unauthorized(), nil
}

func checkIPRegex(ip string) bool {
	// e.g. 1.1.1.1 , not 1.1.1.1/32
	regex := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`
	compiled, _ := regexp.Compile(regex)

	match := compiled.Match([]byte(ip))

	if !match {
		log.Printf("IP Address '%s' does not match regex", ip)
		return false
	}

	return true
}

func wrapAddWhitelist(slackUserName string, ip string, environment string, resource string) (events.APIGatewayProxyResponse, error) {
	err := internal.Init()

	if err != nil {
		return internalError(), nil
	}

	err = internal.CheckExistingEntry(ip, environment)

	if err != nil {
		if err == internal.ErrEntryExists {
			return whitelistAlreadyExists(), nil
		}

		return internalError(), nil
	}

	err = internal.AddEntryToSG(ip, slackUserName, environment, resource)

	if err != nil {
		return internalError(), nil
	}
	return addedIPSuccessfully(), nil
}

func main() {
	lambda.Start(lambdaHanlder)
}

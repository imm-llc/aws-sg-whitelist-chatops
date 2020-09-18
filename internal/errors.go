package internal

import "errors"

var (
	// ErrAWSSession indicates there was an error creating the AWS session
	ErrAWSSession = errors.New("Unable to create AWS session")
	// ErrEntryExists indicates the IP has already been added
	ErrEntryExists = errors.New("IP already whitelisted")
	// ErrDynamoDB indicates there was an issue adding the item to DynamoDB
	ErrDynamoDB = errors.New("Error adding DynamoDB entry")
)

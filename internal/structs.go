package internal

// DynamoWhitelistEntry is our entry into DynamoDB
type DynamoWhitelistEntry struct {
	IP          string
	Username    string
	Expiration  int64
	FromPort    int64
	ToPort      int64
	Protocol    string
	Group       string
	Environment string
	ID          string
}

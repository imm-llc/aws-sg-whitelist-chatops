package main

import "os"

var (
	allowedSlackUserIDS = map[string]string{
		"staging":    os.Getenv("STAGING_ALLOWED_SLACK_USER_IDS"),
		"production": os.Getenv("PRODUCTION_ALLOWED_SLACK_USER_IDS"),
	}
)

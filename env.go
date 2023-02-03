package main

import (
	"log"
	"os"
	"strconv"
)

func getEnvInt(key string) int {
	var variable string = os.Getenv(key)
	if len(variable) == 0 {
		variable = "0"
	}
	val, err := strconv.Atoi(variable)
	if err != nil {
		log.Fatal("Failed to convert env var to integer")
	}
	return val
}

func logEnvs() {
	log.Printf(`Environment variables:
	- PORT: %v
	- SLACK_BOT_OAUTH_TOKEN: %v
	- SLACK_BOT_NOTIFICATION_COLOR: %v
	- SLACK_BOT_NOTIFICATION_GREATINGS: %v
	- GITLAB_WEBHOOK_SECRET_TOKEN: %v
	- SLACK_EVENT_READ_CHANNEL: %v
	- USER_EMAIL_DOMAIN: %v
	- USER_EMAIL_SPACE_REPLACER: %v`,
		PORT,
		SLACK_BOT_OAUTH_TOKEN,
		SLACK_BOT_NOTIFICATION_COLOR,
		SLACK_BOT_NOTIFICATION_GREATINGS,
		GITLAB_WEBHOOK_SECRET_TOKEN,
		SLACK_EVENT_READ_CHANNEL,
		USER_EMAIL_DOMAIN,
		USER_EMAIL_SPACE_REPLACER)
}

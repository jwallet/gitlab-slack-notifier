package main

import "os"

type Environment string

var PORT int = getEnvInt("PORT")

var SLACK_BOT_READ_CHANNEL = os.Getenv("SLACK_BOT_READ_CHANNEL")

var SLACK_BOT_OAUTH_TOKEN = os.Getenv("SLACK_BOT_OAUTH_TOKEN")

var DOMAIN_EMAIL = os.Getenv("DOMAIN_EMAIL") // "gmail.com"

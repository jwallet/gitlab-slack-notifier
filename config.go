package main

import "os"

type Environment string

// Constants

var USER_EMAIL_SPACE_REPLACER = "."

var SLACK_BOT_NOTIFICATION_COLOR = "#0099CC"
var SLACK_BOT_NOTIFICATION_GREATINGS = "You were mentionned on GitLab"

// Environment variables

var PORT int = getEnvInt("PORT")

var SLACK_BOT_OAUTH_TOKEN string = os.Getenv("SLACK_BOT_OAUTH_TOKEN")
var SLACK_BOT_READ_CHANNEL string = os.Getenv("SLACK_BOT_READ_CHANNEL")
var USER_EMAIL_DOMAIN string = os.Getenv("USER_EMAIL_DOMAIN") // "gmail.com"

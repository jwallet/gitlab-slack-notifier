package main

import "os"

type Environment string

var PORT int = getEnvInt("PORT") // optional, default 3000

var SLACK_BOT_OAUTH_TOKEN string = os.Getenv("SLACK_BOT_OAUTH_TOKEN")                       // required, bot token, xoxb-12345678-12345678
var SLACK_BOT_NOTIFICATION_COLOR string = os.Getenv("SLACK_BOT_NOTIFICATION_COLOR")         // optional, color of the notification border, "#0099CC"
var SLACK_BOT_NOTIFICATION_GREATINGS string = os.Getenv("SLACK_BOT_NOTIFICATION_GREATINGS") // required, message sent by the bot, contrains { author, repositoryName, mergeRequestTitle } "*%v* mentionned you on %v merge request **%v**"

var GITLAB_WEBHOOK_SECRET_TOKEN string = os.Getenv("GITLAB_WEBHOOK_SECRET_TOKEN") // required, signature of gitlab webhook

var SLACK_EVENT_READ_CHANNEL string = os.Getenv("SLACK_EVENT_READ_CHANNEL")   // optional, monitor one channel "CH1234567"
var USER_EMAIL_DOMAIN string = os.Getenv("USER_EMAIL_DOMAIN")                 // required, email domain, "business.com"
var USER_EMAIL_SPACE_REPLACER string = os.Getenv("USER_EMAIL_SPACE_REPLACER") // optional, leave empty to replace spaces to nothing, "."

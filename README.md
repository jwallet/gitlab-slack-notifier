# GitLab to Slack Notifier
#### Send a private message when a user is mentionned in comments

![preview](./215661429-dd1b2944-4b9f-46a0-9d87-f06c4f05f5f9.png)

### What it does ?
1. From the GitLab Slack Integration, it reads published ***comments/notes*** to the selected ***Slack channel*** of your Slack bot webhook.
     1. GitLab cannot publish to different channels from the Slack Integration; [opened issue](https://gitlab.com/gitlab-org/gitlab/-/issues/12895)
1. The bot will try to find one or more mentions `@that.guy` in the comment.

    _This Guy (this.guy) commented on merge request !1 in Project / Repo: MR_
    >    _@that.guy I need your review._

1. It will fetch from **GitLab** Open API, the user info to get its fullname. (_See the [**What's missing?**](#whats-missing), to retrieve the email directly_)
   1. It will format his fullname to a user email name using `formatFullnameToUserEmail()` to a lowercase, dot separated format without diacritic, e.g.: Nathan Côté-Dumais → nathan.cote-dumais.
   1. It will then use the environment variable `USER_EMAIL_DOMAIN` to create a valid email, e.g.: nathan.cote-dumais → `nathan.cote-dumais@business.com`
1. It will fetch from Slack API the user info using his email to then extract his userID, e.g.: `UA1BCDEF`.
1. Finally, it will publish a private message notification to Slack to this user using his userID. The user will receive the notification from the bot itself in the Slack app section. 

### Setup
1. Create a [Slack bot](https://api.slack.com/apps) with an active incoming webhook.
1. Copy the webhook URL to your GitLab → Repo → Settings → Integrations → Slack Notifications Integration → Webhook URL
1. Publish this repo and serve it as a web service
1. Set the environment variables in `config.go`
    1. `PORT` Default to `3000`
    1. `SLACK_BOT_READ_CHANNEL` same as the webhook channel, the ID can be found by opening your Slack Workspace in Slack web app in a browser and getting it from the URL
    1. `SLACK_BOT_OAUTH_TOKEN` see your bot **Bot User OAuth Token** under **Install App** section
    1. `USER_EMAIL_DOMAIN` all user emails on the same email domain `@business.com`
1. Go back to your bot, in the **Event subscriptions** and paste where you host this app `https://my.webservice.com/`
1. Then, in the same section, _Subscribe to bot events_ by adding **message.channels** `Scope channels:history` to be able to read the channel where you receive GitLab comments.
1. Go to **OAuth & Permissions**, scroll down to **Scopes**, and select these scopes:
    1. `channels:history` to read the channel
    1. `im:write` to notify a user
    1. `users:read` to fetch user info from Slack API
    1. `users:read.email` to fetch user info from SLACK API
    1. `incoming-webhook` (optional) if this bot is used by GitLab to post to the channel (gitlab) 
### What's missing?
Some other ways to get the user email directly from GitLab
#### Self-hosted
If you are self-hosting GitLab, then this bot can be simplified by getting the user email from GitLab self-hosted API `GET:Users` with an oauth access token, instead of using a formatter `formatFullnameToUserEmail()` and environment variable `USER_EMAIL_DOMAIN` and `USER_EMAIL_SPACE_REPLACER` to generate a business user email. I did not have to implement this because my business email format is based on the user fullname and this info was available through the public API.
#### User.username is the same as my business email
If you are in luck and everyone of your users have set their username `my.name` properly and they match your business email, then you could just concat this value to your domain. I was out of luck and had to fallback to the user fullname `My Name` and format it.
#### Ask your users to set their Slack email as their GitLab public email
If all users set a public email on their profile, you can fetch it one-by-one by using `GET:Users/:id` or by fetching all members of a group (your business) with your private token `GET:Groups/:groupId/members?private_token=ACCESS_TOKEN`.
#### Still not working for you?
Don't forget you can change the formatter that uses the user fullname in the code to make it match your business email. However, if all users are connected without a business email and on your Slack app as well, then it will be really hard for you to match the GitLab user to a Slack UserID. GitLab has some opened issues on that matter. 
##### Make the bot smarter
This is a long shot, but you can improve the bot by asking your Slack users to identify themself to the bot and connect their GitLab account to the bot so the bot can keep a Dictionary of `Dictionary<GitLabUserName, SlackUserId>` on a database. Just like `GitLab` app bot does, it adds an entry in your GitLab profile under `Chat` to have access to your GitLab profile, so no need to have your email since the bot already knows the Slack UserID since the connect-request came from their.

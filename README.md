# GitLab to Slack Notifier
#### Send a private message when a user is mentionned in comments

### What it does ?
1. From the GitLab Slack Integration, it reads published ***comments/notes*** to the selected ***Slack channel*** of your Slack bot webhook.
     1. GitLab cannot publish to different channels from the Slack Integration; [opened issue](https://gitlab.com/gitlab-org/gitlab/-/issues/12895)
1. The bot will try to find one or more mentions `@that.guy` in the comment.

    _This Guy (this.guy) commented on merge request !1 in Project / Repo: MR_
    >    _@that.guy I need your review._

1. It will fetch from **GitLab** Open API, the user info to get its fullname. (_See the [**What's missing?**](#whats-missing), to retrieve the email directly_)
   1. It will format his fullname to a user email name using `formatGitLabUsernameTag()` to a lowercase, dot separated format without diacritic, e.g.: Nathan Côté-Dumais → nathan.cote-dumais.
   1. It will then use the environment variable `DOMAIN_EMAIL` to create a valid email, e.g.: nathan.cote-dumais → `nathan.cote-dumais@business.com`
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
    1. `DOMAIN_EMAIL` all user emails on the same email domain `@business.com`
1. Go back to your bot, in the **Event subscriptions** and paste where you host this app `https://my.webservice.com/`
1. Then, in the same section, _Subscribe to bot events_ by adding **message.channels** `Scope channels:history` to be able to read the channel where you receive GitLab comments.

### What's missing?
This bot can be simplified by getting the user email from GitLab private API `GET:Users`, instead of using a formatter `formatGitLabUsernameTag()` and environment variable `USER_DOMAIN` to generate a user email. I did not implemented this yet since my business email format is based on the user fullname and this info was available through the public API.
1. Before fetching a user info, implement a request to authenticate your bot on [GitLab API](https://docs.gitlab.com/ee/api/rest/)
2. Fetch the `GET:Users` with `email=foo.bar@domain.com` as a query param along with your `access_token`
3. Your response will have more data, including `user.email`
4. Skip the rest of the code until the bot fetches the user from Slack API with an email, pass it there.

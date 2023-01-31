package main

import (
	"fmt"
	"regexp"
	"strings"
)

func handle(event SlackEvent) error {
	if event.Type != SLACK_BOT_READ_CHANNEL {
		return fmt.Errorf("Not monitoring the right channel, stopping.")
	}
	if event.Subtype != "bot_message" {
		return fmt.Errorf("Not a message from a bot, stopping.")
	}
	if len(event.Attachments) == 0 {
		return fmt.Errorf("No message found, no mention needed, stopping.")
	}

	comment := event.Attachments[0].Text
	usernames := getAllUsernameTags(comment)
	fmt.Printf("Got usernames: %v\n", usernames)

	for _, username := range usernames {
		userID, err := getUserID(username)
		if err != nil {
			return err
		}
		fmt.Printf("Got userID %v for username %v\n", userID, username)
		if userID == "" {
			return fmt.Errorf("Did not find any user ID for %v", username)
		}

		botMessage := &BotMessage{
			UserID:  userID,
			Link:    getLink(event.Text),
			Text:    comment,
			Channel: event.Channel,
			EventTS: event.TS,
		}

		notify(botMessage)
	}

	return nil
}

func getLink(text string) string {
	selector := regexp.MustCompile(`<https://gitlab.com/.*?\|commented on merge request`)
	tagURL := selector.FindString(text)
	urlCleaner := strings.NewReplacer("<", "", "|commented on merge request", "")
	return urlCleaner.Replace(tagURL)
}

func getAllUsernameTags(comment string) []string {
	selector := regexp.MustCompile(`@\w+.?\w+`)
	usernames := selector.FindAllString(comment, -1)
	for i, username := range usernames {
		usernames[i] = strings.Replace(username, "@", "", 1)
	}
	return usernames
}

func getUserID(username string) (string, error) {
	// fetch user from GitLab that generate an email
	userEmail, err := fetchGitLabUserToFormattedEmail(username)
	if err != nil {
		return "", err
	}

	if !isEmailValid(userEmail) {
		return "", fmt.Errorf("User email is invalid")
	}

	// send query to Slack with email to retrieve the user ID
	slackUser, err := fetchSlackUser(userEmail)
	if err != nil {
		return "", err
	}

	return slackUser.User.Id, nil
}

func fetchGitLabUserToFormattedEmail(username string) (string, error) {
	// send query to gitlab with tag to retrieve user fullname
	gitLabUser, err := fetchBasicGitLabUser(username)
	fmt.Printf("GitLab user payload: %v\n", gitLabUser)
	if err != nil {
		return "", err
	}

	// user email is public
	if gitLabUser.Email != "" {
		return gitLabUser.Email, nil
	}

	// transform fullname to email username
	usernameTag, err := formatFullnameToUserEmail(gitLabUser.Name)
	if err != nil {
		usernameTag = username
		fmt.Printf("Cannot parse user name to tag: %v\n", err)
	}

	fmt.Printf("Fullname to email username: %v\n", usernameTag)

	// generate full email
	return usernameTag + "@" + USER_EMAIL_DOMAIN, nil
}

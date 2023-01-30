package main

import (
	"fmt"
	"regexp"
	"strings"
)

func handle(event SlackEvent) error {
	if event.Type == "C02HSU2AWN8" { // #random channel
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
		fmt.Printf("Got userID %v for username %v\n", userID, username)
		if err != nil {
			return err
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
	// call gitlab with tag to retrieve email
	gitLabUser, err := fetchBasicGitLabUser(username)
	fmt.Println(gitLabUser)
	if err != nil {
		return "", err
	}
	usernameTag, err := formatGitLabUsernameTag(gitLabUser.Name)
	if err != nil {
		usernameTag = username
		fmt.Printf("Cannot parse user name to tag: %v\n", err)
	}
	fmt.Println(usernameTag)

	// call slack with email to retrieve id
	slackUser, err := fetchSlackUser(usernameTag + "@" + DOMAIN_EMAIL)
	if err != nil {
		return "", err
	}

	return slackUser.User.Id, nil
}

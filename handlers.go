package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Event struct {
	author       string
	message      string
	repository   string
	mergeRequest string
	link         string
}

func handleSlackEvent(slackEvent SlackEvent) error {
	log.Printf("Slack event: %v\n", slackEvent.Type)
	if slackEvent.Type != "message" {
		return fmt.Errorf("Bot can only handle slack event message, stopping.")
	}
	if slackEvent.Subtype != "bot_message" {
		return fmt.Errorf("Not a message from a bot, stopping.")
	}
	if len(slackEvent.Attachments) == 0 {
		return fmt.Errorf("No message found, no mention needed, stopping.")
	}
	if SLACK_EVENT_READ_CHANNEL != "" && slackEvent.Channel != SLACK_EVENT_READ_CHANNEL {
		return fmt.Errorf("Not monitoring the right channel, stopping.")
	}

	event := Event{
		author:       getAliasFromEventText(slackEvent.Text, author),
		message:      slackEvent.Attachments[0].Text,
		repository:   getAliasFromEventText(slackEvent.Text, repo),
		mergeRequest: getAliasFromEventText(slackEvent.Text, mr),
		link:         getAliasFromEventText(slackEvent.Text, link),
	}

	return handle(event)
}

var previousNoteId int32 = 0

func handleGitLabWebhook(gitLabEvent GitLabWebhookEvent) error {
	if gitLabEvent.EventType != "note" {
		return fmt.Errorf("Bot can only handle gitlab webhook notes, stopping.")
	}
	if gitLabEvent.Note.Note == "" && gitLabEvent.Note.Description == "" {
		return fmt.Errorf("No message found, no mention needed, stopping.")
	}
	if previousNoteId == gitLabEvent.Note.Id {
		return fmt.Errorf("Bot prevented to proceed the same gitlab note id, stopping.")
	}
	previousNoteId = gitLabEvent.Note.Id
	log.Printf("GitLab event: %v\n", gitLabEvent)

	event := Event{
		author:       gitLabEvent.User.Name,
		message:      defaults(gitLabEvent.Note.Note, gitLabEvent.Note.Description),
		repository:   defaults(gitLabEvent.Repository.Name, gitLabEvent.Project.Name),
		mergeRequest: defaults(gitLabEvent.MergeRequest.Title, gitLabEvent.MergeRequest.Description),
		link:         gitLabEvent.Note.Url,
	}

	return handle(event)
}

func handle(event Event) error {
	usernames := getAllUsernameTags(event.message)
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

		greatings := strings.Join([]string{
			fmt.Sprintf("<@%s>", userID),
			getGreatings(event),
		}, " ")

		attachment := SlackAttachment{
			Text:   event.message,
			Footer: event.link,
			Color:  SLACK_BOT_NOTIFICATION_COLOR,
		}

		attachmentData := &bytes.Buffer{}
		enc := json.NewEncoder(attachmentData)
		enc.SetEscapeHTML(false)
		err = enc.Encode([]SlackAttachment{attachment})
		if err != nil {
			return err
		}

		botMessage := &BotMessage{
			UserID:      userID,
			Attachments: attachmentData.String(),
			Message:     greatings,
		}

		notify(botMessage)
	}

	return nil
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Event struct {
	author       string
	message      string
	repository   string
	mergeRequest string
	link         string
}

var previousEventTS string = ""

func handleSlackEvent(slackEvent SlackEvent) error {
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
	if slackEvent.TS == previousEventTS {
		return fmt.Errorf("Bot prevented to proceed the same Slack event based on the timestamp %v, stopping.", slackEvent.TS)
	}

	previousEventTS = slackEvent.TS
	fmt.Printf("Slack event passed validation for timestamp: %v\n", slackEvent.TS)

	event := Event{
		author:       getAliasFromEventText(slackEvent.Text, author),
		message:      slackEvent.Attachments[0].Text,
		repository:   getAliasFromEventText(slackEvent.Text, repo),
		mergeRequest: getAliasFromEventText(slackEvent.Text, mr),
		link:         getAliasFromEventText(slackEvent.Text, link),
	}

	return handle(event)
}

var previousNoteId int64 = 0

func handleGitLabWebhook(gitLabEvent GitLabWebhookEvent) error {
	if gitLabEvent.EventType != "note" {
		return fmt.Errorf("Bot can only handle gitlab webhook notes, stopping.")
	}
	if gitLabEvent.Note.Note == "" && gitLabEvent.Note.Description == "" {
		return fmt.Errorf("No message found, no mention needed, stopping.")
	}
	if previousNoteId == gitLabEvent.Note.Id {
		return fmt.Errorf("Bot prevented to proceed the same gitlab note id ''%v'' with note url ''%v'', stopping.", gitLabEvent.Note.Id, gitLabEvent.Note.Url)
	}

	previousNoteId = gitLabEvent.Note.Id
	fmt.Printf("GitLab event passed validation for note: %v\n", gitLabEvent.Note.Url)

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
	fmt.Printf("GitLab note usernames found: %v\n", usernames)

	for _, username := range usernames {
		userID, err := getUserID(username)
		if err != nil {
			return err
		}

		if userID == "" {
			return fmt.Errorf("Bot did not find any user ID for %v", username)
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
		return "", fmt.Errorf("User email %v is invalid", userEmail)
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

	fmt.Printf("GitLab user fullname to email username: %v\n", usernameTag)

	// generate full email
	return usernameTag + "@" + USER_EMAIL_DOMAIN, nil
}

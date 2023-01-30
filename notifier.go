package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type BotMessage struct {
	UserID  string `json:"userId"`
	Link    string `json:"link"`
	Text    string `json:"text"`
	Channel string `json:"channel"`
	EventTS string `json:"event_ts"`
}

type RequestPermalink struct {
	Permalink string `json:"permalink"`
}

func notify(message *BotMessage) error {
	permalink, err := fetchPermalink(message)
	if err != nil {
		return err
	}
	message.Link = permalink
	err = pushNotification(message)
	return err
}

func fetchPermalink(message *BotMessage) (string, error) {
	client := &http.Client{}
	var requestLink RequestPermalink

	endpoint := fmt.Sprintf("https://slack.com/api/chat.getPermalink?token=%s&channel=%s&message_ts=%s",
		SLACK_BOT_OAUTH_TOKEN, message.Channel, message.EventTS)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	fmt.Printf("Fetching permalink: %s\n", endpoint)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GET permalink Failed %v", resp.StatusCode)
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Reading GET permalink request")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Casting GET permalink request")
	json.Unmarshal(body, &requestLink)

	return requestLink.Permalink, nil
}

func pushNotification(message *BotMessage) error {
	alert := strings.Join([]string{
		fmt.Sprintf("<@%s>", message.UserID),
		"Tu as été mentionné dans le channel `#_notifications`",
	}, " ")

	attachment := SlackAttachment{
		Fallback: "Message originale: " + message.Link,
		Text:     message.Text,
		Footer:   message.Link,
	}
	attachmentData := &bytes.Buffer{}
	enc := json.NewEncoder(attachmentData)
	enc.SetEscapeHTML(false)
	err := enc.Encode([]SlackAttachment{attachment})
	if err != nil {
		return err
	}

	client := &http.Client{}

	payload := url.Values{}
	payload.Set("token", SLACK_BOT_OAUTH_TOKEN)
	payload.Set("channel", message.UserID)
	payload.Set("text", alert)
	payload.Set("attachments", attachmentData.String())

	endpoint := "https://slack.com/api/chat.postMessage"
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Printf("Pushing notifiction: %s\n", endpoint)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return fmt.Errorf("POST notification Failed %v", resp.StatusCode)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response SlackReponse
	fmt.Println("Reading POST notification request")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Casting POST notification request")
	json.Unmarshal(body, &response)
	fmt.Println(string(body))

	if response.Ok == false {
		return fmt.Errorf(response.Error)
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type BotMessage struct {
	UserID      string `json:"userId"`
	Message     string `json:"text"`
	Attachments string `json:"attachments"`
}

type RequestPermalink struct {
	Permalink string `json:"permalink"`
}

func notify(message *BotMessage) error {
	return pushNotification(message)
}

func pushNotification(event *BotMessage) error {
	client := getClient()

	payload := url.Values{}
	payload.Set("token", SLACK_BOT_OAUTH_TOKEN)
	payload.Set("channel", event.UserID)
	payload.Set("text", event.Message)
	payload.Set("attachments", event.Attachments)

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	client.CloseIdleConnections()
	json.Unmarshal(body, &response)
	log.Printf("Notifying Slack user id: %v\n", string(event.UserID))

	if response.Ok == false {
		return fmt.Errorf(response.Error)
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SlackPayload struct {
	Type      string     `json:"type"`
	Event     SlackEvent `json:"event,omitempty"`
	Challenge string     `json:"challenge,omitempty"`
}

type SlackAttachment struct {
	Id       int16  `json:"id,omitempty"`
	Color    string `json:"color,omitempty"`
	Fallback string `json:"fallback,omitempty"`
	Text     string `json:"text,omitempty"`
	Footer   string `json:"footer,omitempty"`
}

type SlackBlock struct {
	Type    string `json:"rich_text,omitempty"`
	BlockID string `json:"block_id,omitempty"`
}

type SlackEvent struct {
	Type        string            `json:"type"`
	Subtype     string            `json:"subtype"`
	Text        string            `json:"text"`
	TS          string            `json:"ts,omitempty"`
	BotID       string            `json:"bot_id"`
	Blocks      []SlackBlock      `json:"blocks,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
	Channel     string            `json:"channel"`
	EventTS     string            `json:"event_ts,omitempty"`
	ChannelType string            `json:"channel_type"`
}

type SlackUser struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	User  struct {
		Id      string `json:"id"`
		Profile struct {
			Email string `json:"email"`
		} `json:"profile"`
	} `json:"user,omitempty"`
}

type SlackReponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func fetchSlackUser(userEmail string) (*SlackUser, error) {
	client := &http.Client{}
	var user SlackUser

	endpoint := fmt.Sprintf("https://slack.com/api/users.lookupByEmail?email=%s", userEmail)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+SLACK_BOT_OAUTH_TOKEN)

	fmt.Printf("Fetching user: %s\n", endpoint)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET user Failed %v", resp.StatusCode)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("Reading GET user request")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Casting GET user request")
	json.Unmarshal(body, &user)
	fmt.Println(string(body))

	if user.Ok == false {
		return nil, fmt.Errorf("Did not find a Slack user matching the email, exception %v", user.Error)
	}

	return &user, nil
}

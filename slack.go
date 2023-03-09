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
	TS          string            `json:"ts"`
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
	client := getClient()

	var user SlackUser

	endpoint := fmt.Sprintf("https://slack.com/api/users.lookupByEmail?email=%s", userEmail)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+SLACK_BOT_OAUTH_TOKEN)

	fmt.Printf("Bot is fetching Slack user profile: %s\n", endpoint)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET user Failed %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	client.CloseIdleConnections()
	json.Unmarshal(body, &user)

	if !user.Ok {
		return nil, fmt.Errorf("Did not find a Slack user matching the email, exception %v", user.Error)
	}

	fmt.Printf("Slack userID found ''%v'' for user ''%v'' (OK result %v)\n", user.User.Id, user.User.Profile.Email, user.Ok)

	return &user, nil
}

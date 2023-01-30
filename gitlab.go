package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type GitLabUser struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	State    bool   `json:"state"`
	Email    string `json:"email,omitempty"`
}

func formatGitLabUsernameTag(username string) (string, error) {
	removedAccents, err := deburr(strings.TrimSpace(username))
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ToLower(removedAccents), " ", "."), nil
}

func fetchGitLabUser(username string) (*GitLabUser, error) {
	return nil, fmt.Errorf("Not implemented. Private admin fetch to access email. Needs to authenticate and retrive an oauth token using an admin access token gitlab: https://docs.gitlab.com/ee/api/users.html#for-administrators")
}

func fetchBasicGitLabUser(username string) (*GitLabUser, error) {
	client := &http.Client{}

	endpoint := fmt.Sprintf("https://gitlab.com/api/v4/users?username=%s", username)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

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

	var users []GitLabUser
	fmt.Println("Casting GET user request")
	json.Unmarshal(body, &users)

	if len(users) < 1 {
		return nil, fmt.Errorf("No user found on GitLab with that username tag.")
	}

	return &users[0], nil
}

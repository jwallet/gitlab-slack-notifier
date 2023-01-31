package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type GitLabUser struct {
	Id       int32  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	State    bool   `json:"state"`
	Email    string `json:"email,omitempty"`
}

const defaultQueryParams = "active=true&blocked=false&without_project_bots=true"

func fetchGitLabUser(username string) (*GitLabUser, error) {
	return nil, fmt.Errorf(`Not implemented!
		To retrive a user private email you need to be an admin, so a GitLab staff member for the cloud service.
		You can retrieve the user public email though by using GET:User/:id.
		If you have GitLab self-hosted you can fetch any user private email.
	    Set your personal or project access token when using GET:Users
		https://docs.gitlab.com/ee/api/users.html#for-administrators`)
}

func fetchBasicGitLabUser(username string) (*GitLabUser, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	endpoint := fmt.Sprintf("https://gitlab.com/api/v4/users?%s&username=%s", defaultQueryParams, username)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Fetching user: %s\n", endpoint)

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

	var users []GitLabUser
	json.Unmarshal(body, &users)
	client.CloseIdleConnections()

	if len(users) < 1 {
		return nil, fmt.Errorf("No user found on GitLab with that username tag.")
	}

	user := &users[0]

	fmt.Printf("GitLab user fullname: %v\n", user.Name)

	return user, nil
}

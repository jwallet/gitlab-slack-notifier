package main

import "strings"

func formatFullnameToUserEmail(username string) (string, error) {
	removedAccents, err := deburr(strings.TrimSpace(username))
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ToLower(removedAccents), " ", USER_EMAIL_SPACE_REPLACER), nil
}

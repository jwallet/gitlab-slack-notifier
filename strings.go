package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func deburr(source string) (string, error) {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(transformer, source)
	if err != nil {
		fmt.Printf("Error normalizing username: %v\n", err)
		return source, err
	}
	return output, nil
}

func getMergeRequestLinkToComment(text string) string {
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

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func formatFullnameToUserEmail(username string) (string, error) {
	removedAccents, err := deburr(strings.TrimSpace(username))
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.ToLower(removedAccents), " ", USER_EMAIL_SPACE_REPLACER), nil
}

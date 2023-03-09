package main

import (
	"regexp"
	"strings"
)

type EventTextAlias string

var (
	author   EventTextAlias
	username EventTextAlias
	link     EventTextAlias
	repo     EventTextAlias
	mr       EventTextAlias
)

var (
	selectorUsername           = regexp.MustCompile(`@[^\s]+`)
	selectorAllCaps            = regexp.MustCompile("([a-z0-9])([A-Z])")
	selectorAuthor             = regexp.MustCompile(`.*?(?P<author>{{author\|?(?P<authorFallback>.*?)}}).*?`)
	selectorRepository         = regexp.MustCompile(`.*?(?P<repository>{{repository\|?(?P<repositoryFallback>.*?)}}).*?`)
	selectorMergeRequest       = regexp.MustCompile(`.*?(?P<mergeRequest>{{mergeRequest\|?(?P<mergeRequestFallback>.*?)}}).*?`)
	selectorAliasesFromComment = regexp.MustCompile(`^(?P<author>.*?) \((?P<username>.*?)\) \<(?P<link>https://gitlab\.com/.*?)\|commented on merge request !\d+\> in <https://gitlab\.com/.*?\|.*?\s.*?\s/\s(?P<repo>.*?)>: \*(?P<mr>.*?)\*$`)
)

func getAliasFromEventText(text string, alias EventTextAlias) string {
	matches := selectorAliasesFromComment.FindStringSubmatch(text)
	result := make(map[string]string)
	for i, name := range selectorAliasesFromComment.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}
	return result[string(alias)]
}

func getUsernameTagsToDotCase(usernames []string) []string {
	result := make([]string, 0)
	for _, username := range usernames {
		result = append(result, strings.ToLower(selectorAllCaps.ReplaceAllString(username, "${1}.${2}")))
	}
	return result
}

func getAllUsernameTags(comment string) []string {
	usernames := selectorUsername.FindAllString(comment, -1)
	for i, username := range usernames {
		usernames[i] = strings.Replace(username, "@", "", 1)
	}
	return distinct(usernames)
}

func getGreatings(event Event) string {
	var text = SLACK_BOT_NOTIFICATION_GREATINGS
	author := selectorAuthor.FindStringSubmatch(text)
	repository := selectorRepository.FindStringSubmatch(text)
	mergeRequest := selectorMergeRequest.FindStringSubmatch(text)
	if len(author) > 1 {
		fallback := ternary(len(author) == 3, author[2], "")
		value := ternary(len(event.author) > 0, event.author, fallback)
		text = strings.Replace(text, author[1], value, 1)
	}
	if len(repository) > 1 {
		fallback := ternary(len(repository) == 3, repository[2], "")
		value := ternary(len(event.repository) > 0, event.repository, fallback)
		text = strings.Replace(text, repository[1], value, 1)
	}
	if len(mergeRequest) > 1 {
		fallback := ternary(len(mergeRequest) == 3, mergeRequest[2], "")
		value := ternary(len(event.repository) > 0, event.mergeRequest, fallback)
		text = strings.Replace(text, mergeRequest[1], value, 1)
	}
	return strings.Replace(text, "\\n", "\n", -1)
}

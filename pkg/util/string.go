package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	linkRegex        = regexp.MustCompile(`^\[[^][]+]\((https?://[^()]+)\)$`)
	headCountRegex   = regexp.MustCompile(`\(([[:digit:]]+)\/?([[:digit:]]+)?\)$`)
	footerRegex      = regexp.MustCompile(`^Created by (.+)$`)
	inputSelectRegex = regexp.MustCompile(`^\d+(?: \d+)*$`)
)

// GetLinkFromDeleteDescription extracts a message link from a delete handler
func GetLinkFromDeleteDescription(description string) (string, error) {
	result := linkRegex.FindStringSubmatch(description)
	if len(result) != 2 {
		return "", fmt.Errorf("invalid event description")
	}
	return result[1], nil
}

// GetIDsFromDeleteLink gets IDs in the message link from a delete handler
func GetIDsFromDeleteLink(link string) (guildID string, channelID string, messageID string, err error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", "", "", err
	}

	// Expect [channels, guildID, eventID, messageID]
	pathSlice := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(pathSlice) != 4 {
		return "", "", "", fmt.Errorf("cannot parse link: %s", link)
	}
	l := len(pathSlice)
	return pathSlice[l-3], pathSlice[l-2], pathSlice[l-1], err
}

func ParseFieldHeadCount(fieldName string) (count int, limit int, err error) {
	result := headCountRegex.FindStringSubmatch(fieldName)
	l := len(result)
	if l == 0 {
		return 0, 0, nil
	}

	if l == 3 {
		count, err = strconv.Atoi(result[1])
		if err != nil {
			return 0, 0, err
		}
		if result[2] == "" {
			return count, 0, nil
		}
		limit, err = strconv.Atoi(result[2])
		if err != nil {
			return 0, 0, err
		}
		return count, limit, nil
	}
	return 0, 0, fmt.Errorf("cannot parse field name for count")
}

func AddUserToField(fieldValue, userName string) (string, error) {
	if fieldValue == "-" {
		return fmt.Sprintf("> %s", userName), nil
	}
	names := strings.Split(fieldValue, "\n")
	for _, n := range names {
		if strings.Trim(n, "> ") == userName {
			return fieldValue, nil
		}
	}
	names = append(names, fmt.Sprintf("> %s", userName))
	return strings.Join(names, "\n"), nil
}

func RemoveUser(names []string, userName string) []string {
	var i int
	for _, n := range names {
		if n != userName {
			names[i] = n
			i++
		}
	}
	names = names[:i]
	return names
}

func ContainsUser(names []string, userName string) bool {
	for _, n := range names {
		if n == userName {
			return true
		}
	}
	return false
}

func GetUsersFromValues(fieldValue string) []string {
	result := make([]string, 0)
	if fieldValue == "-" {
		return result
	}
	for _, n := range strings.Split(fieldValue, "\n") {
		result = append(result, strings.Trim(n, "> "))
	}
	return result
}

func NameListToValues(list []string) string {
	if len(list) == 0 {
		return "-"
	}
	result := make([]string, 0)
	for _, n := range list {
		result = append(result, fmt.Sprintf("> %s", n))
	}
	return strings.Join(result, "\n")
}

func GetUserFromFooter(footText string) string {
	match := footerRegex.FindStringSubmatch(footText)
	if len(match) != 2 {
		return ""
	}
	return match[1]
}

func IsInputOption(input string) bool {
	result := inputSelectRegex.FindStringSubmatch(input)
	if len(result) != 1 {
		return false
	}
	return true
}

func MergeSlices(slices ...[]string) []string {
	result := make([]string, 0)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

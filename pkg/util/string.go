package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	linkRegex      = regexp.MustCompile(`^\[[^][]+]\((https?://[^()]+)\)$`)
	headCountRegex = regexp.MustCompile(`\(([[:digit:]]+)\/?([[:digit:]]+)?\)$`)
)

// GetLinkFromDeleteDescription extracts a message link from a delete handler
func GetLinkFromDeleteDescription(description string) (string, error) {
	result := linkRegex.FindStringSubmatch(description)
	if len(result) != 2 {
		return "", fmt.Errorf("invalid event description")
	}
	fmt.Println(description, result[1])
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
		return -1, -1, nil
	}

	if l == 3 {
		count, err = strconv.Atoi(result[1])
		if err != nil {
			return -1, -1, err
		}
		if result[2] == "" {
			return count, -1, nil
		}
		limit, err = strconv.Atoi(result[2])
		if err != nil {
			return -1, -1, err
		}
		return count, limit, nil
	}
	return -1, -1, fmt.Errorf("cannot parse field name for count")
}

func IncrementFieldCounter(fieldName string) (string, error) {
	var result string
	match := headCountRegex.FindStringSubmatch(fieldName)
	l := len(match)

	if l == 0 {
		result = fieldName + " (1)"
		return result, nil
	}

	if l == 3 {
		c, err := strconv.Atoi(match[1])
		if err != nil {
			return result, err
		}
		if match[2] == "" {
			result = headCountRegex.ReplaceAllString(fieldName, "") + fmt.Sprintf("(%d)", c+1)
			return result, nil
		}
		limit, err := strconv.Atoi(match[2])
		if err != nil {
			return result, err
		}
		if c >= limit {
			return result, fmt.Errorf("cannot exceed field limit")
		}
		result = headCountRegex.ReplaceAllString(fieldName, "") + fmt.Sprintf("(%d/%d)", c+1, limit)
		return result, nil
	}
	return result, fmt.Errorf("cannot parse field name")
}

func DecrementFieldCounter(fieldName string) (string, error) {
	var result string
	match := headCountRegex.FindStringSubmatch(fieldName)
	l := len(match)

	if l == 0 {
		return fieldName, nil
	}

	if l == 3 {
		c, err := strconv.Atoi(match[1])
		if err != nil {
			return result, err
		}
		c--
		if match[2] == "" {
			if c == 0 {
				result = headCountRegex.ReplaceAllString(fieldName, "")
				return strings.Trim(result, " "), nil
			}
			result = headCountRegex.ReplaceAllString(fieldName, "") + fmt.Sprintf("(%d)", c)
			return result, nil
		}
		limit, err := strconv.Atoi(match[2])
		if err != nil {
			return result, err
		}
		if c < 0 {
			return result, fmt.Errorf("cannot be negative")
		}
		result = headCountRegex.ReplaceAllString(fieldName, "") + fmt.Sprintf("(%d/%d)", c, limit)
		return result, nil
	}
	return result, fmt.Errorf("cannot parse field name")
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

func RemoveUserFromField(fieldValue, userName string) (string, error) {
	if fieldValue == "-" {
		return fieldValue, nil
	}
	names := strings.Split(fieldValue, "\n")
	var i int
	for _, n := range names {
		if strings.Trim(n, "> ") != userName {
			names[i] = n
			i++
		}
	}
	names = names[:i]
	if len(names) == 0 {
		return "-", nil
	}
	return strings.Join(names, "\n"), nil
}

func ContainsUserInField(fieldValue, userName string) bool {
	if fieldValue == "-" {
		return false
	}
	for _, n := range strings.Split(fieldValue, "\n") {
		if strings.Trim(n, "> ") == userName {
			return true
		}
	}
	return false
}

func GetUsersFromValues(fieldValue string) []string {
	result := []string{}
	if fieldValue == "-" {
		return result
	}
	for _, n := range strings.Split(fieldValue, "\n") {
		result = append(result, strings.Trim(n, "> "))
	}
	return result
}

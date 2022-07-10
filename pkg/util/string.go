package util

import (
	"encoding/base64"
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

	LineFeed = "ლ(´ڡ`ლ)"
)

// GetLinkFromDeleteDescription extracts a message link from a delete handler
func GetLinkFromDeleteDescription(description string) (string, error) {
	result := linkRegex.FindStringSubmatch(description)
	if len(result) != 2 {
		return "", fmt.Errorf("invalid event description")
	}
	return result[1], nil
}

// GetIDsFromDiscordLink gets IDs in the message link from a delete handler
func GetIDsFromDiscordLink(link string) (guildID string, channelID string, messageID string, err error) {
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

// GetDiscordLinkFromCalendarDescription gets the Discord link from a Google calendar event description
func GetDiscordLinkFromCalendarDescription(description string) (string, error) {
	result := strings.Split(description, LineFeed)
	if len(result) != 2 {
		return "", fmt.Errorf("unexpected description format")
	}
	discordLink := strings.Trim(result[len(result)-1], "\n")
	return discordLink, nil
}

// ParseFieldHeadCount extracts the attendees and limits of an event from an embed message field
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

// AddUserToField adds a user to a field value
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

// IsInputOption validates input for one or more choices
func IsInputOption(input string) bool {
	result := inputSelectRegex.FindStringSubmatch(input)
	if len(result) != 1 {
		return false
	}
	return true
}

// EncodeToGoogleCalendarBase64 converts an event ID and calendar ID to a base64 encoded string
// See https://stackoverflow.com/questions/53928044/how-do-i-construct-a-link-to-a-google-calendar-event
func EncodeToGoogleCalendarBase64(eventID, calendarID string) string {
	calendarID = strings.Replace(calendarID, "@group.calendar.google.com", "@g", 1)
	result := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s %s", eventID, calendarID)))
	return strings.Trim(result, "=")
}

// ParseEventID parses the event ID from a calendar event link
func ParseEventID(link string) (string, error) {
	result := linkRegex.FindStringSubmatch(link)
	if len(result) != 2 {
		return "", fmt.Errorf("invalid link")
	}
	u, err := url.Parse(result[1])
	if err != nil {
		return "", err
	}
	q := u.Query()
	id := q.Get("eid")
	if id == "" {
		return "", fmt.Errorf("cannot find eid")
	}
	return id, nil
}

// DecodeToGoogleEventID converts a base64 encoded ID to an event ID and calendar ID
func DecodeToGoogleEventID(encodedID string) (eid, calendarID string, err error) {
	result, err := base64.RawStdEncoding.DecodeString(encodedID)
	if err != nil {
		return "", "", err
	}
	ids := strings.Split(string(result), " ")
	if len(ids) != 2 {
		return "", "", fmt.Errorf("unknown encoded ID")
	}
	cid := strings.Replace(ids[1], "@g", "@group.calendar.google.com", 1)
	return ids[0], cid, nil
}

package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func PrintTime(start, end time.Time) string {
	base := fmt.Sprintf("<t:%d:F>", start.Unix())
	relative := fmt.Sprintf("\n🕔<t:%d:R>", start.Unix())
	if end.IsZero() {
		return base + relative
	}

	if end.Sub(start) < time.Hour*24 {
		base = base + fmt.Sprintf(" - <t:%d:t>", end.Unix())
	} else {
		base = base + fmt.Sprintf(" - <t:%d:F>", end.Unix())
	}
	return base + relative
}

// PrintHumanReadableTime prints time to a local format
func PrintHumanReadableTime(start, end time.Time) string {
	if end.Equal(time.Time{}) || end.Before(start) {
		return ""
	}
	return fmt.Sprintf("%s", end.Sub(start).Round(time.Minute))
}

func PrintBlockValues(input string) string {
	if input == "" {
		return "```-```"
	}
	return fmt.Sprintf("```%s```", input)
}

// PrintAddGoogleCalendarLink returns a Google calendar link with event params
func PrintAddGoogleCalendarLink(title, description string, startTime, endTime time.Time) string {
	if endTime.IsZero() {
		endTime = startTime
	}

	s, e := startTime.UTC().Format(GoogleCalendarTimeFormat), endTime.UTC().Format(GoogleCalendarTimeFormat)

	u, _ := url.Parse("https://www.google.com/calendar/event?action=TEMPLATE&text=&details=&location=")
	q := u.Query()
	q.Set("text", title)
	q.Set("details", description)
	u.RawQuery = q.Encode()

	// TODO: Encode multiple dates rather than constructing
	re := regexp.MustCompile("[[:punct:]]")
	s, e = re.ReplaceAllString(s, ""), re.ReplaceAllString(e, "")

	link := u.String() + "&dates=" + s + "/" + e

	return fmt.Sprintf("[Add to Google Calendar](%s)", link)
}

// PrintGoogleCalendarEventLink prints a link to a Google calendar event using base64 encoding
func PrintGoogleCalendarEventLink(ID string) string {
	u, _ := url.Parse("https://www.google.com/calendar/event?eid=")
	q := u.Query()
	q.Set("eid", ID)
	u.RawQuery = q.Encode()
	return fmt.Sprintf("[View in Calendar](%s)", u.String())
}

func PrintGoogleCalendarDescription(description string, discordLink string) string {
	result := description + fmt.Sprintf("\n%s\n%s", LineFeed, discordLink)
	return result
}

// GetTimesFromLink gets start and end times from a markdown calendar link via query params
func GetTimesFromLink(link string) (start, end time.Time, err error) {
	result := linkRegex.FindStringSubmatch(link)
	if len(result) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid value")
	}
	u, err := url.Parse(result[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	q := u.Query()
	times := strings.Split(q.Get("dates"), "/")
	if len(times) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("cannot parse time")
	}
	start, err = time.Parse(GoogleCalendarTimeFormat, times[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err = time.Parse(GoogleCalendarTimeFormat, times[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end, nil
}

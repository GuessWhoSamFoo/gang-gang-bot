package util

import (
	"fmt"
	"net/url"
	"regexp"
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

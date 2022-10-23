package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTimesFromLink(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name:          "valid link",
			input:         "[Add to Google Calendar](https://www.google.com/calendar/event?action=TEMPLATE&details=lol&location=&text=my+new+title&dates=20240101T080000Z/20240101T100000Z)",
			expectedStart: time.Date(2024, time.January, 1, 8, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2024, time.January, 1, 10, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			start, end, err := GetTimesFromLink(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStart, start)
			assert.Equal(t, tc.expectedEnd, end)
		})
	}
}

func TestPrintHumanReadableTime(t *testing.T) {
	cases := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
		expected  string
	}{
		{
			name:      "end before start",
			startTime: time.Now().Add(time.Minute),
			endTime:   time.Now(),
		},
		{
			name:      "end is zero time",
			startTime: time.Now(),
			endTime:   time.Time{},
		},
		{
			name:      "valid",
			startTime: time.Now(),
			endTime:   time.Now().Add(time.Minute),
			expected:  "1m0s",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PrintHumanReadableTime(tc.startTime, tc.endTime)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestPrintAddGoogleCalendarLink(t *testing.T) {
	fake := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name      string
		title     string
		desc      string
		startTime time.Time
		endTime   time.Time
		expected  string
	}{
		{
			name:      "valid",
			title:     "event title",
			desc:      "event description",
			startTime: fake,
			endTime:   fake.Add(time.Minute),
			expected:  "[Add to Google Calendar](https://www.google.com/calendar/event?action=TEMPLATE&details=event+description&location=&text=event+title&dates=20220101T000000Z/20220101T000100Z)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PrintAddGoogleCalendarLink(tc.title, tc.desc, tc.startTime, tc.endTime)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestPrintTime(t *testing.T) {
	testTime := time.Date(2022, 01, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected string
	}{
		{
			name:     "end after start",
			start:    testTime,
			end:      testTime.Add(time.Minute),
			expected: "<t:1640995200:F> - <t:1640995260:t>\n🕔<t:1640995200:R>",
		},
		{
			name:     "end before start",
			start:    testTime,
			end:      time.Time{},
			expected: "<t:1640995200:F>\n🕔<t:1640995200:R>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PrintTime(tc.start, tc.end)
			assert.Equal(t, tc.expected, got)
		})
	}
}

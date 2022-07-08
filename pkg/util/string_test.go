package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLinkFromDeleteDescription(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
		isErr    bool
	}{
		{
			name:     "valid link",
			input:    "[testing](https://discord.com/channels/530946869023604746/530949686962421770/755612707146825839)",
			expected: "https://discord.com/channels/530946869023604746/530949686962421770/755612707146825839",
		},
		{
			name:     "calendar link",
			input:    "[Add to Google Calendar](https://www.google.com/calendar/event?action=TEMPLATE&text=12321&details=&location=&dates=20220618T040000Z/20220618T050000Z)",
			expected: "https://www.google.com/calendar/event?action=TEMPLATE&text=12321&details=&location=&dates=20220618T040000Z/20220618T050000Z",
		},
		{
			name:  "invalid link",
			input: "https://google.com",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetLinkFromDeleteDescription(tc.input)
			if err == nil && tc.isErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestParseFieldHeadCount(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expectedCount int
		expectedLimit int
	}{
		{
			name:          "no limit with count",
			input:         "✅ Accepted (17)",
			expectedCount: 17,
			expectedLimit: 0,
		},
		{
			name:          "limit",
			input:         "✅ Accepted (1/10)",
			expectedCount: 1,
			expectedLimit: 10,
		},
		{
			name:          "limit zero",
			input:         "✅ Accepted (0/1)",
			expectedCount: 0,
			expectedLimit: 1,
		},
		{
			name:          "limit declined",
			input:         "❌ Declined (12/25)",
			expectedCount: 12,
			expectedLimit: 25,
		},
		{
			name:          "no limit without count",
			input:         "❔ Tentative",
			expectedCount: 0,
			expectedLimit: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotCount, gotLimit, err := ParseFieldHeadCount(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCount, gotCount)
			assert.Equal(t, tc.expectedLimit, gotLimit)
		})
	}
}

func TestAddUserToField(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		userName string
		expected string
	}{
		{
			name:     "empty",
			value:    "-",
			userName: "foo",
			expected: "> foo",
		},
		{
			name:     "add user",
			value:    "> foo",
			userName: "bar",
			expected: "> foo\n> bar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := AddUserToField(tc.value, tc.userName)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestRemoveUser(t *testing.T) {
	cases := []struct {
		name     string
		names    []string
		userName string
		expected []string
	}{
		{
			name:     "empty",
			names:    []string{},
			userName: "foo",
			expected: []string{},
		},
		{
			name:     "should remove",
			names:    []string{"hello", "world"},
			userName: "world",
			expected: []string{"hello"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := RemoveUser(tc.names, tc.userName)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestGetUsersFromValues(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty",
			input:    "-",
			expected: []string{},
		},
		{
			name:     "one",
			input:    "> foo",
			expected: []string{"foo"},
		},
		{
			name:     "many",
			input:    "> foo\n> bar\n> baz",
			expected: []string{"foo", "bar", "baz"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, GetUsersFromValues(tc.input))
		})
	}
}

func TestNameListToValues(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty",
			input:    []string{},
			expected: "-",
		},
		{
			name:     "one name",
			input:    []string{"foo"},
			expected: "> foo",
		},
		{
			name:     "many",
			input:    []string{"foo", "bar", "baz"},
			expected: "> foo\n> bar\n> baz",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, NameListToValues(tc.input))
		})
	}
}

func TestGetUserFromFooter(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "has match",
			input:    "Created by a funky dude",
			expected: "a funky dude",
		},
		{
			name:     "invalid input",
			input:    "not a real footer",
			expected: "",
		},
		{
			name:     "non english chars",
			input:    "Created by Weirdo ハロー・ワールド",
			expected: "Weirdo ハロー・ワールド",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, GetUserFromFooter(tc.input))
		})
	}
}

func TestIsInputOption(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "is valid",
			input:    "1 3 5",
			expected: true,
		},
		{
			name:     "invalid",
			input:    "testing",
			expected: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsInputOption(tc.input))
		})
	}
}

func TestEncodeToGoogleCalendarBase64(t *testing.T) {
	cases := []struct {
		name       string
		eventID    string
		calendarID string
		expected   string
	}{
		{
			name:       "base",
			eventID:    "2vpae347k2a70gbhnmgmve9fll",
			calendarID: "s8lsr7ohbqi5u52b89vmvmq1ak@group.calendar.google.com",
			expected:   "MnZwYWUzNDdrMmE3MGdiaG5tZ212ZTlmbGwgczhsc3I3b2hicWk1dTUyYjg5dm12bXExYWtAZw",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncodeToGoogleCalendarBase64(tc.eventID, tc.calendarID)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestDecodeToGoogleEventID(t *testing.T) {
	cases := []struct {
		name               string
		input              string
		expectedEventID    string
		expectedCalendarID string
	}{
		{
			name:               "base",
			input:              "MnZwYWUzNDdrMmE3MGdiaG5tZ212ZTlmbGwgczhsc3I3b2hicWk1dTUyYjg5dm12bXExYWtAZw",
			expectedEventID:    "2vpae347k2a70gbhnmgmve9fll",
			expectedCalendarID: "s8lsr7ohbqi5u52b89vmvmq1ak@group.calendar.google.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			eventID, calendarID, err := DecodeToGoogleEventID(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedEventID, eventID)
			assert.Equal(t, tc.expectedCalendarID, calendarID)
		})
	}
}

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

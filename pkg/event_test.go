package pkg

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func Test_getEventFromMessage(t *testing.T) {
	cases := []struct {
		name     string
		input    *discordgo.Message
		expected *Event
	}{
		{
			name: "simple",
			input: &discordgo.Message{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "title",
						Description: "desc",
						Color:       1234,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  acceptedBase,
								Value: "-",
							},
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Created by test",
						},
					},
				},
			},
			expected: &Event{
				Title:         "title",
				Description:   "desc",
				Limit:         -1,
				Accepted:      0,
				AcceptedNames: []string{},
				Owner:         "test",
				Color:         1234,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetEventFromMessage(tc.input)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("expected: %v, got: %v", tc.expected, got)
			}
		})
	}
}

func Test_convertEventToMessageEmbed(t *testing.T) {
	cases := []struct {
		name     string
		input    *Event
		expected *discordgo.MessageEmbed
	}{
		{
			name: "base",
			input: &Event{
				Title:          "testing",
				Description:    "hello world",
				Limit:          -1,
				Accepted:       1,
				AcceptedNames:  []string{"user"},
				DeclinedNames:  []string{},
				TentativeNames: []string{},
				WaitlistNames:  []string{},
				Color:          1234,
				Owner:          "foo",
			},
			expected: &discordgo.MessageEmbed{
				Title:       "testing",
				Description: "hello world",
				Color:       1234,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Time",
						Value: util.PrintTime(time.Time{}, time.Time{}),
					},
					{
						Name:  "Links",
						Value: util.PrintAddGoogleCalendarLink("testing", "hello world", time.Time{}, time.Time{}),
					},
					{
						Name:   acceptedBase + " (1)",
						Value:  "> user",
						Inline: true,
					},
					{
						Name:   declinedBase,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   tentativeBase,
						Value:  "-",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Created by foo",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ConvertEventToMessageEmbed(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assert.Equal(t, got, tc.expected)
		})
	}
}

package pkg

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
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
							{
								Name:  string(WaitlistField),
								Value: "> test",
							},
							{
								Name:   "Calendar",
								Value:  util.PrintGoogleCalendarEventLink("MnZwYWUzNDdrMmE3MGdiaG5tZ212ZTlmbGwgczhsc3I3b2hicWk1dTUyYjg5dm12bXExYWtAZw"),
								Inline: true,
							},
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Created by test",
						},
					},
				},
			},
			expected: &Event{
				Title:       "title",
				Description: "desc",
				RoleGroup: &RoleGroup{
					Roles: []*Role{
						{
							Icon:      AcceptedIcon,
							FieldName: AcceptedField,
							Users:     []string{},
						},
					},
					Waitlist: map[FieldType]*Role{
						AcceptedField: {
							Icon:      "",
							FieldName: WaitlistField,
							Users:     []string{"test"},
							Count:     1,
						},
					},
				},
				Owner: "test",
				Color: 1234,
				ID:    "MnZwYWUzNDdrMmE3MGdiaG5tZ212ZTlmbGwgczhsc3I3b2hicWk1dTUyYjg5dm12bXExYWtAZw",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetEventFromMessage(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_convertEventToMessageEmbed(t *testing.T) {
	rg := NewDefaultRoleGroup()
	err := rg.ToggleRole(AcceptedField, "user")
	assert.NoError(t, err)
	cases := []struct {
		name     string
		input    *Event
		expected *discordgo.MessageEmbed
	}{
		{
			name: "base",
			input: &Event{
				Title:       "testing",
				Description: "hello world",
				RoleGroup:   rg,
				Color:       1234,
				Owner:       "foo",
				ID:          "MnZwYWUzNDdrMmE3MGdiaG5tZ212ZTlmbGwgczhsc3I3b2hicWk1dTUyYjg5dm12bXExYWtAZw",
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
						Name:   "Links",
						Value:  util.PrintAddGoogleCalendarLink("testing", "hello world", time.Time{}, time.Time{}),
						Inline: true,
					},
					{
						Name:  "Calendar",
						Value: util.PrintGoogleCalendarEventLink("MnZwYWUzNDdrMmE3MGdiaG5tZ212ZTlmbGwgczhsc3I3b2hicWk1dTUyYjg5dm12bXExYWtAZw"),
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
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestEvent_SetDuration(t *testing.T) {
	now := time.Now()

	cases := []struct {
		input    string
		expected time.Time
		isErr    bool
	}{
		{
			input:    "1 hour",
			expected: now.Add(time.Hour),
		},
		{
			input:    "1 hour and 30 minutes",
			expected: now.Add(time.Minute * 90),
		},
		{
			input:    "1 hour 30 minutes",
			expected: now.Add(time.Minute * 90),
		},
		{
			input: "1hr 30 minutes",
			isErr: true,
		},
		{
			input:    "invalid",
			expected: now,
		},
		// TODO: Find a more configurable parser
		//{
		//	input:    "1h30m",
		//	expected: now.Add(time.Minute * 90),
		//},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			event := &Event{Start: now}
			err := event.SetDuration(tc.input)
			if !tc.isErr {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, event.End)
		})
	}
}

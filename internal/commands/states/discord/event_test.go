package discord

import (
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEvent_AddTitle(t *testing.T) {
	e := Event{}
	expected := "test"
	e.AddTitle(expected)
	assert.Equal(t, expected, e.Title)
}

func TestFromFSMToEvent(t *testing.T) {
	cases := []struct {
		name     string
		key      MetadataKey
		val      interface{}
		expected *Event
	}{
		{
			name:     "title",
			key:      Title,
			val:      "title",
			expected: &Event{Title: "title"},
		},
		{
			name:     "description",
			key:      Description,
			val:      "desc",
			expected: &Event{Description: "desc"},
		},
		{
			name:     "location",
			key:      Location,
			val:      "Seattle",
			expected: &Event{Location: "Seattle"},
		},
		{
			name:     "start time",
			key:      StartTime,
			val:      time.Time{},
			expected: &Event{Start: time.Time{}},
		},
		{
			name:     "duration",
			key:      Duration,
			val:      time.Time{},
			expected: &Event{End: time.Time{}},
		},
		{
			name: "attendee",
			key:  Attendee,
			val:  role.NewDefaultRoleGroup(),
			expected: &Event{
				RoleGroup: role.NewDefaultRoleGroup(),
			},
		},
		{
			name:     "owner",
			key:      Owner,
			val:      "owner",
			expected: &Event{Owner: "owner"},
		},
		{
			name:     "color",
			key:      Color,
			val:      Purple,
			expected: &Event{Color: Purple},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := fsm.NewFSM("idle", fsm.Events{}, fsm.Callbacks{})
			f.SetMetadata(tc.key.String(), tc.val)

			got, err := FromFSMToEvent(f)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestGetEventFromMessage(t *testing.T) {
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
								Name:  AcceptedBase,
								Value: "-",
							},
							{
								Name:  string(role.WaitlistField),
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
				RoleGroup: &role.RoleGroup{
					Roles: []*role.Role{
						{
							Icon:      role.AcceptedIcon,
							FieldName: role.AcceptedField,
							Users:     []string{},
						},
					},
					Waitlist: map[role.FieldType]*role.Role{
						role.AcceptedField: {
							Icon:      "",
							FieldName: role.WaitlistField,
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

func TestConvertEventToMessageEmbed(t *testing.T) {
	rg := role.NewDefaultRoleGroup()
	err := rg.ToggleRole(role.AcceptedField, "user")
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
						Name:   AcceptedBase + " (1)",
						Value:  "> user",
						Inline: true,
					},
					{
						Name:   DeclinedBase,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   TentativeBase,
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
			input:    "invalid",
			expected: now,
		},
		// TODO: Find a more configurable parser
		//{
		//	input:    "1h30m",
		//	expected: now.Add(time.Minute * 90),
		//},
		//{
		//	input: "1hr 30 minutes",
		//	isErr: true,
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

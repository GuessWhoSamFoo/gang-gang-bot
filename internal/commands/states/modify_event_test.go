package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/bwmarrin/discordgo"
	"github.com/ewohltman/discordgo-mock/mockconstants"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewModifyEventState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewModifyEventState(*opts)
	assert.NotNil(t, s)
}

func TestModifyEventState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	event := discord.Event{
		Title:       "event",
		Description: "description",
		Location:    "Seattle",
		Start:       time.Time{},
		End:         time.Time{},
		RoleGroup:   role.NewDefaultRoleGroup(),
		Owner:       mockconstants.TestUser,
		Color:       discord.Purple,
		ID:          "id",
		DiscordLink: "example.com",
	}
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "add title",
			input:    "1",
			expected: AddTitle.String(),
		},
		{
			name:     "add description",
			input:    "2",
			expected: AddDescription.String(),
		},
		{
			name:     "set date",
			input:    "3",
			expected: SetDate.String(),
		},
		{
			name:     "set duration",
			input:    "4",
			expected: SetDuration.String(),
		},
		{
			name:     "set location",
			input:    "5",
			expected: SetLocation.String(),
		},
		{
			name:     "invalid",
			input:    "invalid",
			expected: ModifyEventRetry.String(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewModifyEventState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: ModifyEvent.String(),
						Src:  []string{"idle"},
						Dst:  ModifyEvent.String(),
					},
					{
						Name: AddTitle.String(),
						Src:  []string{ModifyEvent.String()},
						Dst:  AddTitle.String(),
					},
					{
						Name: AddDescription.String(),
						Src:  []string{ModifyEvent.String()},
						Dst:  AddDescription.String(),
					},
					{
						Name: SetDate.String(),
						Src:  []string{ModifyEvent.String()},
						Dst:  SetDate.String(),
					},
					{
						Name: SetDuration.String(),
						Src:  []string{ModifyEvent.String()},
						Dst:  SetDuration.String(),
					},
					{
						Name: SetLocation.String(),
						Src:  []string{ModifyEvent.String()},
						Dst:  SetLocation.String(),
					},
					{
						Name: ModifyEventRetry.String(),
						Src:  []string{ModifyEvent.String()},
						Dst:  ModifyEventRetry.String(),
					},
				},
				fsm.Callbacks{
					ModifyEvent.String(): s.OnState,
				},
			)
			f.SetMetadata(discord.EventObject.String(), event)

			go func() {
				s.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					s.input <- tc.input
				}
				s.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), ModifyEvent.String())
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

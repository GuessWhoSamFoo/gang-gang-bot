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

func TestNewRemoveResponseState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewRemoveResponseState(*opts)
	assert.NotNil(t, s)
}

func TestRemoveResponseState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	cases := []struct {
		name     string
		input    string
		expected string
		isErr    bool
	}{
		{
			name:     "remove users",
			input:    "1 2",
			expected: RemoveResponse.String(),
		},
		{
			name:     "invalid",
			input:    "invalid",
			expected: RemoveResponseRetry.String(),
		},
		{
			name:     "cancel",
			input:    "cancel",
			expected: Cancel.String(),
			isErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rg := role.NewDefaultRoleGroup()
			rg.Roles[0].Users = []string{mockconstants.TestUser, "leo"}

			event := discord.Event{
				Title:       "event",
				Description: "description",
				Location:    "Seattle",
				Start:       time.Time{},
				End:         time.Time{},
				RoleGroup:   rg,
				Owner:       mockconstants.TestUser,
				Color:       discord.Purple,
				ID:          "id",
				DiscordLink: "example.com",
			}

			s := NewRemoveResponseState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: RemoveResponse.String(),
						Src:  []string{"idle"},
						Dst:  RemoveResponse.String(),
					},
					{
						Name: RemoveResponseRetry.String(),
						Src:  []string{RemoveResponse.String()},
						Dst:  RemoveResponseRetry.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{RemoveResponse.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					RemoveResponse.String(): s.OnState,
				},
			)
			f.SetMetadata(discord.EventObject.String(), event)

			go func() {
				s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					s.inputHandler.inputChan <- tc.input
				}
				s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), RemoveResponse.String())
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

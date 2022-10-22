package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/ewohltman/discordgo-mock/mockconstants"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewStartEditState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewStartEditState(*opts)
	assert.NotNil(t, s)
}

func TestStartEditState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	opts.InteractionCreate.Interaction.Member.Permissions = 8589934592
	opts.InteractionCreate.Interaction.Message = &discordgo.Message{
		Embeds: []*discordgo.MessageEmbed{
			{
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
						Name:   discord.AcceptedBase + " (1)",
						Value:  "> user",
						Inline: true,
					},
					{
						Name:   discord.DeclinedBase,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   discord.TentativeBase,
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
	cases := []struct {
		name     string
		input    string
		expected string
		isErr    bool
	}{
		{
			name:     "modify event",
			input:    "1",
			expected: ModifyEvent.String(),
		},
		{
			name:     "remove responses",
			input:    "2",
			expected: RemoveResponse.String(),
		},
		{
			name:     "add response",
			input:    "3",
			expected: AddResponse.String(),
		},
		{
			name:     "invalid",
			input:    "invalid",
			expected: StartEditRetry.String(),
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
			s := NewStartEditState(*opts)
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
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: StartEdit.String(),
						Src:  []string{"idle"},
						Dst:  StartEdit.String(),
					},
					{
						Name: ModifyEvent.String(),
						Src:  []string{StartEdit.String()},
						Dst:  ModifyEvent.String(),
					},
					{
						Name: RemoveResponse.String(),
						Src:  []string{StartEdit.String()},
						Dst:  RemoveResponse.String(),
					},
					{
						Name: AddResponse.String(),
						Src:  []string{StartEdit.String()},
						Dst:  AddResponse.String(),
					},
					{
						Name: StartEditRetry.String(),
						Src:  []string{StartEdit.String()},
						Dst:  StartEditRetry.String(),
					},
					{
						Name: ContinueEdit.String(),
						Src:  []string{StartEdit.String(), ModifyEvent.String()},
						Dst:  ContinueEdit.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{StartEdit.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					StartEdit.String(): s.OnState,
				},
			)
			f.SetMetadata(discord.EventObject.String(), event)
			s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
				s.inputHandler.inputChan <- tc.input
			}

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
				wg.Done()
			}()

			go func() {
				err = f.Event(context.TODO(), StartEdit.String())
				if tc.isErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				wg.Done()
			}()

			wg.Wait()
			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

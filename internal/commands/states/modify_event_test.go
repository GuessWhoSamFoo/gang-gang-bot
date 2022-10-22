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
	"sync"
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
			expected: ContinueEdit.String(),
		},
		{
			name:     "add description",
			input:    "2",
			expected: ContinueEdit.String(),
		},
		{
			name:     "set date",
			input:    "3",
			expected: ContinueEdit.String(),
		},
		{
			name:     "set duration",
			input:    "4",
			expected: ContinueEdit.String(),
		},
		{
			name:     "set location",
			input:    "5",
			expected: ContinueEdit.String(),
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
						Name: ContinueEdit.String(),
						Src: []string{
							ModifyEvent.String(),
							AddTitle.String(),
							AddDescription.String(),
							SetDate.String(),
							SetDuration.String(),
							SetLocation.String(),
						},
						Dst: ContinueEdit.String(),
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
				err = f.Event(context.TODO(), ModifyEvent.String())
				assert.NoError(t, err)
				wg.Done()
			}()
			wg.Wait()
			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

func TestNewModifyEventRetryState(t *testing.T) {
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
			expected: ContinueEdit.String(),
		},
		{
			name:     "add description",
			input:    "2",
			expected: ContinueEdit.String(),
		},
		{
			name:     "set date",
			input:    "3",
			expected: ContinueEdit.String(),
		},
		{
			name:     "set duration",
			input:    "4",
			expected: ContinueEdit.String(),
		},
		{
			name:     "set location",
			input:    "5",
			expected: ContinueEdit.String(),
		},
		{
			name:     "invalid",
			input:    "invalid",
			expected: SelfTransition.String(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewModifyEventRetryState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: ModifyEventRetry.String(),
						Src:  []string{"idle"},
						Dst:  ModifyEventRetry.String(),
					},
					{
						Name: ContinueEdit.String(),
						Src: []string{
							ModifyEventRetry.String(),
							AddTitle.String(),
							AddDescription.String(),
							SetDate.String(),
							SetDuration.String(),
							SetLocation.String(),
						},
						Dst: ContinueEdit.String(),
					},
					{
						Name: AddTitle.String(),
						Src:  []string{ModifyEventRetry.String()},
						Dst:  AddTitle.String(),
					},
					{
						Name: AddDescription.String(),
						Src:  []string{ModifyEventRetry.String()},
						Dst:  AddDescription.String(),
					},
					{
						Name: SetDate.String(),
						Src:  []string{ModifyEventRetry.String()},
						Dst:  SetDate.String(),
					},
					{
						Name: SetDuration.String(),
						Src:  []string{ModifyEventRetry.String()},
						Dst:  SetDuration.String(),
					},
					{
						Name: SetLocation.String(),
						Src:  []string{ModifyEventRetry.String()},
						Dst:  SetLocation.String(),
					},
					{
						Name: SelfTransition.String(),
						Src:  []string{ModifyEventRetry.String()},
						Dst:  SelfTransition.String(),
					},
				},
				fsm.Callbacks{
					ModifyEventRetry.String(): s.OnState,
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
				err = f.Event(context.TODO(), ModifyEventRetry.String())
				assert.NoError(t, err)
				wg.Done()
			}()
			wg.Wait()
			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

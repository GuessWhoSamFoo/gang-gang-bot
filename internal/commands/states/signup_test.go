package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/bwmarrin/discordgo"
	"github.com/ewohltman/discordgo-mock/mockconstants"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewSignUpState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewSignUpState(*opts)
	assert.NotNil(t, s)
}

func TestSignUpState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
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
		isErr    bool
	}{
		{
			name:     "add accepted",
			input:    "1",
			expected: SignUp.String(),
		},
		{
			name:     "add declined",
			input:    "2",
			expected: SignUp.String(),
		},
		{
			name:     "add tentative",
			input:    "3",
			expected: SignUp.String(),
		},
		{
			name:     "invalid",
			input:    "invalid",
			expected: SignUpRetry.String(),
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
			s := NewSignUpState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: SignUp.String(),
						Src:  []string{"idle"},
						Dst:  SignUp.String(),
					},
					{
						Name: SignUpRetry.String(),
						Src:  []string{SignUp.String()},
						Dst:  SignUpRetry.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{SignUp.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					SignUp.String(): s.OnState,
				},
			)
			f.SetMetadata(discord.Username.String(), "leo")
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
				err = f.Event(context.TODO(), SignUp.String())
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

func TestSignUpRetryState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
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
		isErr    bool
	}{
		{
			name:     "add accepted",
			input:    "1",
			expected: SignUpRetry.String(),
		},
		{
			name:     "add declined",
			input:    "2",
			expected: SignUpRetry.String(),
		},
		{
			name:     "add tentative",
			input:    "3",
			expected: SignUpRetry.String(),
		},
		{
			name:     "invalid",
			input:    "invalid",
			expected: SelfTransition.String(),
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
			s := NewSignUpRetryState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: SignUpRetry.String(),
						Src:  []string{"idle"},
						Dst:  SignUpRetry.String(),
					},
					{
						Name: SelfTransition.String(),
						Src:  []string{SignUpRetry.String()},
						Dst:  SelfTransition.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{SignUpRetry.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					SignUpRetry.String(): s.OnState,
				},
			)
			f.SetMetadata(discord.Username.String(), "leo")
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
				err = f.Event(context.TODO(), SignUpRetry.String())
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

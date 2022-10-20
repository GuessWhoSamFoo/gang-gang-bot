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

func TestNewAddResponseState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewAddResponseState(*opts)
	assert.NotNil(t, s)
}

func TestNewUnknownUserState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewUnknownUserState(*opts)
	assert.NotNil(t, s)
}

func TestAddResponseState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	rg := role.NewDefaultRoleGroup()
	rg.Roles[0].Users = []string{"leo"}

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

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "user in guild",
			input:    "testUserBot",
			expected: SignUp.String(),
		},
		{
			name:     "multiple results matched",
			input:    "test",
			expected: SelfTransition.String(),
		},
		{
			name:     "user does not exist in guild",
			input:    "invalid",
			expected: UnknownUser.String(),
		},
		{
			name:     "user already signed up",
			input:    "leo",
			expected: Cancel.String(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAddResponseState(*opts)

			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: AddResponse.String(),
						Src:  []string{"idle"},
						Dst:  AddResponse.String(),
					},
					{
						Name: SignUp.String(),
						Src:  []string{AddResponse.String()},
						Dst:  SignUp.String(),
					},
					{
						Name: UnknownUser.String(),
						Src:  []string{AddResponse.String()},
						Dst:  UnknownUser.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{AddResponse.String()},
						Dst:  Cancel.String(),
					},
					{
						Name: SelfTransition.String(),
						Src:  []string{AddResponse.String()},
						Dst:  SelfTransition.String(),
					},
				},
				fsm.Callbacks{
					AddResponse.String(): s.OnState,
				},
			)
			f.SetMetadata(discord.EventObject.String(), event)
			go func() {
				s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					s.inputHandler.inputChan <- tc.input
				}
				s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), AddResponse.String())
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

func TestUnknownUserState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "to add response",
			input:    "1",
			expected: AddResponse.String(),
		},
		{
			name:     "to sign up",
			input:    "2",
			expected: SignUp.String(),
		},
		{
			name:     "to cancel",
			input:    "3",
			expected: Cancel.String(),
		},
		{
			name:     "invalid",
			input:    "invalid",
			expected: UnknownUserRetry.String(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUnknownUserState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: UnknownUser.String(),
						Src:  []string{"idle"},
						Dst:  UnknownUser.String(),
					},
					{
						Name: AddResponse.String(),
						Src:  []string{UnknownUser.String()},
						Dst:  AddResponse.String(),
					},
					{
						Name: SignUp.String(),
						Src:  []string{UnknownUser.String()},
						Dst:  SignUp.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{UnknownUser.String()},
						Dst:  Cancel.String(),
					},
					{
						Name: UnknownUserRetry.String(),
						Src:  []string{UnknownUser.String()},
						Dst:  UnknownUserRetry.String(),
					},
				},
				fsm.Callbacks{
					UnknownUser.String(): s.OnState,
				},
			)
			f.SetMetadata(discord.Username.String(), mockconstants.TestUser)

			go func() {
				s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					s.inputHandler.inputChan <- tc.input
				}
				s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), UnknownUser.String())
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, f.Current())
		})
	}

}

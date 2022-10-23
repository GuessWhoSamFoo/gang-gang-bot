package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewSetAttendeeState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewSetAttendeeState(*opts)
	assert.NotNil(t, s)
}

func TestSetAttendeeState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	ctx := context.Background()

	cases := []struct {
		name          string
		input         string
		expectedState string
		isErr         bool
	}{
		{
			name:          "within range",
			input:         "50",
			expectedState: SetAttendeeLimit.String(),
		},
		{
			name:          "out of bounds",
			input:         "10000",
			expectedState: SetAttendeeRetry.String(),
		},
		{
			name:          "none",
			input:         "none",
			expectedState: SetAttendeeLimit.String(),
		},
		{
			name:          "invalid",
			input:         "invalid",
			expectedState: SetAttendeeRetry.String(),
		},
		{
			name:          "cancel",
			input:         "cancel",
			expectedState: Cancel.String(),
			isErr:         true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSetAttendeeState(*opts)

			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: SetAttendeeLimit.String(),
						Src:  []string{"idle"},
						Dst:  SetAttendeeLimit.String(),
					},
					{
						Name: SetAttendeeRetry.String(),
						Src:  []string{SetAttendeeLimit.String()},
						Dst:  SetAttendeeRetry.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{SetAttendeeLimit.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					SetAttendeeLimit.String(): s.OnState,
				},
			)
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
				err = f.Event(ctx, SetAttendeeLimit.String())
				if tc.isErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				wg.Done()
			}()
			wg.Wait()
			assert.Equal(t, tc.expectedState, f.Current())
		})
	}
}

func TestSetAttendeeRetryState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	ctx := context.Background()

	cases := []struct {
		name          string
		input         string
		expectedState string
		isErr         bool
	}{
		{
			name:          "within range",
			input:         "50",
			expectedState: SetAttendeeRetry.String(),
		},
		{
			name:          "out of bounds",
			input:         "10000",
			expectedState: SelfTransition.String(),
		},
		{
			name:          "none",
			input:         "none",
			expectedState: SetAttendeeRetry.String(),
		},
		{
			name:          "invalid",
			input:         "invalid",
			expectedState: SelfTransition.String(),
		},
		{
			name:          "cancel",
			input:         "cancel",
			expectedState: Cancel.String(),
			isErr:         true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSetAttendeeRetryState(*opts)

			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: SetAttendeeRetry.String(),
						Src:  []string{"idle"},
						Dst:  SetAttendeeRetry.String(),
					},
					{
						Name: SelfTransition.String(),
						Src:  []string{SetAttendeeRetry.String()},
						Dst:  SelfTransition.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{SetAttendeeRetry.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					SetAttendeeRetry.String(): s.OnState,
				},
			)
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
				err = f.Event(ctx, SetAttendeeRetry.String())
				if tc.isErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				wg.Done()
			}()
			wg.Wait()
			assert.Equal(t, tc.expectedState, f.Current())
		})
	}
}

func Test_validateAttendee(t *testing.T) {
	cases := []struct {
		name  string
		input string
		isErr bool
	}{
		{
			name:  "within bounds",
			input: "100",
		},
		{
			name:  "none",
			input: "None",
		},
		{
			name:  "out of bounds higher",
			input: "251",
			isErr: true,
		},
		{
			name:  "out of bounds negative",
			input: "-1",
			isErr: true,
		},
		{
			name:  "out of bounds zero",
			input: "0",
			isErr: true,
		},
		{
			name:  "not a number",
			input: "invalid",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := fsm.NewFSM("", fsm.Events{}, fsm.Callbacks{})
			f.SetMetadata(discord.Attendee.String(), tc.input)
			err := validateAttendee(f, discord.Attendee)
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

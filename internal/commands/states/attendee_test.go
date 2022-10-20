package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSetAttendeeState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewSetAttendeeState(*opts)
	assert.NotNil(t, s)
}

func TestSetAttendeeState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
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
			go func() {
				s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					s.inputHandler.inputChan <- tc.input
				}
				s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(ctx, SetAttendeeLimit.String())
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedState, f.Current())
		})
	}
}

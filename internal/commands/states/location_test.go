package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSetLocationState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	s := NewSetLocationState(*opts)
	assert.NotNil(t, s)
}

func TestSetLocationState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	cases := []struct {
		name          string
		input         string
		expectedState string
		isErr         bool
	}{
		{
			name:          "location",
			input:         "Seattle",
			expectedState: SetLocation.String(),
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
			s := NewSetLocationState(*opts)

			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: SetLocation.String(),
						Src:  []string{"idle"},
						Dst:  SetLocation.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{SetLocation.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					SetLocation.String(): s.OnState,
				},
			)
			go func() {
				s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					s.inputHandler.inputChan <- tc.input
				}
				s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), SetLocation.String())
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedState, f.Current())
		})
	}
}

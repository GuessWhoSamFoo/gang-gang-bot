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

func TestNewAddDescriptionState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewAddDescriptionState(*opts)
	assert.NotNil(t, s)
}

func TestAddDescriptionState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)

	cases := []struct {
		name          string
		input         string
		expectedDesc  string
		expectedState string
		isErr         bool
	}{
		{
			name:          "title",
			input:         "hello",
			expectedDesc:  "hello",
			expectedState: AddDescription.String(),
		},
		{
			name:          "cancel",
			input:         "cancel",
			expectedDesc:  "",
			expectedState: Cancel.String(),
			isErr:         true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewAddDescriptionState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: AddDescription.String(),
						Src:  []string{"idle"},
						Dst:  AddDescription.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{AddDescription.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					AddDescription.String(): s.OnState,
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
				err = f.Event(context.TODO(), AddDescription.String())
				if tc.isErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
				wg.Done()
			}()
			wg.Wait()
			got, err := Get(f, discord.Description)
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedState, f.Current())
			assert.Equal(t, tc.expectedDesc, got)
		})
	}
}

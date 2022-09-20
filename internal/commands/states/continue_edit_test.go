package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContinueEditState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	s := NewContinueEditState(*opts)
	assert.NotNil(t, s)
}

func TestNewContinueEditState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "to process edit",
			input:    "1",
			expected: ProcessEdit.String(),
		},
		{
			name:     "to modify event",
			input:    "2",
			expected: ModifyEvent.String(),
		},
		{
			name:     "invalid input",
			input:    "invalid",
			expected: ContinueEditRetry.String(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewContinueEditState(*opts)

			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: ContinueEdit.String(),
						Src:  []string{"idle"},
						Dst:  ContinueEdit.String(),
					},
					{
						Name: ProcessEdit.String(),
						Src:  []string{ContinueEdit.String()},
						Dst:  ProcessEdit.String(),
					},
					{
						Name: ModifyEvent.String(),
						Src:  []string{ContinueEdit.String()},
						Dst:  ModifyEvent.String(),
					},
					{
						Name: ContinueEditRetry.String(),
						Src:  []string{ContinueEdit.String()},
						Dst:  ContinueEditRetry.String(),
					},
				},
				fsm.Callbacks{
					ContinueEdit.String(): c.OnState,
				},
			)

			go func() {
				c.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					c.input <- tc.input
				}
				c.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), ContinueEdit.String())
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

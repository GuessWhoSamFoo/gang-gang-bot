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
			expected: ContinueEdit.String(),
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
				c.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					c.inputHandler.inputChan <- tc.input
				}
				c.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), ContinueEdit.String())
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

func TestContinueEditRetryState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	cases := []struct {
		name     string
		input    string
		expected string
		isErr    bool
	}{
		{
			name:     "to process edit",
			input:    "1",
			expected: ContinueEditRetry.String(),
		},
		{
			name:     "to modify event",
			input:    "2",
			expected: ModifyEvent.String(),
		},
		{
			name:     "invalid input",
			input:    "invalid",
			expected: SelfTransition.String(),
		},
		{
			name:     "cancel",
			input:    "Cancel",
			expected: Cancel.String(),
			isErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewContinueEditRetryState(*opts)

			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: ContinueEditRetry.String(),
						Src:  []string{"idle"},
						Dst:  ContinueEditRetry.String(),
					},
					{
						Name: ProcessEdit.String(),
						Src:  []string{ContinueEditRetry.String()},
						Dst:  ProcessEdit.String(),
					},
					{
						Name: ModifyEvent.String(),
						Src:  []string{ContinueEditRetry.String()},
						Dst:  ModifyEvent.String(),
					},
					{
						Name: SelfTransition.String(),
						Src:  []string{ContinueEditRetry.String()},
						Dst:  SelfTransition.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{ContinueEditRetry.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					ContinueEditRetry.String(): c.OnState,
				},
			)

			go func() {
				c.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
					c.inputHandler.inputChan <- tc.input
				}
				c.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
			}()

			err = f.Event(context.TODO(), ContinueEditRetry.String())
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, f.Current())
		})
	}
}

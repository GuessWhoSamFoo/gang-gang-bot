package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewDurationState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	s := NewDurationState(*opts)
	assert.NotNil(t, s)
}

func TestSetDurationState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	now := time.Now()

	cases := []struct {
		name          string
		input         string
		expectedState string
		expectedTime  time.Time
		isErr         bool
	}{
		{
			name:          "valid",
			input:         "1 hour",
			expectedState: SetDuration.String(),
			expectedTime:  now.Add(time.Minute * 60),
		},
		{
			name:          "invalid",
			input:         "invalid",
			expectedState: SetDurationRetry.String(),
			expectedTime:  time.Time{},
		},
		{
			name:          "none",
			input:         "none",
			expectedState: SetDuration.String(),
			expectedTime:  time.Time{},
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
			d := NewDurationState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: SetDuration.String(),
						Src:  []string{"idle"},
						Dst:  SetDuration.String(),
					},
					{
						Name: SetDurationRetry.String(),
						Src:  []string{SetDuration.String()},
						Dst:  SetDurationRetry.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{SetDuration.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					SetDuration.String(): d.OnState,
				},
			)
			f.SetMetadata(discord.StartTime.String(), now)
			d.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
				d.inputHandler.inputChan <- tc.input
			}
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				d.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
				wg.Done()
			}()

			go func() {
				err = f.Event(context.TODO(), SetDuration.String())
				if tc.isErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					got, err := Get(f, discord.Duration)
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedTime, got)
				}
				wg.Done()
			}()
			wg.Wait()
			assert.Equal(t, tc.expectedState, f.Current())
		})
	}
}

func TestSetDurationRetryState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	now := time.Now()

	cases := []struct {
		name          string
		input         string
		expectedState string
		expectedTime  time.Time
		isErr         bool
	}{
		{
			name:          "valid",
			input:         "1 hour",
			expectedState: SetDurationRetry.String(),
			expectedTime:  now.Add(time.Minute * 60),
		},
		{
			name:          "invalid",
			input:         "invalid",
			expectedState: SelfTransition.String(),
			expectedTime:  time.Time{},
		},
		{
			name:          "none",
			input:         "none",
			expectedState: SetDurationRetry.String(),
			expectedTime:  time.Time{},
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
			d := NewDurationRetryState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: SetDurationRetry.String(),
						Src:  []string{"idle"},
						Dst:  SetDurationRetry.String(),
					},
					{
						Name: SelfTransition.String(),
						Src:  []string{SetDurationRetry.String()},
						Dst:  SelfTransition.String(),
					},
					{
						Name: Cancel.String(),
						Src:  []string{SetDurationRetry.String()},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					SetDurationRetry.String(): d.OnState,
				},
			)
			f.SetMetadata(discord.StartTime.String(), now)
			d.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
				d.inputHandler.inputChan <- tc.input
			}
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				d.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
				wg.Done()
			}()
			go func() {
				err = f.Event(context.TODO(), SetDurationRetry.String())
				if tc.isErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					got, err := Get(f, discord.Duration)
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedTime, got)
				}
				wg.Done()
			}()
			wg.Wait()
			assert.Equal(t, tc.expectedState, f.Current())
		})
	}
}

func Test_validateDuration(t *testing.T) {
	now := time.Now()
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	d := NewDurationState(*opts)

	f := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: SetDuration.String(),
				Src:  []string{"idle"},
				Dst:  SetDuration.String(),
			},
		},
		fsm.Callbacks{
			SetDuration.String(): d.OnState,
		},
	)

	e := &fsm.Event{
		FSM: f,
	}
	f.SetMetadata(discord.StartTime.String(), now)

	cases := []struct {
		input    string
		expected time.Time
		isErr    bool
	}{
		{
			input:    "1 hour",
			expected: now.Add(time.Hour),
		},
		{
			input:    "1 hour and 30 minutes",
			expected: now.Add(time.Minute * 90),
		},
		{
			input:    "1 hour 30 minutes",
			expected: now.Add(time.Minute * 90),
		},
		{
			input:    "none",
			expected: time.Time{},
		},
		{
			input:    "invalid",
			expected: time.Time{},
			isErr:    true,
		},
		// TODO: Find a more configurable parser
		//{
		//	input:    "1h30m",
		//	expected: now.Add(time.Minute * 90),
		//},
		//{
		//	input:    "1hr 30 minutes",
		//	expected: time.Time{},
		//	isErr:    true,
		//},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			f.SetMetadata(discord.Duration.String(), tc.input)
			err = validateDuration(e, discord.Duration)
			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			actual, exists := f.Metadata(discord.Duration.String())
			assert.True(t, exists)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

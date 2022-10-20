package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
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
	now := time.Now()
	f.SetMetadata(discord.StartTime.String(), now)

	go func() {
		d.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
			d.inputHandler.inputChan <- "1 hour"
		}
		d.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
	}()

	err = f.Event(context.TODO(), SetDuration.String())
	assert.NoError(t, err)

	actual, err := Get(f, discord.Duration)
	assert.NoError(t, err)

	expected := now.Add(time.Hour)
	assert.Equal(t, expected, actual)
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
			input:    "invalid",
			expected: now,
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

			if tc.isErr {
				assert.Equal(t, tc.input, actual)
			} else {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

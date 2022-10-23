package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCancelState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)

	s := NewCancelState(*opts)
	assert.NotNil(t, s)
}

func TestCancelState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)

	cases := []struct {
		name   string
		action string
	}{
		{
			name:   "event create",
			action: CreateAction,
		},
		{
			name:   "event edit",
			action: EditAction,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCancelState(*opts)
			f := fsm.NewFSM(
				"idle",
				fsm.Events{
					{
						Name: Cancel.String(),
						Src:  []string{"idle"},
						Dst:  Cancel.String(),
					},
				},
				fsm.Callbacks{
					Cancel.String(): c.OnState,
				},
			)
			f.SetMetadata(discord.Action.String(), tc.action)
			err = f.Event(context.TODO(), Cancel.String())
			assert.NoError(t, err)

			got, err := Get(f, discord.Action)
			assert.NoError(t, err)

			assert.Equal(t, tc.action, got)
		})
	}
}

package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSelfTransitionState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewSelfTransitionState(*opts)
	assert.NotNil(t, s)
}

func TestSelfTransitionState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewSelfTransitionState(*opts)

	f := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: SelfTransition.String(),
				Src:  []string{"idle"},
				Dst:  SelfTransition.String(),
			},
			{
				Name: "idle",
				Src:  []string{SelfTransition.String()},
				Dst:  "idle",
			},
		},
		fsm.Callbacks{
			SelfTransition.String(): s.OnState,
		},
	)
	err = f.Event(context.TODO(), SelfTransition.String())
	assert.NoError(t, err)
	assert.Equal(t, "idle", f.Current())
}

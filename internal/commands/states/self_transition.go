package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
)

// SelfTransitionState is an intermediary state that loops back to the previous state
type SelfTransitionState struct {
	*discord.Options
}

func NewSelfTransitionState(o discord.Options) *SelfTransitionState {
	return &SelfTransitionState{
		&o,
	}
}

func (s *SelfTransitionState) OnState(ctx context.Context, e *fsm.Event) {
	err := e.FSM.Event(ctx, e.Src)
	if err != nil {
		e.Err = err
		return
	}
}

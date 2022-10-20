package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
)

type TimeoutState struct {
	*discord.Options
	inputHandler *InputHandler
}

func NewTimeoutState(o discord.Options) *TimeoutState {
	return &TimeoutState{
		Options:      &o,
		inputHandler: NewInputHandler(&o),
	}
}

func (t *TimeoutState) OnState(_ context.Context, e *fsm.Event) {
	_, err := t.Session.ChannelMessageSend(t.Channel.ID, "I'm not sure where you went. We can try this again later.")
	if err != nil {
		e.Err = err
	}
}

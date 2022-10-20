package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
)

type CancelState struct {
	*discord.Options
}

func NewCancelState(o discord.Options) *CancelState {
	return &CancelState{
		&o,
	}
}

func (c *CancelState) OnState(_ context.Context, e *fsm.Event) {
	action, err := Get(e.FSM, discord.Action)
	if err != nil {
		e.Err = err
		return
	}

	_, err = c.Session.ChannelMessageSend(c.Channel.ID, fmt.Sprintf("Event %s has been canceled", action))
	if err != nil {
		e.Err = err
	}
}

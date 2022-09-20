package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
)

type CancelState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel
}

func NewCancelState(o discord.Options) *CancelState {
	return &CancelState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
	}
}

func (c *CancelState) OnState(_ context.Context, e *fsm.Event) {
	action, err := Get(e.FSM, discord.Action)
	if err != nil {
		e.Err = err
		return
	}

	_, err = c.session.ChannelMessageSend(c.channel.ID, fmt.Sprintf("Event %s has been canceled", action))
	if err != nil {
		e.Err = err
	}
}

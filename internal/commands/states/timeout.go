package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
)

type TimeoutState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel
}

func NewTimeoutState(o discord.Options) *TimeoutState {
	return &TimeoutState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
	}
}

func (t *TimeoutState) OnState(_ context.Context, e *fsm.Event) {
	_, err := t.session.ChannelMessageSend(t.channel.ID, "I'm not sure where you went. We can try this again later.")
	if err != nil {
		e.Err = err
	}
}

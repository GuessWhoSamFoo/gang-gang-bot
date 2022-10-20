package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"time"
)

type SetLocationState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewSetLocationState(o discord.Options) *SetLocationState {
	return &SetLocationState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (l *SetLocationState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := l.session.ChannelMessageSendEmbed(l.channel.ID, &discord.EnterLocationMessage)
	if err != nil {
		e.Err = err
		return
	}

	err = l.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.Location, 60*time.Second)
	if err != nil {
		e.Err = err
		return
	}
}

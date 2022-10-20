package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"time"
)

type AddDescriptionState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewAddDescriptionState(o discord.Options) *AddDescriptionState {
	return &AddDescriptionState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (a *AddDescriptionState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := a.session.ChannelMessageSendEmbed(a.channel.ID, &discord.EnterDescriptionMessage)
	if err != nil {
		e.Err = err
		return
	}

	err = a.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.Description, 60*time.Second)
	if err != nil {
		e.Err = err
	}
}

package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"time"
)

type AddTitleState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewAddTitleState(o discord.Options) *AddTitleState {
	return &AddTitleState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (a *AddTitleState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := a.session.ChannelMessageSendEmbed(a.channel.ID, &discord.EnterTitleMessage)
	if err != nil {
		e.Err = err
		return
	}

	err = a.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.Title, 60*time.Second)
	if err != nil {
		e.Err = err
	}
}

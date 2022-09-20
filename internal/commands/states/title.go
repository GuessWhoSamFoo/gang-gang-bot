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

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewAddTitleState(o discord.Options) *AddTitleState {
	i := make(chan string)

	return &AddTitleState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		input:             i,
		handlerFunc: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.ChannelID == o.Channel.ID {
				i <- m.Content
			}
		},
	}
}

func (a *AddTitleState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := a.session.ChannelMessageSendEmbed(a.channel.ID, &discord.EnterTitleMessage)
	if err != nil {
		e.Err = err
		return
	}

	err = AwaitInputOrTimeout(ctx, 60*time.Second, a.session, a.input, e.FSM, a.handlerFunc, discord.Title)
	if err != nil {
		e.Err = err
	}
}

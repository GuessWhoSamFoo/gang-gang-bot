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

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewSetLocationState(o discord.Options) *SetLocationState {
	i := make(chan string)

	return &SetLocationState{
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

func (l *SetLocationState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := l.session.ChannelMessageSendEmbed(l.channel.ID, &discord.EnterLocationMessage)
	if err != nil {
		e.Err = err
		return
	}

	err = AwaitInputOrTimeout(ctx, 60*time.Second, l.session, l.input, e.FSM, l.handlerFunc, discord.Location)
	if err != nil {
		e.Err = err
		return
	}
}

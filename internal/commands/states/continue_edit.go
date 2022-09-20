package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"time"
)

type ContinueEditState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewContinueEditState(o discord.Options) *ContinueEditState {
	i := make(chan string)

	return &ContinueEditState{
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

func (c *ContinueEditState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := c.session.ChannelMessageSendEmbed(c.channel.ID, &discord.EditConfirmationMessage); err != nil {
		e.Err = err
		return
	}

	if err := AwaitInputOrTimeout(ctx, 60*time.Second, c.session, c.input, e.FSM, c.handlerFunc, discord.MenuOption); err != nil {
		e.Err = err
		return
	}
	state, err := ConfirmSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, ContinueEditRetry.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
		return
	}
	err = e.FSM.Event(ctx, state)
	if err != nil {
		e.Err = err
		return
	}
}

type ContinueEditRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewContinueEditRetryState(o discord.Options) *ContinueEditRetryState {
	i := make(chan string)

	return &ContinueEditRetryState{
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

func (c *ContinueEditRetryState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := c.session.ChannelMessageSend(c.channel.ID, discord.InvalidEntryText); err != nil {
		e.Err = err
		return
	}

	if err := AwaitInputOrTimeout(ctx, 60*time.Second, c.session, c.input, e.FSM, c.handlerFunc, discord.MenuOption); err != nil {
		e.Err = err
		return
	}

	state, err := ConfirmSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
	err = e.FSM.Event(ctx, state)
	if err != nil {
		e.Err = err
		return
	}
}

func ConfirmSelect(e *fsm.Event) (string, error) {
	val, err := Get(e.FSM, discord.MenuOption)
	if err != nil {
		return "", err
	}

	opts := map[string]chatState{
		"1": ProcessEdit,
		"2": ModifyEvent,
	}
	option, ok := opts[val.(string)]
	if !ok {
		return "", fmt.Errorf("cannot find response")
	}
	return option.String(), nil
}

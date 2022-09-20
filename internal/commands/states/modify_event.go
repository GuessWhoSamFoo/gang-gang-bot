package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"time"
)

type ModifyEventState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel
	input             chan string
	handlerFunc       func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewModifyEventState(o discord.Options) *ModifyEventState {
	i := make(chan string)

	return &ModifyEventState{
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

func (m *ModifyEventState) OnState(ctx context.Context, e *fsm.Event) {
	obj, err := Get(e.FSM, discord.EventObject)
	if err != nil {
		e.Err = err
		return
	}
	event, ok := obj.(discord.Event)
	if !ok {
		e.Err = fmt.Errorf("cannot get event")
		return
	}

	if _, err = m.session.ChannelMessageSendEmbed(m.channel.ID, &discordgo.MessageEmbed{
		Title: "What would you like to modify?",
		Color: discord.Purple,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "1 ⋅ Title",
				Value: util.PrintBlockValues(event.Title),
			},
			{
				Name:  "2 ⋅ Description",
				Value: util.PrintBlockValues(event.Description),
			},
			{
				Name:   "3 ⋅ Start Time",
				Value:  fmt.Sprintf("```%s```", event.Start.In(time.Local).Format(util.HumanTimeFormat)),
				Inline: true,
			},
			{
				Name:   "4 ⋅ Duration",
				Value:  util.PrintBlockValues(util.PrintHumanReadableTime(event.Start, event.End)),
				Inline: true,
			},
			{
				Name:  "5 ⋅ Location",
				Value: fmt.Sprintf("```%s```", event.Location),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: discord.OptionText + "\n" + discord.CancelText,
		},
	}); err != nil {
		e.Err = err
		return
	}

	if err = AwaitInputOrTimeout(ctx, 60*time.Second, m.session, m.input, e.FSM, m.handlerFunc, discord.MenuOption); err != nil {
		e.Err = err
		return
	}

	state, err := EditFieldSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, ModifyEventRetry.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
		return
	}
	if err = e.FSM.Event(ctx, state); err != nil {
		e.Err = err
		return
	}
}

type ModifyEventRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel
	input             chan string
	handlerFunc       func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewModifyEventRetryState(o discord.Options) *ModifyEventRetryState {
	i := make(chan string)

	return &ModifyEventRetryState{
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

func (m *ModifyEventRetryState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := m.session.ChannelMessageSend(m.channel.ID, discord.InvalidEntryText); err != nil {
		e.Err = err
		return
	}
	if err := AwaitInputOrTimeout(ctx, 60*time.Second, m.session, m.input, e.FSM, m.handlerFunc, discord.MenuOption); err != nil {
		e.Err = err
		return
	}

	state, err := EditFieldSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
	if err = e.FSM.Event(ctx, state); err != nil {
		e.Err = err
		return
	}
}

func EditFieldSelect(e *fsm.Event) (string, error) {
	val, err := Get(e.FSM, discord.MenuOption)
	if err != nil {
		return "", err
	}

	opts := map[string]chatState{
		"1": AddTitle,
		"2": AddDescription,
		"3": SetDate,
		"4": SetDuration,
		"5": SetLocation,
	}
	option, ok := opts[val.(string)]
	if !ok {
		return "", fmt.Errorf("cannot find response")
	}
	return option.String(), nil
}

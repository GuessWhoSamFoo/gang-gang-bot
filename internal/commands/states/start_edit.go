package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"time"
)

const EditAction = "modification"

type StartEditState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewStartEditState(o discord.Options) *StartEditState {
	i := make(chan string)
	return &StartEditState{
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

func (s *StartEditState) OnState(ctx context.Context, e *fsm.Event) {
	event, err := discord.GetEventFromMessage(s.interactionCreate.Interaction.Message)
	if err != nil {
		e.Err = err
		return
	}
	if s.interactionCreate.Interaction.Member.User.Username != event.Owner && s.interactionCreate.Interaction.Member.Permissions&discordgo.PermissionManageEvents == 0 {
		if _, err := s.session.ChannelMessageSendEmbed(s.channel.ID, discord.EditInsufficientPermissionMessage); err != nil {
			e.Err = fmt.Errorf("failed to send message: %v", err)
			return
		}
		e.Err = fmt.Errorf("insufficient permissions to edit %s", s.interactionCreate.Interaction.Message.ID)
		return
	}
	e.FSM.SetMetadata(discord.Action.String(), EditAction)
	event.DiscordLink = fmt.Sprintf("https://discord.com/channels/%s/%s/%s", s.interactionCreate.GuildID, s.interactionCreate.Interaction.ChannelID, s.interactionCreate.Interaction.Message.ID)

	if _, err := s.session.ChannelMessageSendEmbed(s.channel.ID, &discord.EnterEditOptionMessage); err != nil {
		e.Err = fmt.Errorf("failed to send message: %v", err)
		return
	}

	if err = AwaitInputOrTimeout(ctx, 60*time.Second, s.session, s.input, e.FSM, s.handlerFunc, discord.MenuOption); err != nil {
		e.Err = err
		return
	}

	state, err := EditTypeSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, StartEditRetry.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
		}
		return
	}
	e.FSM.SetMetadata(discord.EventObject.String(), *event)
	if err = e.FSM.Event(ctx, state); err != nil {
		e.Err = err
		return
	}
}

type StartEditRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewStartEditRetryState(o discord.Options) *StartEditRetryState {
	i := make(chan string)
	return &StartEditRetryState{
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

func (r *StartEditRetryState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := r.session.ChannelMessageSend(r.channel.ID, discord.InvalidEntryText); err != nil {
		e.Err = err
		return
	}

	if err := AwaitInputOrTimeout(ctx, 60*time.Second, r.session, r.input, e.FSM, r.handlerFunc, discord.MenuOption); err != nil {
		e.Err = err
		return
	}

	state, err := EditTypeSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
		}
		return
	}

	err = e.FSM.Event(ctx, state)
	if err = e.FSM.Event(ctx, state); err != nil {
		e.Err = err
		return
	}
}

func EditTypeSelect(e *fsm.Event) (string, error) {
	val, err := Get(e.FSM, discord.MenuOption)
	if err != nil {
		return "", err
	}

	opts := map[string]chatState{
		"1": ModifyEvent,
		"2": RemoveResponse,
		"3": AddResponse,
	}
	option, ok := opts[val.(string)]
	if !ok {
		return "", fmt.Errorf("cannot find response")
	}
	return option.String(), nil
}

package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"time"
)

type SignUpState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewSignUpState(o discord.Options) *SignUpState {
	return &SignUpState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (s *SignUpState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := s.session.ChannelMessageSendEmbed(s.channel.ID, &discordgo.MessageEmbed{
		Title:       "Which signup option should we add the user to?",
		Description: fmt.Sprintf("**1**⠀%s\n**2**⠀%s\n**3**⠀%s", discord.AcceptedBase, discord.DeclinedBase, discord.TentativeBase),
		Color:       discord.Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: discord.CancelText,
		},
	}); err != nil {
		e.Err = err
		return
	}
	user, err := Get(e.FSM, discord.Username)
	if err != nil {
		e.Err = err
		return
	}
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

	if err = s.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
		e.Err = err
		return
	}
	err = SelectRole(e, s.session, s.interactionCreate, &event, fmt.Sprintf("%v", user))
	if err != nil {
		eventErr := e.FSM.Event(ctx, SignUpRetry.String())
		if eventErr != nil {
			e.Err = eventErr
			return
		}
	}
}

type SignUpRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewSignUpRetryState(o discord.Options) *SignUpRetryState {
	return &SignUpRetryState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (r *SignUpRetryState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := r.session.ChannelMessageSend(r.channel.ID, discord.InvalidEntryText); err != nil {
		e.Err = err
		return
	}
	user, err := Get(e.FSM, discord.Username)
	if err != nil {
		e.Err = err
		return
	}
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

	if err = r.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
		e.Err = err
		return
	}
	err = SelectRole(e, r.session, r.interactionCreate, &event, fmt.Sprintf("%v", user))
	if err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
}

func SelectRole(e *fsm.Event, s *discordgo.Session, ic *discordgo.InteractionCreate, event *discord.Event, name string) error {
	val, err := Get(e.FSM, discord.MenuOption)
	if err != nil {
		return err
	}

	opts := map[string]func(s *discordgo.Session, ic *discordgo.InteractionCreate, name string) error{
		"1": event.ToggleAccept,
		"2": event.ToggleDecline,
		"3": event.ToggleTentative,
	}
	option, ok := opts[val.(string)]
	if !ok {
		return fmt.Errorf("cannot find %s response", e.FSM.Current())
	}
	return option(s, ic, name)
}

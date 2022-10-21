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

	inputHandler *InputHandler
}

func NewContinueEditState(o discord.Options) *ContinueEditState {
	return &ContinueEditState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (c *ContinueEditState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := c.session.ChannelMessageSendEmbed(c.channel.ID, &discord.EditConfirmationMessage); err != nil {
		e.Err = err
		return
	}

	if err := c.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
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
	// No-op as all edits will be processed
	if state == ProcessEdit.String() {
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

	inputHandler *InputHandler
}

func NewContinueEditRetryState(o discord.Options) *ContinueEditRetryState {
	return &ContinueEditRetryState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (c *ContinueEditRetryState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := c.session.ChannelMessageSend(c.channel.ID, discord.InvalidEntryText); err != nil {
		e.Err = err
		return
	}

	if err := c.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
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
		return
	}
	// No-op as all edits will be processed
	if state == ProcessEdit.String() {
		return
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
		return "", fmt.Errorf("cannot find %s response", e.FSM.Current())
	}
	return option.String(), nil
}

package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/tj/go-naturaldate"
	"time"
)

type SetDurationState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewDurationState(o discord.Options) *SetDurationState {
	i := make(chan string)

	return &SetDurationState{
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

func (d *SetDurationState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := d.session.ChannelMessageSendEmbed(d.channel.ID, &discord.EnterDurationMessage)
	if err != nil {
		e.Err = err
		return
	}

	if err = AwaitInputOrTimeout(ctx, 60*time.Second, d.session, d.input, e.FSM, d.handlerFunc, discord.Duration); err != nil {
		e.Err = err
		return
	}
	if err = validateDuration(e, discord.Duration); err != nil {
		eventErr := e.FSM.Event(ctx, SetDurationRetry.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
}

type SetDurationRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewDurationRetryState(o discord.Options) *SetDurationRetryState {
	i := make(chan string)

	return &SetDurationRetryState{
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

func (d *SetDurationRetryState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := d.session.ChannelMessageSend(d.channel.ID, discord.InvalidDurationText)
	if err != nil {
		e.Err = err
		return
	}

	if err = AwaitInputOrTimeout(ctx, 60*time.Second, d.session, d.input, e.FSM, d.handlerFunc, discord.Duration); err != nil {
		e.Err = err
		return
	}
	if err = validateDuration(e, discord.Duration); err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
}

func validateDuration(e *fsm.Event, key discord.MetadataKey) error {
	val, err := Get(e.FSM, key)
	if err != nil {
		return err
	}
	input := fmt.Sprintf("%v", val)
	start, err := Get(e.FSM, discord.StartTime)
	if err != nil {
		return err
	}
	startTime := start.(time.Time)

	endTime, err := naturaldate.Parse(input, startTime, naturaldate.WithDirection(naturaldate.Future))
	if err != nil {
		return err
	}

	if input == "" || endTime.Before(startTime) || endTime.Equal(startTime) {
		return fmt.Errorf("invalid end time")
	}
	e.FSM.SetMetadata(key.String(), endTime)
	return nil
}

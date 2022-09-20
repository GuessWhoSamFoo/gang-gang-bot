package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/tj/go-naturaldate"
	"strings"
	"time"
)

type SetDateState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewSetDateState(o discord.Options) *SetDateState {
	i := make(chan string)

	return &SetDateState{
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

func (d *SetDateState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := d.session.ChannelMessageSendEmbed(d.channel.ID, &discord.EnterDateStartMessage)
	if err != nil {
		e.Err = err
		return
	}

	if err = AwaitInputOrTimeout(ctx, 60*time.Second, d.session, d.input, e.FSM, d.handlerFunc, discord.StartTime); err != nil {
		e.Err = err
		return
	}

	if err = validateTime(e, d.session, d.channel, discord.StartTime); err != nil {
		eventErr := e.FSM.Event(ctx, SetDateRetry.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
}

type SetDateRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	input       chan string
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
}

func NewSetDateRetryState(o discord.Options) *SetDateRetryState {
	i := make(chan string)

	return &SetDateRetryState{
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

func (r *SetDateRetryState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := r.session.ChannelMessageSend(r.channel.ID, discord.InvalidEventTimeText)
	if err != nil {
		e.Err = err
		return
	}

	if err = AwaitInputOrTimeout(ctx, 60*time.Second, r.session, r.input, e.FSM, r.handlerFunc, discord.StartTime); err != nil {
		e.Err = err
		return
	}
	if err = validateTime(e, r.session, r.channel, discord.StartTime); err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
}

func validateTime(e *fsm.Event, s *discordgo.Session, c *discordgo.Channel, key discord.MetadataKey) error {
	val, err := Get(e.FSM, key)
	if err != nil {
		return err
	}
	input := fmt.Sprintf("%v", val)
	now := time.Now()

	var startTime time.Time
	startTime, err = naturaldate.Parse(input, now, naturaldate.WithDirection(naturaldate.Future))
	if err != nil {
		startTime, err = dateparse.ParseLocal(input)
		if err != nil {
			_, msgErr := s.ChannelMessageSend(c.ID, discord.InvalidStartTimeText)
			return fmt.Errorf("%v: %v", err, msgErr)
		}
	}
	if !strings.EqualFold(input, "now") && startTime.Equal(now) {
		return fmt.Errorf("failed to parse time")
	}

	if startTime.Before(now) {
		return fmt.Errorf("start time cannot be in the past")
	}
	e.FSM.SetMetadata(discord.StartTime.String(), startTime)
	return nil
}

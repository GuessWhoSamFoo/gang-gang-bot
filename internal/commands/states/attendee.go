package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

type SetAttendeeState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewSetAttendeeState(o discord.Options) *SetAttendeeState {
	return &SetAttendeeState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (s *SetAttendeeState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := s.session.ChannelMessageSendEmbed(s.channel.ID, &discord.EnterAttendeeLimitMessage)
	if err != nil {
		e.Err = err
		return
	}

	if err = s.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.Attendee, 60*time.Second); err != nil {
		e.Err = err
		return
	}

	if err = validateAttendee(e, discord.Attendee); err != nil {
		err = e.FSM.Event(ctx, SetAttendeeRetry.String())
		if err != nil {
			e.Err = err
			return
		}
	}
}

type SetAttendeeRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewSetAttendeeRetryState(o discord.Options) *SetAttendeeRetryState {
	return &SetAttendeeRetryState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (r *SetAttendeeRetryState) OnState(ctx context.Context, e *fsm.Event) {
	_, err := r.session.ChannelMessageSend(r.channel.ID, discord.InvalidEventLimitText)
	if err != nil {
		e.Err = err
		return
	}
	if err = r.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.Attendee, 60*time.Second); err != nil {
		e.Err = err
		return
	}
	if err = validateAttendee(e, discord.Attendee); err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
	}
}

func validateAttendee(e *fsm.Event, key discord.MetadataKey) error {
	val, err := Get(e.FSM, key)
	if err != nil {
		return err
	}

	rg := role.NewDefaultRoleGroup()

	var n int
	if strings.EqualFold(fmt.Sprintf("%v", val), "none") {
		n = 0
		rg.SetLimit(role.AcceptedField, n)
		e.FSM.SetMetadata(discord.Attendee.String(), rg)
		return nil
	}
	n, err = strconv.Atoi(fmt.Sprintf("%v", val))
	if err != nil {
		return err
	}

	if !InRange(n, 250, 1) {
		return fmt.Errorf("attendee range out of bounds")
	}
	rg.SetLimit(role.AcceptedField, n)
	e.FSM.SetMetadata(discord.Attendee.String(), rg)
	return nil
}

func InRange(n, high, low int) bool {
	return n <= high && n >= low
}

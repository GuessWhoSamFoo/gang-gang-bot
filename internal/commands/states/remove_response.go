package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

type RemoveResponseState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewRemoveResponseState(o discord.Options) *RemoveResponseState {
	return &RemoveResponseState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
		inputHandler:      NewInputHandler(&o),
	}
}

func (r *RemoveResponseState) OnState(ctx context.Context, e *fsm.Event) {
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

	var count int
	for _, rol := range event.RoleGroup.Roles {
		count += len(rol.Users)
	}
	if count == 0 {
		_, err = r.session.ChannelMessageSendEmbed(r.channel.ID, &discordgo.MessageEmbed{
			Title: "Event doesn't have any responses",
		})
		if err != nil {
			e.Err = err
			return
		}
		if err = e.FSM.Event(ctx, Cancel.String()); err != nil {
			e.Err = err
			return
		}
		e.Err = fmt.Errorf("event has no responses")
		return
	}

	var desc string
	var counter int
	users := make([]string, 0)
	// Braille space is used instead because hard spaces in embeds are not documented
	for _, r := range event.RoleGroup.Roles {
		users = append(users, r.Users...)
		for _, n := range r.Users {
			counter++
			desc = desc + fmt.Sprintf("**%d**⠀%s %s\n", counter, r.Icon, n)
		}
	}
	if _, err = r.session.ChannelMessageSendEmbed(r.channel.ID, &discordgo.MessageEmbed{
		Title:       "Which responses would you like to remove?",
		Description: desc,
		Color:       discord.Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Enter the number(s) of the desired option(s), separated by spaces\n" + discord.CancelText,
		},
	}); err != nil {
		e.Err = err
		return
	}

	wl, ok := event.RoleGroup.Waitlist[role.AcceptedField]
	if !ok {
		e.Err = fmt.Errorf("cannot find accepted field")
		return
	}
	nameMap := map[int]string{}
	for index, name := range append(users, wl.Users...) {
		nameMap[index+1] = name
	}

	if err = r.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
		e.Err = err
		return
	}

	names, err := selectMultiple(e, nameMap)
	if err != nil {
		eventErr := e.FSM.Event(ctx, RemoveResponseRetry.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
		return
	}

	for _, n := range names {
		if err = event.RemoveFromAllLists(r.session, r.interactionCreate, n); err != nil {
			e.Err = err
			return
		}
	}

	offWaitlistUsers, err := event.PromoteFromWaitlists()
	if err != nil {
		e.Err = err
		return
	}

	for _, u := range offWaitlistUsers {
		if err = event.NotifyUserOffWaitlist(r.session, r.interactionCreate.Interaction, u); err != nil {
			e.Err = err
			return
		}
	}
}

type RemoveResponseRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewRemoveResponseRetryState(o discord.Options) *RemoveResponseRetryState {
	return &RemoveResponseRetryState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,

		inputHandler: NewInputHandler(&o),
	}
}

func (r *RemoveResponseRetryState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := r.session.ChannelMessageSend(r.channel.ID, discord.InvalidRemoveResponseText); err != nil {
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

	wl, ok := event.RoleGroup.Waitlist[role.AcceptedField]
	if !ok {
		e.Err = fmt.Errorf("cannot find accepted field")
		return
	}

	var counter int
	users := make([]string, 0)
	// Braille space is used instead because hard spaces in embeds are not documented
	for _, r := range event.RoleGroup.Roles {
		counter += len(r.Users)
		users = append(users, r.Users...)
	}
	nameMap := map[int]string{}
	for index, name := range append(users, wl.Users...) {
		nameMap[index+1] = name
	}

	if err = r.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
		e.Err = err
		return
	}

	names, err := selectMultiple(e, nameMap)
	if err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		if eventErr != nil {
			e.Err = fmt.Errorf("%v: %v", err, eventErr)
			return
		}
		return
	}

	for _, n := range names {
		if err = event.RemoveFromAllLists(r.session, r.interactionCreate, n); err != nil {
			e.Err = err
			return
		}
	}

	offWaitlistUsers, err := event.PromoteFromWaitlists()
	if err != nil {
		e.Err = err
		return
	}

	for _, u := range offWaitlistUsers {
		if err = event.NotifyUserOffWaitlist(r.session, r.interactionCreate.Interaction, u); err != nil {
			e.Err = err
			return
		}
	}
}

func selectMultiple(e *fsm.Event, nameMap map[int]string) ([]string, error) {
	val, err := Get(e.FSM, discord.MenuOption)
	if err != nil {
		return nil, err
	}
	if !util.IsInputOption(val.(string)) {
		return nil, fmt.Errorf("invalid select input")
	}
	names := make([]string, 0)
	for _, n := range strings.Split(val.(string), " ") {
		option, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}

		n, ok := nameMap[option]
		if !ok {
			return nil, fmt.Errorf("cannot find name")
		}
		names = append(names, n)
	}
	return names, nil
}

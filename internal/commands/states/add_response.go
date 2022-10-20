package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"time"
)

type AddResponseState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewAddResponseState(o discord.Options) *AddResponseState {
	return &AddResponseState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,

		inputHandler: NewInputHandler(&o),
	}
}

func (a *AddResponseState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := a.session.ChannelMessageSendEmbed(a.channel.ID, &discord.EnterUserNameMessage); err != nil {
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

	// TODO: Handle guilds with more than 1000 members
	members, err := a.session.GuildMembers(a.interactionCreate.Interaction.GuildID, "0", 1000)
	if err != nil {
		e.Err = fmt.Errorf("failed to get guild members: %v", err)
		return
	}
	names := make([]string, 0)
	for _, m := range members {
		names = append(names, m.User.Username)
	}

	if err = a.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.Username, 60*time.Second); err != nil {
		e.Err = err
		return
	}

	state, err := a.findUser(names, e, &event)
	if err != nil {
		e.Err = err
		return
	}

	err = e.FSM.Event(ctx, state)
	if err != nil {
		e.Err = err
		return
	}
}

func (a *AddResponseState) findUser(names []string, e *fsm.Event, event *discord.Event) (string, error) {
	result, err := Get(e.FSM, discord.Username)
	if err != nil {
		return "", err
	}
	matches := fuzzy.Find(result.(string), names)
	numMatches := len(matches)
	if numMatches == 0 {
		return UnknownUser.String(), nil
	}
	if numMatches > 1 {
		if _, err = a.session.ChannelMessageSend(a.channel.ID, discord.FoundMultipleText); err != nil {
			return "", err
		}
		return SelfTransition.String(), nil
	}
	user := matches[0]
	for _, r := range event.RoleGroup.Roles {
		if util.ContainsUser(r.Users, user) {
			if _, err = a.session.ChannelMessageSend(a.channel.ID, discord.UserSignedUpText); err != nil {
				return "", err
			}
			return Cancel.String(), nil
		}
	}
	e.FSM.SetMetadata(discord.Username.String(), user)
	return SignUp.String(), nil
}

type UnknownUserState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewUnknownUserState(o discord.Options) *UnknownUserState {
	return &UnknownUserState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,

		inputHandler: NewInputHandler(&o),
	}
}

func (u *UnknownUserState) OnState(ctx context.Context, e *fsm.Event) {
	name, found := e.FSM.Metadata(discord.Username.String())
	if !found {
		e.Err = fmt.Errorf("cannot find username")
	}

	if _, err := u.session.ChannelMessageSendEmbed(u.channel.ID, &discordgo.MessageEmbed{
		Title:       "We couldn't find a Discord user with that name",
		Color:       discord.Purple,
		Description: fmt.Sprintf("**1** Try another name\n**2** Add **%s** as a non Discord user\n**3** Cancel", name),
		Footer: &discordgo.MessageEmbedFooter{
			Text: discord.OptionText,
		},
	}); err != nil {
		e.Err = err
		return
	}

	if err := u.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
		e.Err = err
		return
	}

	state, err := userAddSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, UnknownUserRetry.String())
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

type UnknownUserRetryState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	inputHandler *InputHandler
}

func NewUnknownUserRetryState(o discord.Options) *UnknownUserRetryState {
	return &UnknownUserRetryState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,

		inputHandler: NewInputHandler(&o),
	}
}

func (u *UnknownUserRetryState) OnState(ctx context.Context, e *fsm.Event) {
	if _, err := u.session.ChannelMessageSend(u.channel.ID, discord.InvalidEntryText); err != nil {
		e.Err = err
		return
	}

	if err := u.inputHandler.AwaitInputOrTimeout(ctx, e.FSM, discord.MenuOption, 60*time.Second); err != nil {
		e.Err = err
		return
	}

	state, err := userAddSelect(e)
	if err != nil {
		eventErr := e.FSM.Event(ctx, SelfTransition.String())
		e.Err = fmt.Errorf("%v: %v", err, eventErr)
		return
	}
	err = e.FSM.Event(ctx, state)
	if err != nil {
		e.Err = err
		return
	}
}

func userAddSelect(e *fsm.Event) (string, error) {
	val, err := Get(e.FSM, discord.MenuOption)
	if err != nil {
		return "", err
	}

	opts := map[string]chatState{
		"1": AddResponse,
		"2": SignUp,
		"3": Cancel,
	}
	option, ok := opts[val.(string)]
	if !ok {
		return "", fmt.Errorf("cannot find response")
	}
	return option.String(), nil
}

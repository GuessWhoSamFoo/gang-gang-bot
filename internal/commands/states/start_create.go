package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
)

const CreateAction string = "creation"

type StartCreateState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel

	responseFunc func(*discordgo.Interaction, *discordgo.InteractionResponse) error
}

func NewStartCreateState(o discord.Options) *StartCreateState {
	return &StartCreateState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,

		responseFunc: o.Session.InteractionRespond,
	}
}

func (s *StartCreateState) OnState(_ context.Context, e *fsm.Event) {
	err := s.responseFunc(s.interactionCreate.Interaction, discord.CreateEventMessage(s.interactionCreate.Interaction.GuildID, s.channel.ID))
	if err != nil {
		e.Err = err
		return
	}

	e.FSM.SetMetadata(discord.Action.String(), CreateAction)
	e.FSM.SetMetadata(discord.GuildID.String(), s.interactionCreate.Interaction.GuildID)
	e.FSM.SetMetadata(discord.Owner.String(), s.interactionCreate.Member.User.Username)
	e.FSM.SetMetadata(discord.Color.String(), discord.Purple)
}

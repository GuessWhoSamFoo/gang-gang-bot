package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
)

type ProcessEditState struct {
	*discord.Options
}

func NewProcessEditState(o discord.Options) *ProcessEditState {
	return &ProcessEditState{
		&o,
	}
}

func (p *ProcessEditState) OnState(_ context.Context, e *fsm.Event) {
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
	embed, err := discord.ConvertEventToMessageEmbed(&event)
	if err != nil {
		e.Err = err
		return
	}
	if _, err = p.Options.Session.ChannelMessageEditEmbed(p.Options.InteractionCreate.Interaction.ChannelID, p.Options.InteractionCreate.Interaction.Message.ID, embed); err != nil {
		e.Err = err
		return
	}
	if err = p.Options.UpdateEvent(&event); err != nil {
		e.Err = err
		return
	}
	msg := p.Options.InteractionCreate.Interaction.Message
	if _, err = p.Options.Session.ChannelMessageSendEmbed(p.Options.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Event has been updated!",
		Color:       discord.Purple,
		Description: fmt.Sprintf("[Click here to view the event](https://discord.com/channels/%s/%s/%s)", p.Options.InteractionCreate.GuildID, p.Options.InteractionCreate.Interaction.ChannelID, msg.ID),
	}); err != nil {
		e.Err = fmt.Errorf("failed to send message: %v", err)
		return
	}
}

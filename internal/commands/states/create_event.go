package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"log"
)

type CreateEventState struct {
	*discord.Options
}

func NewCreateEventState(o discord.Options) *CreateEventState {
	return &CreateEventState{
		&o,
	}
}

func (c *CreateEventState) OnState(_ context.Context, e *fsm.Event) {
	if c.Options.InteractionCreate.Interaction.Member == nil {
		e.Err = fmt.Errorf("cannot find user who created event")
		return
	}

	event, err := discord.FromFSMToEvent(e.FSM)
	if err != nil {
		e.Err = err
		return
	}

	if err = c.Options.CreateGoogleEvent(event); err != nil {
		e.Err = err
		return
	}

	acceptedField := discord.AcceptedBase
	for _, r := range event.RoleGroup.Roles {
		if r.FieldName == role.AcceptedField && r.Limit > 0 {
			acceptedField = acceptedField + fmt.Sprintf(" (0/%d)", r.Limit)
		}
	}

	msg, err := c.Options.Session.ChannelMessageSendComplex(c.Options.InteractionCreate.Interaction.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       event.Title,
				Description: event.Description,
				Color:       discord.Purple,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "Time",
						// https://discord.com/developers/docs/reference#message-formatting-timestamp-styles
						Value: util.PrintTime(event.Start, event.End),
					},
					{
						Name:   "Links",
						Value:  util.PrintAddGoogleCalendarLink(event.Title, event.Description, event.Start, event.End),
						Inline: true,
					},
					{
						Name:   "Location",
						Value:  event.Location,
						Inline: true,
					},
					{
						Name:  "Calendar",
						Value: util.PrintGoogleCalendarEventLink(event.ID),
					},
					{
						Name:   acceptedField,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   discord.DeclinedBase,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   discord.TentativeBase,
						Value:  "-",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Created by %v", c.Options.InteractionCreate.Interaction.Member.User.Username),
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discord.AcceptButton,
					discord.DeclineButton,
					discord.TentativeButton,
					discord.EditButton,
					discord.DeleteButton,
				},
			},
		},
	})
	if err != nil {
		e.Err = err
		return
	}

	event.DiscordLink = fmt.Sprintf("https://discord.com/channels/%s/%s/%s", c.Options.InteractionCreate.GuildID, c.Options.InteractionCreate.Interaction.ChannelID, msg.ID)
	_, err = c.Session.ChannelMessageSendEmbed(c.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Event has been created",
		Color:       discord.Purple,
		Description: fmt.Sprintf("[Click here to view the event](%s)", event.DiscordLink),
	})
	if err != nil {
		e.Err = fmt.Errorf("failed to send message: %v", err)
		return
	}

	if err = c.UpdateEvent(event); err != nil {
		e.Err = err
		return
	}

	guildEvent := &discordgo.GuildScheduledEventParams{
		Name:               event.Title,
		Description:        event.DiscordLink,
		ScheduledStartTime: &event.Start,
		ScheduledEndTime:   &event.End,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: event.Location,
		},
		Status: 1,
	}

	_, err = c.Session.GuildScheduledEventCreate(c.Options.InteractionCreate.GuildID, guildEvent)
	if err != nil {
		e.Err = err
		return
	}

	log.Println("Successfully created event")
	return
}

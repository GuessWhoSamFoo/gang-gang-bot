package internal

import (
	"context"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "event",
			Description:              "Create a new event",
			DefaultMemberPermissions: &pkg.EventPermission,
			DMPermission:             &pkg.DMPermission,
		},
		{
			Name:        "my_events",
			Description: "View a list of upcoming events you've organized or signed up for",
		},
		{
			Name:        "upcoming_events",
			Description: "View a list of upcoming events",
		},
		//{
		//	Name:        "edit",
		//	Description: "Modify an existing event",
		//},
	}
)

func (sm *StateManager) CreateEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if sm.HasUser(i.Member.User.ID) {
		discord.NotifyCommandInProgress(s, i)
		return
	}
	sm.AddUser(i.Member.User.ID)
	defer sm.RemoveUser(i.Member.User.ID)

	c, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("cannot create channel: %v", err)
		return
	}
	ctx := context.Background()
	opts := discord.Options{
		Session:           s,
		InteractionCreate: i,
		Channel:           c,
		CalendarClient:    sm.CalendarClient,
	}

	f, err := NewDefaultStateFactory(opts).Factory(commands.CreateType)
	if err != nil {
		return
	}

	if err := f.Event(ctx, states.StartCreate.String()); err != nil {
		log.Println(err)
		return
	}
	if err := f.Event(ctx, states.AddTitle.String()); err != nil {
		log.Println(err)
		return
	}
	if err := f.Event(ctx, states.AddDescription.String()); err != nil {
		log.Println(err)
		return
	}
	if err := f.Event(ctx, states.SetAttendeeLimit.String()); err != nil {
		log.Println(err)
		return
	}
	if err := f.Event(ctx, states.SetDate.String()); err != nil {
		log.Println(err)
		return
	}
	if err := f.Event(ctx, states.SetLocation.String()); err != nil {
		log.Println(err)
		return
	}
	if err := f.Event(ctx, states.SetDuration.String()); err != nil {
		log.Println(err)
		return
	}
	if err := f.Event(ctx, states.CreateEvent.String()); err != nil {
		log.Println(err)
		return
	}
}

func (sm *StateManager) ListMyEventsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if sm == nil {
		log.Printf("commands manager is nil")
		return
	}
	if sm.CalendarClient == nil {
		log.Println("calendar client is nil")
		return
	}
	if i.Member == nil {
		log.Println("cannot find user")
		return
	}

	events, err := sm.CalendarClient.ListEvents()
	if err != nil {
		log.Printf("cannot get events from calendar: %v", err)
		return
	}

	var desc string
	for _, e := range events {
		link, err := util.GetDiscordLinkFromCalendarDescription(e.Description)
		if err != nil {
			log.Printf("failed to get discord link from calendar: %v", err)
			return
		}
		_, channelID, messageID, err := util.GetIDsFromDiscordLink(link)
		if err != nil {
			log.Printf("failed to get id: %v", err)
			return
		}

		msg, err := s.ChannelMessage(channelID, messageID)
		if err != nil {
			log.Printf("failed to get message: %v", err)
			return
		}

		event, err := discord.GetEventFromMessage(msg)
		if err != nil {
			log.Printf("failed to convert message: %v", err)
			return
		}

		if event.Owner == i.Member.User.Username || event.RoleGroup.HasUser(i.Member.User.Username, role.AcceptedField) {
			desc += util.PrintEventListItem(event.Start, event.Title, link)
		}
	}

	if desc == "" {
		desc = "No events found!"
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "My Events",
					Color:       discord.Purple,
					Description: desc,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("failed to reponsd: %v", err)
		return
	}
}

func (sm *StateManager) ListUpcomingEventsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if sm == nil {
		log.Printf("commands manager is nil")
		return
	}
	if sm.CalendarClient == nil {
		log.Println("calendar client is nil")
		return
	}
	events, err := sm.CalendarClient.ListEvents()
	if err != nil {
		log.Printf("cannot list events from Google calendar")
		return
	}

	var desc string
	for _, e := range events {
		link, err := util.GetDiscordLinkFromCalendarDescription(e.Description)
		if err != nil {
			log.Printf("failed to get discord link from calendar: %v", err)
			return
		}
		_, channelID, messageID, err := util.GetIDsFromDiscordLink(link)
		if err != nil {
			log.Printf("failed to get id: %v", err)
			return
		}

		msg, err := s.ChannelMessage(channelID, messageID)
		if err != nil {
			log.Printf("failed to get message: %v", err)
			return
		}

		event, err := discord.GetEventFromMessage(msg)
		if err != nil {
			log.Printf("failed to convert message: %v", err)
			return
		}
		desc += util.PrintEventListItem(event.Start, event.Title, link)
	}

	if desc == "" {
		desc = "No events found!"
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Upcoming Events",
					Color:       discord.Purple,
					Description: desc,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("failed to reponsd: %v", err)
		return
	}
}

//func EditEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	// TODO: Send a DM with edit message sequence
//}

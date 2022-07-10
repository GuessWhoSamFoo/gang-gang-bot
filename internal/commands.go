package internal

import (
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
		//{
		//	Name:        "edit",
		//	Description: "Modify an existing event",
		//},
	}
)

func (sm *StateManager) CreateEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if sm == nil {
		log.Printf("state manager is nil")
		return
	}
	if sm.CalendarClient == nil {
		log.Println("calendar client is nil")
		return
	}

	c, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("cannot create channel: %v", err)
		return
	}

	eb, err := pkg.NewEventBuilder(s, c, i)
	if err != nil {
		log.Printf("cannot create builder: %v", err)
		return
	}

	if sm.HasUser(i.Member.User.ID) {
		pkg.NotifyCommandInProgress(s, i)
		return
	}
	sm.AddUser(i.Member.User.ID)
	defer sm.RemoveUser(i.Member.User.ID)

	if err := eb.StartCreate(); err != nil {
		log.Printf("failed to start event create chat: %v", err)
		return
	}

	if err := eb.AddTitle(); err != nil {
		log.Printf("failed to set event title: %v", err)
		return
	}

	if err := eb.AddDescription(); err != nil {
		log.Printf("failed to set event description: %v", err)
		return
	}

	if err := eb.SetAttendeeLimit(); err != nil {
		log.Printf("failed to set attendee limit: %v", err)
		return
	}

	if err := eb.SetDate(); err != nil {
		log.Printf("failed to set starting date: %v", err)
		return
	}

	if err := eb.SetLocation(); err != nil {
		log.Printf("failed to set location: %v", err)
		return
	}

	if err := eb.SetDuration(); err != nil {
		log.Printf("failed to set duration: %v", err)
		return
	}

	if err := sm.CalendarClient.CreateGoogleEvent(eb.Event); err != nil {
		log.Printf("failed to add event to calendar: %v", err)
		return
	}

	if err := eb.CreateEvent(); err != nil {
		log.Printf("failed to create event: %v", err)
		return
	}

	if err := sm.CalendarClient.UpdateEvent(eb.Event); err != nil {
		log.Printf("failed to add discord link to calendar: %v", err)
		return
	}
}

func (sm *StateManager) ListEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if sm == nil {
		log.Printf("state manager is nil")
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
	if sm.HasUser(i.Member.User.ID) {
		pkg.NotifyCommandInProgress(s, i)
		return
	}

	sm.AddUser(i.Member.User.ID)
	defer sm.RemoveUser(i.Member.User.ID)

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

		event, err := pkg.GetEventFromMessage(msg)
		if err != nil {
			log.Printf("failed to convert message: %v", err)
			return
		}

		if event.Owner == i.Member.User.Username || event.RoleGroup.HasUser(i.Member.User.Username, pkg.AcceptedField) {
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
					Color:       pkg.Purple,
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

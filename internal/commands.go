package internal

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg"
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
		//{
		//	Name:        "my_events",
		//	Description: "View a list of upcoming events you've organized or signed up for",
		//},
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

	if err := eb.SetDuration(); err != nil {
		log.Printf("failed to set duration: %v", err)
		return
	}

	if err := eb.CreateEvent(); err != nil {
		log.Printf("failed to create event: %v", err)
		return
	}
}

//func ListEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	// TODO: Return an ephemeral message of all events one is organizing or signed up for
//}
//
//func EditEventHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	// TODO: Send a DM with edit message sequence
//}

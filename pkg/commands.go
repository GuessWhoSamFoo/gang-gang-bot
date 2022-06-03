package pkg

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	eventPermission int64 = discordgo.PermissionManageEvents
	dmPermission          = true
	purple                = 10181046

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:                     "event",
			Description:              "Create a new event",
			DefaultMemberPermissions: &eventPermission,
			DMPermission:             &dmPermission,
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"event": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			c, err := s.UserChannelCreate(i.Member.User.ID)
			if err != nil {
				log.Printf("cannot create channel: %v", err)
				return
			}

			eb, err := NewEventBuilder(s, c, i)
			if err != nil {
				log.Printf("cannot create builder: %v", err)
				return
			}
			if err := eb.StartChat(); err != nil {
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
		},
	}

	ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		// TODO: ComponentHandlers need to be implemented
		"accept": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("Accept interaction")
			return
		},
		"decline": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("Decline interaction")
			return
		},
		"tentative": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("Tentative interaction")
			return
		},
		"edit": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("Edit interaction")
			return
		},
		"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			log.Println("Delete interaction")
			return
		},
	}
)

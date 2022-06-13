package pkg

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

var (
	ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"accept":        AcceptHandler,
		"decline":       DeclineHandler,
		"tentative":     TentativeHandler,
		"edit":          EditHandler,
		"delete":        DeleteHandler,
		"confirmDelete": ConfirmDeleteHandler,
	}

	AcceptButton = discordgo.Button{
		Label:    "✅",
		Style:    discordgo.SecondaryButton,
		CustomID: "accept",
	}
	DeclineButton = discordgo.Button{
		Label:    "❌",
		Style:    discordgo.SecondaryButton,
		CustomID: "decline",
	}
	TentativeButton = discordgo.Button{
		Label:    "❔",
		Style:    discordgo.SecondaryButton,
		CustomID: "tentative",
	}
	EditButton = discordgo.Button{
		Label:    "Edit",
		Style:    discordgo.PrimaryButton,
		CustomID: "edit",
	}
	DeleteButton = discordgo.Button{
		Label:    "Delete",
		Style:    discordgo.DangerButton,
		CustomID: "delete",
	}
)

func AcceptHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	event, err := getEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}
	newEmbed := &discordgo.MessageEmbed{}
	for _, embed := range i.Message.Embeds {
		for _, f := range embed.Fields {
			if strings.Contains(f.Name, "Accepted") {
				if !util.ContainsUserInField(f.Value, i.Member.User.Username) {
					newName, err := util.IncrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot accept: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.AddUserToField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot accept: %v", err)
						return
					}
					f.Value = newVal
				} else {
					newName, err := util.DecrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot accept: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.RemoveUserFromField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot accept: %v", err)
						return
					}
					f.Value = newVal
				}
			} else if strings.Contains(f.Name, "Declined") || strings.Contains(f.Name, "Tentative") {
				if util.ContainsUserInField(f.Value, i.Member.User.Username) {
					newName, err := util.DecrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot decrease field: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.RemoveUserFromField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot remove user: %v", err)
						return
					}
					f.Value = newVal
				}
			}
		}
		newEmbed = embed
	}

	if event.limit == -1 || event.accepted < event.limit {
		if _, err := s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, newEmbed); err != nil {
			log.Printf("cannot accept: %v", err)
		}
	} else {
		// TODO: Implement waitlist
	}
	return
}

func DeclineHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	newEmbed := &discordgo.MessageEmbed{}
	for _, embed := range i.Message.Embeds {
		for _, f := range embed.Fields {
			if strings.Contains(f.Name, "Declined") {
				if !util.ContainsUserInField(f.Value, i.Member.User.Username) {
					newName, err := util.IncrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot decline: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.AddUserToField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot decline: %v", err)
						return
					}
					f.Value = newVal
				} else {
					newName, err := util.DecrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot decline: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.RemoveUserFromField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot decline: %v", err)
						return
					}
					f.Value = newVal
				}
			} else if strings.Contains(f.Name, "Accepted") || strings.Contains(f.Name, "Tentative") {
				if util.ContainsUserInField(f.Value, i.Member.User.Username) {
					newName, err := util.DecrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot decrease field: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.RemoveUserFromField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot remove user: %v", err)
						return
					}
					f.Value = newVal
				}
			}
		}
		newEmbed = embed
	}

	if _, err := s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, newEmbed); err != nil {
		log.Printf("cannot decline: %v", err)
	}
	return
}

func TentativeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	newEmbed := &discordgo.MessageEmbed{}
	for _, embed := range i.Message.Embeds {
		for _, f := range embed.Fields {
			if strings.Contains(f.Name, "Tentative") {
				if !util.ContainsUserInField(f.Value, i.Member.User.Username) {
					newName, err := util.IncrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot set tentative: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.AddUserToField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot set tentative: %v", err)
						return
					}
					f.Value = newVal
				} else {
					newName, err := util.DecrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot set tentative: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.RemoveUserFromField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot set tentative: %v", err)
						return
					}
					f.Value = newVal
				}
			} else if strings.Contains(f.Name, "Declined") || strings.Contains(f.Name, "Accepted") {
				if util.ContainsUserInField(f.Value, i.Member.User.Username) {
					newName, err := util.DecrementFieldCounter(f.Name)
					if err != nil {
						log.Printf("cannot decrease field: %v", err)
						return
					}
					f.Name = newName
					newVal, err := util.RemoveUserFromField(f.Value, i.Member.User.Username)
					if err != nil {
						log.Printf("cannot remove user: %v", err)
						return
					}
					f.Value = newVal
				}
			}
		}
		newEmbed = embed
	}

	if _, err := s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, newEmbed); err != nil {
		log.Printf("cannot set tenative: %v", err)
	}
	return
}

func EditHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Check server permissions before edit
	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func DeleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Check server permissions before delete
	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	c, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("failed to get channel: %v", err)
		return
	}

	event, err := getEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}

	if _, err := s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Confirm event deletion",
				Description: fmt.Sprintf("[%s](https://discord.com/channels/%s/%s/%s)", event.title, i.GuildID, i.ChannelID, i.Message.ID),
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Delete this event",
						Style:    discordgo.DangerButton,
						CustomID: "confirmDelete",
					},
				},
			},
		},
	}); err != nil {
		log.Printf("failed to send message: %v", err)
		return
	}
}

func ConfirmDeleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	event, err := getEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to parse event from interaction: %v", err)
		return
	}

	url, err := util.GetLinkFromDeleteDescription(event.description)
	if err != nil {
		log.Printf("failed to parse link from description: %v", err)
		return
	}

	_, channelID, messageID, err := util.GetIDsFromDeleteLink(url)
	if err != nil {
		log.Printf("failed to get ID from link: %v", err)
		return
	}

	if err := s.ChannelMessageDelete(channelID, messageID); err != nil {
		log.Printf("failed to delete message: %v", err)
		return
	}

	if _, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Event deleted",
			},
		},
		ID:         i.Message.ID,
		Channel:    i.Message.ChannelID,
		Components: []discordgo.MessageComponent{},
	}); err != nil {
		log.Printf("failed to edit message: %v", err)
		return
	}

}

func getEventFromMessage(msg *discordgo.Message) (*Event, error) {
	if len(msg.Embeds) != 1 {
		return nil, fmt.Errorf("expected 1 embed: got %d", len(msg.Embeds))
	}
	e := &Event{}
	for _, embed := range msg.Embeds {
		e.title = embed.Title
		e.description = embed.Description
		for _, f := range embed.Fields {
			count, limit, err := util.ParseFieldHeadCount(f.Name)
			if err != nil {
				return nil, err
			}
			switch {
			case strings.Contains(f.Name, "Accepted"):
				if count > 0 {
					e.accepted = count
				}
				e.limit = limit
				e.acceptedNames = util.GetUsersFromValues(f.Value)
			case strings.Contains(f.Name, "Declined"):
				if count > 0 {
					e.declined = count
				}
				e.declined = count
				e.declinedNames = util.GetUsersFromValues(f.Value)
			case strings.Contains(f.Name, "Tentative"):
				if count > 0 {
					e.tentative = count
				}
				e.tentativeNames = util.GetUsersFromValues(f.Value)
			default:
				// TODO: time, links, and waitlist
			}
		}
	}
	return e, nil
}

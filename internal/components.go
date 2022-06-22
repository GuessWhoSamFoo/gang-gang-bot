package internal

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"log"
)

func (sm *StateManager) AcceptHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}); err != nil {
		log.Println(err)
	}
	e, err := pkg.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}

	if err := e.ToggleAccept(s, i, i.Member.User.Username); err != nil {
		log.Printf("toggle accept: %v", err)
		return
	}

	embed, err := pkg.ConvertEventToMessageEmbed(e)
	if err != nil {
		log.Printf("failed to convert event: %v", err)
		return
	}

	if _, err := s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed); err != nil {
		log.Printf("failed to edit embed: %v", err)
		return
	}
	log.Printf("User: %s accepted event %s", i.Member.User.Username, fmt.Sprintf("%s/%s", i.ChannelID, i.Message.ID))
	return
}

func (sm *StateManager) DeclineHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}); err != nil {
		log.Println(err)
	}
	e, err := pkg.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}
	if err := e.ToggleDecline(s, i, i.Member.User.Username); err != nil {
		log.Printf("toggle decline: %v", err)
		return
	}

	embed, err := pkg.ConvertEventToMessageEmbed(e)
	if err != nil {
		log.Printf("failed to convert event: %v", err)
		return
	}

	if _, err := s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed); err != nil {
		log.Printf("cannot decline: %v", err)
		return
	}
	log.Printf("User: %s declined event %s", i.Member.User.Username, fmt.Sprintf("%s/%s", i.ChannelID, i.Message.ID))
	return
}

func (sm *StateManager) TentativeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}); err != nil {
		log.Println(err)
	}
	e, err := pkg.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}

	if err := e.ToggleTentative(s, i, i.Member.User.Username); err != nil {
		log.Printf("toggle tentative: %v", err)
		return
	}

	embed, err := pkg.ConvertEventToMessageEmbed(e)
	if err != nil {
		log.Printf("failed to convert event: %v", err)
		return
	}

	if _, err := s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, embed); err != nil {
		log.Printf("cannot set tenative: %v", err)
		return
	}
	log.Printf("User: %s marked tentative %s", i.Member.User.Username, fmt.Sprintf("%s/%s", i.ChannelID, i.Message.ID))
	return
}

func (sm *StateManager) EditHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Check server permissions before edit
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}); err != nil {
		log.Println(err)
	}

	c, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("cannot create channel: %v", err)
		return
	}

	if sm.ActiveMap.HasUser(i.Member.User.ID) {
		pkg.NotifyCommandInProgress(s, i)
		return
	}
	sm.AddUser(i.Member.User.ID)
	defer sm.RemoveUser(i.Member.User.ID)

	eb, err := pkg.NewEventBuilder(s, c, i)
	if err != nil {
		log.Printf("cannot creater builder: %v", err)
		return
	}

	if err := eb.StartEdit(); err != nil {
		log.Printf("failed to edit event: %v", err)
		return
	}

	if err := eb.ProcessEdit(); err != nil {
		log.Printf("failed to confirm edit: %v", err)
		return
	}
}

func (sm *StateManager) DeleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Check server permissions before delete
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}); err != nil {
		log.Println(err)
	}

	c, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("failed to get channel: %v", err)
		return
	}

	e, err := pkg.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}

	if _, err := s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Confirm event deletion",
				Description: fmt.Sprintf("[%s](https://discord.com/channels/%s/%s/%s)", e.Title, i.GuildID, i.ChannelID, i.Message.ID),
				Color:       pkg.Purple,
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

func (sm *StateManager) ConfirmDeleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}); err != nil {
		log.Println(err)
	}

	e, err := pkg.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to parse event from interaction: %v", err)
		return
	}

	url, err := util.GetLinkFromDeleteDescription(e.Description)
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
				Color: pkg.Purple,
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

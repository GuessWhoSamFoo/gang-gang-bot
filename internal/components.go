package internal

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
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
	e, err := discord.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}

	if err := e.ToggleAccept(s, i, i.Member.User.Username); err != nil {
		log.Printf("toggle accept: %v", err)
		return
	}

	embed, err := discord.ConvertEventToMessageEmbed(e)
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
	e, err := discord.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}
	if err := e.ToggleDecline(s, i, i.Member.User.Username); err != nil {
		log.Printf("toggle decline: %v", err)
		return
	}

	embed, err := discord.ConvertEventToMessageEmbed(e)
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
	e, err := discord.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}

	if err := e.ToggleTentative(s, i, i.Member.User.Username); err != nil {
		log.Printf("toggle tentative: %v", err)
		return
	}

	embed, err := discord.ConvertEventToMessageEmbed(e)
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
	if sm.CalendarClient == nil {
		log.Println("calendar client is nil")
		return
	}

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
		discord.NotifyCommandInProgress(s, i)
		return
	}
	sm.AddUser(i.Member.User.ID)
	defer sm.RemoveUser(i.Member.User.ID)
	ctx := context.Background()
	opts := discord.Options{
		Session:           s,
		InteractionCreate: i,
		Channel:           c,
		CalendarClient:    sm.CalendarClient,
	}

	f, err := NewDefaultStateFactory(opts).Factory(commands.EditType)
	if err != nil {
		return
	}

	if err = f.Event(ctx, states.StartEdit.String()); err != nil {
		log.Println(err)
		return
	}
	if err = f.Event(ctx, states.ProcessEdit.String()); err != nil {
		log.Println(err)
		return
	}
	log.Printf("User: %s edited event %s", i.Member.User.Username, fmt.Sprintf("%s/%s", i.ChannelID, i.Message.ID))
}

func (sm *StateManager) DeleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: Check server permissions before delete
	if sm == nil || sm.CalendarClient == nil {
		log.Printf("cannot find commands manager email client")
		return
	}

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

	e, err := discord.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return
	}
	if i.Member.User.Username != e.Owner && i.Interaction.Member.Permissions&discordgo.PermissionManageEvents == 0 {
		if _, err := s.ChannelMessageSendEmbed(c.ID, discord.DeleteInsufficientPermissionMessage); err != nil {
			log.Printf("failed to send message: %v", err)
			return
		}
		log.Printf("insufficient permissions to delete %s", i.Interaction.Message.ID)
		return
	}

	if _, err := s.ChannelMessageSendComplex(c.ID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Confirm event deletion",
				Description: fmt.Sprintf("[%s](https://discord.com/channels/%s/%s/%s)", e.Title, i.GuildID, i.ChannelID, i.Message.ID),
				Color:       discord.Purple,
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
	log.Printf("User: %s deleted event %s", i.Member.User.Username, fmt.Sprintf("%s/%s", i.ChannelID, i.Message.ID))
}

func (sm *StateManager) ConfirmDeleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if sm == nil || sm.CalendarClient == nil {
		log.Printf("cannot find commands manager email client")
		return
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
	}); err != nil {
		log.Println(err)
	}

	e, err := discord.GetEventFromMessage(i.Message)
	if err != nil {
		log.Printf("failed to parse event from interaction: %v", err)
		return
	}

	url, err := util.GetLinkFromDeleteDescription(e.Description)
	if err != nil {
		log.Printf("failed to parse link from description: %v", err)
		return
	}

	_, channelID, messageID, err := util.GetIDsFromDiscordLink(url)
	if err != nil {
		log.Printf("failed to get ID from link: %v", err)
		return
	}

	msg, err := s.ChannelMessage(channelID, messageID)
	if err != nil {
		return
	}
	cEvent, err := discord.GetEventFromMessage(msg)
	if err != nil {
		return
	}

	if err := sm.CalendarClient.DeleteEvent(cEvent); err != nil {
		log.Printf("cannot delete google event: %v", err)
		return
	}

	if err = s.ChannelMessageDelete(channelID, messageID); err != nil {
		log.Printf("failed to delete message: %v", err)
		return
	}

	if _, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Event deleted",
				Color: discord.Purple,
			},
		},
		ID:         i.Message.ID,
		Channel:    i.Message.ChannelID,
		Components: []discordgo.MessageComponent{},
	}); err != nil {
		log.Printf("failed to edit message: %v", err)
		return
	}

	events, err := s.GuildScheduledEvents(sm.Config.Discord.GuildID, false)
	if err != nil {
		log.Printf("cannot get guild events: %v", err)
		return
	}

	for _, guildEvent := range events {
		if guildEvent.Description == cEvent.Description {
			err = s.GuildScheduledEventDelete(sm.Config.Discord.GuildID, guildEvent.ID)
			if err != nil {
				log.Printf("failed to delete guild event: %v", err)
				return
			}
		}
	}
}

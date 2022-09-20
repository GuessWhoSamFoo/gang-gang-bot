package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Options struct {
	Session           *discordgo.Session
	InteractionCreate *discordgo.InteractionCreate
	Channel           *discordgo.Channel

	*CalendarClient
}

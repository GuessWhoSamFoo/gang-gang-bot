package services

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

// NewDiscordSession creates a new discord session
func NewDiscordSession(token string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	return session, nil
}

func ReadyEvent(_ *discordgo.Session, _ *discordgo.Ready) {
	log.Println("Discord session is up")
}

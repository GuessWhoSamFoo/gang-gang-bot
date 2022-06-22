package services

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type Bot struct {
	Session    *discordgo.Session
	CommandMap map[string]string // id:name
	Config     *internal.Config
}

func NewBot(c *internal.Config) (*Bot, error) {
	return &Bot{
		CommandMap: map[string]string{},
		Config:     c,
	}, nil
}

func (b *Bot) Start() error {
	if b.Session != nil {
		return fmt.Errorf("session exists")
	}
	var err error
	b.Session, err = NewDiscordSession(b.Config.Secret.Token)
	if err != nil {
		return err
	}

	l, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}
	time.Local = l
	b.Session.AddHandler(readyEvent)

	sm := internal.NewStateManager()
	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := sm.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			} else {
				log.Fatalln("cannot add command handler")
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := sm.ComponentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			} else {
				log.Fatalln("cannot add component handler")
			}
		default:
			log.Println("unknown handler type")
		}
	})
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("cannot open session: %v", err)
	}

	for _, v := range internal.Commands {
		c, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.Config.Discord.GuildID, v)
		if err != nil {
			return err
		}
		b.CommandMap[c.ID] = c.Name
	}
	return nil
}

func (b *Bot) Close() error {
	if b.Session == nil {
		return fmt.Errorf("nil session")
	}
	for id, name := range b.CommandMap {
		log.Println("removing command /" + name)
		if err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, b.Config.Discord.GuildID, id); err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
	return b.Session.Close()
}

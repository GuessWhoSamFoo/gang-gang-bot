package internal

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/services"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type Bot struct {
	Session    *discordgo.Session
	CommandMap map[string]string // id:name
	Config     *Config
}

func NewBot(c *Config) (*Bot, error) {
	return &Bot{
		CommandMap: map[string]string{},
		Config:     c,
	}, nil
}

func (b *Bot) Start() error {
	ctx := context.Background()
	sm := NewStateManager(b.Config)

	if b.Session != nil {
		return fmt.Errorf("session exists")
	}

	var err error
	c, err := services.NewCalendarClient(ctx, b.Config.Google.CalendarID, b.Config.Google.Credentials)
	if err != nil {
		return err
	}
	sm.CalendarClient = c

	b.Session, err = services.NewDiscordSession(b.Config.Secret.Token)
	if err != nil {
		return err
	}

	l, err := time.LoadLocation(util.StaticLocation)
	if err != nil {
		return err
	}
	time.Local = l
	b.Session.AddHandler(services.ReadyEvent)
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

	for _, v := range Commands {
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

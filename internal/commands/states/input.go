package states

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

type InputHandler struct {
	*discord.Options
	handlerFunc func(*discordgo.Session, *discordgo.MessageCreate)
	inputChan   chan string
}

func NewInputHandler(o *discord.Options) *InputHandler {
	i := make(chan string)
	return &InputHandler{
		Options: o,
		handlerFunc: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.ChannelID == o.Channel.ID {
				if m.Content == "" || m.Content == "\n" {
					return
				}
				i <- m.Content
			}
		},
		inputChan: i,
	}
}

func (ih *InputHandler) AwaitInputOrTimeout(ctx context.Context, f *fsm.FSM, key discord.MetadataKey, wait time.Duration) error {
	cancelFunc := ih.Options.Session.AddHandler(ih.handlerFunc)
	defer cancelFunc()
	// NewTimer is used instead because time.After can leak memory if the timer doesn't fire
	timer := time.NewTimer(wait)
	select {
	case result := <-ih.inputChan:
		if strings.EqualFold(result, "cancel") {
			err := f.Event(ctx, Cancel.String())
			return fmt.Errorf("event action canceled: %v", err)
		}
		f.SetMetadata(key.String(), result)
	case <-timer.C:
		err := f.Event(ctx, Timeout.String())
		return fmt.Errorf("event action timed out: %v", err)
	}
	return nil
}

func Get(f *fsm.FSM, key discord.MetadataKey) (interface{}, error) {
	val, exists := f.Metadata(key.String())
	if !exists {
		return "", fmt.Errorf("key %s does not exist", key.String())
	}
	return val, nil
}

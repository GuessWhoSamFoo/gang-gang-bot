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

func AwaitInputOrTimeout(ctx context.Context, wait time.Duration, s *discordgo.Session, input chan string, f *fsm.FSM, hf func(*discordgo.Session, *discordgo.MessageCreate), key discord.MetadataKey) error {
	cancelFunc := s.AddHandler(hf)
	defer cancelFunc()
	// NewTimer is used instead because time.After can leak memory if the timer doesn't fire
	timer := time.NewTimer(wait)
	select {
	case result := <-input:
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

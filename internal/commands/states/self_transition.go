package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
)

// SelfTransitionState is an intermediary state that loops back to the previous state
type SelfTransitionState struct {
	session           *discordgo.Session
	interactionCreate *discordgo.InteractionCreate
	channel           *discordgo.Channel
}

func NewSelfTransitionState(o discord.Options) *SelfTransitionState {
	return &SelfTransitionState{
		session:           o.Session,
		interactionCreate: o.InteractionCreate,
		channel:           o.Channel,
	}
}

func (s *SelfTransitionState) OnState(ctx context.Context, e *fsm.Event) {
	err := e.FSM.Event(ctx, e.Src)
	if err != nil {
		e.Err = err
		return
	}
}

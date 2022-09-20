package commands

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states"
)

type FSMState interface {
	OnState(ctx context.Context, e *fsm.Event)
}

type ActionType string

const (
	CreateType = ActionType(states.CreateAction)
	EditType   = ActionType(states.EditAction)
)

func (a ActionType) String() string {
	return string(a)
}

func InitState() string {
	return states.Idle.String()
}

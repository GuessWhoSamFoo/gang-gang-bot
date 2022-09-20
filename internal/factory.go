package internal

import (
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
)

type StateFactory interface {
	Factory(action commands.ActionType) (*fsm.FSM, error)
}

type DefaultStateFactory struct {
	typeFactory map[commands.ActionType]*fsm.FSM
}

var _ StateFactory = DefaultStateFactory{}

func NewDefaultStateFactory(o discord.Options) *DefaultStateFactory {
	return &DefaultStateFactory{
		typeFactory: map[commands.ActionType]*fsm.FSM{
			commands.CreateType: fsm.NewFSM(
				commands.InitState(),
				commands.CreateEvents(),
				commands.CreateTransitions(o),
			),
			commands.EditType: fsm.NewFSM(
				commands.InitState(),
				commands.EditEvents(),
				commands.EditTransitions(o),
			),
		},
	}
}

func (d DefaultStateFactory) Factory(action commands.ActionType) (*fsm.FSM, error) {
	stateMachine, ok := d.typeFactory[action]
	if !ok {
		return nil, fmt.Errorf("factory type not found")
	}
	return stateMachine, nil
}

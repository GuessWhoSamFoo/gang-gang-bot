package commands

import (
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
)

func CreateEvents() fsm.Events {
	return fsm.Events{
		{
			Name: states.StartCreate.String(),
			Src:  []string{states.Idle.String()},
			Dst:  states.StartCreate.String(),
		},
		{
			Name: states.AddTitle.String(),
			Src:  []string{states.StartCreate.String()},
			Dst:  states.AddTitle.String(),
		},
		{
			Name: states.AddDescription.String(),
			Src:  []string{states.AddTitle.String()},
			Dst:  states.AddDescription.String(),
		},
		{
			Name: states.SetAttendeeLimit.String(),
			Src:  []string{states.AddDescription.String()},
			Dst:  states.SetAttendeeLimit.String(),
		},
		{
			Name: states.SetAttendeeRetry.String(),
			Src:  []string{states.SetAttendeeLimit.String(), states.SelfTransition.String()},
			Dst:  states.SetAttendeeRetry.String(),
		},
		{
			Name: states.SetDate.String(),
			Src:  []string{states.SetAttendeeLimit.String(), states.SetAttendeeRetry.String()},
			Dst:  states.SetDate.String(),
		},
		{
			Name: states.SetDateRetry.String(),
			Src:  []string{states.SetDate.String(), states.SelfTransition.String()},
			Dst:  states.SetDateRetry.String(),
		},
		{
			Name: states.SetLocation.String(),
			Src:  []string{states.SetDate.String(), states.SetDateRetry.String()},
			Dst:  states.SetLocation.String(),
		},
		{
			Name: states.SetDuration.String(),
			Src:  []string{states.SetLocation.String()},
			Dst:  states.SetDuration.String(),
		},
		{
			Name: states.SetDurationRetry.String(),
			Src:  []string{states.SetDuration.String(), states.SelfTransition.String()},
			Dst:  states.SetDurationRetry.String(),
		},
		{
			Name: states.CreateEvent.String(),
			Src:  []string{states.SetDuration.String(), states.SetDurationRetry.String()},
			Dst:  states.CreateEvent.String(),
		},
		{
			Name: states.SelfTransition.String(),
			Src: []string{
				states.SetAttendeeRetry.String(),
				states.SetDateRetry.String(),
				states.SetDurationRetry.String(),
			},
			Dst: states.SelfTransition.String(),
		},
		{
			Name: states.Cancel.String(),
			Src: []string{
				states.AddTitle.String(),
				states.AddDescription.String(),
				states.SetAttendeeLimit.String(),
				states.SetAttendeeRetry.String(),
				states.SetDate.String(),
				states.SetDateRetry.String(),
				states.SetLocation.String(),
				states.SetDuration.String(),
				states.SetDurationRetry.String(),
			},
			Dst: states.Cancel.String(),
		},
		{
			Name: states.Timeout.String(),
			Src: []string{
				states.AddTitle.String(),
				states.AddDescription.String(),
				states.SetAttendeeLimit.String(),
				states.SetAttendeeRetry.String(),
				states.SetDate.String(),
				states.SetDateRetry.String(),
				states.SetLocation.String(),
				states.SetDuration.String(),
				states.SetDurationRetry.String(),
			},
			Dst: states.Timeout.String(),
		},
	}
}

func CreateEventStates(o discord.Options) map[string]FSMState {
	return map[string]FSMState{
		states.Cancel.String():           states.NewCancelState(o),
		states.Timeout.String():          states.NewTimeoutState(o),
		states.StartCreate.String():      states.NewStartCreateState(o),
		states.AddTitle.String():         states.NewAddTitleState(o),
		states.AddDescription.String():   states.NewAddDescriptionState(o),
		states.SetAttendeeLimit.String(): states.NewSetAttendeeState(o),
		states.SetAttendeeRetry.String(): states.NewSetAttendeeRetryState(o),
		states.SetDate.String():          states.NewSetDateState(o),
		states.SetDateRetry.String():     states.NewSetDateRetryState(o),
		states.SetLocation.String():      states.NewSetLocationState(o),
		states.SetDuration.String():      states.NewDurationState(o),
		states.SetDurationRetry.String(): states.NewDurationRetryState(o),
		states.CreateEvent.String():      states.NewCreateEventState(o),
		states.SelfTransition.String():   states.NewSelfTransitionState(o),
	}
}

func CreateTransitions(o discord.Options) fsm.Callbacks {
	all := CreateEventStates(o)
	callbacks := fsm.Callbacks{}
	for k, v := range all {
		callbacks[k] = v.OnState
	}
	return callbacks
}

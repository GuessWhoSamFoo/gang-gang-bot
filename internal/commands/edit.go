package commands

import (
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
)

func EditEvents() fsm.Events {
	return fsm.Events{
		{
			Name: states.StartEdit.String(),
			Src:  []string{states.Idle.String(), states.ContinueEdit.String()},
			Dst:  states.StartEdit.String(),
		},
		{
			Name: states.StartEditRetry.String(),
			Src:  []string{states.StartEdit.String(), states.SelfTransition.String()},
			Dst:  states.StartEditRetry.String(),
		},
		{
			Name: states.ModifyEvent.String(),
			Src: []string{
				states.StartEdit.String(),
				states.ContinueEdit.String(),
				states.ContinueEditRetry.String(),
			},
			Dst: states.ModifyEvent.String(),
		},
		{
			Name: states.ModifyEventRetry.String(),
			Src:  []string{states.ModifyEvent.String(), states.SelfTransition.String()},
			Dst:  states.ModifyEventRetry.String(),
		},
		{
			Name: states.RemoveResponse.String(),
			Src:  []string{states.StartEdit.String(), states.StartEditRetry.String()},
			Dst:  states.RemoveResponse.String(),
		},
		{
			Name: states.RemoveResponseRetry.String(),
			Src:  []string{states.RemoveResponse.String(), states.SelfTransition.String()},
			Dst:  states.RemoveResponseRetry.String(),
		},
		{
			Name: states.AddResponse.String(),
			Src: []string{
				states.StartEdit.String(),
				states.SelfTransition.String(),
				states.UnknownUser.String(),
				states.UnknownUserRetry.String(),
			},
			Dst: states.AddResponse.String(),
		},
		{
			Name: states.UnknownUser.String(),
			Src:  []string{states.AddResponse.String()},
			Dst:  states.UnknownUser.String(),
		},
		{
			Name: states.UnknownUserRetry.String(),
			Src:  []string{states.UnknownUser.String(), states.SelfTransition.String()},
			Dst:  states.UnknownUserRetry.String(),
		},
		{
			Name: states.SignUp.String(),
			Src: []string{
				states.AddResponse.String(),
				states.UnknownUser.String(),
				states.UnknownUserRetry.String(),
			},
			Dst: states.SignUp.String(),
		},
		{
			Name: states.SignUpRetry.String(),
			Src: []string{
				states.SignUp.String(),
				states.SelfTransition.String(),
			},
			Dst: states.SignUpRetry.String(),
		},
		{
			Name: states.AddTitle.String(),
			Src:  []string{states.ModifyEvent.String()},
			Dst:  states.AddTitle.String(),
		},
		{
			Name: states.AddDescription.String(),
			Src:  []string{states.ModifyEvent.String()},
			Dst:  states.AddDescription.String(),
		},
		{
			Name: states.SetDate.String(),
			Src:  []string{states.ModifyEvent.String()},
			Dst:  states.SetDate.String(),
		},
		{
			Name: states.SetLocation.String(),
			Src:  []string{states.ModifyEvent.String()},
			Dst:  states.SetLocation.String(),
		},
		{
			Name: states.ContinueEdit.String(),
			Src: []string{
				states.AddTitle.String(),
				states.AddDescription.String(),
				states.SetDate.String(),
				states.SetLocation.String(),
			},
			Dst: states.ContinueEdit.String(),
		},
		{
			Name: states.ContinueEditRetry.String(),
			Src:  []string{states.ContinueEdit.String(), states.SelfTransition.String()},
			Dst:  states.ContinueEditRetry.String(),
		},
		{
			Name: states.ProcessEdit.String(),
			Src: []string{
				states.ContinueEdit.String(),
				states.ContinueEditRetry.String(),
				states.SignUp.String(),
				states.SignUpRetry.String(),
				states.RemoveResponse.String(),
				states.RemoveResponseRetry.String(),
			},
			Dst: states.ProcessEdit.String(),
		},
		{
			Name: states.SelfTransition.String(),
			Src: []string{
				states.AddResponse.String(),
				states.StartEditRetry.String(),
				states.ModifyEventRetry.String(),
				states.ContinueEditRetry.String(),
				states.SignUpRetry.String(),
				states.RemoveResponseRetry.String(),
				states.UnknownUserRetry.String(),
			},
			Dst: states.SelfTransition.String(),
		},
		{
			Name: states.Cancel.String(),
			Src: []string{
				states.StartEdit.String(),
				states.ModifyEvent.String(),
				states.RemoveResponse.String(),
				states.RemoveResponseRetry.String(),
				states.AddResponse.String(),
				states.AddTitle.String(),
				states.AddDescription.String(),
				states.SetDate.String(),
				states.SetLocation.String(),
				states.ContinueEdit.String(),
				states.ContinueEditRetry.String(),
				states.UnknownUser.String(),
				states.UnknownUserRetry.String(),
				states.SignUp.String(),
				states.SignUpRetry.String(),
			},
			Dst: states.Cancel.String(),
		},
		{
			Name: states.Timeout.String(),
			Src: []string{
				states.StartEdit.String(),
				states.ModifyEvent.String(),
				states.RemoveResponse.String(),
				states.RemoveResponseRetry.String(),
				states.AddResponse.String(),
				states.AddTitle.String(),
				states.AddDescription.String(),
				states.SetDate.String(),
				states.SetLocation.String(),
				states.ContinueEdit.String(),
				states.ContinueEditRetry.String(),
				states.UnknownUser.String(),
				states.UnknownUserRetry.String(),
				states.SignUp.String(),
				states.SignUpRetry.String(),
			},
			Dst: states.Timeout.String(),
		},
	}
}

func EditEventStates(o discord.Options) map[string]FSMState {
	return map[string]FSMState{
		states.Cancel.String():              states.NewCancelState(o),
		states.Timeout.String():             states.NewTimeoutState(o),
		states.StartEdit.String():           states.NewStartEditState(o),
		states.StartEditRetry.String():      states.NewStartEditRetryState(o),
		states.ModifyEvent.String():         states.NewModifyEventState(o),
		states.ModifyEventRetry.String():    states.NewModifyEventRetryState(o),
		states.AddTitle.String():            states.NewAddTitleState(o),
		states.AddDescription.String():      states.NewAddDescriptionState(o),
		states.SetDate.String():             states.NewSetDateState(o),
		states.SetLocation.String():         states.NewSetLocationState(o),
		states.ContinueEdit.String():        states.NewContinueEditState(o),
		states.ContinueEditRetry.String():   states.NewContinueEditRetryState(o),
		states.RemoveResponse.String():      states.NewRemoveResponseState(o),
		states.RemoveResponseRetry.String(): states.NewRemoveResponseRetryState(o),
		states.ProcessEdit.String():         states.NewProcessEditState(o),
		states.AddResponse.String():         states.NewAddResponseState(o),
		states.UnknownUser.String():         states.NewUnknownUserState(o),
		states.UnknownUserRetry.String():    states.NewUnknownUserRetryState(o),
		states.SignUp.String():              states.NewSignUpState(o),
		states.SignUpRetry.String():         states.NewSignUpRetryState(o),
		states.SelfTransition.String():      states.NewSelfTransitionState(o),
	}
}

func EditTransitions(o discord.Options) fsm.Callbacks {
	callbacks := fsm.Callbacks{}
	for k, v := range EditEventStates(o) {
		callbacks[k] = v.OnState
	}
	return callbacks
}

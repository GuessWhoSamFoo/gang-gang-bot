package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSetAttendeeState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewSetAttendeeState(*opts)
	assert.NotNil(t, s)
}

func TestSetAttendeeState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	ctx := context.Background()

	s := NewSetAttendeeState(*opts)

	f := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: SetAttendeeLimit.String(),
				Src:  []string{"idle"},
				Dst:  SetAttendeeLimit.String(),
			},
		},
		fsm.Callbacks{
			SetAttendeeLimit.String(): s.OnState,
		},
	)

	expected := role.NewDefaultRoleGroup()
	limit := "50"
	expected.SetLimit(role.AcceptedField, 50)

	go func() {
		s.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
			s.input <- limit
		}
		s.handlerFunc(opts.Session, &discordgo.MessageCreate{})
	}()

	err = f.Event(ctx, SetAttendeeLimit.String())
	assert.NoError(t, err)

	actual, err := Get(f, discord.Attendee)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

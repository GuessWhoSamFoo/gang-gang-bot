package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/ewohltman/discordgo-mock/mockconstants"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStartCreateState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)

	s := NewStartCreateState(*opts)
	assert.NotNil(t, s)
}

func TestNewStartCreateState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)

	s := NewStartCreateState(*opts)
	s.responseFunc = mock.NewInteractionResponse

	f := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: StartCreate.String(),
				Src:  []string{"idle"},
				Dst:  StartCreate.String(),
			},
		},
		fsm.Callbacks{
			StartCreate.String(): s.OnState,
		},
	)
	err = f.Event(context.TODO(), StartCreate.String())
	assert.NoError(t, err)

	action, err := Get(f, discord.Action)
	assert.NoError(t, err)
	guildID, err := Get(f, discord.GuildID)
	assert.NoError(t, err)
	owner, err := Get(f, discord.Owner)
	assert.NoError(t, err)
	color, err := Get(f, discord.Color)
	assert.NoError(t, err)

	assert.Equal(t, CreateAction, action)
	assert.Equal(t, mockconstants.TestGuild, guildID)
	assert.Equal(t, mockconstants.TestUser, owner)
	assert.Equal(t, discord.Purple, color)
}

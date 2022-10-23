package states

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCreateEventState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)

	s := NewCreateEventState(*opts)
	assert.NotNil(t, s)
}

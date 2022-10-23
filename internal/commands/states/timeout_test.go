package states

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTimeoutState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewTimeoutState(*opts)
	assert.NotNil(t, s)
}

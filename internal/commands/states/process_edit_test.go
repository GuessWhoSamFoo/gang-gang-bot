package states

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProcessEditState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewProcessEditState(*opts)
	assert.NotNil(t, s)
}

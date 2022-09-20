package states

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTimeoutState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewTimeoutState(*opts)
	assert.NotNil(t, s)
}

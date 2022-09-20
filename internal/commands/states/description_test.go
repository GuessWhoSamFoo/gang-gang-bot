package states

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAddDescriptionState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)
	s := NewAddDescriptionState(*opts)
	assert.NotNil(t, s)
}

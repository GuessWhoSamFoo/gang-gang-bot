package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAwaitInputOrTimeout(t *testing.T) {
	ctx := context.TODO()
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	ts := NewTimeoutState(*opts)

	f := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: Timeout.String(),
				Src:  []string{"idle"},
				Dst:  Timeout.String(),
			},
		},
		fsm.Callbacks{
			Timeout.String(): ts.OnState,
		})
	err = ts.inputHandler.AwaitInputOrTimeout(ctx, f, "", 10*time.Millisecond)
	assert.Errorf(t, err, "event creation timed out: %v", nil)
}

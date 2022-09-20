package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAwaitInputOrTimeout(t *testing.T) {
	ctx := context.TODO()
	handlerFunc := func(s *discordgo.Session, m *discordgo.MessageCreate) {}
	opts, err := mock.NewOptions()
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
	err = AwaitInputOrTimeout(ctx, 10*time.Millisecond, opts.Session, make(chan string), f, handlerFunc, "")
	assert.Errorf(t, err, "event creation timed out: %v", nil)
}

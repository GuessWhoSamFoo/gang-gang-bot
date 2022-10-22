package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewSetDateState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	s := NewSetDateState(*opts)
	assert.NotNil(t, s)
}

func TestSetDateState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	s := NewSetDateState(*opts)

	f := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: SetDate.String(),
				Src:  []string{"idle"},
				Dst:  SetDate.String(),
			},
		},
		fsm.Callbacks{
			SetDate.String(): s.OnState,
		},
	)
	s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
		s.inputHandler.inputChan <- "tomorrow"
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
		wg.Done()
	}()
	go func() {
		err = f.Event(context.TODO(), SetDate.String())
		assert.NoError(t, err)
		wg.Done()
	}()
	wg.Wait()
	actual, err := Get(f, discord.StartTime)
	assert.NoError(t, err)

	cur := time.Now()
	expected := time.Date(cur.Year(), cur.Month(), cur.Day()+1, 0, 0, 0, 0, cur.Location())
	assert.Equal(t, expected, actual)
}

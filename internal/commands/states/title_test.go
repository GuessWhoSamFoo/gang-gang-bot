package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewAddTitleState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)

	s := NewAddTitleState(*opts)
	assert.NotNil(t, s)
}

func TestAddTitleState_OnState(t *testing.T) {
	opts, err := discord.NewMockOptions()
	assert.NoError(t, err)
	s := NewAddTitleState(*opts)

	f := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: AddTitle.String(),
				Src:  []string{"idle"},
				Dst:  AddTitle.String(),
			},
		},
		fsm.Callbacks{
			AddTitle.String(): s.OnState,
		},
	)

	expected := "hello world"
	s.inputHandler.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
		s.inputHandler.inputChan <- expected
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		s.inputHandler.handlerFunc(opts.Session, &discordgo.MessageCreate{})
		wg.Done()
	}()

	go func() {
		err = f.Event(context.TODO(), AddTitle.String())
		assert.NoError(t, err)
		wg.Done()
	}()

	wg.Wait()
	actual, err := Get(f, discord.Title)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

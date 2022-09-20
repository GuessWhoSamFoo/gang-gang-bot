package states

import (
	"context"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/discord"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAddTitleState(t *testing.T) {
	opts, err := mock.NewOptions()
	assert.NoError(t, err)

	s := NewAddTitleState(*opts)
	assert.NotNil(t, s)
}

func TestAddTitleState_OnState(t *testing.T) {
	opts, err := mock.NewOptions()
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

	go func() {
		s.handlerFunc = func(session *discordgo.Session, create *discordgo.MessageCreate) {
			s.input <- expected
		}
		s.handlerFunc(opts.Session, &discordgo.MessageCreate{})
	}()

	err = f.Event(context.TODO(), AddTitle.String())
	assert.NoError(t, err)

	actual, err := Get(f, discord.Title)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

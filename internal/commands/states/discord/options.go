package discord

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/mock"
	"github.com/bwmarrin/discordgo"
	"github.com/ewohltman/discordgo-mock/mockchannel"
	"github.com/ewohltman/discordgo-mock/mockconstants"
)

type Options struct {
	Session           *discordgo.Session
	InteractionCreate *discordgo.InteractionCreate
	Channel           *discordgo.Channel

	*CalendarClient
}

// NewMockOptions returns a mocked Discord user session
func NewMockOptions() (*Options, error) {
	session, err := mock.NewSession()
	if err != nil {
		return nil, err
	}

	channel := mockchannel.New(
		mockchannel.WithID(mockconstants.TestChannel),
	)

	ic := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			GuildID: mockconstants.TestGuild,
			Member: &discordgo.Member{
				User: &discordgo.User{
					Username: mockconstants.TestUser,
				},
			},
		},
	}

	return &Options{
		Session:           session,
		InteractionCreate: ic,
		Channel:           channel,
	}, nil
}

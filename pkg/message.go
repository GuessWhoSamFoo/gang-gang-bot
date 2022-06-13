package pkg

import "github.com/bwmarrin/discordgo"

var (
	cancelText = "To exit, type 'cancel'"

	EnterTitleMessage = discordgo.MessageEmbed{
		Title:       "Enter the event title",
		Color:       purple,
		Description: "Up to 200 characters are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}
	EnterDescriptionMessage = discordgo.MessageEmbed{
		Title:       "Enter the event description",
		Color:       purple,
		Description: "Type `None` for no description. Up to 1600 characters are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EnterAttendeeLimitMessage = discordgo.MessageEmbed{
		Title:       "Enter the maximum number of attendees",
		Color:       purple,
		Description: "Type `None` for no limit. Up to 250 attendees are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EnterDateStartMessage = discordgo.MessageEmbed{
		Title: "When should the event start",
		Color: purple,
		// TODO: Support various time input formats
		// Description: "> Friday at 9pm\n> Tomorrow at 18:00\n> Now\n> In 1 hour\n> YYYY-MM-DD 7:00 PM",
		Description: "> YYYY-MM-DD 7:00 PM",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EnterDurationMessage = discordgo.MessageEmbed{
		Title:       "What is the duration of this event?",
		Color:       purple,
		Description: "Type `None` for no duration.\n> 2 hours\n> 45 minutes\n> 1 hour and 30 minutes",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}
)

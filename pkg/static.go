package pkg

import "github.com/bwmarrin/discordgo"

// Discord Button Components
var (
	AcceptButton = discordgo.Button{
		Label:    "✅",
		Style:    discordgo.SecondaryButton,
		CustomID: "accept",
	}
	DeclineButton = discordgo.Button{
		Label:    "❌",
		Style:    discordgo.SecondaryButton,
		CustomID: "decline",
	}
	TentativeButton = discordgo.Button{
		Label:    "❔",
		Style:    discordgo.SecondaryButton,
		CustomID: "tentative",
	}
	EditButton = discordgo.Button{
		Label:    "Edit",
		Style:    discordgo.PrimaryButton,
		CustomID: "edit",
	}
	DeleteButton = discordgo.Button{
		Label:    "Delete",
		Style:    discordgo.DangerButton,
		CustomID: "delete",
	}
)

// Discord Static Responses
var (
	cancelText                = "To exit, type 'cancel'"
	invalidEventLimitText     = "Entry must be between 1 and 250 (or `None` for no limit). Try again:"
	invalidEntryText          = "Invalid entry. Please select a number from the list above."
	invalidEventTimeText      = "Event start time cannot be in the past. Try again:"
	invalidRemoveResponseText = "Invalid selection. Enter the number(s) of the desired option(s), separated by spaces. \n\nFor example: `1 3 5`"
	foundMultipleText         = "We've found more than one user for the search term. Try something more specific:"
	foundNoneText             = "We couldn't find a user with that name. Try again:"
	userSignedUpText          = "That user is already signed up for this event."
	optionText                = "Enter a number to select an option"

	EnterTitleMessage = discordgo.MessageEmbed{
		Title:       "Enter the event title",
		Color:       Purple,
		Description: "Up to 200 characters are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}
	EnterDescriptionMessage = discordgo.MessageEmbed{
		Title:       "Enter the event description",
		Color:       Purple,
		Description: "Type `None` for no description. Up to 1600 characters are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EnterAttendeeLimitMessage = discordgo.MessageEmbed{
		Title:       "Enter the maximum number of attendees",
		Color:       Purple,
		Description: "Type `None` for no limit. Up to 250 attendees are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EnterDateStartMessage = discordgo.MessageEmbed{
		Title: "When should the event start",
		Color: Purple,
		// TODO: Support various time input formats
		// Description: "> Friday at 9pm\n> Tomorrow at 18:00\n> Now\n> In 1 hour\n> YYYY-MM-DD 7:00 PM",
		Description: "> tomorrow at 10:15am\n> now\n> YYYY-MM-DD 7:00 PM",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EnterDurationMessage = discordgo.MessageEmbed{
		Title:       "What is the duration of this event?",
		Color:       Purple,
		Description: "Type `None` for no duration.\n> 2 hours\n> 45 minutes\n> 1 hour and 30 minutes",
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EnterLocationMessage = discordgo.MessageEmbed{
		Title: "Where does this event take place?",
		Color: Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	CommandInProcessMessage = discordgo.MessageEmbed{
		Title:       "You have another command in process",
		Color:       Purple,
		Description: "Check your direct messages with me",
	}

	EnterEditOptionMessage = discordgo.MessageEmbed{
		Title:       "What would you like to do?",
		Color:       Purple,
		Description: "**1** Modify the event\n**2** Remove responses\n**3** Add a response",
		Footer: &discordgo.MessageEmbedFooter{
			Text: optionText + "\n" + cancelText,
		},
	}

	EnterUserNameMessage = discordgo.MessageEmbed{
		Title:       "Enter the name of the user you'd like to add",
		Description: "An exact match isn't needed. A few characters of their name will suffice!",
		Color:       Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}

	EditConfirmationMessage = discordgo.MessageEmbed{
		Title:       "Would you like to keep editing?",
		Description: "**1** No, I'm all done\n**2** Yes, keep editing",
		Color:       Purple,
	}
)

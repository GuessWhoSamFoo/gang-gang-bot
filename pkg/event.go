package pkg

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/tj/go-naturaldate"
	"log"
	"strconv"
	"time"
)

type Event struct {
	title          string
	description    string
	limit          int
	start          time.Time
	end            time.Time
	accepted       int
	declined       int
	tentative      int
	waitlist       int
	acceptedNames  []string
	declinedNames  []string
	tentativeNames []string
	waitlistNames  []string
	// TODO: image, frequency, localization
}

// NewEvent creates a new event
func NewEvent() *Event {
	return &Event{
		limit: -1,
		start: time.Now(),
	}
}

func (e *Event) AddTitle(title string) {
	e.title = title
}

func (e *Event) AddDescription(description string) {
	e.description = description
}

func (e *Event) SetMaximumAttendees(number int) {
	e.limit = number
}

func (e *Event) SetStartTime(startTime time.Time) {
	e.start = startTime
}

func (e *Event) SetDuration(length string) error {
	if length == "" {
		return nil
	}
	end, err := naturaldate.Parse(length, e.start, naturaldate.WithDirection(naturaldate.Future))
	if err != nil {
		return err
	}
	e.end = end
	return nil
}

type EventBuilder struct {
	Event             *Event
	Session           *discordgo.Session
	Channel           *discordgo.Channel
	GuildID           string
	InteractionCreate *discordgo.InteractionCreate
}

// NewEventBuilder manages the lifecycle of an event creation
func NewEventBuilder(s *discordgo.Session, c *discordgo.Channel, ic *discordgo.InteractionCreate) (*EventBuilder, error) {
	if s == nil || c == nil || ic == nil {
		return nil, fmt.Errorf("missing builder resource %v, %v, %v", s, c, ic)
	}
	return &EventBuilder{
		Event:             NewEvent(),
		Session:           s,
		Channel:           c,
		InteractionCreate: ic,
	}, nil
}

func (eb *EventBuilder) StartChat() error {
	messages, err := eb.Session.ChannelMessages(eb.Channel.ID, 1, "", "", "")
	if err != nil {
		return fmt.Errorf("failed to check last message from channel")
	}
	for _, m := range messages {
		for _, embed := range m.Embeds {
			if embed.Title == EnterTitleMessage.Title {
				if err := eb.Session.InteractionRespond(eb.InteractionCreate.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title:       "You have another command in process",
								Color:       purple,
								Description: "Check your direct messages with me",
							},
						},
						Flags: discordgo.MessageFlagsEphemeral,
					},
				}); err != nil {
					return fmt.Errorf("respond to command in process: %v", err)
				}
				return fmt.Errorf("existing command in process")
			}
		}
	}

	msg, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterTitleMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	if err := eb.Session.InteractionRespond(eb.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Let's create an event",
					Color:       purple,
					Description: fmt.Sprintf("I've sent you a [direct message](https://discordapp.com/channels/%s/%s) with next steps.", eb.InteractionCreate.GuildID, msg.ChannelID),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		return fmt.Errorf("command handler response: %v", err)
	}
	eb.GuildID = eb.InteractionCreate.Interaction.GuildID
	return nil
}

func (eb *EventBuilder) AddTitle() error {
	result, err := eb.waitForInput(time.Second * 60)
	if err != nil {
		return err
	}
	eb.Event.AddTitle(result.(string))
	return nil
}

func (eb *EventBuilder) AddDescription() error {
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterDescriptionMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	result, err := eb.waitForInput(time.Second * 60)
	if err != nil {
		return err
	}
	eb.Event.AddDescription(result.(string))
	return nil
}

func (eb *EventBuilder) SetAttendeeLimit() error {
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterAttendeeLimitMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	for i := 0; i < 3; i++ {
		result, err := eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}

		s, ok := result.(string)
		if ok && s == "" {
			return nil
		}

		val, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		if val > 1 && val < 250 {
			eb.Event.SetMaximumAttendees(val)
			return nil
		}
		if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, "Entry must be between 1 and 250 (or `None` for no limit). Try again:"); err != nil {
			return err
		}
	}
	return fmt.Errorf("unable to set attendee limit")
}

func (eb *EventBuilder) SetDate() error {
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterDateStartMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	result, err := eb.waitForInput(time.Second * 60)
	if err != nil {
		return err
	}

	startTime, err := dateparse.ParseLocal(result.(string))
	if err != nil {
		return fmt.Errorf("cannot parse time: %v", err)
	}

	eb.Event.SetStartTime(startTime)
	return nil
}

func (eb *EventBuilder) SetDuration() error {
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterDurationMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	result, err := eb.waitForInput(time.Second * 60)
	if err != nil {
		return fmt.Errorf("cannot parse duration: %v", err)
	}

	if err := eb.Event.SetDuration(result.(string)); err != nil {
		return err
	}
	return nil
}

func (eb *EventBuilder) CreateEvent() error {
	if eb.InteractionCreate.Interaction.Member == nil {
		return fmt.Errorf("cannot find user who created event")
	}

	acceptedField := "✅ Accepted"
	if eb.Event.limit > 0 {
		acceptedField = acceptedField + fmt.Sprintf(" (0/%d)", eb.Event.limit)
	}
	msg, err := eb.Session.FollowupMessageCreate(eb.InteractionCreate.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       eb.Event.title,
				Description: eb.Event.description,
				Color:       purple,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "Time",
						// https://discord.com/developers/docs/reference#message-formatting-timestamp-styles
						Value: util.PrintTime(eb.Event.start, eb.Event.end),
					},
					{
						Name:  "Links",
						Value: util.PrintAddGoogleCalendarLink(eb.Event.title, eb.Event.description, eb.Event.start, eb.Event.end),
					},
					{
						Name:   acceptedField,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   "❌ Declined",
						Value:  "-",
						Inline: true,
					},
					{
						Name:   "❔ Tentative",
						Value:  "-",
						Inline: true,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Created by %v", eb.InteractionCreate.Interaction.Member.User.Username),
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					AcceptButton,
					DeclineButton,
					TentativeButton,
					EditButton,
					DeleteButton,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Event has been created",
		Color:       purple,
		Description: fmt.Sprintf("[Click here to view the event](https://discord.com/channels/%s/%s/%s)", eb.InteractionCreate.GuildID, eb.InteractionCreate.Interaction.ChannelID, msg.ID),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	log.Println("Successfully created event")
	return nil
}

func (eb *EventBuilder) waitForInput(timeout time.Duration) (interface{}, error) {
	input := make(chan string)
	cancelFunc := eb.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Only interested in inputs via DM
		if m.ChannelID == eb.Channel.ID {
			input <- m.Content
		}
	})
	defer cancelFunc()

	for {
		select {
		case result := <-input:
			if result == "" {
				continue
			}
			if result == "cancel" {
				if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, "Event creation has been canceled"); err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("canceled event creation")
			}
			if result == "None" {
				result = ""
			}
			return result, nil
		case <-time.After(timeout):
			input <- ""
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, "I'm not sure where you went. We can try this again later."); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("create event timed out")
		}
	}
}

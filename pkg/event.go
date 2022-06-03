package pkg

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/tj/go-naturaldate"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type Event struct {
	title       string
	description string
	limit       int
	start       time.Time
	end         time.Time
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
	msg, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Enter the event title",
		Color:       purple,
		Description: "Up to 200 characters are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "To exit, type 'cancel'",
		},
	})
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
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Enter the event description",
		Color:       purple,
		Description: "Type `None` for no description. Up to 1600 characters are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "To exit, type 'cancel'",
		},
	})
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
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Enter the maximum number of attendees",
		Color:       purple,
		Description: "Type `None` for no limit. Up to 250 attendees are permitted",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "To exit, type 'cancel'",
		},
	})
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
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title: "When should the event start",
		Color: purple,
		// TODO: Support various time input formats
		// Description: "> Friday at 9pm\n> Tomorrow at 18:00\n> Now\n> In 1 hour\n> YYYY-MM-DD 7:00 PM",
		Description: "> YYYY-MM-DD 7:00 PM",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "To exit, type 'cancel'",
		},
	})
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
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "What is the duration of this event?",
		Color:       purple,
		Description: "Type `None` for no duration.\n> 2 hours\n> 45 minutes\n> 1 hour and 30 minutes",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "To exit, type 'cancel'",
		},
	})
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

func (eb *EventBuilder) printTime(start, end time.Time) string {
	base := fmt.Sprintf("<t:%d:F>", start.Unix())
	relative := fmt.Sprintf("\n🕔<t:%d:R>", start.Unix())
	if end.IsZero() {
		return base + relative
	}

	if end.Sub(start) < time.Hour*24 {
		base = base + fmt.Sprintf(" - <t:%d:t>", end.Unix())
	} else {
		base = base + fmt.Sprintf(" - <t:%d:F>", end.Unix())
	}
	return base + relative
}

func (eb *EventBuilder) printAddGoogleCalendarLink(title, description string, startTime, endTime time.Time) string {
	if endTime.IsZero() {
		endTime = startTime
	}

	s, e := startTime.UTC().Format(GoogleCalendarTimeFormat), endTime.UTC().Format(GoogleCalendarTimeFormat)

	u, _ := url.Parse("https://www.google.com/calendar/event?action=TEMPLATE&text=&details=&location=")
	q := u.Query()
	q.Set("text", title)
	q.Set("details", description)
	u.RawQuery = q.Encode()

	// TODO: Encode multiple dates rather than constructing
	re := regexp.MustCompile("[[:punct:]]")
	s, e = re.ReplaceAllString(s, ""), re.ReplaceAllString(e, "")

	link := u.String() + "&dates=" + s + "/" + e

	return fmt.Sprintf("[Add to Google Calendar](%s)", link)
}

func (eb *EventBuilder) CreateEvent() error {
	if eb.InteractionCreate.Interaction.Member == nil {
		return fmt.Errorf("cannot find user who created event")
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
						Value: eb.printTime(eb.Event.start, eb.Event.end),
					},
					{
						Name:  "Links",
						Value: eb.printAddGoogleCalendarLink(eb.Event.title, eb.Event.description, eb.Event.start, eb.Event.end),
					},
					{
						Name:   "✅ Accepted",
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
					discordgo.Button{
						Label:    "✅",
						Style:    discordgo.SecondaryButton,
						CustomID: "accept",
					},
					discordgo.Button{
						Label:    "❌",
						Style:    discordgo.SecondaryButton,
						CustomID: "decline",
					},
					discordgo.Button{
						Label:    "❔",
						Style:    discordgo.SecondaryButton,
						CustomID: "tentative",
					},
					discordgo.Button{
						Label:    "Edit",
						Style:    discordgo.PrimaryButton,
						CustomID: "edit",
					},
					discordgo.Button{
						Label:    "Delete",
						Style:    discordgo.DangerButton,
						CustomID: "delete",
					},
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

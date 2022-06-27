package pkg

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/tj/go-naturaldate"
	"log"
	"strconv"
	"strings"
	"time"
)

type ActionType string

const (
	CreateType ActionType = "creation"
	EditType   ActionType = "modification"
)

type EventBuilder struct {
	Event             *Event
	Session           *discordgo.Session
	Channel           *discordgo.Channel
	GuildID           string
	InteractionCreate *discordgo.InteractionCreate
	Action            ActionType
	CalendarID        string
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

// StartCreate starts the sequence of prompts to create a new event
func (eb *EventBuilder) StartCreate() error {
	if err := eb.Session.InteractionRespond(eb.InteractionCreate.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Let's create an event",
					Color:       Purple,
					Description: fmt.Sprintf("I've sent you a [direct message](https://discordapp.com/channels/%s/%s) with next steps.", eb.InteractionCreate.GuildID, eb.Channel.ID),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		return fmt.Errorf("command handler response: %v", err)
	}
	eb.GuildID = eb.InteractionCreate.Interaction.GuildID
	eb.Action = CreateType
	return nil
}

func (eb *EventBuilder) StartEdit() error {
	event, err := GetEventFromMessage(eb.InteractionCreate.Interaction.Message)
	if err != nil {
		return err
	}
	if eb.InteractionCreate.Member.User.Username != event.Owner && eb.InteractionCreate.Interaction.Member.Permissions&discordgo.PermissionManageEvents == 0 {
		if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, EditInsufficientPermissionMessage); err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}
		return fmt.Errorf("insufficient permissions to edit %s", eb.InteractionCreate.Interaction.Message.ID)
	}
	eb.Event = event
	eb.Action = EditType
	eb.Event.DiscordLink = fmt.Sprintf("https://discord.com/channels/%s/%s/%s", eb.InteractionCreate.GuildID, eb.InteractionCreate.Interaction.ChannelID, eb.InteractionCreate.Interaction.Message.ID)
	if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterEditOptionMessage); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

Loop:
	for {
		result, err := eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}
		option := result.(string)
		switch {
		case option == "1":
			if err := eb.modifyEvent(); err != nil {
				return err
			}
			break Loop
		case option == "2":
			if err := eb.removeResponses(); err != nil {
				return err
			}
			break Loop
		case option == "3":
			if err := eb.addResponse(); err != nil {
				return err
			}
			break Loop
		default:
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, invalidEntryText); err != nil {
				return err
			}
		}
	}
	return nil
}

func (eb *EventBuilder) modifyEvent() error {
	if eb.Event == nil {
		return fmt.Errorf("event is nil")
	}

	for {
		if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
			Title: "What would you like to modify?",
			Color: Purple,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "1 ⋅ Title",
					Value: util.PrintBlockValues(eb.Event.Title),
				},
				{
					Name:  "2 ⋅ Description",
					Value: util.PrintBlockValues(eb.Event.Description),
				},
				{
					Name:   "3 ⋅ Start Time",
					Value:  fmt.Sprintf("```%s```", eb.Event.Start.In(time.Local).Format(util.HumanTimeFormat)),
					Inline: true,
				},
				{
					Name:   "4 ⋅ Duration",
					Value:  util.PrintBlockValues(util.PrintHumanReadableTime(eb.Event.Start, eb.Event.End)),
					Inline: true,
				},
				{
					Name:  "5 ⋅ Location",
					Value: fmt.Sprintf("```%s```", eb.Event.Location),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: optionText + "\n" + cancelText,
			},
		}); err != nil {
			return err
		}
		result, err := eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}

		switch option := result.(string); {
		case option == "1":
			if err := eb.AddTitle(); err != nil {
				return err
			}
		case option == "2":
			if err := eb.AddDescription(); err != nil {
				return err
			}
		case option == "3":
			if err := eb.SetDate(); err != nil {
				return err
			}
		case option == "4":
			if err := eb.SetDuration(); err != nil {
				return err
			}
		case option == "5":
			if err := eb.SetLocation(); err != nil {
				return err
			}
		default:
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, invalidEntryText); err != nil {
				return err
			}
		}
		if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EditConfirmationMessage); err != nil {
			return err
		}

		result, err = eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}
		if result.(string) == "1" {
			break
		}
	}
	return nil
}

func (eb *EventBuilder) removeResponses() error {
	if eb.Event == nil {
		return fmt.Errorf("event is nil")
	}

	var count int
	for _, r := range eb.Event.RoleGroup.Roles {
		count += r.Count
	}
	if count == 0 {
		_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
			Title: "Event doesn't have any responses",
		})
		if err != nil {
			return err
		}
	}

	var desc string
	var counter int
	users := make([]string, 0)
	// Braille space is used instead because hard spaces in embeds are not documented
	for _, role := range eb.Event.RoleGroup.Roles {
		users = append(users, role.Users...)
		for _, n := range role.Users {
			counter++
			desc = desc + fmt.Sprintf("**%d**⠀%s %s\n", counter, role.Icon, n)
		}
	}
	if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Which responses would you like to remove?",
		Description: desc,
		Color:       Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Enter the number(s) of the desired option(s), separated by spaces\n" + cancelText,
		},
	}); err != nil {
		return err
	}

	wl, ok := eb.Event.RoleGroup.Waitlist[AcceptedField]
	if !ok {
		return fmt.Errorf("cannot find accepted field")
	}
	nameMap := map[int]string{}
	for index, name := range append(users, wl.Users...) {
		nameMap[index+1] = name
	}

	for {
		result, err := eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}
		if !util.IsInputOption(result.(string)) {
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, invalidRemoveResponseText); err != nil {
				return err
			}
		}
		for _, n := range strings.Split(result.(string), " ") {
			option, err := strconv.Atoi(n)
			if err != nil {
				return err
			}

			val, ok := nameMap[option]
			if !ok {
				return fmt.Errorf("cannot find name")
			}
			if err := eb.Event.RemoveFromAllLists(eb.Session, eb.InteractionCreate, val); err != nil {
				return err
			}
		}
		break
	}
	return nil
}

func (eb *EventBuilder) addResponse() error {
	if eb.Event == nil {
		return fmt.Errorf("event is nil")
	}
	if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterUserNameMessage); err != nil {
		return err
	}
	// TODO: Handle guilds with more than 1000 members
	members, err := eb.Session.GuildMembers(eb.InteractionCreate.Interaction.GuildID, "0", 1000)
	if err != nil {
		return fmt.Errorf("failed to get guild members: %v", err)
	}
	names := make([]string, 0)
	for _, m := range members {
		names = append(names, m.User.Username)
	}

	var user string
	for {
		result, err := eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}
		matches := fuzzy.Find(result.(string), names)
		numMatches := len(matches)
		if numMatches == 0 {
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, foundNoneText); err != nil {
				return err
			}
			continue
		}
		if numMatches > 1 {
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, foundMultipleText); err != nil {
				return err
			}
			continue
		}
		user = matches[0]
		for _, r := range eb.Event.RoleGroup.Roles {
			if util.ContainsUser(r.Users, user) {
				if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, userSignedUpText); err != nil {
					return err
				}
				return fmt.Errorf("user already signed up")
			}
		}
		break
	}

	if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Which signup option should we add the user to?",
		Description: fmt.Sprintf("**1**⠀%s\n**2**⠀%s\n**3**⠀%s", acceptedBase, declinedBase, tentativeBase),
		Color:       Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text: cancelText,
		},
	}); err != nil {
		return err
	}

Loop:
	for {
		result, err := eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}
		option := result.(string)
		switch {
		case option == "1":
			if err := eb.Event.ToggleAccept(eb.Session, eb.InteractionCreate, user); err != nil {
				return err
			}
			break Loop
		case option == "2":
			if err := eb.Event.ToggleDecline(eb.Session, eb.InteractionCreate, user); err != nil {
				return err
			}
			break Loop
		case option == "3":
			if err := eb.Event.ToggleTentative(eb.Session, eb.InteractionCreate, user); err != nil {
				return err
			}
			break Loop
		default:
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, invalidEntryText); err != nil {
				return err
			}
		}
	}
	return nil
}

func (eb *EventBuilder) ProcessEdit() error {
	if eb.Event == nil {
		return fmt.Errorf("event is nil")
	}
	embed, err := ConvertEventToMessageEmbed(eb.Event)
	if err != nil {
		return err
	}

	if _, err := eb.Session.ChannelMessageEditEmbed(eb.InteractionCreate.Interaction.ChannelID, eb.InteractionCreate.Interaction.Message.ID, embed); err != nil {
		return err
	}

	msg := eb.InteractionCreate.Interaction.Message
	if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Event has been updated!",
		Color:       Purple,
		Description: fmt.Sprintf("[Click here to view the event](https://discord.com/channels/%s/%s/%s)", eb.InteractionCreate.GuildID, eb.InteractionCreate.Interaction.ChannelID, msg.ID),
	}); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

// AddTitle adds a title to the event
func (eb *EventBuilder) AddTitle() error {
	if _, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterTitleMessage); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	result, err := eb.waitForInput(time.Second * 60)
	if err != nil {
		return err
	}
	eb.Event.AddTitle(result.(string))
	return nil
}

// AddDescription adds a description to the event
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

// SetAttendeeLimit sets the maximum number of attendees
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

		if val >= 1 && val <= 250 {
			eb.Event.SetMaximumAttendees(val)
			return nil
		}
		if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, invalidEventLimitText); err != nil {
			return err
		}
	}
	return fmt.Errorf("unable to set attendee limit")
}

// SetDate sets the starting time of an event. It attempts to parse time in natural language
func (eb *EventBuilder) SetDate() error {
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterDateStartMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	for {
		result, err := eb.waitForInput(time.Second * 60)
		if err != nil {
			return err
		}

		var startTime time.Time
		now := time.Now()
		startTime, err = naturaldate.Parse(result.(string), now, naturaldate.WithDirection(naturaldate.Future))
		if err != nil {
			startTime, err = dateparse.ParseLocal(result.(string))
			if err != nil {
				return fmt.Errorf("cannot parse time: %v", err)
			}
		}
		if startTime.Before(now) {
			if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, invalidEventTimeText); err != nil {
				return err
			}
			continue
		}
		eb.Event.SetStartTime(startTime)
		break
	}
	return nil
}

// SetDuration sets the length of an event
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

// SetLocation sets the location of an event
func (eb *EventBuilder) SetLocation() error {
	_, err := eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &EnterLocationMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	result, err := eb.waitForInput(time.Second * 60)
	if err != nil {
		return fmt.Errorf("cannot parse location: %v", err)
	}

	if err := eb.Event.SetLocation(result.(string)); err != nil {
		return err
	}
	return nil
}

// CreateEvent converts the internal event object to a Discord message
func (eb *EventBuilder) CreateEvent() error {
	if eb.InteractionCreate.Interaction.Member == nil {
		return fmt.Errorf("cannot find user who created event")
	}

	acceptedField := acceptedBase
	for _, r := range eb.Event.RoleGroup.Roles {
		if r.FieldName == AcceptedField && r.Limit > 0 {
			acceptedField = acceptedField + fmt.Sprintf(" (0/%d)", r.Limit)
		}
	}

	msg, err := eb.Session.ChannelMessageSendComplex(eb.InteractionCreate.Interaction.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       eb.Event.Title,
				Description: eb.Event.Description,
				Color:       Purple,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "Time",
						// https://discord.com/developers/docs/reference#message-formatting-timestamp-styles
						Value: util.PrintTime(eb.Event.Start, eb.Event.End),
					},
					{
						Name:   "Links",
						Value:  util.PrintAddGoogleCalendarLink(eb.Event.Title, eb.Event.Description, eb.Event.Start, eb.Event.End),
						Inline: true,
					},
					{
						Name:   "Location",
						Value:  eb.Event.Location,
						Inline: true,
					},
					{
						Name:  "Calendar",
						Value: util.PrintGoogleCalendarEventLink(eb.Event.ID),
					},
					{
						Name:   acceptedField,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   declinedBase,
						Value:  "-",
						Inline: true,
					},
					{
						Name:   tentativeBase,
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

	eb.Event.DiscordLink = fmt.Sprintf("https://discord.com/channels/%s/%s/%s", eb.InteractionCreate.GuildID, eb.InteractionCreate.Interaction.ChannelID, msg.ID)
	_, err = eb.Session.ChannelMessageSendEmbed(eb.Channel.ID, &discordgo.MessageEmbed{
		Title:       "Event has been created",
		Color:       Purple,
		Description: fmt.Sprintf("[Click here to view the event](%s)", eb.Event.DiscordLink),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	log.Println("Successfully created event")
	return nil
}

// waitForInput waits for the user to reply in the same channel
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
				if _, err := eb.Session.ChannelMessageSend(eb.Channel.ID, fmt.Sprintf("Event %s has been canceled", eb.Action)); err != nil {
					return nil, err
				}
				return nil, fmt.Errorf("canceled event %s", eb.Action)
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

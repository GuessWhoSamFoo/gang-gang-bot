package pkg

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/tj/go-naturaldate"
	"golang.org/x/exp/slices"
	"log"
	"strings"
	"time"
)

var (
	acceptedBase  = "✅ Accepted"
	declinedBase  = "❌ Declined"
	tentativeBase = "❔ Tentative"
)

// Event is an internal representation of a formatted Discord Embed Message
type Event struct {
	Title       string
	Description string
	Limit       int
	Start       time.Time
	End         time.Time
	Accepted    int
	Declined    int
	Tentative   int
	Waitlist    int
	// TODO: Generalize to N signup options
	AcceptedNames  []string
	DeclinedNames  []string
	TentativeNames []string
	WaitlistNames  []string
	Owner          string
	Color          int
	// TODO: image, frequency, localization
}

// NewEvent creates a new event
func NewEvent() *Event {
	return &Event{
		Limit:          -1,
		Start:          time.Now(),
		AcceptedNames:  []string{},
		DeclinedNames:  []string{},
		TentativeNames: []string{},
		WaitlistNames:  []string{},
	}
}

func (e *Event) AddTitle(title string) {
	e.Title = title
}

func (e *Event) AddDescription(description string) {
	e.Description = description
}

func (e *Event) SetMaximumAttendees(number int) {
	e.Limit = number
}

func (e *Event) SetStartTime(startTime time.Time) {
	e.Start = startTime
}

func (e *Event) SetDuration(length string) error {
	if length == "" {
		return nil
	}
	end, err := naturaldate.Parse(length, e.Start, naturaldate.WithDirection(naturaldate.Future))
	if err != nil {
		return err
	}
	e.End = end
	return nil
}

func (e *Event) ToggleAccept(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if slices.Contains(e.DeclinedNames, name) {
		e.Declined--
		e.DeclinedNames = util.RemoveUser(e.DeclinedNames, name)
	}
	if slices.Contains(e.TentativeNames, name) {
		e.Tentative--
		e.TentativeNames = util.RemoveUser(e.TentativeNames, name)
	}
	if slices.Contains(e.WaitlistNames, name) {
		e.Waitlist--
		e.WaitlistNames = util.RemoveUser(e.WaitlistNames, name)
	}

	if e.Limit != -1 && e.Accepted >= e.Limit && !slices.Contains(e.AcceptedNames, name) {
		e.Waitlist++
		e.WaitlistNames = append(e.WaitlistNames, name)
	} else if !slices.Contains(e.AcceptedNames, name) {
		e.Accepted++
		e.AcceptedNames = append(e.AcceptedNames, name)
	} else {
		e.Accepted--
		e.AcceptedNames = util.RemoveUser(e.AcceptedNames, name)
		if err := e.BumpWaitlist(s, i.Interaction); err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) ToggleDecline(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if slices.Contains(e.AcceptedNames, name) {
		e.Accepted--
		e.AcceptedNames = util.RemoveUser(e.AcceptedNames, name)
		if err := e.BumpWaitlist(s, i.Interaction); err != nil {
			return err
		}
	}
	if slices.Contains(e.TentativeNames, name) {
		e.Tentative--
		e.TentativeNames = util.RemoveUser(e.TentativeNames, name)
	}
	if slices.Contains(e.WaitlistNames, name) {
		e.Waitlist--
		e.WaitlistNames = util.RemoveUser(e.WaitlistNames, name)
	}

	if !slices.Contains(e.DeclinedNames, name) {
		e.Declined++
		e.DeclinedNames = append(e.DeclinedNames, name)
	} else {
		e.Declined--
		e.DeclinedNames = util.RemoveUser(e.DeclinedNames, name)
	}
	return nil
}

func (e *Event) ToggleTentative(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if slices.Contains(e.AcceptedNames, name) {
		e.Accepted--
		e.AcceptedNames = util.RemoveUser(e.AcceptedNames, name)
		if err := e.BumpWaitlist(s, i.Interaction); err != nil {
			return err
		}
	}
	if slices.Contains(e.DeclinedNames, name) {
		e.Declined--
		e.DeclinedNames = util.RemoveUser(e.DeclinedNames, name)
	}
	if slices.Contains(e.WaitlistNames, name) {
		e.Waitlist--
		e.WaitlistNames = util.RemoveUser(e.WaitlistNames, name)
	}

	if !slices.Contains(e.TentativeNames, name) {
		e.Tentative++
		e.TentativeNames = append(e.TentativeNames, name)
	} else {
		e.Tentative--
		e.TentativeNames = util.RemoveUser(e.TentativeNames, name)
	}
	return nil
}

func (e *Event) RemoveFromAllLists(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if util.ContainsUser(e.WaitlistNames, name) {
		e.WaitlistNames = util.RemoveUser(e.WaitlistNames, name)
		e.Waitlist--
	}
	if util.ContainsUser(e.AcceptedNames, name) {
		e.AcceptedNames = util.RemoveUser(e.AcceptedNames, name)
		e.Accepted--
		if err := e.BumpWaitlist(s, i.Interaction); err != nil {
			return err
		}
	}
	if util.ContainsUser(e.DeclinedNames, name) {
		e.DeclinedNames = util.RemoveUser(e.DeclinedNames, name)
		e.Declined--
	}
	if util.ContainsUser(e.TentativeNames, name) {
		e.TentativeNames = util.RemoveUser(e.TentativeNames, name)
		e.Tentative--
	}
	return nil
}

func (e *Event) BumpWaitlist(s *discordgo.Session, i *discordgo.Interaction) error {
	if (e.Limit != -1 && e.Accepted >= e.Limit) || e.Waitlist <= 0 {
		return nil
	}
	e.Waitlist--
	name := e.WaitlistNames[0]
	e.WaitlistNames = e.WaitlistNames[1:]
	e.Accepted++
	e.AcceptedNames = append(e.AcceptedNames, name)

	c, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		return err
	}
	// TODO: Handle guilds with more than 1000 members
	members, err := s.GuildMembersSearch(i.GuildID, name, 1000)
	if err != nil {
		return err
	}
	for _, m := range members {
		if m.User.Username == name {
			if _, err := s.ChannelMessageSendEmbed(c.ID, &discordgo.MessageEmbed{
				Title:       "You have been moved off the waitlist!",
				Color:       Purple,
				Description: fmt.Sprintf("[Click here to view the event](https://discord.com/channels/%s/%s/%s)", i.GuildID, i.ChannelID, i.Message.ID),
			}); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func NotifyCommandInProgress(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				&CommandInProcessMessage,
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		log.Printf("failed to send message: %v", err)
	}
}

func GetEventFromMessage(msg *discordgo.Message) (*Event, error) {
	if len(msg.Embeds) != 1 {
		return nil, fmt.Errorf("expected 1 embed: got %d", len(msg.Embeds))
	}
	e := &Event{}
	embed := msg.Embeds[0]
	e.Title = embed.Title
	e.Description = embed.Description
	e.Color = embed.Color
	for _, f := range embed.Fields {
		switch {
		case strings.Contains(f.Name, acceptedBase):
			_, limit, err := util.ParseFieldHeadCount(f.Name)
			if err != nil {
				return nil, err
			}
			e.Limit = limit
			e.AcceptedNames = util.GetUsersFromValues(f.Value)
			e.Accepted = len(e.AcceptedNames)
		case strings.Contains(f.Name, declinedBase):
			e.DeclinedNames = util.GetUsersFromValues(f.Value)
			e.Declined = len(e.DeclinedNames)
		case strings.Contains(f.Name, tentativeBase):
			e.TentativeNames = util.GetUsersFromValues(f.Value)
			e.Tentative = len(e.TentativeNames)
		case f.Name == "Links":
			var err error
			e.Start, e.End, err = util.GetTimesFromLink(f.Value)
			if err != nil {
				return nil, err
			}
		case f.Name == "Waitlist":
			e.WaitlistNames = util.GetUsersFromValues(f.Value)
			e.Waitlist = len(e.WaitlistNames)
		case f.Name == "Time":
			// no-op since start/end times comes from Links
		default:
			return nil, fmt.Errorf("unknown field: %s", f.Name)
		}
	}

	if embed.Footer != nil {
		e.Owner = util.GetUserFromFooter(embed.Footer.Text)
	}
	return e, nil
}

func ConvertEventToMessageEmbed(event *Event) (*discordgo.MessageEmbed, error) {
	msg := &discordgo.MessageEmbed{}
	acceptedName, declinedName, tentativeName := acceptedBase, declinedBase, tentativeBase
	if event.Limit == -1 && event.Accepted > 0 {
		acceptedName = fmt.Sprintf(acceptedName+" (%d)", event.Accepted)
	}
	if event.Limit > 0 {
		acceptedName = fmt.Sprintf(acceptedName+" (%d/%d)", event.Accepted, event.Limit)
	}
	if event.Declined > 0 {
		declinedName = fmt.Sprintf(declinedBase+" (%d)", event.Declined)
	}
	if event.Tentative > 0 {
		tentativeName = fmt.Sprintf(tentativeBase+" (%d)", event.Tentative)
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Time",
			Value: util.PrintTime(event.Start, event.End),
		},
		{
			Name:  "Links",
			Value: util.PrintAddGoogleCalendarLink(event.Title, event.Description, event.Start, event.End),
		},
		{
			Name:   acceptedName,
			Value:  util.NameListToValues(event.AcceptedNames),
			Inline: true,
		},
		{
			Name:   declinedName,
			Value:  util.NameListToValues(event.DeclinedNames),
			Inline: true,
		},
		{
			Name:   tentativeName,
			Value:  util.NameListToValues(event.TentativeNames),
			Inline: true,
		},
	}
	if event.Waitlist > 0 || event.Accepted == event.Limit {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Waitlist",
			Value: util.NameListToValues(event.WaitlistNames),
		})
	}

	msg.Title = event.Title
	msg.Description = event.Description
	msg.Color = event.Color
	msg.Fields = fields
	msg.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Created by %v", event.Owner),
	}
	return msg, nil
}

package discord

import (
	"fmt"
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands/states/role"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/tj/go-naturaldate"
	"log"
	"strings"
	"time"
)

type MetadataKey string

const (
	Action      MetadataKey = "action"
	GuildID     MetadataKey = "guildID"
	Title       MetadataKey = "title"
	Description MetadataKey = "description"
	Attendee    MetadataKey = "attendee"
	Location    MetadataKey = "location"
	StartTime   MetadataKey = "start"
	Duration    MetadataKey = "duration"
	Owner       MetadataKey = "owner"
	Color       MetadataKey = "color"
	ID          MetadataKey = "id"
	MenuOption  MetadataKey = "menuOption"

	EventObject MetadataKey = "eventObject"
	Username    MetadataKey = "username"
)

func (m MetadataKey) String() string {
	return string(m)
}

var (
	AcceptedBase  = fmt.Sprintf("%s %s", role.AcceptedIcon, role.AcceptedField)
	DeclinedBase  = fmt.Sprintf("%s %s", role.DeclinedIcon, role.DeclinedField)
	TentativeBase = fmt.Sprintf("%s %s", role.TentativeIcon, role.TentativeField)
)

// Event is an internal representation of a formatted Discord Embed Message
type Event struct {
	Title       string
	Description string
	Location    string
	Start       time.Time
	End         time.Time
	RoleGroup   *role.RoleGroup
	Owner       string
	Color       int
	ID          string // base64 encoded eventID + calendarID
	DiscordLink string
}

func (e *Event) AddTitle(title string) {
	e.Title = title
}

func (e *Event) AddDescription(description string) {
	e.Description = description
}

func (e *Event) SetMaximumAttendees(field role.FieldType, number int) {
	for _, r := range e.RoleGroup.Roles {
		if r.FieldName == field {
			r.Limit = number
			return
		}
	}
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
	if end.Before(e.Start) {
		return fmt.Errorf("cannot set end before start time")
	}
	e.End = end
	return nil
}

func (e *Event) SetLocation(location string) error {
	if location == "" {
		return nil
	}
	e.Location = location
	return nil
}

func (e *Event) ToggleAccept(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if e.RoleGroup == nil {
		return fmt.Errorf("missing role group")
	}
	prev := e.RoleGroup.PeekWaitlist(role.AcceptedField)
	if err := e.RoleGroup.ToggleRole(role.AcceptedField, name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(role.AcceptedField) {
		if err := e.NotifyUserOffWaitlist(s, i.Interaction, prev); err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) ToggleDecline(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if e.RoleGroup == nil {
		return fmt.Errorf("missing role group")
	}
	prev := e.RoleGroup.PeekWaitlist(role.AcceptedField)
	if err := e.RoleGroup.ToggleRole(role.DeclinedField, name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(role.AcceptedField) {
		if err := e.NotifyUserOffWaitlist(s, i.Interaction, prev); err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) ToggleTentative(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if e.RoleGroup == nil {
		return fmt.Errorf("missing role group")
	}
	prev := e.RoleGroup.PeekWaitlist(role.AcceptedField)
	if err := e.RoleGroup.ToggleRole(role.TentativeField, name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(role.AcceptedField) {
		if err := e.NotifyUserOffWaitlist(s, i.Interaction, prev); err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) RemoveFromAllLists(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if e.RoleGroup == nil {
		return fmt.Errorf("missing role group")
	}
	prev := e.RoleGroup.PeekWaitlist(role.AcceptedField)
	if err := e.RoleGroup.RemoveFromAllLists(name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(role.AcceptedField) {
		if err := e.NotifyUserOffWaitlist(s, i.Interaction, prev); err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) NotifyUserOffWaitlist(s *discordgo.Session, i *discordgo.Interaction, name string) error {
	if name == "" {
		return nil
	}
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

// NotifyCommandInProgress notifies a user if another interaction is pending input
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

func FromFSMToEvent(f *fsm.FSM) (*Event, error) {
	e := &Event{}

	if title, found := f.Metadata(Title.String()); found {
		e.Title = fmt.Sprintf("%s", title)
	}
	if description, found := f.Metadata(Description.String()); found {
		e.Description = fmt.Sprintf("%s", description)
	}
	if location, found := f.Metadata(Location.String()); found {
		e.Location = fmt.Sprintf("%s", location)
	}
	if start, found := f.Metadata(StartTime.String()); found {
		val, ok := start.(time.Time)
		if !ok {
			return nil, fmt.Errorf("cannot cast key: %s", StartTime.String())
		}
		e.Start = val
	}
	if end, found := f.Metadata(Duration.String()); found {
		val, ok := end.(time.Time)
		if !ok {
			return nil, fmt.Errorf("cannot cast key: %s", Duration.String())
		}
		e.End = val
	}
	if rg, found := f.Metadata(Attendee.String()); found {
		val, ok := rg.(*role.RoleGroup)
		if !ok {
			return nil, fmt.Errorf("cannot cast key: %s", Attendee.String())
		}
		e.RoleGroup = val
	}
	if owner, found := f.Metadata(Owner.String()); found {
		e.Owner = fmt.Sprintf("%s", owner)
	}
	if color, found := f.Metadata(Color.String()); found {
		val, ok := color.(int)
		if ok {
			e.Color = val
		}
	}
	return e, nil
}

// GetEventFromMessage converts a Discord message into an event type
func GetEventFromMessage(msg *discordgo.Message) (*Event, error) {
	if len(msg.Embeds) != 1 {
		return nil, fmt.Errorf("expected 1 embed: got %d", len(msg.Embeds))
	}
	e := &Event{}
	embed := msg.Embeds[0]
	e.Title = embed.Title
	e.Description = embed.Description
	e.Color = embed.Color
	e.RoleGroup = &role.RoleGroup{
		Roles: []*role.Role{},
		Waitlist: map[role.FieldType]*role.Role{
			role.AcceptedField: {
				Icon:      "",
				FieldName: role.WaitlistField,
				Users:     []string{},
			},
		},
	}
	var err error
	for _, f := range embed.Fields {
		switch {
		case strings.Contains(f.Name, AcceptedBase):
			_, limit, err := util.ParseFieldHeadCount(f.Name)
			if err != nil {
				return nil, err
			}
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Roles = append(e.RoleGroup.Roles, &role.Role{
				Icon:      role.AcceptedIcon,
				FieldName: role.AcceptedField,
				Users:     users,
				Count:     len(users),
				Limit:     limit,
			})
		case strings.Contains(f.Name, DeclinedBase):
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Roles = append(e.RoleGroup.Roles, &role.Role{
				Icon:      role.DeclinedIcon,
				FieldName: role.DeclinedField,
				Users:     users,
				Count:     len(users),
			})
		case strings.Contains(f.Name, TentativeBase):
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Roles = append(e.RoleGroup.Roles, &role.Role{
				Icon:      role.TentativeIcon,
				FieldName: role.TentativeField,
				Users:     users,
				Count:     len(users),
			})
		case f.Name == "Links":
			e.Start, e.End, err = util.GetTimesFromLink(f.Value)
			if err != nil {
				return nil, err
			}
		case f.Name == "Calendar":
			e.ID, err = util.ParseEventID(f.Value)
			if err != nil {
				return nil, err
			}
		case f.Name == "Location":
			e.Location = f.Value
		case f.Name == string(role.WaitlistField):
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Waitlist[role.AcceptedField] = &role.Role{
				Icon:      "",
				FieldName: role.WaitlistField,
				Users:     users,
				Count:     len(users),
			}
		case f.Name == "Time":
			// no-op since start/end times comes from Links
		default:
			return nil, fmt.Errorf("unknown field: %s", f.Name)
		}
	}
	if msg.GuildID != "" && msg.ChannelID != "" && msg.ID != "" {
		e.DiscordLink = fmt.Sprintf("https://discord.com/channels/%s/%s/%s", msg.GuildID, msg.ChannelID, msg.ID)
	}

	if embed.Footer != nil {
		e.Owner = util.GetUserFromFooter(embed.Footer.Text)
	}
	return e, nil
}

// ConvertEventToMessageEmbed converts an internal Event into a Discord Embed message
func ConvertEventToMessageEmbed(event *Event) (*discordgo.MessageEmbed, error) {
	msg := &discordgo.MessageEmbed{}
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Time",
			Value: util.PrintTime(event.Start, event.End),
		},
		{
			Name:   "Links",
			Value:  util.PrintAddGoogleCalendarLink(event.Title, event.Description, event.Start, event.End),
			Inline: true,
		},
	}

	if event.Location != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Location",
			Value:  event.Location,
			Inline: true,
		})
	}
	if event.ID != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Calendar",
			Value: util.PrintGoogleCalendarEventLink(event.ID),
		})
	}

	for _, r := range event.RoleGroup.Roles {
		name := fmt.Sprintf("%s %s", r.Icon, r.FieldName)
		if r.Limit == 0 && r.Count > 0 {
			name = fmt.Sprintf("%s %s (%d)", r.Icon, r.FieldName, r.Count)
		}
		if r.Limit > 0 {
			name = fmt.Sprintf("%s %s (%d/%d)", r.Icon, r.FieldName, r.Count, r.Limit)
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  util.NameListToValues(r.Users),
			Inline: true,
		})
	}

	for _, r := range event.RoleGroup.Roles {
		if wl, ok := event.RoleGroup.Waitlist[r.FieldName]; ok && wl.Count > 0 {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:  string(wl.FieldName),
				Value: util.NameListToValues(wl.Users),
			})
		}
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

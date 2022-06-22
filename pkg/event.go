package pkg

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"github.com/bwmarrin/discordgo"
	"github.com/tj/go-naturaldate"
	"log"
	"strings"
	"time"
)

var (
	acceptedBase  = fmt.Sprintf("%s %s", AcceptedIcon, AcceptedField)
	declinedBase  = fmt.Sprintf("%s %s", DeclinedIcon, DeclinedField)
	tentativeBase = fmt.Sprintf("%s %s", TentativeIcon, TentativeField)
)

// Event is an internal representation of a formatted Discord Embed Message
type Event struct {
	Title       string
	Description string
	Start       time.Time
	End         time.Time
	RoleGroup   *RoleGroup
	Owner       string
	Color       int
	// TODO: image, frequency, localization
}

// NewEvent creates a new event
func NewEvent() *Event {
	return &Event{
		Start:     time.Now(),
		RoleGroup: NewDefaultRoleGroup(),
	}
}

func (e *Event) AddTitle(title string) {
	e.Title = title
}

func (e *Event) AddDescription(description string) {
	e.Description = description
}

func (e *Event) SetMaximumAttendees(number int) {
	for _, r := range e.RoleGroup.Roles {
		if r.FieldName == AcceptedField {
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
	e.End = end
	return nil
}

func (e *Event) ToggleAccept(s *discordgo.Session, i *discordgo.InteractionCreate, name string) error {
	if e.RoleGroup == nil {
		return fmt.Errorf("missing role group")
	}
	prev := e.RoleGroup.PeekWaitlist(AcceptedField)
	if err := e.RoleGroup.ToggleRole(AcceptedField, name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(AcceptedField) {
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
	prev := e.RoleGroup.PeekWaitlist(AcceptedField)
	if err := e.RoleGroup.ToggleRole(DeclinedField, name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(AcceptedField) {
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
	prev := e.RoleGroup.PeekWaitlist(AcceptedField)
	if err := e.RoleGroup.ToggleRole(TentativeField, name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(AcceptedField) {
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
	prev := e.RoleGroup.PeekWaitlist(AcceptedField)
	if err := e.RoleGroup.RemoveFromAllLists(name); err != nil {
		return err
	}
	if prev != "" && prev != e.RoleGroup.PeekWaitlist(AcceptedField) {
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
	e.RoleGroup = &RoleGroup{
		Roles: []*Role{},
		Waitlist: map[FieldType]*Role{
			AcceptedField: {
				Icon:      "",
				FieldName: WaitlistField,
				Users:     []string{},
			},
		},
	}
	for _, f := range embed.Fields {
		switch {
		case strings.Contains(f.Name, acceptedBase):
			_, limit, err := util.ParseFieldHeadCount(f.Name)
			if err != nil {
				return nil, err
			}
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Roles = append(e.RoleGroup.Roles, &Role{
				Icon:      AcceptedIcon,
				FieldName: AcceptedField,
				Users:     users,
				Count:     len(users),
				Limit:     limit,
			})
		case strings.Contains(f.Name, declinedBase):
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Roles = append(e.RoleGroup.Roles, &Role{
				Icon:      DeclinedIcon,
				FieldName: DeclinedField,
				Users:     users,
				Count:     len(users),
			})
		case strings.Contains(f.Name, tentativeBase):
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Roles = append(e.RoleGroup.Roles, &Role{
				Icon:      TentativeIcon,
				FieldName: TentativeField,
				Users:     users,
				Count:     len(users),
			})
		case f.Name == "Links":
			var err error
			e.Start, e.End, err = util.GetTimesFromLink(f.Value)
			if err != nil {
				return nil, err
			}
		case f.Name == string(WaitlistField):
			users := util.GetUsersFromValues(f.Value)
			e.RoleGroup.Waitlist[AcceptedField] = &Role{
				Icon:      "",
				FieldName: WaitlistField,
				Users:     users,
				Count:     len(users),
			}
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
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Time",
			Value: util.PrintTime(event.Start, event.End),
		},
		{
			Name:  "Links",
			Value: util.PrintAddGoogleCalendarLink(event.Title, event.Description, event.Start, event.End),
		},
	}
	for _, role := range event.RoleGroup.Roles {
		name := fmt.Sprintf("%s %s", role.Icon, role.FieldName)
		if role.Limit == 0 && role.Count > 0 {
			name = fmt.Sprintf("%s %s (%d)", role.Icon, role.FieldName, role.Count)
		}
		if role.Limit > 0 {
			name = fmt.Sprintf("%s %s (%d/%d)", role.Icon, role.FieldName, role.Count, role.Limit)
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  util.NameListToValues(role.Users),
			Inline: true,
		})
	}

	for _, role := range event.RoleGroup.Roles {
		if wl, ok := event.RoleGroup.Waitlist[role.FieldName]; ok && wl.Count > 0 {
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

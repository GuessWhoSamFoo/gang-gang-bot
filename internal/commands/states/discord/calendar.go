package discord

import (
	"context"
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"os"
	"strings"
	"time"
)

type CalendarClient struct {
	ctx        context.Context
	service    *calendar.Service
	calendarID string
}

// NewCalendarClient creates a new client to query a Google Calendar API
func NewCalendarClient(ctx context.Context, calendarID string, credentials []byte) (*CalendarClient, error) {
	temp, err := os.CreateTemp(os.TempDir(), "bot-")
	if err != nil {
		return nil, err
	}
	if _, err := temp.Write(credentials); err != nil {
		return nil, err
	}
	defer os.Remove(temp.Name())
	svc, err := calendar.NewService(ctx, option.WithCredentialsFile(temp.Name()))
	if err != nil {
		return nil, err
	}
	return &CalendarClient{
		ctx:        ctx,
		service:    svc,
		calendarID: calendarID,
	}, nil
}

func (c *CalendarClient) CreateGoogleEvent(event *Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}
	gEvent := toGoogleEvent(event)
	cEvent, err := c.service.Events.Insert(c.calendarID, gEvent).Do()
	if err != nil {
		return err
	}
	event.ID = util.EncodeToGoogleCalendarBase64(cEvent.Id, c.calendarID)
	return nil
}

func (c *CalendarClient) ListEvents() ([]*calendar.Event, error) {
	events, err := c.service.Events.List(c.calendarID).SingleEvents(true).TimeMin(time.Now().Format(time.RFC3339)).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}

func (c *CalendarClient) UpdateEvent(event *Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}
	gEvent := toGoogleEvent(event)
	eventID, _, err := util.DecodeToGoogleEventID(event.ID)
	if err != nil {
		return err
	}
	_, err = c.service.Events.Patch(c.calendarID, eventID, gEvent).Do()
	if err != nil {
		return err
	}
	return nil
}

func (c *CalendarClient) DeleteEvent(event *Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}
	eventID, _, err := util.DecodeToGoogleEventID(event.ID)
	if err != nil {
		return fmt.Errorf("decoded %s: %v", eventID, err)
	}
	if err := c.service.Events.Delete(c.calendarID, eventID).Do(); err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}
	return nil
}

// toGoogleEvent converts the event type to a calendar event
func toGoogleEvent(event *Event) *calendar.Event {
	if event == nil {
		return nil
	}
	gEvent := &calendar.Event{
		Summary:     event.Title,
		Description: event.Description,
		Location:    event.Location,
	}
	if event.DiscordLink != "" && !strings.Contains(event.Description, util.LineFeed) {
		gEvent.Description = util.PrintGoogleCalendarDescription(event.Description, event.DiscordLink)
	}

	if !event.Start.IsZero() {
		gEvent.Start = &calendar.EventDateTime{
			DateTime: event.Start.Format(time.RFC3339),
			TimeZone: util.StaticLocation,
		}
		gEvent.End = &calendar.EventDateTime{
			DateTime: event.End.Format(time.RFC3339),
			TimeZone: util.StaticLocation,
		}
	}
	return gEvent
}

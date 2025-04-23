package googlecalendar

import (
	"log"
	googleapi "mvp-2-spms/integrations/google-api"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const DAYS_PERIOD = 7
const HOURS_IN_DAY = 24
const EVENT_DURATION_HOURS = 1

type googleCalendarApi struct {
	googleapi.Google
	api *calendar.Service
}

func InitCalendarApi(googleAPI googleapi.GoogleAPI) googleCalendarApi {
	c := googleCalendarApi{Google: googleapi.InintGoogle(googleAPI)}
	return c
}

func (c *googleCalendarApi) AuthentificateService(token *oauth2.Token) error {
	if err := c.Authentificate(token); err != nil {
		log.Printf("Ошибка при аутентификации в claendar: %v", err)
	}

	api, err := calendar.NewService(c.GetContext(), option.WithHTTPClient(c.GetClient()))
	if err != nil {
		return err
	}

	c.api = api
	return nil
}

// startTime should be UTC+0!!!
func (c *googleCalendarApi) AddEvent(startTime time.Time, summary string, desc string, calendarId string) (*calendar.Event, error) {
	endTime := strings.Split(startTime.Add(EVENT_DURATION_HOURS*time.Hour).Format(time.RFC3339), "Z")[0]
	event := &calendar.Event{
		Summary:     summary,
		Description: desc,
		Start: &calendar.EventDateTime{
			TimeZone: "Etc/GMT-5",
			DateTime: strings.Split(startTime.Format(time.RFC3339), "Z")[0],
		},
		End: &calendar.EventDateTime{
			TimeZone: "Etc/GMT-5", //////////////////////////////////////?????????????????????????????????????
			DateTime: endTime,
		},
		Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=1"},
	}
	result, err := c.api.Events.Insert(calendarId, event).Do()
	if err == nil {
		return result, nil
	}
	return nil, err
}

func (c *googleCalendarApi) DeleteSlot(eventId string, calendarId string) (*calendar.Event, error) {
	err := c.api.Events.Delete(calendarId, eventId).Do()
	return nil, err
}

func (c *googleCalendarApi) AddSlot(startTime time.Time, duration int, desc string, calendarId string) (*calendar.Event, error) {
	endTime := strings.Split(startTime.Add(time.Duration(duration)*time.Minute).Format(time.RFC3339), "Z")[0]
	event := &calendar.Event{
		Summary:     "SPAMS Slot",
		Description: desc,
		Start: &calendar.EventDateTime{
			TimeZone: "Etc/GMT-5",
			DateTime: strings.Split(startTime.Format(time.RFC3339), "Z")[0],
		},
		End: &calendar.EventDateTime{
			TimeZone: "Etc/GMT-5",
			DateTime: endTime,
		},
		Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=1"},
	}
	result, err := c.api.Events.Insert(calendarId, event).Do()
	if err == nil {
		return result, nil
	}
	return nil, err
}

func (c *googleCalendarApi) UpdateEvent(startTime time.Time, duration int, calendarId string, eventId string) (*calendar.Event, error) {
	endTime := strings.Split(startTime.Add(time.Duration(duration)*time.Minute).Format(time.RFC3339), "Z")[0]
	event, err := c.GetEventById(eventId, calendarId)
	if err != nil {
		log.Printf("event not found")
	}
	event.Start = &calendar.EventDateTime{
		TimeZone: "Etc/GMT-5",
		DateTime: strings.Split(startTime.Format(time.RFC3339), "Z")[0],
	}
	event.End = &calendar.EventDateTime{
		TimeZone: "Etc/GMT-5",
		DateTime: endTime,
	}
	_, err = c.api.Events.Update(calendarId, eventId, event).Do()
	return nil, err
}

func (c *googleCalendarApi) GetEventById(eventId string, calendarId string) (*calendar.Event, error) {
	event, err := c.api.Events.Get(calendarId, eventId).Do()
	if err == nil {
		return event, nil
	}
	return nil, err
}

func (c *googleCalendarApi) GetSchedule(startTime time.Time, calendarId string) (*calendar.Events, error) {
	events, err := c.api.Events.List(calendarId).ShowDeleted(false).SingleEvents(true).TimeMin(startTime.Format(time.RFC3339)).OrderBy("startTime").Do()
	if err == nil {
		return events, nil
	}
	return nil, err
}

func (c *googleCalendarApi) GetAllCalendars() (*calendar.CalendarList, error) {
	calendars, err := c.api.CalendarList.List().Do()
	if err == nil {
		return calendars, nil
	}
	return nil, err
}

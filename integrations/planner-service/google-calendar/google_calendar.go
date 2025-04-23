package googlecalendar

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/api/calendar/v3"
	entities "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/manage-meetings/inputdata"
	"mvp-2-spms/services/models"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type GoogleCalendar struct {
	api googleCalendarApi
}

func InitGoogleCalendar(api googleCalendarApi) *GoogleCalendar {
	return &GoogleCalendar{api: api}
}

func (c *GoogleCalendar) AddMeeting(meeting entities.Meeting, plannerInfo models.PlannerIntegration) (models.PlannerMeeting, error) {
	event, err := c.api.AddEvent(meeting.Time, meeting.Name, meeting.Description, plannerInfo.PlannerData.Id)
	if err != nil {
		return models.PlannerMeeting{}, err
	}

	return models.PlannerMeeting{
		Meeting:          meeting,
		MeetingPlannerId: event.Id,
	}, nil
}

func (c *GoogleCalendar) DeleteSlot(eventId string, plannerInfo models.PlannerIntegration) (models.PlannerSlot, error) {
	_, err := c.api.DeleteSlot(eventId, plannerInfo.PlannerData.Id)
	return models.PlannerSlot{}, err
}

func (c *GoogleCalendar) UpdateSlot(eventId string, slot inputdata.AddSlot, plannerInfo models.PlannerIntegration) error {
	_, err := c.api.UpdateEvent(slot.MeetingTime, slot.Duration, plannerInfo.PlannerData.Id, eventId)
	if err != nil {
		return err
	}
	return nil
}

func (c *GoogleCalendar) AddSlot(slot inputdata.AddSlot, plannerInfo models.PlannerIntegration) (models.PlannerSlot, error) {
	event, err := c.api.AddSlot(slot.MeetingTime, slot.Duration, slot.Description, plannerInfo.PlannerData.Id)
	if err != nil {
		return models.PlannerSlot{}, err
	}
	entity := entities.Slot{
		Description: event.Description,
		ProfessorId: strconv.Itoa(int(slot.ProfessorId)),
		EventId:     event.Id,
		IsOnline:    slot.IsOnline,
	}

	return models.PlannerSlot{
		Slot:             entity,
		MeetingPlannerId: event.Id,
	}, nil
}

func (c *GoogleCalendar) GetAllPlanners() ([]models.PlannerData, error) {
	planners, err := c.api.GetAllCalendars()
	if err != nil {
		return []models.PlannerData{}, err
	}

	result := []models.PlannerData{}
	for _, pl := range planners.Items {
		result = append(result, models.PlannerData{
			Id:   pl.Id,
			Name: pl.Summary,
		})
	}

	return result, nil
}

func (c *GoogleCalendar) GetScheduleMeetingIds(from time.Time, plannerInfo models.PlannerIntegration) ([]string, error) {
	events, err := c.api.GetSchedule(from, plannerInfo.PlannerData.Id)
	if err != nil {
		return []string{}, err
	}

	result := []string{}
	for _, event := range events.Items {
		result = append(result, strings.Split(event.Id, "_")[0])
	}

	return result, nil
}

func (c *GoogleCalendar) FindMeetingById(meetId string, plannerInfo models.PlannerIntegration) (*calendar.Event, error) {
	event, err := c.api.GetEventById(meetId, plannerInfo.PlannerData.Id)
	if err != nil || event.Id == "" {
		return nil, err
	}

	return event, nil
}

func (c *GoogleCalendar) GetAuthLink(redirectURI string, accountId int, returnURL string) (string, error) {
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// encode as JSON!
	statestr := base64.URLEncoding.EncodeToString([]byte(fmt.Sprint(accountId, ",", returnURL)))

	url, err := c.api.GetAuthLink(redirectURI, statestr)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (c *GoogleCalendar) GetToken(code string) (*oauth2.Token, error) {
	token, err := c.api.GetToken(code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (c *GoogleCalendar) Authentificate(token *oauth2.Token) error {
	err := c.api.AuthentificateService(token)
	if err != nil {
		return err
	}

	return nil
}

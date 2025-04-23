package interfaces

import (
	"google.golang.org/api/calendar/v3"
	entities "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/manage-meetings/inputdata"
	"mvp-2-spms/services/models"
	"time"
)

type IPlannerService interface {
	IIntegration
	AddMeeting(meeting entities.Meeting, plannerInfo models.PlannerIntegration) (models.PlannerMeeting, error)
	AddSlot(slot inputdata.AddSlot, plannerInfo models.PlannerIntegration) (models.PlannerSlot, error)
	UpdateSlot(eventId string, slot inputdata.AddSlot, plannerInfo models.PlannerIntegration) error
	DeleteSlot(eventId string, plannerInfo models.PlannerIntegration) (models.PlannerSlot, error)
	FindMeetingById(meetId string, plannerInfo models.PlannerIntegration) (*calendar.Event, error)
	GetScheduleMeetingIds(from time.Time, plannerInfo models.PlannerIntegration) ([]string, error)
	GetAllPlanners() ([]models.PlannerData, error)
}

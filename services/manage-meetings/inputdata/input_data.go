package inputdata

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"time"
)

type GetProfessorMeetings struct {
	ProfessorId uint
	From        *time.Time
	To          *time.Time
}

type AddMeeting struct {
	ProfessorId uint
	ProjectId   uint
	Name        string
	Description string
	MeetingTime time.Time
	StudentId   int
	IsOnline    bool
}

type AddSlot struct {
	ProfessorId uint
	Description string
	MeetingTime time.Time
	Duration    int
	IsOnline    bool
}

func (am *AddMeeting) MapToMeetingEntity() entities.Meeting {
	return entities.Meeting{
		OrganizerId:   fmt.Sprint(am.ProfessorId),
		Name:          am.Name,
		Description:   am.Description,
		ParticipantId: fmt.Sprint(am.StudentId),
		ProjectId:     fmt.Sprint(am.ProjectId),
		Time:          am.MeetingTime,
		IsOnline:      am.IsOnline,
		Status:        entities.MeetingStatus(entities.MeetingPlanned),
	}
}

func (am *AddSlot) MapToSlotEntity() entities.Slot {
	return entities.Slot{
		ProfessorId: fmt.Sprint(am.ProfessorId),
		Description: am.Description,
		EventId:     "",
		IsOnline:    am.IsOnline,
	}
}

package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
	"time"
)

type GetProfesorMeetings struct {
	Meetings []getProfMeetingsData `json:"meetings"`
}

func MapToGetProfesorMeetings(meetings []GetProfesorMeetingsEntities) GetProfesorMeetings {
	outputProjects := []getProfMeetingsData{}
	for _, meet := range meetings {
		id, _ := strconv.Atoi(meet.Meeting.Id)
		projId, _ := strconv.Atoi(meet.Project.Id)
		outputProjects = append(outputProjects,
			getProfMeetingsData{
				Id:          id,
				Name:        meet.Meeting.Name,
				Description: meet.Meeting.Description,
				MeetingTime: meet.Meeting.Time,
				Participant: getProfMeetingsParticipantData{
					FullName: meet.Student.FullNameToString(),
					Cource:   meet.Student.Course,

					ProjectTheme: meet.Project.Theme,
					ProjectId:    projId,
				},
				HasPlannerMeeting: meet.HasPlannerMeeting,
				Status:            meet.Meeting.Status.String(),
				IsOnline:          meet.Meeting.IsOnline,
			})
	}
	return GetProfesorMeetings{
		Meetings: outputProjects,
	}
}

type GetProfesorMeetingsEntities struct {
	Meeting           entities.Meeting
	Student           entities.Student
	Project           entities.Project
	HasPlannerMeeting bool
}

type getProfMeetingsData struct {
	Id                int                            `json:"id"`
	Name              string                         `json:"name"`
	Description       string                         `json:"description"`
	MeetingTime       time.Time                      `json:"time"`
	Participant       getProfMeetingsParticipantData `json:"student"`
	IsOnline          bool                           `json:"is_online"`
	Status            string                         `json:"status"`
	HasPlannerMeeting bool                           `json:"has_planner_meeting"`
}

type getProfMeetingsParticipantData struct {
	FullName     string `json:"name"`
	Cource       uint   `json:"cource"`
	ProjectId    int    `json:"project_id"`
	ProjectTheme string `json:"project_theme"`
}

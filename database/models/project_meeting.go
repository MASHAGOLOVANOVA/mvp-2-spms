package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type ProjectMeeting struct {
	Id            uint `gorm:"column:id"`
	StudMeetingId uint `gorm:"column:stud_meeting_id"`
	ProjectId     uint `gorm:"column:project_id"`
}

func (*ProjectMeeting) TableName() string {
	return "project_meeting"
}

func (pj *ProjectMeeting) MapToThis(slotId int, projectId int) {
	pj.ProjectId = uint(projectId)
	pj.StudMeetingId = uint(slotId)
}

func (pj *ProjectMeeting) MapToEntity() entities.ProjectMeeting {
	return entities.ProjectMeeting{
		Id:            fmt.Sprint(pj.Id),
		StudMeetingId: fmt.Sprint(pj.StudMeetingId),
		ProjectId:     fmt.Sprint(pj.ProjectId),
	}
}

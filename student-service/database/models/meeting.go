package models

import (
	"database/sql"
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
	"time"
)

type Meeting struct {
	Id                   uint           `gorm:"column:id"`
	Name                 string         `gorm:"column:name"`
	Description          string         `gorm:"column:description"`
	MeetingTime          time.Time      `gorm:"column:meeting_time"`
	StudentParticipantId uint           `gorm:"column:student_id"`
	ProfessorId          uint           `gorm:"column:professor_id"`
	ProjectId            sql.NullInt64  `gorm:"column:project_id"`
	IsOnline             bool           `gorm:"column:is_online"`
	Status               uint           `gorm:"column:status"`
	PlannerId            sql.NullString `gorm:"column:planner_id;default:null"`
}

func (*Meeting) TableName() string {
	return "meeting"
}

func (pj *Meeting) MapToEntity() entities.Meeting {
	return entities.Meeting{
		Id:            fmt.Sprint(pj.Id),
		Name:          pj.Name,
		Description:   pj.Description,
		Time:          pj.MeetingTime,
		OrganizerId:   fmt.Sprint(pj.ProfessorId),
		ParticipantId: fmt.Sprint(pj.StudentParticipantId),
		IsOnline:      pj.IsOnline,
		Status:        entities.MeetingStatus(pj.Status),
	}
}

func (pj *Meeting) MapEntityToThis(entity entities.Meeting) {
	mId, _ := strconv.Atoi(entity.Id)
	prId, _ := strconv.Atoi(entity.OrganizerId)
	stId, _ := strconv.Atoi(entity.ParticipantId)
	pj.Id = uint(mId)
	pj.Name = entity.Name
	pj.Description = entity.Description
	pj.MeetingTime = entity.Time
	pj.StudentParticipantId = uint(stId)
	pj.ProfessorId = uint(prId)
	pj.IsOnline = entity.IsOnline
	pj.Status = uint(entity.Status)
	pid, _ := strconv.Atoi(entity.ProjectId)
	if pid != 0 {
		pj.ProjectId.Valid = true
		pj.ProjectId.Int64 = int64(pid)
	}
}

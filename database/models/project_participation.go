package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type ProjectParticipation struct {
	Id        uint `gorm:"column:id"`
	ProjectId uint `gorm:"column:project_id"`
	StudentId uint `gorm:"column:student_id"`
}

func (ProjectParticipation) TableName() string {
	return "project_participation"
}

func (pj ProjectParticipation) MapToEntity() entities.ProjectParticipation {
	return entities.ProjectParticipation{
		Id:        fmt.Sprint(pj.Id),
		StudentId: fmt.Sprint(pj.StudentId),
		ProjectId: fmt.Sprint(pj.ProjectId),
	}
}

func (p *ProjectParticipation) MapEntityToThis(entity entities.ProjectParticipation) {
	prId, _ := strconv.Atoi(entity.Id)
	prStudentId, _ := strconv.Atoi(entity.StudentId)
	pProjectId, _ := strconv.Atoi(entity.ProjectId)
	p.Id = uint(prId)
	p.StudentId = uint(prStudentId)
	p.ProjectId = uint(pProjectId)
}

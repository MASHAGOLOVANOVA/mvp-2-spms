package models

import (
	"database/sql"
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type Project struct {
	Id                 uint            `gorm:"column:id"`
	Theme              string          `gorm:"column:theme"`
	Year               uint            `gorm:"column:year"`
	SupervisorId       uint            `gorm:"column:supervisor_id"`
	DefenceGrade       sql.NullFloat64 `gorm:"column:defence_grade;default:null"`
	FinalGrade         sql.NullFloat64 `gorm:"column:final_grade;default:null"`
	SupervisorReviewId sql.NullInt64   `gorm:"column:supervisor_review_id;default:null"`
	RepoId             sql.NullInt64   `gorm:"column:repo_id;default:null"`
	CloudId            sql.NullString  `gorm:"column:cloud_id;default:null"`
	StageId            uint            `gorm:"column:stage_id"`
	StatusId           uint            `gorm:"column:status_id"`
}

func (Project) TableName() string {
	return "project"
}

func (pj Project) MapToEntity(studentIds []string) entities.Project {
	return entities.Project{
		Id:           fmt.Sprint(pj.Id),
		Theme:        pj.Theme,
		SupervisorId: fmt.Sprint(pj.SupervisorId),
		Year:         pj.Year,
		Stage:        entities.ProjectStage(pj.StageId),
		Status:       entities.ProjectStatus(pj.StatusId),
		StudentIds:   studentIds,
	}
}

func (p *Project) MapEntityToThis(entity entities.Project) {
	prId, _ := strconv.Atoi(entity.Id)
	pId, _ := strconv.Atoi(entity.SupervisorId)
	p.Id = uint(prId)
	p.Theme = entity.Theme
	p.Year = entity.Year
	p.SupervisorId = uint(pId)
	p.StageId = uint(entity.Stage)
	p.StatusId = uint(entity.Status)
}

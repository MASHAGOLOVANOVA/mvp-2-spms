package models

import (
	"database/sql"
	entities "mvp-2-spms/domain-aggregate"
	"time"
)

type SupervisorReview struct {
	Id           sql.NullInt64 `gorm:"column:id"`
	CreationDate time.Time     `gorm:"column:creation_date"`
}

func (*SupervisorReview) TableName() string {
	return "supervisor_review"
}

func (sr *SupervisorReview) MapToEntity(criterias []entities.Criteria) entities.SupervisorReview {
	return entities.SupervisorReview{
		Id:           uint(sr.Id.Int64),
		CreationDate: sr.CreationDate,
		Criterias:    criterias,
	}
}

func (sr *SupervisorReview) MapEntityToThis(e entities.SupervisorReview) {
	sr.CreationDate = e.CreationDate
	if e.Id != 0 {
		sr.Id = sql.NullInt64{
			Int64: int64(e.Id),
			Valid: true,
		}
	}
}

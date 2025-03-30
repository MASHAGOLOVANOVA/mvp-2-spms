package models

import (
	"database/sql"
	entities "mvp-2-spms/domain-aggregate"
)

type ReviewCriteria struct {
	Id                 sql.NullInt64   `gorm:"column:id"`
	Description        string          `gorm:"column:description"`
	Grade              sql.NullFloat64 `gorm:"column:grade"`
	GradeWeight        float32         `gorm:"column:grade_weight"`
	SupervieorReviewId uint            `gorm:"column:supervisor_review_id"`
}

func (*ReviewCriteria) TableName() string {
	return "review_criteria"
}

func (r *ReviewCriteria) MapToEntity() entities.Criteria {
	return entities.Criteria{
		Description: r.Description,
		Grade:       float32(r.Grade.Float64),
		Weight:      r.GradeWeight,
	}
}

func (r *ReviewCriteria) MapEntityToThis(e entities.Criteria) {
	r.Description = e.Description
	if e.Grade != 0 {
		r.Grade = sql.NullFloat64{
			Float64: float64(e.Grade),
			Valid:   true,
		}
	}
	r.GradeWeight = e.Weight
}

package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type Department struct {
	Id    uint   `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	UniId uint   `gorm:"column:uni_id"`
}

func (Department) TableName() string {
	return "department"
}

func (d Department) MapToEntity() entities.Department {
	return entities.Department{
		Id:           fmt.Sprint(d.Id),
		Name:         d.Name,
		UniversityId: (fmt.Sprint(d.UniId)),
	}
}

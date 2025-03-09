package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type University struct {
	Id     uint   `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	CityId uint   `gorm:"column:city_id"`
}

func (University) TableName() string {
	return "university"
}

func (u University) MapToEntity() entities.University {
	return entities.University{
		Id:   fmt.Sprint(u.Id),
		Name: u.Name,
		City: (fmt.Sprint(u.CityId)),
	}
}

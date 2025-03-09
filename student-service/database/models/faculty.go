package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type Faculty struct {
	Id     uint   `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	DeptId uint   `gorm:"column:dept_id"`
}

func (Faculty) TableName() string {
	return "faculty"
}

func (d Faculty) MapToEntity() entities.Faculty {
	return entities.Faculty{
		Id:           fmt.Sprint(d.Id),
		Name:         d.Name,
		DepartmentId: (fmt.Sprint(d.DeptId)),
	}
}

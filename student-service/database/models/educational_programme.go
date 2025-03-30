package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type EducationalProgramme struct {
	Id        uint   `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	FacultyId uint   `gorm:"column:faculty_id"`
	EdLevel   uint   `gorm:"column:ed_level"`
}

func (EducationalProgramme) TableName() string {
	return "educational_programme"
}

func (ep EducationalProgramme) MapToEntity() entities.EducationalProgramme {
	return entities.EducationalProgramme{
		Id:               fmt.Sprint(ep.Id),
		Name:             ep.Name,
		EducationalLevel: entities.EducationalLevel(ep.EdLevel),
		FacultyId:        fmt.Sprint(ep.FacultyId),
	}
}

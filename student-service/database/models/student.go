package models

import (
	"database/sql"
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
	"time"
)

type Student struct {
	Id                     uint          `gorm:"column:id"`
	Name                   string        `gorm:"column:name"`
	Surname                string        `gorm:"column:surname"`
	Middlename             string        `gorm:"column:middlename"`
	EnrollmentYear         uint          `gorm:"column:enrollment_year"`
	EducationalProgrammeId sql.NullInt64 `gorm:"column:educational_programme_id;default:null"`
}

func (*Student) TableName() string {
	return "student"
}

func (s *Student) MapToEntity() entities.Student {
	return entities.Student{
		Person: entities.Person{
			Id:         fmt.Sprint(s.Id),
			Name:       s.Name,
			Surname:    s.Surname,
			Middlename: s.Middlename,
		},
		EducationalProgrammeId: fmt.Sprint(s.EducationalProgrammeId.Int64),
		Cource:                 s.GetCource(),
		//EnrollmentYear: s.EnrollmentYear,
	}
}

func (s *Student) MapEntityToThis(entity entities.Student) {
	sId, _ := strconv.Atoi(entity.Id)
	epId, _ := strconv.Atoi(entity.EducationalProgrammeId)
	s.Id = uint(sId)
	s.Name = entity.Name
	s.Surname = entity.Surname
	s.Middlename = entity.Middlename
	s.EnrollmentYear = getStudentEnrollmentYear(entity.Cource)
	if epId != 0 {
		s.EducationalProgrammeId = sql.NullInt64{Int64: int64(epId), Valid: true}
	}
}

func (s *Student) GetCource() uint {
	currentDate := time.Now()
	if currentDate.Month() > 9 {
		return uint(currentDate.Year()) - s.EnrollmentYear + 1
	}
	return uint(currentDate.Year()) - s.EnrollmentYear
}

func getStudentEnrollmentYear(cource uint) uint {
	currentDate := time.Now()
	if currentDate.Month() > 9 {
		return uint(currentDate.Year()) - cource + 1
	}
	return uint(currentDate.Year()) - cource
}

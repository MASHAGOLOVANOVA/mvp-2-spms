package models

import (
	"database/sql"
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"
	"strconv"
	"time"
)

type Student struct {
	Id                   uint   `gorm:"column:id"`
	Name                 string `gorm:"column:name"`
	Surname              string `gorm:"column:surname"`
	Middlename           string `gorm:"column:middlename"`
	EnrollmentYear       uint   `gorm:"column:enrollment_year"`
	EducationalProgramme string `gorm:"column:ed_program"`
	University           string `gorm:"column:university"`
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
		EducationalProgramme: s.EducationalProgramme,
		Course:               s.GetCource(),
		University:           s.University,
	}
}

func (s *Student) MapEntityToThis(entity entities.Student) {
	sId, _ := strconv.Atoi(entity.Id)
	s.Id = uint(sId)
	s.Name = entity.Name
	s.Surname = entity.Surname
	s.Middlename = entity.Middlename
	s.EnrollmentYear = getStudentEnrollmentYear(entity.Course)
	s.EducationalProgramme = entity.EducationalProgramme
	s.University = entity.University
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

type StudentAccount struct {
	Id        uint          `gorm:"column:id"`
	Login     string        `gorm:"column:login"`
	StudentId sql.NullInt64 `gorm:"column:student_id;default:null"`
}

func (*StudentAccount) TableName() string {
	return "student_account"
}

func (sa *StudentAccount) MapToEntity() entities.StudentAccount {
	return entities.StudentAccount{
		Id:        fmt.Sprint(sa.Id),
		Login:     sa.Login,
		StudentId: fmt.Sprint(sa.StudentId.Int64),
	}
}

func (sa *StudentAccount) MapEntityToThis(entity entities.StudentAccount) {
	sId, _ := strconv.Atoi(entity.Id)
	sa.Id = uint(sId)
	sa.Login = entity.Login
	if entity.StudentId != "" {
		studentId, _ := strconv.Atoi(entity.StudentId)
		sa.StudentId = sql.NullInt64{Int64: int64(studentId), Valid: true}
	} else {
		sa.StudentId = sql.NullInt64{Valid: false}
	}
}

func (sa *StudentAccount) MapModelToThis(entity usecasemodels.StudentAccount) {
	sId, _ := strconv.Atoi(entity.Id)
	sa.Id = uint(sId)
	sa.Login = entity.Login
	if entity.StudentId != "" {
		studentId, _ := strconv.Atoi(entity.StudentId)
		sa.StudentId = sql.NullInt64{Int64: int64(studentId), Valid: true}
	} else {
		sa.StudentId = sql.NullInt64{Valid: false}
	}
}

func (sa *StudentAccount) MapToUseCaseModel() usecasemodels.StudentAccount {
	return usecasemodels.StudentAccount{
		Login:     sa.Login,
		Id:        fmt.Sprint(sa.Id),
		StudentId: fmt.Sprint(sa.Id),
	}
}

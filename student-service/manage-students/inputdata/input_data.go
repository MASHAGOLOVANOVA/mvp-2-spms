package inputdata

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type GetStudents struct {
	ProfessorId uint
}

type AddStudent struct {
	ProfessorId            uint
	Name                   string
	Surname                string
	Middlename             string
	EducationalProgrammeId uint
	Cource                 uint
}

func (as AddStudent) MapToStudentEntity() entities.Student {
	return entities.Student{
		Person: entities.Person{
			Name:       as.Name,
			Surname:    as.Surname,
			Middlename: as.Middlename,
		},
		EducationalProgrammeId: fmt.Sprint(as.EducationalProgrammeId),
		Cource:                 as.Cource,
	}
}

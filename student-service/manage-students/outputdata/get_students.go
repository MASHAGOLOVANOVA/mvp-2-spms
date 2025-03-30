package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetStudents struct {
	Projects []getStudentsData `json:"students"`
}

func MapToGetStudents(studentEntities []GetStudentsEntities) GetStudents {
	outputStudents := []getStudentsData{}
	for _, studenttEntitiy := range studentEntities {
		id, _ := strconv.Atoi(studenttEntitiy.Student.Id)
		pid, _ := strconv.Atoi(studenttEntitiy.PtojectId)
		outputStudents = append(outputStudents,
			getStudentsData{
				Id:                   id,
				Name:                 studenttEntitiy.Student.Name,
				Surname:              studenttEntitiy.Student.Surname,
				Middlename:           studenttEntitiy.Student.Middlename,
				EducationalProgramme: studenttEntitiy.EducationalProgramme,
				Cource:               int(studenttEntitiy.Student.Cource),
				PtojectTheme:         studenttEntitiy.ProjectTheme,
				PtojectId:            pid,
			})
	}
	return GetStudents{
		Projects: outputStudents,
	}
}

type GetStudentsEntities struct {
	PtojectId            string
	ProjectTheme         string
	Student              entities.Student
	EducationalProgramme string
}

type getStudentsData struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Surname              string `json:"surname"`
	Middlename           string `json:"middlename"`
	EducationalProgramme string `json:"education_programme"`
	Cource               int    `json:"cource"`
	PtojectTheme         string `json:"project_theme"`
	PtojectId            int    `json:"project_id"`
}

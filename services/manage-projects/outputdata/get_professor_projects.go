package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetProfessorProjects struct {
	Projects []getProfProjProjectData `json:"projects"`
}

type ProjStudentData struct {
	StudentName       string `json:"student_name"`
	StudentCourse     int    `json:"student_course"`
	StudentEdProgram  string `json:"student_ed_program"`
	StudentUniversity string `json:"student_university"`
}

func MapToGetProfessorProjects(projectEntities []GetProfessorProjectsEntities) GetProfessorProjects {
	outputProjects := []getProfProjProjectData{}
	for _, projectEntitiy := range projectEntities {
		id, _ := strconv.Atoi(projectEntitiy.Project.Id)

		studentData := make([]ProjStudentData, len(projectEntitiy.Students))
		for i, student := range projectEntitiy.Students {
			studentData[i] = ProjStudentData{
				StudentName:       student.FullNameToString(),
				StudentCourse:     int(student.Course),
				StudentEdProgram:  student.EducationalProgramme,
				StudentUniversity: student.University,
			}
		}

		outputProjects = append(outputProjects,
			getProfProjProjectData{
				Id:       id,
				Theme:    projectEntitiy.Project.Theme,
				Status:   projectEntitiy.Project.Status.String(),
				Stage:    projectEntitiy.Project.Stage.String(),
				Year:     int(projectEntitiy.Project.Year),
				Students: studentData,
			})
	}
	return GetProfessorProjects{
		Projects: outputProjects,
	}
}

type GetProfessorProjectsEntities struct {
	Project  entities.Project
	Students []entities.Student
}

type getProfProjProjectData struct {
	Id       int               `json:"id"`
	Theme    string            `json:"theme"`
	Students []ProjStudentData `json:"students"`
	Status   string            `json:"status"`
	Stage    string            `json:"stage"`
	Year     int               `json:"year"`
}

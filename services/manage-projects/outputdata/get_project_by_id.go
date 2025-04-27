package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetProjectById struct {
	Id              int                         `json:"id"`
	Theme           string                      `json:"theme"`
	Students        []GetProjectByIdStudentData `json:"students"`
	Status          string                      `json:"status"`
	Stage           string                      `json:"stage"`
	Year            int                         `json:"year"`
	CloudFolderLink string                      `json:"cloud_folder_link"`
}

func MapToGetProjectsById(project entities.Project, students []entities.Student, folderLink string) GetProjectById {
	pId, _ := strconv.Atoi(project.Id)

	// Convert students slice to GetProjectByIdStudentData slice
	studentData := make([]GetProjectByIdStudentData, len(students))
	for i, student := range students {
		sId, _ := strconv.Atoi(student.Id)
		studentData[i] = GetProjectByIdStudentData{
			Id:         sId,
			Name:       student.Name,
			Surname:    student.Surname,
			Middlename: student.Middlename,
			Cource:     int(student.Course),
		}
	}

	return GetProjectById{
		Id:              pId,
		Theme:           project.Theme,
		Students:        studentData,
		Status:          project.Status.String(),
		Stage:           project.Stage.String(),
		Year:            int(project.Year),
		CloudFolderLink: folderLink,
	}
}

type GetProjectByIdStudentData struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Middlename string `json:"middlename"`
	Cource     int    `json:"cource"`
}

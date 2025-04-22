package interfaces

import (
	"mvp-2-spms/services/manage-professors/inputdata"
	"mvp-2-spms/services/manage-professors/outputdata"
)

type IProfessorInteractor interface {
	GetProfessors() (outputdata.GetProfessors, error)
	GetStudentApplications(studentId string) (outputdata.GetApplications, error)
	GetProfessorApplications(profId string) (outputdata.GetApplications, error)
	Apply(apply inputdata.Apply) error
	UpdateApplicationStatus(apply inputdata.ApplicationStatus) error
}

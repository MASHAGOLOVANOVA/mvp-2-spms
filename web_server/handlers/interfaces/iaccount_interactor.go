package interfaces

import (
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/manage-accounts/outputdata"
)

type IAccountInteractor interface {
	GetDriveBaseFolderName(folderId, profId string, cloudDrive interfaces.ICloudDrive) (string, error)
	GetPlannerIntegration(input inputdata.GetPlannerIntegration) (outputdata.GetPlannerIntegration, error)
	GetDriveIntegration(input inputdata.GetDriveIntegration) (outputdata.GetDriveIntegration, error)
	GetRepoHubIntegration(input inputdata.GetRepoHubIntegration) (outputdata.GetRepoHubIntegration, error)
	SetPlannerIntegration(input inputdata.SetPlannerIntegration, planner interfaces.IPlannerService) (outputdata.SetPlannerIntegration, error)
	SetDriveIntegration(input inputdata.SetDriveIntegration, planner interfaces.ICloudDrive) (outputdata.SetDriveIntegration, error)
	SetRepoHubIntegration(input inputdata.SetRepoHubIntegration, planner interfaces.IGitRepositoryHub) (outputdata.SetRepoHubIntegration, error)
	GetAccountIntegrations(input inputdata.GetAccountIntegrations) (outputdata.GetAccountIntegrations, error)
	GetProfessorInfo(input inputdata.GetProfessorInfo) (outputdata.GetProfessorInfo, error)
	GetStudentInfo(input inputdata.GetStudentInfo) (outputdata.GetStudentInfo, error)
	CheckCredsValidity(input inputdata.CheckCredsValidity) (bool, error)
	CheckUsernameExists(input inputdata.CheckUsernameExists) (bool, error)
	CheckStudentExists(input inputdata.CheckStudentExists) (bool, error)
	SignUp(input inputdata.SignUp) (outputdata.SignUp, error)
	StudentSignUp(input inputdata.StudentSignUp) (outputdata.SignUp, error)
	GetAccountProfessorId(login string) (string, error)
	GetAccountStudentId(login string) (string, error)
	SetProfessorPlanner(plannerId, profId string) error
	GetProfessorIntegrPlanners(profId string, planner interfaces.IPlannerService) (outputdata.GetProfessorIntegrPlanners, error)
}

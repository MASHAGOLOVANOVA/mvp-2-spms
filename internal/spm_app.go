package internal

import (
	mngInterfaces "mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/models"
	hInterfaces "mvp-2-spms/web_server/handlers/interfaces"
)

type GitRepositoryHubs map[models.GetRepoHubName](mngInterfaces.IGitRepositoryHub)
type CloudDrives map[models.CloudDriveName](mngInterfaces.ICloudDrive)
type Planners map[models.PlannerName](mngInterfaces.IPlannerService)

type StudentsProjectsManagementApp struct {
	Intercators  Intercators
	Integrations Integrations
}

type Intercators struct {
	AccountManager   hInterfaces.IAccountInteractor
	ProjectManager   hInterfaces.IProjetInteractor
	StudentManager   hInterfaces.IStudentInteractor
	MeetingManager   hInterfaces.IMeetingInteractor
	TaskManager      hInterfaces.ITaskInteractor
	UnversityManager hInterfaces.IUniversityInteractor
	ProfessorManager hInterfaces.IProfessorInteractor
}

type Integrations struct {
	GitRepositoryHubs
	CloudDrives
	Planners
}

type Repositories struct {
	Projects     mngInterfaces.IProjetRepository
	Students     mngInterfaces.IStudentRepository
	Universities mngInterfaces.IUniversityRepository
	Meetings     mngInterfaces.IMeetingRepository
	Accounts     mngInterfaces.IAccountRepository
	Tasks        mngInterfaces.ITaskRepository
	Professors   mngInterfaces.IProfessorRepository
}

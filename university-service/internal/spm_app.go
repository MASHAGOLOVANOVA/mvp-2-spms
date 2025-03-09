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
	UniversityManager hInterfaces.IUniversityInteractor
}

type Integrations struct {
}

type Repositories struct {
	Universities mngInterfaces.IUniversityRepository
}

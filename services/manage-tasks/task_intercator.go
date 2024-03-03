package managetasks

import (
	"fmt"
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-tasks/inputdata"
	"mvp-2-spms/services/manage-tasks/outputdata"
)

type TaskInteractor struct {
	projectRepo interfaces.IProjetRepository
	cloudDrive  interfaces.ICloudDrive
	//accountRepo interfaces.IAccountRepository
	taskRepo interfaces.ITaskRepository
}

func InitTaskInteractor(projRepo interfaces.IProjetRepository, cloudDrive interfaces.ICloudDrive) *TaskInteractor {
	return &TaskInteractor{
		projectRepo: projRepo,
		cloudDrive:  cloudDrive,
	}
}

func (p *TaskInteractor) AddTask(input inputdata.AddTask) outputdata.AddTask {
	// add to db
	task := p.taskRepo.CreateTask(input.MapToTaskEntity())
	// get project folder id
	projFolder := p.projectRepo.GetProjectCloudFolderId(fmt.Sprint(input.ProjectId))
	// add folder to cloud (create folder and task file)
	driveTask := p.cloudDrive.AddTaskToDrive(task, projFolder)
	// add folder id and file id from drive
	p.taskRepo.AssignDriveTask(driveTask)
	// returning id
	output := outputdata.MapToAddTask(task)
	return output
}

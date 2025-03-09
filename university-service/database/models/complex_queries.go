package models

import (
	entities "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/models"
)

type ProjectTaskInfo struct {
	Statuses []statusCount
}

type statusCount struct {
	Status int `gorm:"column:status"`
	Count  int `gorm:"column:count"`
}

func (pti *ProjectTaskInfo) MapToUseCaseModel() models.TasksInfo {
	result := models.TasksInfo{}
	for _, s := range pti.Statuses {
		switch s.Status {
		case int(entities.NotStarted):
			{
				result.NotStartedCount = s.Count
			}
		case int(entities.InProgress):
			{
				result.InProgressCount = s.Count
			}
		case int(entities.Finished):
			{
				result.FinishedCount = s.Count
			}
		}
	}
	return result
}

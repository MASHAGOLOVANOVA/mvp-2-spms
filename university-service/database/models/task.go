package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
	"time"
)

type Task struct {
	Id          uint      `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Deadline    time.Time `gorm:"column:deadline"`
	ProjectId   uint      `gorm:"column:project_id"`
	Status      uint      `gorm:"column:status"`
	FolderId    string    `gorm:"column:folder_id;default:null"`
	TaskFileId  string    `gorm:"column:task_file_id;default:null"`
}

func (*Task) TableName() string {
	return "task"
}

func (t *Task) MapToEntity() entities.Task {
	return entities.Task{
		Id:          fmt.Sprint(t.Id),
		Name:        t.Name,
		Description: t.Description,
		Deadline:    t.Deadline,
		ProjectId:   fmt.Sprint(t.ProjectId),
		Status:      entities.TaskStatus(t.Status),
	}
}

func (t *Task) MapEntityToThis(entity entities.Task) {
	mId, _ := strconv.Atoi(entity.Id)
	prId, _ := strconv.Atoi(entity.ProjectId)
	t.Id = uint(mId)
	t.Name = entity.Name
	t.Description = entity.Description
	t.ProjectId = uint(prId)
	t.Deadline = entity.Deadline
	t.Status = uint(entity.Status)
}

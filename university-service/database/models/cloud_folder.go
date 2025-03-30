package models

import (
	"mvp-2-spms/services/models"
)

type CloudFolder struct {
	Id   string `gorm:"column:id"`
	Link string `gorm:"column:link"`
}

func (*CloudFolder) TableName() string {
	return "cloud_folder"
}

func (cf *CloudFolder) MapToUseCaseModel() models.DriveFolder {
	return models.DriveFolder{
		Id:   cf.Id,
		Link: cf.Link,
	}
}

func (cf *CloudFolder) MapUseCaseModelToThis(model models.DriveFolder) {
	cf.Id = model.Id
	cf.Link = model.Link
}

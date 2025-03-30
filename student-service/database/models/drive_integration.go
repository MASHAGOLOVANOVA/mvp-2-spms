package models

import (
	"fmt"
	"mvp-2-spms/services/models"
	"strconv"
)

type DriveIntegration struct {
	AccountId    uint   `gorm:"column:account_id"`
	ApiKey       string `gorm:"column:api_key"`
	BaseFolderId string `gorm:"column:base_folder_id"`
	Type         int    `gorm:"column:type"`
}

func (*DriveIntegration) TableName() string {
	return "drive_integration"
}

func (di *DriveIntegration) MapToUseCaseModel() models.CloudDriveIntegration {
	return models.CloudDriveIntegration{
		BaseIntegration: models.BaseIntegration{
			AccountId: fmt.Sprint(di.AccountId),
			ApiKey:    di.ApiKey,
			Type:      di.Type,
		},
		DriveData: models.DriveData{
			BaseFolderId: di.BaseFolderId,
		},
	}
}

func (di *DriveIntegration) MapUseCaseModelToThis(model models.CloudDriveIntegration) {
	sId, _ := strconv.Atoi(model.AccountId)
	di.AccountId = uint(sId)
	di.BaseFolderId = model.BaseFolderId
	di.ApiKey = model.ApiKey
	di.Type = model.Type
}

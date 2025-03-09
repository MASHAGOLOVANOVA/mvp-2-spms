package models

import (
	"fmt"
	"mvp-2-spms/services/models"
	"strconv"
)

type GitRepositoryIntegration struct {
	AccountId uint   `gorm:"column:account_id"`
	ApiKey    string `gorm:"column:api_key"`
	Type      int    `gorm:"column:type"`
}

func (*GitRepositoryIntegration) TableName() string {
	return "git_repository_integration"
}

func (gr *GitRepositoryIntegration) MapToUseCaseModel() models.BaseIntegration {
	return models.BaseIntegration{
		AccountId: fmt.Sprint(gr.AccountId),
		ApiKey:    gr.ApiKey,
		Type:      gr.Type,
	}
}

func (gr *GitRepositoryIntegration) MapUseCaseModelToThis(model models.BaseIntegration) {
	sId, _ := strconv.Atoi(model.AccountId)
	gr.AccountId = uint(sId)
	gr.ApiKey = model.ApiKey
	gr.Type = model.Type
}

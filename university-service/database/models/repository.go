package models

import (
	"mvp-2-spms/services/models"
)

type Repository struct {
	Id          uint   `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	OwnerName   string `gorm:"column:owner_name"`
	IsPublic    bool   `gorm:"column:is_public"`
	RepoHubType int    `gorm:"column:repo_hub_type"`
}

func (Repository) TableName() string {
	return "repository"
}

func (r Repository) MapToUseCaseModel() models.Repository {
	return models.Repository{
		RepoId:    r.Name,
		OwnerName: r.OwnerName,
	}
}

func (r *Repository) MapModelToThis(model models.Repository) {
	r.Name = model.RepoId
	r.OwnerName = model.OwnerName
	r.IsPublic = true
	r.RepoHubType = model.RepoType
}

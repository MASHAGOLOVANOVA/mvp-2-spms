package models

import (
	entities "mvp-2-spms/domain-aggregate"
	"time"
)

type Repository struct {
	RepoId    string
	OwnerName string
	Commits   []Commit
	Branches  []branch
	RepoType  int
	ProjectId int
}

type branch struct {
	Name string
	Head Commit
}

type Commit struct {
	SHA         string
	Description string
	Date        time.Time
	Author      string
	// id repo
}

type ProjectInRepository struct {
	Repository
	entities.Project
}

package domainaggregate

import (
	"fmt"
)

type ProjectStatus int

const (
	ProjectNotConfirmed ProjectStatus = iota
	ProjectInProgress
	ProjectFinished
	ProjectCancelled
)

type ProjectStage int

const (
	Analysis ProjectStage = iota
	Design
	Development
	Testing
	Deployment
)

type Project struct {
	Id           string
	Theme        string
	SupervisorId string
	StudentIds   []string
	Year         uint
	Stage        ProjectStage
	Status       ProjectStatus
}

func (s ProjectStatus) String() string {
	switch s {
	case ProjectStatus(ProjectNotConfirmed):
		return "NotConfirmed"
	case ProjectStatus(ProjectInProgress):
		return "InProgress"
	case ProjectStatus(ProjectFinished):
		return "Finished"
	case ProjectStatus(ProjectCancelled):
		return "Cancelled"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}

func (s ProjectStage) String() string {
	switch s {
	case ProjectStage(Analysis):
		return "Analysis"
	case ProjectStage(Design):
		return "Design"
	case ProjectStage(Development):
		return "Development"
	case ProjectStage(Testing):
		return "Testing"
	case ProjectStage(Deployment):
		return "Deployment"
	default:
		return fmt.Sprintf("%d", int(s))
	}
}

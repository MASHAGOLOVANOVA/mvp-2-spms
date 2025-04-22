package inputdata

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type ApplicationStatus struct {
	Id          int
	ProfessorId int  `json:"professor_id,omitempty"`
	StudentId   int  `json:"student_id,omitempty"`
	Status      bool `json:"status,omitempty"`
}

func (a ApplicationStatus) MapToApplyEntity() entities.Apply {
	return entities.Apply{
		Id:          fmt.Sprint(a.Id),
		StudentId:   fmt.Sprint(a.StudentId),
		ProfessorId: fmt.Sprint(a.ProfessorId),
		Status:      fmt.Sprint(a.Status),
	}
}

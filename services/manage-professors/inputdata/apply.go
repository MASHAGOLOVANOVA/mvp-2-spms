package inputdata

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type Apply struct {
	ProfessorId uint
	StudentId   uint
}

func (a Apply) MapToApplyEntity() entities.Apply {
	return entities.Apply{
		StudentId:   fmt.Sprint(a.StudentId),
		ProfessorId: fmt.Sprint(a.ProfessorId),
	}
}

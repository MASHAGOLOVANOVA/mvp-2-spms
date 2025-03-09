package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type AddStudent struct {
	Id int `json:"student_id"`
}

func MapToAddStudent(student entities.Student) AddStudent {
	sId, _ := strconv.Atoi(student.Id)
	return AddStudent{
		Id: sId,
	}
}

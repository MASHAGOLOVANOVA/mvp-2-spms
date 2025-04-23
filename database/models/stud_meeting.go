package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
)

type StudMeeting struct {
	Id        uint `gorm:"column:id"`
	SlotId    uint `gorm:"column:slot_id"`
	StudentId uint `gorm:"column:student_id"`
}

func (*StudMeeting) TableName() string {
	return "student_meeting"
}

func (pj *StudMeeting) MapToThis(slotId int, studId int) {
	pj.StudentId = uint(studId)
	pj.SlotId = uint(slotId)
}

func (pj *StudMeeting) MapToEntity() entities.StudMeeting {
	return entities.StudMeeting{
		Id:        fmt.Sprint(pj.Id),
		SlotId:    fmt.Sprint(pj.SlotId),
		StudentId: fmt.Sprint(pj.StudentId),
	}
}

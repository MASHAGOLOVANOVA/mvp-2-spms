package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type Slot struct {
	Id          uint   `gorm:"column:id"`
	Description string `gorm:"column:description"`
	EventId     string `gorm:"column:event_id"`
	ProfessorId uint   `gorm:"column:professor_id"`
	IsOnline    bool   `gorm:"column:is_online"`
	PlannerId   string `gorm:"column:planner_id"`
	Status      int    `gorm:"column:status"`
}

func (*Slot) TableName() string {
	return "slot"
}

func (pj *Slot) MapToEntity() entities.Slot {
	return entities.Slot{
		Id:          fmt.Sprint(pj.Id),
		Description: pj.Description,
		ProfessorId: fmt.Sprint(pj.ProfessorId),
		EventId:     fmt.Sprint(pj.EventId),
		IsOnline:    pj.IsOnline,
	}
}

func (pj *Slot) MapEntityToThis(entity entities.Slot, plannerId string) {
	mId, _ := strconv.Atoi(entity.Id)
	prId, _ := strconv.Atoi(entity.ProfessorId)
	pj.Id = uint(mId)
	pj.Description = entity.Description
	pj.EventId = entity.EventId
	pj.ProfessorId = uint(prId)
	pj.IsOnline = entity.IsOnline
	pj.PlannerId = plannerId
	pj.Status = 0
}

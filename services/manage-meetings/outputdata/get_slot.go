package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetSlot struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	ProfessorId string `json:"professor_id"`
	EventId     string `json:"event_id"`
	IsOnline    bool   `json:"is_online"`
}

func MapToGetSlot(slot entities.Slot) GetSlot {
	sId, _ := strconv.Atoi(slot.Id)
	return GetSlot{
		Id:          sId,
		Description: slot.Description,
		ProfessorId: slot.ProfessorId,
		EventId:     slot.EventId,
		IsOnline:    slot.IsOnline,
	}
}

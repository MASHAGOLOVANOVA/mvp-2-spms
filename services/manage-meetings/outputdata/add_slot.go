package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type AddSlot struct {
	Id int `json:"id"`
}

func MapToAddSlot(meeting entities.Slot) AddSlot {
	sId, _ := strconv.Atoi(meeting.Id)
	return AddSlot{
		Id: sId,
	}
}

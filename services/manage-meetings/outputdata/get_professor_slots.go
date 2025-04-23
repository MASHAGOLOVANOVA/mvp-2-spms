package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
	"time"
)

type GetProfessorSlots struct {
	Slots []GetProfSLotsData `json:"slots"`
}

func MapToGetProfesorSlots(meetings []GetProfesorSlotsEntities) GetProfessorSlots {
	outputProjects := []GetProfSLotsData{}
	for _, meet := range meetings {
		parseInt, _ := strconv.Atoi(meet.Slot.Id)
		outputProjects = append(outputProjects,
			GetProfSLotsData{
				Id:          parseInt,
				Description: meet.Description,
				StartTime:   meet.StartTime,
				EndTime:     meet.EndTime,
			})
	}
	return GetProfessorSlots{
		Slots: outputProjects,
	}
}

type GetProfesorSlotsEntities struct {
	Slot        entities.Slot
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type GetProfSLotsData struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

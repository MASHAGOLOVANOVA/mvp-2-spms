package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
	"time"
)

type GetStudentSlots struct {
	Slots []GetStudSlotsData `json:"slots"`
}

func MapToGetStudentSlots(meetings []GetStudentSlotsEntities) GetStudentSlots {
	outputProjects := []GetStudSlotsData{}
	for _, meet := range meetings {
		parseInt, _ := strconv.Atoi(meet.Slot.Id)
		profId, _ := strconv.Atoi(meet.Slot.ProfessorId)
		studId, _ := strconv.Atoi(meet.StudMeeting.StudentId)
		outputProjects = append(outputProjects,
			GetStudSlotsData{
				Id:            parseInt,
				Description:   meet.Description,
				StartTime:     meet.StartTime,
				EndTime:       meet.EndTime,
				ProfessorId:   profId,
				StudentId:     studId,
				ProfessorName: meet.ProfessorName,
				StudentName:   meet.StudentName,
				IsOnline:      meet.Slot.IsOnline,
			})
	}
	return GetStudentSlots{
		Slots: outputProjects,
	}
}

type GetStudentSlotsEntities struct {
	StudMeeting   entities.StudMeeting
	Slot          entities.Slot
	Description   string
	StartTime     time.Time
	EndTime       time.Time
	ProfessorName string
	StudentName   string
	IsOnline      bool
}

type GetStudSlotsData struct {
	Id            int       `json:"id"`
	Description   string    `json:"description"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	ProfessorId   int       `json:"professor_id"`
	StudentId     int       `json:"student_id"`
	ProfessorName string    `json:"professor_name"`
	StudentName   string    `json:"student_name"`
	IsOnline      bool      `json:"is_online"`
}

package models

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type Apply struct {
	Id          uint  `gorm:"column:id"`
	StudentId   uint  `gorm:"column:student_id"`
	ProfessorId uint  `gorm:"column:professor_id"`
	Status      *bool `gorm:"column:status"`
}

func (*Apply) TableName() string {
	return "application"
}

// Метод для преобразования структуры Apply в сущность
func (a *Apply) MapToEntity() entities.Apply {
	var statusStr string
	if a.Status != nil {
		statusStr = strconv.FormatBool(*a.Status)
	} else {
		statusStr = "null"
	}

	return entities.Apply{
		Id:          strconv.Itoa(int(a.Id)),
		StudentId:   strconv.Itoa(int(a.StudentId)),
		ProfessorId: strconv.Itoa(int(a.ProfessorId)),
		Status:      statusStr,
	}
}

// Метод для преобразования сущности в структуру Apply
func (a *Apply) MapEntityToThis(entity entities.Apply) {
	studentId, err := strconv.ParseUint(entity.StudentId, 10, 32)
	if err != nil {
		return
	}
	a.StudentId = uint(studentId)

	professorId, err := strconv.ParseUint(entity.ProfessorId, 10, 32)
	if err != nil {
		return
	}
	a.ProfessorId = uint(professorId)

	if entity.Status == "null" {
		a.Status = nil // Устанавливаем статус в nil, если он "null"
	} else {
		status, err := strconv.ParseBool(entity.Status)
		if err != nil {
			return
		}
		a.Status = &status // Присваиваем указатель на статус
	}
}

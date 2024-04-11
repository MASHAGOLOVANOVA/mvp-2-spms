package models

import (
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type Professor struct {
	Id            uint   `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	Surname       string `gorm:"column:surname"`
	Middlename    string `gorm:"column:middlename"`
	ScienceDegree string `gorm:"column:science_degree"`
	UniversityId  uint   `gorm:"column:university_id"`
}

func (*Professor) TableName() string {
	return "professor"
}

func (p *Professor) MapToEntity() entities.Professor {
	return entities.Professor{
		Person: entities.Person{
			Id:         fmt.Sprint(p.Id),
			Name:       p.Name,
			Surname:    p.Surname,
			Middlename: p.Middlename,
		},
		ScienceDegree: p.ScienceDegree,
		UniversityId:  fmt.Sprint(p.UniversityId),
	}
}

func (p *Professor) MapEntityToThis(entity entities.Professor) {
	pId, _ := strconv.Atoi(entity.Id)
	uId, _ := strconv.Atoi(entity.UniversityId)
	p.Id = uint(pId)
	p.Name = entity.Name
	p.Surname = entity.Surname
	p.Middlename = entity.Middlename
	p.ScienceDegree = entity.ScienceDegree
	p.UniversityId = uint(uId)
}

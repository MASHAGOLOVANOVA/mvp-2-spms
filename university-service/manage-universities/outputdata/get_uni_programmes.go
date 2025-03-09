package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	"strconv"
)

type GetUniEducationalProgrammes struct {
	Programmes []getUniEdProgData `json:"programmes"`
}

func MapToGetUniEducationalProgramme(programmes []entities.EducationalProgramme) GetUniEducationalProgrammes {
	output := []getUniEdProgData{}
	for _, prog := range programmes {
		id, _ := strconv.Atoi(prog.Id)
		output = append(output,
			getUniEdProgData{
				Id:   id,
				Name: prog.Name,
			})
	}
	return GetUniEducationalProgrammes{
		Programmes: output,
	}
}

type getUniEdProgData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	// add ed level?
}

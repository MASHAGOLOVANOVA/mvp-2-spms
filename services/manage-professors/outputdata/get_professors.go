package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
)

type GetProfessors struct {
	Professors []getProfData `json:"professors"`
}

func MapToGetProfessors(profEntities []GetProfessorsEntities) GetProfessors {
	outputProfessors := make([]getProfData, len(profEntities))
	for i, profEntity := range profEntities {
		outputProfessors[i] = getProfData{
			Id:   profEntity.Professor.Id,
			Name: profEntity.Professor.FullNameToString(),
		}
	}
	return GetProfessors{
		Professors: outputProfessors,
	}
}

type GetProfessorsEntities struct {
	Professor entities.Professor
}

type getProfData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

package outputdata

import entities "mvp-2-spms/domain-aggregate"

type GetApplications struct {
	Applications []getApplicationData `json:"applications"`
}

type getApplicationData struct {
	Id            string `json:"id"`
	StudentId     string `json:"student_id"`
	ProfessorId   string `json:"professor_id"`
	StudentName   string `json:"student_name"`
	ProfessorName string `json:"professor_name"`
	Status        string `json:"status"`
}

func MapToGetApplications(applEntities []GetApplicationEntities) GetApplications {
	outputApplications := make([]getApplicationData, len(applEntities))
	for i, applEntity := range applEntities {
		outputApplications[i] = getApplicationData{
			Id:            applEntity.Application.Id,
			StudentName:   applEntity.Student.FullNameToString(),
			ProfessorName: applEntity.Professor.FullNameToString(),
			StudentId:     applEntity.Student.Id,
			ProfessorId:   applEntity.Professor.Id,
			Status:        applEntity.Application.Status,
		}
	}
	return GetApplications{
		Applications: outputApplications,
	}
}

type GetApplicationEntities struct {
	Application entities.Apply
	Professor   entities.Professor
	Student     entities.Student
}

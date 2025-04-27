package outputdata

import (
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"
	"strconv"
)

type GetProfessorInfo struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type GetStudentInfo struct {
	Id         int    `json:"id"`
	Login      string `json:"login"`
	Name       string `json:"name"`
	University string `json:"university"`
	EdProgName string `json:"ed_prog_name"`
}

func MapToGetStudentAccountInfo(stud entities.Student, studAcc entities.StudentAccount) GetStudentInfo {
	sId, _ := strconv.Atoi(stud.Id)
	return GetStudentInfo{
		Id:    sId,
		Login: "", ////////////////////////////
		Name:  stud.FullNameToString(),
	}
}
func MapModelToGetStudentAccountInfo(stud entities.Student, studAcc usecasemodels.StudentAccount) GetStudentInfo {
	sId, _ := strconv.Atoi(stud.Id)
	return GetStudentInfo{
		Id:    sId,
		Login: "", ////////////////////////////
		Name:  stud.FullNameToString(),
	}
}

func MapToGetAccountInfo(prof entities.Professor) GetProfessorInfo {
	pId, _ := strconv.Atoi(prof.Id)

	return GetProfessorInfo{
		Id:    pId,
		Login: "", ////////////////////////////
		Name:  prof.FullNameToString(),
	}
}

package managetasks

import (
	"fmt"
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-universities/inputdata"
	"mvp-2-spms/services/manage-universities/outputdata"
)

type UniversityInteractor struct {
	uniRepo interfaces.IUniversityRepository
}

func InitUniversityInteractor(uRepo interfaces.IUniversityRepository) *UniversityInteractor {
	return &UniversityInteractor{
		uniRepo: uRepo,
	}
}

func (p *UniversityInteractor) GetUniEdProgrammes(input inputdata.GetUniEducationalProgrammes) (outputdata.GetUniEducationalProgrammes, error) {
	// get progs from db
	progs, _ := p.uniRepo.GetUniversityEducationalProgrammes(fmt.Sprint(input.UniversityId))
	output := outputdata.MapToGetUniEducationalProgramme(progs)
	return output, nil
}

package manage_professors

import (
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-professors/inputdata"
	"mvp-2-spms/services/manage-professors/outputdata"
)

type ApplicationInteractor struct {
	professorRepo interfaces.IProfessorRepository
	accountRepo   interfaces.IAccountRepository
	studentRepo   interfaces.IStudentRepository
}

func InitApplicationInteractor(profRepo interfaces.IProfessorRepository,
	accRepo interfaces.IAccountRepository, studRepo interfaces.IStudentRepository) *ApplicationInteractor {
	return &ApplicationInteractor{
		professorRepo: profRepo,
		accountRepo:   accRepo,
		studentRepo:   studRepo,
	}
}

func (p *ApplicationInteractor) GetStudentApplications(studentId string) (outputdata.GetApplications, error) {
	applEntities, err := p.professorRepo.GetApplicationsByStudent(studentId)
	if err != nil {
		return outputdata.GetApplications{}, err // Возвращаем ошибку, если она возникла
	}

	// Преобразуем полученные сущности профессоров в нужный формат
	var outputEntities []outputdata.GetApplicationEntities
	for _, appl := range applEntities {
		resChan := <-p.accountRepo.GetProfessorById(appl.ProfessorId)
		res, _ := p.studentRepo.GetStudentById(appl.StudentId)
		outputEntities = append(outputEntities, outputdata.GetApplicationEntities{
			Student:     res,
			Professor:   resChan.Professor,
			Application: appl,
		})
	}

	output := outputdata.MapToGetApplications(outputEntities)
	return output, nil
}

func (p *ApplicationInteractor) GetProfessorApplications(profId string) (outputdata.GetApplications, error) {
	applEntities, err := p.professorRepo.GetApplicationsByProfessor(profId)
	if err != nil {
		return outputdata.GetApplications{}, err // Возвращаем ошибку, если она возникла
	}

	// Преобразуем полученные сущности профессоров в нужный формат
	var outputEntities []outputdata.GetApplicationEntities
	for _, appl := range applEntities {
		resChan := <-p.accountRepo.GetProfessorById(appl.ProfessorId)
		res, _ := p.studentRepo.GetStudentById(appl.StudentId)
		outputEntities = append(outputEntities, outputdata.GetApplicationEntities{
			Student:     res,
			Professor:   resChan.Professor,
			Application: appl,
		})
	}

	output := outputdata.MapToGetApplications(outputEntities)
	return output, nil
}

func (p *ApplicationInteractor) GetProfessors() (outputdata.GetProfessors, error) {
	// Получаем профессоров из репозитория
	profEntities, err := p.professorRepo.GetProfessors()
	if err != nil {
		return outputdata.GetProfessors{}, err // Возвращаем ошибку, если она возникла
	}

	// Преобразуем полученные сущности профессоров в нужный формат
	var outputEntities []outputdata.GetProfessorsEntities
	for _, prof := range profEntities {
		outputEntities = append(outputEntities, outputdata.GetProfessorsEntities{
			Professor: prof,
		})
	}

	output := outputdata.MapToGetProfessors(outputEntities)
	return output, nil
}

func (p *ApplicationInteractor) Apply(apply inputdata.Apply) error {
	return p.professorRepo.Apply(apply.MapToApplyEntity())
}

func (p *ApplicationInteractor) UpdateApplicationStatus(apply inputdata.ApplicationStatus) error {
	return p.professorRepo.UpdateApplication(apply.MapToApplyEntity())
}

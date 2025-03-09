package managestudents

import (
	"errors"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-students/inputdata"
	"mvp-2-spms/services/manage-students/outputdata"
	"mvp-2-spms/services/models"
)

type StudentInteractor struct {
	studentRepo interfaces.IStudentRepository
	projetRepo  interfaces.IProjetRepository
	uniRepo     interfaces.IUniversityRepository
}

func InitStudentInteractor(stRepo interfaces.IStudentRepository, pjRepo interfaces.IProjetRepository, uRepo interfaces.IUniversityRepository) *StudentInteractor {
	return &StudentInteractor{
		studentRepo: stRepo,
		projetRepo:  pjRepo,
		uniRepo:     uRepo,
	}
}

func (p *StudentInteractor) AddStudent(input inputdata.AddStudent) (outputdata.AddStudent, error) {
	// adding student to db, returns created student (with id)
	student, err := p.studentRepo.CreateStudent(input.MapToStudentEntity())
	if err != nil {
		return outputdata.AddStudent{}, err
	}

	// returning id
	output := outputdata.MapToAddStudent(student)
	return output, nil
}

func (p *StudentInteractor) GetStudents(input inputdata.GetStudents) (outputdata.GetStudents, error) {
	// get from database
	stEntities := []outputdata.GetStudentsEntities{}

	students, err := p.studentRepo.GetStudents()
	if err != nil {
		return outputdata.GetStudents{}, err
	}

	for _, student := range students {
		project, err := p.projetRepo.GetStudentCurrentProject(student.Id)
		if err != nil {
			if !errors.Is(err, models.ErrStudentHasNoCurrentProject) {
				return outputdata.GetStudents{}, err
			}
			project = domainaggregate.Project{} // поменять на нил
		}

		edProg, err := p.uniRepo.GetEducationalProgrammeById(student.EducationalProgrammeId)
		if err != nil {
			return outputdata.GetStudents{}, err
		}

		stEntities = append(stEntities, outputdata.GetStudentsEntities{
			ProjectTheme:         project.Theme,
			Student:              student,
			EducationalProgramme: edProg.Name,
			PtojectId:            project.Id,
		})
	}

	output := outputdata.MapToGetStudents(stEntities)
	return output, nil
}

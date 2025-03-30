package universityrepository

import (
	"errors"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"

	"gorm.io/gorm"
)

type UniversityRepository struct {
	dbContext database.Database
}

func InitUniversityRepository(dbcxt database.Database) *UniversityRepository {
	return &UniversityRepository{
		dbContext: dbcxt,
	}
}

func (u *UniversityRepository) GetEducationalProgrammeById(epId string) (entities.EducationalProgramme, error) {
	var edProg models.EducationalProgramme

	result := u.dbContext.DB.Select("*").Where("id = ?", epId).Take(&edProg)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.EducationalProgramme{}, usecasemodels.ErrEdProgrammmeNotFound
		}
		return entities.EducationalProgramme{}, result.Error
	}

	return edProg.MapToEntity(), nil
}

func (u *UniversityRepository) GetEducationalProgrammeFullById(epId string) (usecasemodels.EdProg, error) {
	var edProg models.EducationalProgramme

	result := u.dbContext.DB.Select("*").Where("id = ?", epId).Take(&edProg)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.EdProg{}, usecasemodels.ErrEdProgrammmeNotFound
		}
		return usecasemodels.EdProg{}, result.Error
	}

	var f models.Faculty
	result = u.dbContext.DB.Select("*").Where("id = ?", edProg.FacultyId).Take(&f)
	if result.Error != nil {
		return usecasemodels.EdProg{}, result.Error
	}

	var d models.Department
	result = u.dbContext.DB.Select("*").Where("id = ?", f.DeptId).Take(&d)
	if result.Error != nil {
		return usecasemodels.EdProg{}, result.Error
	}

	return usecasemodels.EdProg{
		EducationalProgramme: edProg.MapToEntity(),
		Dept:                 d.Name,
		Faculty:              f.Name,
	}, nil
}

func (u *UniversityRepository) GetUniversityById(uId string) (entities.University, error) {
	var uni models.University

	result := u.dbContext.DB.Select("*").Where("id = ?", uId).Take(&uni)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.University{}, usecasemodels.ErrUniNoFound
		}
		return entities.University{}, result.Error
	}

	return uni.MapToEntity(), nil
}

func (u *UniversityRepository) GetUniversityEducationalProgrammes(uniId string) ([]entities.EducationalProgramme, error) {
	var edProgs []models.EducationalProgramme

	result := u.dbContext.DB.Raw(
		`SELECT educational_programme.* 
		FROM (SELECT * 
		FROM department
		WHERE uni_id = ?) as depts
		LEFT JOIN faculty ON faculty.dept_id = depts.id
		LEFT JOIN educational_programme ON educational_programme.faculty_id = faculty.id;`,
		uniId).Scan(&edProgs)

	if result.Error != nil {
		return []entities.EducationalProgramme{}, result.Error
	}

	edProgrammes := []entities.EducationalProgramme{}
	for _, p := range edProgs {
		edProgrammes = append(edProgrammes, p.MapToEntity())
	}
	return edProgrammes, nil
}

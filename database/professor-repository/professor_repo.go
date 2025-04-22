package professor_repository

import (
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
)

type ProfessorRepository struct {
	dbContext database.Database
}

func InitProfessorRepository(dbcxt database.Database) *ProfessorRepository {
	return &ProfessorRepository{
		dbContext: dbcxt,
	}
}

func (r *ProfessorRepository) GetApplicationsByStudent(studentId string) ([]entities.Apply, error) {
	var applyDb []models.Apply
	result := r.dbContext.DB.Select("*").Where("student_id = ?", studentId).Find(&applyDb)
	if result.Error != nil {
		return nil, result.Error
	}

	// Маппим профессоров из базы данных в сущности
	applications := make([]entities.Apply, len(applyDb))
	for i, ap := range applyDb {
		applications[i] = ap.MapToEntity()
	}
	return applications, nil
}

func (r *ProfessorRepository) GetApplicationsByProfessor(profId string) ([]entities.Apply, error) {
	var applyDb []models.Apply
	result := r.dbContext.DB.Select("*").Where("professor_id = ?", profId).Find(&applyDb)
	if result.Error != nil {
		return nil, result.Error
	}

	// Маппим профессоров из базы данных в сущности
	applications := make([]entities.Apply, len(applyDb))
	for i, ap := range applyDb {
		applications[i] = ap.MapToEntity()
	}
	return applications, nil
}

// GetProfessors получает список профессоров из базы данных.
func (r *ProfessorRepository) GetProfessors() ([]entities.Professor, error) {
	var professorsDb []models.Professor
	result := r.dbContext.DB.Select("*").Find(&professorsDb)
	if result.Error != nil {
		return nil, result.Error
	}

	// Маппим профессоров из базы данных в сущности
	professors := make([]entities.Professor, len(professorsDb))
	for i, pj := range professorsDb {
		professors[i] = pj.MapToEntity()
	}
	return professors, nil
}

func (r *ProfessorRepository) Apply(apply entities.Apply) error {
	dbapply := models.Apply{}
	dbapply.MapEntityToThis(apply)
	dbapply.Status = nil
	result := r.dbContext.DB.Create(&dbapply)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ProfessorRepository) UpdateApplication(apply entities.Apply) error {
	applDb := models.Apply{}
	result := r.dbContext.DB.Where("id = ?", apply.Id).Take(&applDb)
	if result.Error != nil {
		return result.Error
	}
	applDb.MapEntityToThis(apply)

	result = r.dbContext.DB.Where("id = ?", apply.Id).Save(&applDb)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

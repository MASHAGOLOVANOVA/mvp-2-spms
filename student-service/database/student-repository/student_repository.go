package studentrepository

import (
	"errors"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"

	"gorm.io/gorm"
)

type StudentRepository struct {
	dbContext database.Database
}

func InitStudentRepository(dbcxt database.Database) *StudentRepository {
	return &StudentRepository{
		dbContext: dbcxt,
	}
}

func (r *StudentRepository) GetStudentById(studId string) (entities.Student, error) {
	var student models.Student

	result := r.dbContext.DB.Select("*").Where("id = ?", studId).Take(&student)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Student{}, usecasemodels.ErrStudentNotFound
		}
		return entities.Student{}, result.Error
	}

	return student.MapToEntity(), nil
}

func (r *StudentRepository) CreateStudent(student entities.Student) (entities.Student, error) {
	dbstudent := models.Student{}
	dbstudent.MapEntityToThis(student)

	result := r.dbContext.DB.Create(&dbstudent)
	if result.Error != nil {
		return entities.Student{}, result.Error
	}

	return dbstudent.MapToEntity(), nil
}

func (r *StudentRepository) DeleteStudent(id int) error {
	dbStudent := models.Student{Id: uint(id)}

	result := r.dbContext.DB.Delete(&dbStudent)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *StudentRepository) GetStudents() ([]entities.Student, error) {
	var studentsDb []models.Student

	result := r.dbContext.DB.Select("*").Find(&studentsDb)
	if result.Error != nil {
		return []entities.Student{}, result.Error
	}

	students := []entities.Student{}
	for _, s := range studentsDb {
		students = append(students, s.MapToEntity())
	}
	return students, nil
}

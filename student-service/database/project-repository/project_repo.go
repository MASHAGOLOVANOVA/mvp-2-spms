package projectrepository

import (
	"database/sql"
	"errors"
	"fmt"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
	usecasemodels "mvp-2-spms/services/models"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	dbContext database.Database
}

func InitProjectRepository(dbcxt database.Database) *ProjectRepository {
	return &ProjectRepository{
		dbContext: dbcxt,
	}
}

func (r *ProjectRepository) GetProfessorProjectsWithFilters(profId string, statusFilter int) ([]entities.Project, error) {
	var projectsDb []models.Project
	result := r.dbContext.DB.Select("*").Where("supervisor_id = ? and status_id = ?", profId, statusFilter).Find(&projectsDb)
	if result.Error != nil {
		return []entities.Project{}, result.Error
	}
	projects := []entities.Project{}
	for _, pj := range projectsDb {
		// вынести в маппер
		projects = append(projects, pj.MapToEntity())
	}
	return projects, nil
}

func (r *ProjectRepository) GetProfessorProjects(profId string) ([]entities.Project, error) {
	var projectsDb []models.Project
	result := r.dbContext.DB.Select("*").Where("supervisor_id = ?", profId).Find(&projectsDb)
	if result.Error != nil {
		return []entities.Project{}, result.Error
	}
	projects := []entities.Project{}
	for _, pj := range projectsDb {
		// вынести в маппер
		projects = append(projects, pj.MapToEntity())
	}
	return projects, nil
}

func (r *ProjectRepository) GetProjectRepository(projId string) (usecasemodels.Repository, error) {
	var project models.Project
	result := r.dbContext.DB.Select("repo_id").Where("id = ?", projId).Take(&project)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.Repository{}, usecasemodels.ErrProjectNotFound
		}
		return usecasemodels.Repository{}, result.Error
	}

	if !project.RepoId.Valid {
		return usecasemodels.Repository{}, usecasemodels.ErrProjectRepoNotFound
	}

	var repo models.Repository
	result = r.dbContext.DB.Select("*").Where("id = ?", project.RepoId).Take(&repo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return usecasemodels.Repository{}, usecasemodels.ErrProjectRepoNotFound
		}
		return usecasemodels.Repository{}, result.Error
	}

	return repo.MapToUseCaseModel(), nil
}

func (r *ProjectRepository) GetProjectById(projId string) (entities.Project, error) {
	var project models.Project
	result := r.dbContext.DB.Select("*").Where("id = ?", projId).Take(&project)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Project{}, usecasemodels.ErrProjectNotFound
		}
		return entities.Project{}, result.Error
	}
	return project.MapToEntity(), nil
}

func (r *ProjectRepository) CreateProject(project entities.Project) (entities.Project, error) {
	dbProject := models.Project{}
	dbProject.MapEntityToThis(project)

	result := r.dbContext.DB.Create(&dbProject)
	if result.Error != nil {
		return entities.Project{}, result.Error
	}
	return dbProject.MapToEntity(), nil
}

func (r *ProjectRepository) CreateProjectWithRepository(project entities.Project, repo usecasemodels.Repository) (usecasemodels.ProjectInRepository, error) {
	dbRepo := models.Repository{}
	dbRepo.MapModelToThis(repo)

	dbProject := models.Project{}
	dbProject.MapEntityToThis(project)

	err := r.dbContext.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&dbRepo)
		if result.Error != nil {
			return result.Error
		}

		dbProject.RepoId = sql.NullInt64{Int64: int64(dbRepo.Id), Valid: true}

		result = tx.Create(&dbProject)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return usecasemodels.ProjectInRepository{}, err
	}
	return usecasemodels.ProjectInRepository{
		Project: dbProject.MapToEntity(),
	}, nil
}

func (r *ProjectRepository) AssignDriveFolder(project usecasemodels.DriveProject) error {
	dbCloudFolder := models.CloudFolder{}
	dbCloudFolder.MapUseCaseModelToThis(project.DriveFolder)

	err := r.dbContext.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&dbCloudFolder)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&models.Project{}).Where("id = ?", project.Project.Id).Update("cloud_id", project.DriveFolder.Id)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return usecasemodels.ErrProjectNotFound
		}
		return nil
	})

	return err
}

func (r *ProjectRepository) GetProjectCloudFolderId(projId string) (string, error) {
	proj := models.Project{}
	result := r.dbContext.DB.Select("cloud_id").Where("id = ?", projId).Take(&proj)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", usecasemodels.ErrProjectNotFound
		}
		return "", result.Error
	}

	if !proj.CloudId.Valid {
		return "", usecasemodels.ErrProjectCloudFolderNotFound
	}
	return fmt.Sprint(proj.CloudId.String), nil
}

func (r *ProjectRepository) GetProjectFolderLink(projId string) (string, error) {
	folder := models.CloudFolder{}
	folderid, err := r.GetProjectCloudFolderId(projId)
	if err != nil {
		return "", err
	}

	result := r.dbContext.DB.Select("link").Where("id = ?", folderid).Take(&folder)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", usecasemodels.ErrProjectCloudFolderLinkNotFound
		}
		return "", result.Error
	}
	return folder.Link, nil
}

func (r *ProjectRepository) GetStudentCurrentProject(studId string) (entities.Project, error) {
	proj := models.Project{}
	result := r.dbContext.DB.Select("*").Where("student_id = ? AND status_id IN(?, ?)",
		studId, entities.ProjectInProgress,
		entities.ProjectNotConfirmed).Order("year desc").Take(&proj)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.Project{}, usecasemodels.ErrStudentHasNoCurrentProject
		}
		return entities.Project{}, result.Error
	}
	return proj.MapToEntity(), nil
}

func (r *ProjectRepository) GetProjectGradingById(projId string) (entities.ProjectGrading, error) {
	var defenceGrade sql.NullFloat64

	result := r.dbContext.DB.Model(models.Project{}).Select("defence_grade").Where("id=?", projId).Take(&defenceGrade)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.ProjectGrading{}, usecasemodels.ErrProjectNotFound
		}
		return entities.ProjectGrading{}, result.Error
	}

	var supReview models.SupervisorReview
	result = r.dbContext.DB.Model(models.Project{}).Select("supervisor_review_id").Where("id=?", projId).Take(&supReview.Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entities.ProjectGrading{}, usecasemodels.ErrProjectNotFound
		}
		return entities.ProjectGrading{}, result.Error
	}

	grading := entities.ProjectGrading{
		ProjectId: projId,
	}
	if defenceGrade.Valid {
		grading.DefenceGrade = float32(defenceGrade.Float64)
	}
	if supReview.Id.Valid {
		result = r.dbContext.DB.Select("*").Where("id=?", supReview.Id).Take(&supReview)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return entities.ProjectGrading{}, usecasemodels.ErrSupervisorReviewNotFound
			}
			return entities.ProjectGrading{}, result.Error
		}

		var dbcriterias []models.ReviewCriteria
		r.dbContext.DB.Select("*").Where("supervisor_review_id=?", supReview.Id).Find(&dbcriterias)

		criterias := []entities.Criteria{}
		for _, c := range dbcriterias {
			criterias = append(criterias, c.MapToEntity())
		}
		grading.SupervisorReview = supReview.MapToEntity(criterias)
	}
	return grading, nil
}

func (r *ProjectRepository) GetProjectTaskInfoById(projId string) (usecasemodels.TasksInfo, error) {
	taskInfo := models.ProjectTaskInfo{}

	result := r.dbContext.DB.Raw(` 
	SELECT status, COUNT(status) as count
	FROM task
	WHERE project_id = ?
	GROUP BY status`, projId).Scan(&taskInfo.Statuses)
	if result.Error != nil {
		return usecasemodels.TasksInfo{}, result.Error
	}

	return taskInfo.MapToUseCaseModel(), nil
}

func (r *ProjectRepository) GetProjectMeetingInfoById(projId string) (usecasemodels.MeetingInfo, error) {
	var meetCount int

	result := r.dbContext.DB.Raw(`
	SELECT COUNT(id) as count
	FROM meeting
	WHERE project_id = ? and status = 1`, projId).Scan(&meetCount)
	if result.Error != nil {
		return usecasemodels.MeetingInfo{}, result.Error
	}

	return usecasemodels.MeetingInfo{
		PassedCount: meetCount,
	}, nil
}

func (r *ProjectRepository) UpdateProject(proj entities.Project) error {
	projDb := models.Project{}
	result := r.dbContext.DB.Where("id = ?", proj.Id).Find(&projDb)
	if result.Error != nil {
		return result.Error
	}
	projDb.MapEntityToThis(proj)

	result = r.dbContext.DB.Where("id = ?", proj.Id).Save(&projDb)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ProjectRepository) UpdateProjectDefenceGrade(projId string, grade float32) error {
	result := r.dbContext.DB.Model(&models.Project{}).Where("id = ?", projId).Update("defence_grade", grade)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ProjectRepository) UpdateProjectSupRew(projId string, sr entities.SupervisorReview) error {
	srDb := models.SupervisorReview{}
	srDb.MapEntityToThis(sr)

	err := r.dbContext.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Save(&srDb)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Where("supervisor_review_id = ?", srDb.Id.Int64).Delete(&models.ReviewCriteria{})
		if result.Error != nil {
			return result.Error
		}

		for _, c := range sr.Criterias {
			cDb := models.ReviewCriteria{}
			cDb.MapEntityToThis(c)
			cDb.SupervieorReviewId = uint(srDb.Id.Int64)

			result = tx.Create(&cDb)
			if result.Error != nil {
				return result.Error
			}
		}

		result = tx.Model(&models.Project{}).Where("id = ?", projId).Update("supervisor_review_id", srDb.Id.Int64)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

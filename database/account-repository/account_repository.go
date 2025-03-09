package accountrepository

import (
	"errors"
	"mvp-2-spms/database"
	"mvp-2-spms/database/models"
	entities "mvp-2-spms/domain-aggregate"
	interfaces "mvp-2-spms/services/interfaces"
	usecasemodels "mvp-2-spms/services/models"
	"strconv"

	"gorm.io/gorm"
)

type AccountRepository struct {
	dbContext database.Database
}

func InitAccountRepository(dbcxt database.Database) *AccountRepository {
	return &AccountRepository{
		dbContext: dbcxt,
	}
}

func (r *AccountRepository) GetAccountByLogin(login string) <-chan interfaces.ResultAccount {
	resultChan := make(chan interfaces.ResultAccount)

	go func() {
		defer close(resultChan)
		acc := models.Account{}

		result := r.dbContext.DB.Select("*").Where("login = ?", login).Take(&acc)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				resultChan <- interfaces.ResultAccount{Account: usecasemodels.Account{}, Err: usecasemodels.ErrAccountNotFound}
				return
			}
			resultChan <- interfaces.ResultAccount{Account: usecasemodels.Account{}, Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultAccount{Account: acc.MapToUseCaseModel(), Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) DeleteAccountByLogin(login string) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		accResult := <-r.GetAccountByLogin(login)
		if accResult.Err != nil {
			resultChan <- interfaces.ResultError{Err: accResult.Err}
			return
		}

		profId, err := strconv.Atoi(accResult.Account.Id)
		if err != nil {
			resultChan <- interfaces.ResultError{Err: err}
			return
		}

		result := r.dbContext.DB.Delete(&models.Professor{Id: uint(profId)})
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				resultChan <- interfaces.ResultError{Err: usecasemodels.ErrProfessorNotFound}
				return
			}
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}
		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) AddProfessor(prof entities.Professor) <-chan interfaces.ResultProfessor {
	resultChan := make(chan interfaces.ResultProfessor)

	go func() {
		defer close(resultChan)

		dbProf := models.Professor{}
		dbProf.MapEntityToThis(prof)

		result := r.dbContext.DB.Create(&dbProf)
		if result.Error != nil {
			resultChan <- interfaces.ResultProfessor{Professor: entities.Professor{}, Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultProfessor{Professor: dbProf.MapToEntity(), Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) DeleteProfessor(profId int) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		dbProf := models.Professor{Id: uint(profId)}
		result := r.dbContext.DB.Delete(&dbProf)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) AddAccount(account usecasemodels.Account) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		dbAcc := models.Account{}
		dbAcc.MapUseCaseModelToThis(account)

		result := r.dbContext.DB.Create(&dbAcc)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) GetProfessorById(id string) <-chan interfaces.ResultProfessor {
	resultChan := make(chan interfaces.ResultProfessor)

	go func() {
		defer close(resultChan)

		prof := models.Professor{}
		result := r.dbContext.DB.Select("*").Where("id = ?", id).Take(&prof)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				resultChan <- interfaces.ResultProfessor{Professor: entities.Professor{}, Err: usecasemodels.ErrProfessorNotFound}
				return
			}
			resultChan <- interfaces.ResultProfessor{Professor: entities.Professor{}, Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultProfessor{Professor: prof.MapToEntity(), Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) GetAccountPlannerData(id string) <-chan interfaces.ResultPlannerIntegration {
	resultChan := make(chan interfaces.ResultPlannerIntegration)

	go func() {
		defer close(resultChan)

		dbPlanner := models.PlannerIntegration{}
		result := r.dbContext.DB.Select("*").Where("account_id = ?", id).Take(&dbPlanner)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				resultChan <- interfaces.ResultPlannerIntegration{PlannerIntegration: usecasemodels.PlannerIntegration{}, Err: usecasemodels.ErrAccountPlannerDataNotFound}
				return
			}
			resultChan <- interfaces.ResultPlannerIntegration{PlannerIntegration: usecasemodels.PlannerIntegration{}, Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultPlannerIntegration{PlannerIntegration: dbPlanner.MapToUseCaseModel(), Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) GetAccountDriveData(id string) <-chan interfaces.ResultCloudDriveIntegration {
	resultChan := make(chan interfaces.ResultCloudDriveIntegration)

	go func() {
		defer close(resultChan)

		dbDrive := models.DriveIntegration{}
		result := r.dbContext.DB.Select("*").Where("account_id = ?", id).Take(&dbDrive)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				resultChan <- interfaces.ResultCloudDriveIntegration{CloudDriveIntegration: usecasemodels.CloudDriveIntegration{}, Err: usecasemodels.ErrAccountDriveDataNotFound}
				return
			}
			resultChan <- interfaces.ResultCloudDriveIntegration{CloudDriveIntegration: usecasemodels.CloudDriveIntegration{}, Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultCloudDriveIntegration{CloudDriveIntegration: dbDrive.MapToUseCaseModel(), Err: nil}
	}()

	return resultChan
}

// can return multiple for 1 account, should consider this
func (r *AccountRepository) GetAccountRepoHubData(id string) <-chan interfaces.ResultBaseIntegration {
	resultChan := make(chan interfaces.ResultBaseIntegration)

	go func() {
		defer close(resultChan)

		dbRHub := models.GitRepositoryIntegration{}
		result := r.dbContext.DB.Select("*").Where("account_id = ?", id).Take(&dbRHub)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				resultChan <- interfaces.ResultBaseIntegration{BaseIntegration: usecasemodels.BaseIntegration{}, Err: usecasemodels.ErrAccountRepoHubDataNotFound}
				return
			}
			resultChan <- interfaces.ResultBaseIntegration{BaseIntegration: usecasemodels.BaseIntegration{}, Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultBaseIntegration{BaseIntegration: dbRHub.MapToUseCaseModel(), Err: nil}
	}()

	return resultChan
}

// func (r *AccountRepository) DeleteAccountPlannerData(id int) error {
// 	dbPl := models.PlannerIntegration{AccountId: uint(id)}

// 	result := r.dbContext.DB.Where("account_id = ?", id).Delete(&dbPl)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (r *AccountRepository) DeleteAccountDriveData(id int) error {
// 	dbDrive := models.DriveIntegration{AccountId: uint(id)}

// 	result := r.dbContext.DB.Where("account_id = ?", id).Delete(&dbDrive)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (r *AccountRepository) DeleteAccountRepoHubData(id int) error {
// 	dbRepo := models.GitRepositoryIntegration{AccountId: uint(id)}

// 	result := r.dbContext.DB.Where("account_id = ?", id).Delete(&dbRepo)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

func (r *AccountRepository) AddAccountPlannerIntegration(integr usecasemodels.PlannerIntegration) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		dbPlanner := models.PlannerIntegration{}
		dbPlanner.MapUseCaseModelToThis(integr)

		result := r.dbContext.DB.Create(&dbPlanner)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) AddAccountDriveIntegration(integr usecasemodels.CloudDriveIntegration) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		dbDrive := models.DriveIntegration{}
		dbDrive.MapUseCaseModelToThis(integr)

		result := r.dbContext.DB.Create(&dbDrive)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) AddAccountRepoHubIntegration(integr usecasemodels.BaseIntegration) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		dbRepoHub := models.GitRepositoryIntegration{}
		dbRepoHub.MapUseCaseModelToThis(integr)

		result := r.dbContext.DB.Create(&dbRepoHub)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) UpdateAccountPlannerIntegration(integr usecasemodels.PlannerIntegration) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		plannerDb := models.PlannerIntegration{}
		plannerDb.MapUseCaseModelToThis(integr)

		result := r.dbContext.DB.Where("account_id = ?", integr.AccountId).Save(&plannerDb)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) UpdateAccountDriveIntegration(integr usecasemodels.CloudDriveIntegration) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		result := r.dbContext.DB.Model(&models.DriveIntegration{}).Where("account_id = ?", integr.AccountId).Update("api_key", integr.ApiKey)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		if result.RowsAffected == 0 {
			resultChan <- interfaces.ResultError{Err: usecasemodels.ErrAccountDriveDataNotFound}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

func (r *AccountRepository) UpdateAccountRepoHubIntegration(integr usecasemodels.BaseIntegration) <-chan interfaces.ResultError {
	resultChan := make(chan interfaces.ResultError)

	go func() {
		defer close(resultChan)

		result := r.dbContext.DB.Model(&models.GitRepositoryIntegration{}).Where("account_id = ?", integr.AccountId).Update("api_key", integr.ApiKey)
		if result.Error != nil {
			resultChan <- interfaces.ResultError{Err: result.Error}
			return
		}

		if result.RowsAffected == 0 {
			resultChan <- interfaces.ResultError{Err: usecasemodels.ErrAccountDriveDataNotFound}
			return
		}

		resultChan <- interfaces.ResultError{Err: nil}
	}()

	return resultChan
}

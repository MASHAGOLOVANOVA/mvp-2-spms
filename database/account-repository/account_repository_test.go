package accountrepository

import (
	"mvp-2-spms/database"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/models"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dsn = "root:root@tcp(0.0.0.0:3308)/student_project_management_test?parseTime=true"

func connectDB() *database.Database {
	gdb, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	db := database.InitDatabade(gdb)
	return db
}

func TestAccountRepo_GetAccountByLogin(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent account", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		resChan := ar.GetAccountByLogin("123")
		resAcc := <-resChan
		// assert
		assert.ErrorIs(t, resAcc.Err, models.ErrAccountNotFound)
	})

	t.Run("ok, get existing account", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		login := "test"
		err := addTestingAccount(login, ar)
		assert.NoError(t, err)

		// act
		resChan := ar.GetAccountByLogin(login)
		resAcc := <-resChan

		// assert
		assert.NoError(t, resAcc.Err)

		// cleanup
		err = deleteTestingAccount(login, ar)
		assert.NoError(t, err)
	})
}
func TestAccountRepo_AddProfessor(t *testing.T) {
	db := connectDB()

	t.Run("ok", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		prof := domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "Test",
				Surname:    "1",
				Middlename: "2",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		}

		// act
		resChan := ar.AddProfessor(prof)
		resProf := <-resChan
		// assert
		assert.NoError(t, resProf.Err)
		resChan1 := ar.GetProfessorById(resProf.Professor.Id)
		resProf1 := <-resChan1

		assert.NoError(t, resProf1.Err)
		assert.Equal(t, resProf.Professor.Id, resProf1.Professor.Id)
		assert.Equal(t, resProf.Professor.Name, resProf1.Professor.Name)
		assert.Equal(t, resProf.Professor.Surname, resProf1.Professor.Surname)
		assert.Equal(t, resProf.Professor.Middlename, resProf1.Professor.Middlename)
		assert.Equal(t, resProf.Professor.ScienceDegree, resProf1.Professor.ScienceDegree)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)
	})

	t.Run("fail, uni id doesnt exist", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		prof := domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "Test",
				Surname:    "1",
				Middlename: "2",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
			UniversityId:  "2",
		}

		// act
		resChan := ar.AddProfessor(prof)
		resProf := <-resChan
		// assert
		assert.Error(t, resProf.Err)
	})
}

func TestAccountRepo_AddAccount(t *testing.T) {
	db := connectDB()

	t.Run("ok", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		prof := domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "Test",
				Surname:    "1",
				Middlename: "2",
			},
			ScienceDegree: "sd",
		}
		resChan := ar.AddProfessor(prof)
		resProf := <-resChan
		assert.NoError(t, resProf.Err)

		login := time.Now().Format(time.RFC3339)
		acc := models.Account{
			Login: login,
			Hash:  []byte{5, 6, 2},
			Salt:  "123232434",
			Id:    resProf.Professor.Id,
		}

		// act
		resErrChan := ar.AddAccount(acc)
		resErr := <-resErrChan
		// assert
		assert.NoError(t, resErr.Err)

		resAccChan1 := ar.GetAccountByLogin(login)
		resAcc := <-resAccChan1
		assert.NoError(t, resAcc.Err)
		assert.Equal(t, acc.Id, resAcc.Account.Id)
		assert.Equal(t, acc.Login, resAcc.Account.Login)
		assert.Equal(t, acc.Hash, resAcc.Account.Hash)
		assert.Equal(t, acc.Salt, resAcc.Account.Salt)

		// cleanup
		err := deleteTestingAccount(login, ar)
		assert.NoError(t, err)
	})

	t.Run("fail, prof id doesnt exist", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)
		login := time.Now().Format(time.RFC3339)
		acc := models.Account{
			Login: login,
			Hash:  []byte{5, 6, 2},
			Salt:  "123232434",
			Id:    "0",
		}

		// act
		errChan := ar.AddAccount(acc)
		err := <-errChan
		// assert
		assert.Error(t, err.Err)
	})
}

func TestAccountRepo_GetProfessorById(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent prof", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		resProfChan := ar.GetProfessorById("123")
		resProf := <-resProfChan

		// assert
		assert.ErrorIs(t, resProf.Err, models.ErrProfessorNotFound)
	})

	t.Run("ok, get existing prof", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "",
				Surname:    "",
				Middlename: "",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		// act
		resProfChan1 := ar.GetProfessorById(resProf.Professor.Id)
		resProf1 := <-resProfChan1

		// assert
		assert.NoError(t, resProf1.Err)
		assert.Equal(t, resProf.Professor.Id, resProf1.Professor.Id)
		assert.Equal(t, resProf.Professor.Name, resProf1.Professor.Name)
		assert.Equal(t, resProf.Professor.Surname, resProf1.Professor.Surname)
		assert.Equal(t, resProf.Professor.Middlename, resProf1.Professor.Middlename)
		assert.Equal(t, resProf.Professor.ScienceDegree, resProf1.Professor.ScienceDegree)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)
	})
}

func TestAccountRepo_GetAccountPlannerData(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent planner data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		resAccChan := ar.GetAccountPlannerData("123")
		resAcc := <-resAccChan

		// assert
		assert.ErrorIs(t, resAcc.Err, models.ErrAccountPlannerDataNotFound)
	})

	t.Run("ok, get existing planner data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		planner := models.PlannerIntegration{
			BaseIntegration: models.BaseIntegration{
				AccountId: resProf.Professor.Id,
				ApiKey:    "api",
				Type:      int(models.GoogleCalendar),
			},
			PlannerData: models.PlannerData{
				Id: time.Now().Format(time.RFC3339),
			},
		}

		resPlannerChan := ar.AddAccountPlannerIntegration(planner)
		resPlanner := <-resPlannerChan
		assert.NoError(t, resPlanner.Err)

		// act
		resPlannerChan1 := ar.GetAccountPlannerData(planner.AccountId)
		resPlanner1 := <-resPlannerChan1

		// assert
		assert.NoError(t, resPlanner1.Err)
		assert.Equal(t, planner.AccountId, resPlanner1.PlannerIntegration.AccountId)
		assert.Equal(t, planner.ApiKey, resPlanner1.PlannerIntegration.ApiKey)
		assert.Equal(t, planner.Id, resPlanner1.PlannerIntegration.Id)
		assert.Equal(t, planner.Type, resPlanner1.PlannerIntegration.Type)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan := ar.DeleteProfessor(profId)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)
	})
}

func TestAccountRepo_GetAccountDriveData(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent drive data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		resDriveChan := ar.GetAccountDriveData("123")
		resDrive := <-resDriveChan
		// assert
		assert.ErrorIs(t, resDrive.Err, models.ErrAccountDriveDataNotFound)
	})

	t.Run("ok, get existing drive data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		drive := models.CloudDriveIntegration{
			BaseIntegration: models.BaseIntegration{
				AccountId: resProf.Professor.Id,
				ApiKey:    "api",
				Type:      int(models.GoogleDrive),
			},
			DriveData: models.DriveData{
				BaseFolderId: time.Now().Format(time.RFC3339),
			},
		}

		resErrChan := ar.AddAccountDriveIntegration(drive)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		// act
		resDriveChan := ar.GetAccountDriveData(drive.AccountId)
		resDrive := <-resDriveChan
		// assert
		assert.NoError(t, resDrive.Err)
		assert.Equal(t, drive.AccountId, resDrive.CloudDriveIntegration.AccountId)
		assert.Equal(t, drive.ApiKey, resDrive.CloudDriveIntegration.ApiKey)
		assert.Equal(t, drive.BaseFolderId, resDrive.CloudDriveIntegration.BaseFolderId)
		assert.Equal(t, drive.Type, resDrive.CloudDriveIntegration.Type)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan = ar.DeleteProfessor(profId)
		resErr = <-resErrChan
		assert.NoError(t, resErr.Err)
	})
}

func TestAccountRepo_GetAccountRepoHubData(t *testing.T) {
	db := connectDB()

	t.Run("fail, get non existent repo hub data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		// act
		resRepoChan := ar.GetAccountRepoHubData("123")
		resRepo := <-resRepoChan

		// assert
		assert.ErrorIs(t, resRepo.Err, models.ErrAccountRepoHubDataNotFound)
	})

	t.Run("ok, get existing repo hub data", func(t *testing.T) {
		// arrange
		ar := InitAccountRepository(*db)

		resProfChan := ar.AddProfessor(domainaggregate.Professor{
			Person: domainaggregate.Person{
				Name:       "dsf",
				Surname:    "sdf",
				Middlename: "sdf",
			},
			ScienceDegree: time.Now().Format(time.RFC3339),
		})
		resProf := <-resProfChan
		assert.NoError(t, resProf.Err)

		repo := models.BaseIntegration{
			AccountId: resProf.Professor.Id,
			ApiKey:    time.Now().Format(time.RFC3339),
			Type:      int(models.GoogleDrive),
		}

		resErrChan := ar.AddAccountRepoHubIntegration(repo)
		resErr := <-resErrChan
		assert.NoError(t, resErr.Err)

		// act
		resRepoChan := ar.GetAccountRepoHubData(repo.AccountId)
		resRepo := <-resRepoChan
		// assert
		assert.NoError(t, resErr.Err)
		assert.Equal(t, repo.AccountId, resRepo.BaseIntegration.AccountId)
		assert.Equal(t, repo.ApiKey, resRepo.BaseIntegration.ApiKey)
		assert.Equal(t, repo.Type, resRepo.BaseIntegration.Type)

		// cleanup
		profId, err := strconv.Atoi(resProf.Professor.Id)
		assert.NoError(t, err)
		resErrChan = ar.DeleteProfessor(profId)
		resErr = <-resErrChan
		assert.NoError(t, resErr.Err)
	})
}

func addTestingAccount(name string, ar *AccountRepository) error {
	resProfChan := ar.AddProfessor(domainaggregate.Professor{
		Person: domainaggregate.Person{
			Name:       "",
			Surname:    "",
			Middlename: "",
		},
		ScienceDegree: "",
	})
	resProf := <-resProfChan
	if resProf.Err != nil {
		return resProf.Err
	}
	resErrChan := ar.AddAccount(models.Account{
		Login: name,
		Hash:  []byte{},
		Salt:  "",
		Id:    resProf.Professor.Id,
	})
	resErr := <-resErrChan
	if resErr.Err != nil {
		return resErr.Err
	}
	return nil
}

func deleteTestingAccount(name string, ar *AccountRepository) error {
	resErrChan := ar.DeleteAccountByLogin(name)
	resErr := <-resErrChan
	if resErr.Err != nil {
		return resErr.Err
	}
	return nil
}

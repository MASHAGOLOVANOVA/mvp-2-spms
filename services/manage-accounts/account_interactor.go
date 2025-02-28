package managestudents

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"fmt"
	entities "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/services/interfaces"
	"mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/manage-accounts/outputdata"
	"mvp-2-spms/services/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/oauth2"
)

const pbkdf2Iterations int = 4096
const pbkdf2HashSize int = 32

type AccountInteractor struct {
	accountRepo interfaces.IAccountRepository
	uniRepo     interfaces.IUniversityRepository
}

func InitAccountInteractor(accRepo interfaces.IAccountRepository, uniRepo interfaces.IUniversityRepository) *AccountInteractor {
	return &AccountInteractor{
		accountRepo: accRepo,
		uniRepo:     uniRepo,
	}
}

func (a *AccountInteractor) GetAccountProfessorId(login string) (string, error) {
	resChan := a.accountRepo.GetAccountByLogin(login)
	resAccount := <-resChan
	if resAccount.Err != nil {
		return "", resAccount.Err
	}
	return resAccount.Account.Id, nil
}

func (a *AccountInteractor) GetProfessorInfo(input inputdata.GetProfessorInfo) (outputdata.GetProfessorInfo, error) {
	resChan := a.accountRepo.GetProfessorById(fmt.Sprint(input.AccountId))
	resProf := <-resChan
	if resProf.Err != nil {
		return outputdata.GetProfessorInfo{}, resProf.Err
	}

	uni, err := a.uniRepo.GetUniversityById(resProf.Professor.UniversityId)
	if err != nil {
		if !errors.Is(err, models.ErrUniNoFound) {
			return outputdata.GetProfessorInfo{}, err
		}
		uni = entities.University{
			Id:   "-1",
			Name: "Университет неизвестен",
			City: "-",
		}
	}

	// add get account login
	output := outputdata.MapToGetAccountInfo(resProf.Professor, uni)
	return output, nil
}

func (a *AccountInteractor) GetPlannerIntegration(input inputdata.GetPlannerIntegration) (outputdata.GetPlannerIntegration, error) {
	resChan := a.accountRepo.GetAccountPlannerData(fmt.Sprint(input.AccountId))
	resPlanner := <-resChan
	if resPlanner.Err != nil {
		return outputdata.GetPlannerIntegration{}, resPlanner.Err
	}

	output := outputdata.MapToGetPlannerIntegration(resPlanner.PlannerIntegration)
	return output, nil
}

func (a *AccountInteractor) GetDriveIntegration(input inputdata.GetDriveIntegration) (outputdata.GetDriveIntegration, error) {
	resChan := a.accountRepo.GetAccountDriveData(fmt.Sprint(input.AccountId))
	resDrive := <-resChan
	if resDrive.Err != nil {
		return outputdata.GetDriveIntegration{}, resDrive.Err
	}

	output := outputdata.MapToGetDriveIntegration(resDrive.CloudDriveIntegration)
	return output, nil
}

func (a *AccountInteractor) SetProfessorPlanner(plannerId, profId string) error {
	resChan := a.accountRepo.GetAccountPlannerData(profId)
	resPlanner := <-resChan
	if resPlanner.Err != nil {
		return resPlanner.Err
	}

	resPlanner.PlannerIntegration.PlannerData.Id = plannerId

	resChan1 := a.accountRepo.UpdateAccountPlannerIntegration(resPlanner.PlannerIntegration)
	err := <-resChan1
	if err.Err != nil {
		return err.Err
	}

	return nil
}

func (a *AccountInteractor) GetProfessorIntegrPlanners(profId string, planner interfaces.IPlannerService) (outputdata.GetProfessorIntegrPlanners, error) {
	resChan := a.accountRepo.GetAccountPlannerData(profId)
	resPlanner := <-resChan
	if resPlanner.Err != nil {
		return outputdata.GetProfessorIntegrPlanners{}, resPlanner.Err
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////
	// check for access token first????????????????????????????????????????????
	token := &oauth2.Token{
		RefreshToken: resPlanner.PlannerIntegration.ApiKey,
	}
	err := planner.Authentificate(token)
	if err != nil {
		return outputdata.GetProfessorIntegrPlanners{}, err
	}

	planners, err := planner.GetAllPlanners()
	if err != nil {
		return outputdata.GetProfessorIntegrPlanners{}, err
	}

	return outputdata.MapToGetProfessorIntegrPlanners(planners), nil
}

func (a *AccountInteractor) GetDriveBaseFolderName(folderId, profId string, cloudDrive interfaces.ICloudDrive) (string, error) {
	resChan := a.accountRepo.GetAccountDriveData(fmt.Sprint(profId))
	resDrive := <-resChan
	if resDrive.Err != nil {
		return "", resDrive.Err
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////
	// check for access token first????????????????????????????????????????????
	token := &oauth2.Token{
		RefreshToken: resDrive.CloudDriveIntegration.ApiKey,
	}
	err := cloudDrive.Authentificate(token)
	if err != nil {
		return "", err
	}

	folderName, err := cloudDrive.GetFolderNameById(folderId)
	if err != nil {
		return "", err
	}

	return folderName, nil
}

func (a *AccountInteractor) GetRepoHubIntegration(input inputdata.GetRepoHubIntegration) (outputdata.GetRepoHubIntegration, error) {
	resChan := a.accountRepo.GetAccountRepoHubData(fmt.Sprint(input.AccountId))
	resRepo := <-resChan
	if resRepo.Err != nil {
		return outputdata.GetRepoHubIntegration{}, resRepo.Err
	}

	output := outputdata.MapToGetRepoHubIntegration(resRepo.BaseIntegration)
	return output, nil
}

func (a *AccountInteractor) SetPlannerIntegration(input inputdata.SetPlannerIntegration, planner interfaces.IPlannerService) (outputdata.SetPlannerIntegration, error) {
	token, err := planner.GetToken(input.AuthCode)
	if err != nil {
		return outputdata.SetPlannerIntegration{}, err
	}

	refreshTok := token.RefreshToken
	accessTok := token.AccessToken
	expires := token.Expiry

	integr := models.PlannerIntegration{
		BaseIntegration: models.BaseIntegration{
			AccountId: fmt.Sprint(input.AccountId),
			ApiKey:    refreshTok,
			Type:      input.Type,
		},
		PlannerData: models.PlannerData{},
	}

	resChan := a.accountRepo.AddAccountPlannerIntegration(integr)
	resErr := <-resChan
	if resErr.Err != nil {
		return outputdata.SetPlannerIntegration{}, resErr.Err
	}

	return outputdata.SetPlannerIntegration{
		AccessToken: accessTok,
		Expiry:      expires,
	}, nil
}

func (a *AccountInteractor) SetDriveIntegration(input inputdata.SetDriveIntegration, drive interfaces.ICloudDrive) (outputdata.SetDriveIntegration, error) {
	token, _ := drive.GetToken(input.AuthCode)
	refreshTok := token.RefreshToken
	accessTok := token.AccessToken
	expires := token.Expiry

	err := drive.Authentificate(token)
	if err != nil {
		return outputdata.SetDriveIntegration{}, err
	}

	baseFolder, err := drive.AddProfessorBaseFolder()
	if err != nil {
		return outputdata.SetDriveIntegration{}, err
	}

	integr := models.CloudDriveIntegration{
		BaseIntegration: models.BaseIntegration{
			AccountId: fmt.Sprint(input.AccountId),
			ApiKey:    refreshTok,
			Type:      input.Type,
		},
		DriveData: baseFolder,
	}

	resChan := a.accountRepo.AddAccountDriveIntegration(integr)
	resErr := <-resChan
	if resErr.Err != nil {
		return outputdata.SetDriveIntegration{}, resErr.Err
	}

	return outputdata.SetDriveIntegration{
		AccessToken: accessTok,
		Expiry:      expires,
	}, nil
}

func (a *AccountInteractor) SetRepoHubIntegration(input inputdata.SetRepoHubIntegration, planner interfaces.IGitRepositoryHub) (outputdata.SetRepoHubIntegration, error) {
	token, err := planner.GetToken(input.AuthCode)
	if err != nil {
		return outputdata.SetRepoHubIntegration{}, err
	}

	refreshTok := token.RefreshToken
	accessTok := token.AccessToken
	expires := token.Expiry
	integr := models.BaseIntegration{
		AccountId: fmt.Sprint(input.AccountId),
		ApiKey:    refreshTok,
		Type:      input.Type,
	}

	resChan := a.accountRepo.AddAccountRepoHubIntegration(integr)
	resErr := <-resChan
	if resErr.Err != nil {
		return outputdata.SetRepoHubIntegration{}, resErr.Err
	}

	return outputdata.SetRepoHubIntegration{
		AccessToken: accessTok,
		Expiry:      expires,
	}, nil
}

func (a *AccountInteractor) GetAccountIntegrations(input inputdata.GetAccountIntegrations) (outputdata.GetAccountIntegrations, error) {
	var (
		outputDrive   *outputdata.GetAccountIntegrationsDrive   = nil
		outputPlanner *outputdata.GetAccountIntegrationsPlanner = nil
		outputRepos   []outputdata.GetAccountIntegrationsIntegr = []outputdata.GetAccountIntegrationsIntegr{}
	)

	found := true
	resChan := a.accountRepo.GetAccountDriveData(fmt.Sprint(input.AccountId))
	resDrive := <-resChan
	if resDrive.Err != nil {
		if !errors.Is(resDrive.Err, models.ErrAccountDriveDataNotFound) {
			return outputdata.GetAccountIntegrations{}, resDrive.Err
		}
		found = false
	}

	if found {
		outputDrive = &outputdata.GetAccountIntegrationsDrive{
			Type: outputdata.GetAccountIntegrationsIntegr{
				Id:   resDrive.CloudDriveIntegration.Type,
				Name: resDrive.CloudDriveIntegration.GetTypeAsString(),
			},
			BaseFolderId: resDrive.CloudDriveIntegration.BaseFolderId,
		}
	}

	found = true
	resChan1 := a.accountRepo.GetAccountPlannerData(fmt.Sprint(input.AccountId))
	resPlanner := <-resChan1
	if resPlanner.Err != nil {
		if !errors.Is(resPlanner.Err, models.ErrAccountPlannerDataNotFound) {
			return outputdata.GetAccountIntegrations{}, resPlanner.Err
		}
		found = false
	}

	if found {
		outputPlanner = &outputdata.GetAccountIntegrationsPlanner{
			Type: outputdata.GetAccountIntegrationsIntegr{
				Id:   resPlanner.PlannerIntegration.Type,
				Name: resPlanner.PlannerIntegration.GetTypeAsString(),
			},
			PlannerName: resPlanner.PlannerIntegration.PlannerData.Id, ///////////////////////////////////////change
		}
	}

	found = true
	resChan2 := a.accountRepo.GetAccountRepoHubData(fmt.Sprint(input.AccountId))
	resRepo := <-resChan2
	if resRepo.Err != nil {
		if !errors.Is(resRepo.Err, models.ErrAccountRepoHubDataNotFound) {
			return outputdata.GetAccountIntegrations{}, resRepo.Err
		}
		found = false
	}

	if found {
		outputRepos = append(outputRepos, outputdata.GetAccountIntegrationsIntegr{
			Id:   resRepo.BaseIntegration.Type,
			Name: resRepo.BaseIntegration.GetRepoHubTypeAsString(),
		})
	}

	return outputdata.MapToGetAccountIntegrations(outputDrive, outputPlanner, outputRepos), nil
}

func (a *AccountInteractor) CheckCredsValidity(input inputdata.CheckCredsValidity) (bool, error) {
	resChan := a.accountRepo.GetAccountByLogin(input.Login)
	resAccount := <-resChan
	if resAccount.Err != nil {
		return false, resAccount.Err
	}

	key := pbkdf2.Key([]byte(input.Password), []byte(resAccount.Account.Salt), pbkdf2Iterations, pbkdf2HashSize, sha512.New)

	return bytes.Equal(key, resAccount.Account.Hash), nil
}

func (a *AccountInteractor) CheckUsernameExists(input inputdata.CheckUsernameExists) (bool, error) {
	resChan := a.accountRepo.GetAccountByLogin(input.Login)
	resAccount := <-resChan
	if resAccount.Err != nil {
		if errors.Is(resAccount.Err, models.ErrAccountNotFound) {
			return false, nil
		}
		return false, resAccount.Err
	}
	return true, nil
}

func (a *AccountInteractor) SignUp(input inputdata.SignUp) (outputdata.SignUp, error) {
	salt := uuid.NewString()
	passHash := pbkdf2.Key([]byte(input.Password), []byte(salt), pbkdf2Iterations, pbkdf2HashSize, sha512.New)

	prof := entities.Professor{
		Person: entities.Person{
			Name:       input.Name,
			Surname:    input.Surname,
			Middlename: input.Middlename,
		},
		ScienceDegree: input.ScienceDegree,
		UniversityId:  fmt.Sprint(input.UniId),
	}

	resChan := a.accountRepo.AddProfessor(prof)
	resProf := <-resChan
	if resProf.Err != nil {
		return outputdata.SignUp{}, resProf.Err
	}

	account := models.Account{
		Login: input.Login,
		Hash:  passHash,
		Salt:  salt,
		Id:    prof.Id,
	}

	resChan1 := a.accountRepo.AddAccount(account)
	resErr := <-resChan1
	if resErr.Err != nil {
		return outputdata.SignUp{}, resErr.Err
	}

	return outputdata.SignUp{
		Id:    account.Id,
		Login: account.Login,
	}, nil
}

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/internal"
	mngInterfaces "mvp-2-spms/services/interfaces"
	ainputdata "mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/manage-projects/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	requestbodies "mvp-2-spms/web_server/handlers/request-bodies"
	responsebodies "mvp-2-spms/web_server/handlers/response-bodies"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	projectInteractor interfaces.IProjetInteractor
	accountInteractor interfaces.IAccountInteractor
	cloudDrives       internal.CloudDrives
	repoHubs          internal.GitRepositoryHubs
}

func InitProjectHandler(projInteractor interfaces.IProjetInteractor, acc interfaces.IAccountInteractor, cd internal.CloudDrives, rh internal.GitRepositoryHubs) ProjectHandler {
	return ProjectHandler{
		projectInteractor: projInteractor,
		accountInteractor: acc,
		cloudDrives:       cd,
		repoHubs:          rh,
	}
}

func (h *ProjectHandler) GetAllProfProjects(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	input := inputdata.GetProfessorProjects{
		ProfessorId: uint(id),
	}
	filter := r.URL.Query().Get("filter")

	var filterStatus int = -1
	switch filter {
	case domainaggregate.ProjectNotConfirmed.String():
		{
			filterStatus = int(domainaggregate.ProjectNotConfirmed)
		}
	case domainaggregate.ProjectCancelled.String():
		{
			filterStatus = int(domainaggregate.ProjectCancelled)
		}
	case domainaggregate.ProjectInProgress.String():
		{
			filterStatus = int(domainaggregate.ProjectInProgress)
		}

	case domainaggregate.ProjectFinished.String():
		{
			filterStatus = int(domainaggregate.ProjectFinished)
		}
	}

	if filterStatus != -1 {
		input.FilterStatus = &filterStatus
	}

	result, err := h.projectInteractor.GetProfessorProjects(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v", err)
	}
}

func (h *ProjectHandler) GetProjectStatusList(w http.ResponseWriter, r *http.Request) {
	result := responsebodies.ProjectStatuses{
		Statuses: []responsebodies.Status{
			{
				Name:  domainaggregate.ProjectNotConfirmed.String(),
				Value: int(domainaggregate.ProjectNotConfirmed),
			},
			{
				Name:  domainaggregate.ProjectInProgress.String(),
				Value: int(domainaggregate.ProjectInProgress),
			},
			{
				Name:  domainaggregate.ProjectFinished.String(),
				Value: int(domainaggregate.ProjectFinished),
			},
			{
				Name:  domainaggregate.ProjectCancelled.String(),
				Value: int(domainaggregate.ProjectCancelled),
			},
		},
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProjectHandler) GetProjectStageList(w http.ResponseWriter, r *http.Request) {
	result := responsebodies.ProjectStages{
		Stages: []responsebodies.Status{
			{
				Name:  domainaggregate.Analysis.String(),
				Value: int(domainaggregate.Analysis),
			},
			{
				Name:  domainaggregate.Design.String(),
				Value: int(domainaggregate.Design),
			},
			{
				Name:  domainaggregate.Development.String(),
				Value: int(domainaggregate.Development),
			},
			{
				Name:  domainaggregate.Testing.String(),
				Value: int(domainaggregate.Testing),
			},
			{
				Name:  domainaggregate.Deployment.String(),
				Value: int(domainaggregate.Deployment),
			},
		},
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}
}

func (h *ProjectHandler) GetProjectCommits(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	projectId, err := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", r.URL.Query().Get("from"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	integInput := ainputdata.GetRepoHubIntegration{
		AccountId: uint(id),
	}

	found := true
	hubInfo, err := h.accountInteractor.GetRepoHubIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountRepoHubDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	if !found {
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	input := inputdata.GetProjectCommits{
		ProfessorId: uint(id),
		ProjectId:   uint(projectId),
		From:        from,
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	result, err := h.projectInteractor.GetProjectCommits(input, h.repoHubs[models.GetRepoHubName(hubInfo.Type)])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	projectId, err := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	input := inputdata.GetProjectById{
		ProfessorId: uint(id),
		ProjectId:   uint(projectId),
	}

	result, err := h.projectInteractor.GetProjectById(input)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}
}

func (h *ProjectHandler) GetProjectSupReport(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	projectId, err := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	var reqB inputdata.GetProjectSupReport
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqB.ProfessorId = uint(id)
	reqB.ProjectId = uint(projectId)

	result, err := h.projectInteractor.GetProjectSupReport(reqB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, result)
}

func (h *ProjectHandler) GetProjectStatistics(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	projectId, err := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	input := inputdata.GetProjectStatsById{
		ProfessorId: uint(id),
		ProjectId:   uint(projectId),
	}

	result, err := h.projectInteractor.GetProjectStatsById(input)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}
}

func (h *ProjectHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var reqB requestbodies.AddProject
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	integInput := ainputdata.GetDriveIntegration{
		AccountId: uint(id),
	}

	found := true
	driveInfo, err := h.accountInteractor.GetDriveIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountDriveDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var drive mngInterfaces.ICloudDrive
	if found {
		drive = h.cloudDrives[models.CloudDriveName(driveInfo.Type)]
	}

	input := inputdata.AddProject{
		ProfessorId:         uint(id),
		Theme:               reqB.Theme,
		StudentId:           uint(reqB.StudentId),
		Year:                uint(reqB.Year),
		RepositoryOwnerName: reqB.RepoOwner,
		RepositoryName:      reqB.RepositoryName,
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	student_id, err := h.projectInteractor.AddProject(input, drive)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(student_id); err != nil {
		log.Printf("Ошибка при кодировании id студента: %v", err)
	}
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	projectId, err := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var reqB requestbodies.UpdateProject
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := inputdata.UpdateProject{
		Id:                  int(projectId),
		ProfessorId:         &id,
		Theme:               reqB.Theme,
		StudentId:           reqB.StudentId,
		Year:                reqB.Year,
		RepositoryOwnerName: reqB.RepoOwner,
		RepositoryName:      reqB.RepositoryName,
		Status:              reqB.Status,
		Stage:               reqB.Stage,
	}

	err = h.projectInteractor.UpdateProject(input, nil)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotProfessors) {
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ProjectHandler) UpdateProjectGrading(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	projectId, err := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var reqB inputdata.UpdateProjectGrading
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqB.ProjId = int(projectId)
	reqB.ProfessorId = &id

	err = h.projectInteractor.UpdateProjectGrading(reqB)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotProfessors) {
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

package handlers

import (
	"encoding/json"
	"mvp-2-spms/internal"
	ainputdata "mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/manage-projects/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	requestbodies "mvp-2-spms/web_server/handlers/request-bodies"
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
	user := GetSessionUser(r)
	id, _ := strconv.Atoi(user.GetProfId())
	input := inputdata.GetProfessorProjects{
		ProfessorId: uint(id),
	}
	result := h.projectInteractor.GetProfessorProjects(input)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProjectHandler) GetProjectCommits(w http.ResponseWriter, r *http.Request) {
	user := GetSessionUser(r)
	id, _ := strconv.Atoi(user.GetProfId())
	projectId, _ := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", r.URL.Query().Get("from"))

	integInput := ainputdata.GetRepoHubIntegration{
		AccountId: uint(id),
	}
	hubInfo := h.accountInteractor.GetRepoHubIntegration(integInput)
	input := inputdata.GetProjectCommits{
		ProfessorId: uint(id),
		ProjectId:   uint(projectId),
		From:        from,
	}
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	result := h.projectInteractor.GetProjectCommits(input, h.repoHubs[models.GetRepoHubName(hubInfo.Type)])
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	user := GetSessionUser(r)
	id, _ := strconv.Atoi(user.GetProfId())
	projectId, _ := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	input := inputdata.GetProjectById{
		ProfessorId: uint(id),
		ProjectId:   uint(projectId),
	}
	result := h.projectInteractor.GetProjectById(input)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProjectHandler) GetProjectStatistics(w http.ResponseWriter, r *http.Request) {
	user := GetSessionUser(r)
	id, _ := strconv.Atoi(user.GetProfId())
	projectId, _ := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	input := inputdata.GetProjectStatsById{
		ProfessorId: uint(id),
		ProjectId:   uint(projectId),
	}
	result := h.projectInteractor.GetProjectStatsById(input)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProjectHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	user := GetSessionUser(r)
	id, _ := strconv.Atoi(user.GetProfId())

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
	err := decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	integInput := ainputdata.GetDriveIntegration{
		AccountId: uint(id),
	}
	driveInfo := h.accountInteractor.GetDriveIntegration(integInput)

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
	student_id := h.projectInteractor.AddProject(input, h.cloudDrives[models.CloudDriveName(driveInfo.Type)])
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student_id)
}

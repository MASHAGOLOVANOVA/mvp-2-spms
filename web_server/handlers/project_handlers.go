package handlers

import (
	"encoding/json"
	"mvp-2-spms/services/manage-projects/inputdata"
	"mvp-2-spms/web_server/handlers/interfaces"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	projectInteractor interfaces.IProjetInteractor
}

func InitProjectHandler(projInteractor interfaces.IProjetInteractor) ProjectHandler {
	return ProjectHandler{
		projectInteractor: projInteractor,
	}
}

func (h *ProjectHandler) GetAllProfProjects(w http.ResponseWriter, r *http.Request) {
	professorIdCookie, _ := r.Cookie("professor_id")
	professorId, _ := strconv.ParseUint(professorIdCookie.Value, 10, 32)
	input := inputdata.GetPfofessorProjects{
		ProfessorId: uint(professorId),
	}
	result := h.projectInteractor.GetProfessorProjects(input)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProjectHandler) GetProjectCommits(w http.ResponseWriter, r *http.Request) {
	professorIdCookie, _ := r.Cookie("professor_id")
	professorId, _ := strconv.ParseUint(professorIdCookie.Value, 10, 32)
	projectId, _ := strconv.ParseUint(chi.URLParam(r, "projectID"), 10, 32)
	from, _ := time.Parse("2006-01-02T15:04:05.000Z", r.URL.Query().Get("from"))
	input := inputdata.GetProjectCommits{
		ProfessorId: uint(professorId),
		ProjectId:   uint(projectId),
		From:        from,
	}
	result := h.projectInteractor.GetProjectCommits(input)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

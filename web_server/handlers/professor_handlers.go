package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"mvp-2-spms/services/manage-professors/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	"net/http"
	"strconv"
)

type ProfessorHandler struct {
	professorInteractor interfaces.IProfessorInteractor
}

func InitProfessorHandler(profInteractor interfaces.IProfessorInteractor) ProfessorHandler {
	return ProfessorHandler{
		professorInteractor: profInteractor,
	}
}

func (h *ProfessorHandler) GetStudentApplications(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	_, err = strconv.Atoi(user.GetAccId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	result, err := h.professorInteractor.GetStudentApplications(user.GetAccId())
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

func (h *ProfessorHandler) GetProfessorApplications(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	_, err = strconv.Atoi(user.GetAccId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	result, err := h.professorInteractor.GetProfessorApplications(user.GetAccId())
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

func (h *ProfessorHandler) GetProfessors(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	_, err = strconv.Atoi(user.GetAccId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	result, err := h.professorInteractor.GetProfessors()
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

func (h *ProfessorHandler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetAccId())
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

	_, err = strconv.ParseUint(chi.URLParam(r, "applicationID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	// декодируем тело запроса
	var input inputdata.ApplicationStatus
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id != int(input.ProfessorId) {
		if errors.Is(err, models.ErrProjectNotProfessors) {
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
	}

	err = h.professorInteractor.UpdateApplicationStatus(input)
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
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProfessorHandler) Apply(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	id, err := strconv.Atoi(user.GetAccId())
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

	profId, err := strconv.ParseUint(chi.URLParam(r, "professorID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	input := inputdata.Apply{
		ProfessorId: uint(profId),
		StudentId:   uint(id),
	}

	err = h.professorInteractor.Apply(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

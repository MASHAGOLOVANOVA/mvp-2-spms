package handlers

import (
	"encoding/json"
	"log"
	"mvp-2-spms/services/manage-universities/inputdata"
	"mvp-2-spms/web_server/handlers/interfaces"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UniversityHandler struct {
	uniInteractor interfaces.IUniversityInteractor
}

func InitUniversityHandler(uInteractor interfaces.IUniversityInteractor) UniversityHandler {
	return UniversityHandler{
		uniInteractor: uInteractor,
	}
}

func (h *UniversityHandler) GetAllUniEdProgrammes(w http.ResponseWriter, r *http.Request) {
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

	uniId, err := strconv.ParseUint(chi.URLParam(r, "uniID"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	input := inputdata.GetUniEducationalProgrammes{
		ProfessorId:  uint(id),
		UniversityId: uint(uniId),
	}
	result, err := h.uniInteractor.GetUniEdProgrammes(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

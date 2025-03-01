package handlers

import (
	"encoding/json"
	"errors"
	"log"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/internal"
	mngInterfaces "mvp-2-spms/services/interfaces"
	ainputdata "mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/manage-tasks/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	requestbodies "mvp-2-spms/web_server/handlers/request-bodies"
	responsebodies "mvp-2-spms/web_server/handlers/response-bodies"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskInteractor    interfaces.ITaskInteractor
	accountInteractor interfaces.IAccountInteractor
	cloudDrives       internal.CloudDrives
}

func InitTaskHandler(taskInteractor interfaces.ITaskInteractor, acc interfaces.IAccountInteractor, cd internal.CloudDrives) TaskHandler {
	return TaskHandler{
		taskInteractor:    taskInteractor,
		accountInteractor: acc,
		cloudDrives:       cd,
	}
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
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
	var reqB requestbodies.AddTask
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

	input := inputdata.AddTask{
		ProfessorId: uint(id),
		Name:        reqB.Name,
		Description: reqB.Description,
		Deadline:    reqB.Deadline,
		ProjectId:   uint(projectId),
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	task_id, err := h.taskInteractor.AddTask(input, drive)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task_id); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v", err)
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
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

	taskId, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 32)
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
	var input inputdata.UpdateTask
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input.Id = int(taskId)
	input.ProfId = id

	err = h.taskInteractor.UpdateTask(input, nil)
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

func (h *TaskHandler) GetTaskStatusList(w http.ResponseWriter, r *http.Request) {
	result := responsebodies.TaskStatuses{
		Statuses: []responsebodies.Status{
			{
				Name:  domainaggregate.NotStarted.String(),
				Value: int(domainaggregate.NotStarted),
			},
			{
				Name:  domainaggregate.InProgress.String(),
				Value: int(domainaggregate.InProgress),
			},
			{
				Name:  domainaggregate.Finished.String(),
				Value: int(domainaggregate.Finished),
			},
		},
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v", err)
	}
}

func (h *TaskHandler) GetAllProjectTasks(w http.ResponseWriter, r *http.Request) {
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

	input := inputdata.GetProjectTasks{
		ProfessorId: uint(id),
		ProjectId:   uint(projectId),
	}

	result, err := h.taskInteractor.GetProjectTasks(input)
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

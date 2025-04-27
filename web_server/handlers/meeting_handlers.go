package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	domainaggregate "mvp-2-spms/domain-aggregate"
	"mvp-2-spms/internal"
	mngInterfaces "mvp-2-spms/services/interfaces"
	ainputdata "mvp-2-spms/services/manage-accounts/inputdata"
	minputdata "mvp-2-spms/services/manage-meetings/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	requestbodies "mvp-2-spms/web_server/handlers/request-bodies"
	responsebodies "mvp-2-spms/web_server/handlers/response-bodies"
	"net/http"
	"strconv"
	"time"
)

type MeetingHandler struct {
	meetingInteractor interfaces.IMeetingInteractor
	accountInteractor interfaces.IAccountInteractor
	planners          internal.Planners
}

func InitMeetingHandler(meetInteractor interfaces.IMeetingInteractor, acc interfaces.IAccountInteractor, pl internal.Planners) MeetingHandler {
	return MeetingHandler{
		meetingInteractor: meetInteractor,
		accountInteractor: acc,
		planners:          pl,
	}
}

func (h *MeetingHandler) AddMeeting(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	print("foundUser")

	id, err := strconv.Atoi(user.GetAccId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	print("foundProf")

	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var reqB requestbodies.AddMeeting
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	print("decoded_fields")

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(id),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var planner mngInterfaces.IPlannerService
	if found {
		planner = h.planners[models.PlannerName(calendarInfo.Type)]
	}

	meetingInput := minputdata.AddMeeting{
		ProfessorId: uint(id),
		Name:        reqB.Name,
		Description: reqB.Description,
		MeetingTime: reqB.MeetingTime,
		StudentId:   reqB.StudentId,
		IsOnline:    reqB.IsOnline,
		ProjectId:   uint(reqB.ProjectId),
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	meeting_id, err := h.meetingInteractor.AddMeeting(meetingInput, planner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(meeting_id); err != nil {
		log.Printf("Ошибка при кодировании meeting id: %v", err)
	}
}

func (h *MeetingHandler) GetProfessorStudentMeetings(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	print("foundUser")

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

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(id),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var planner mngInterfaces.IPlannerService
	if found {
		planner = h.planners[models.PlannerName(calendarInfo.Type)]
	}

	slots, err := h.meetingInteractor.GetProfessorStudentMeetings(id, planner, minputdata.GetProfessorMeetings{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании результата: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(slots); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}

}

func (h *MeetingHandler) GetStudentMeetings(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	print("foundUser")

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

	slots, err := h.meetingInteractor.GetStudentMeetings(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании результата: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(slots); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}

}

func (h *MeetingHandler) GetProfessorSlots(w http.ResponseWriter, r *http.Request) {
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

	profId, err := strconv.Atoi(chi.URLParam(r, "professorID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(profId),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var planner mngInterfaces.IPlannerService
	if found {
		planner = h.planners[models.PlannerName(calendarInfo.Type)]
	}

	get := r.URL.Query().Get("filter")

	slots, err := h.meetingInteractor.GetProfessorSlots(id, profId, planner, get)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании результата: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(slots); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}

}

func (h *MeetingHandler) DeleteSlot(w http.ResponseWriter, r *http.Request) {
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

	slotId, err := strconv.Atoi(chi.URLParam(r, "slotID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(id),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var planner mngInterfaces.IPlannerService
	if found {
		planner = h.planners[models.PlannerName(calendarInfo.Type)]
	}

	err = h.meetingInteractor.DeleteSlot(id, slotId, planner)
	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *MeetingHandler) ChooseSlot(w http.ResponseWriter, r *http.Request) {
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

	slotId, err := strconv.Atoi(chi.URLParam(r, "slotID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	err = h.meetingInteractor.ChooseSlot(id, slotId)
	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *MeetingHandler) AddSlotProject(w http.ResponseWriter, r *http.Request) {
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

	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var reqB struct {
		ProjectId int `json:"project_id"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	slotId, err := strconv.Atoi(chi.URLParam(r, "slotID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	err = h.meetingInteractor.BindSlotToProject(slotId, reqB.ProjectId)
	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *MeetingHandler) UpdateSlot(w http.ResponseWriter, r *http.Request) {
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

	// декодируем тело запроса
	var reqB struct {
		Description string `json:"description"`
		MeetingTime string `json:"meeting_time"` // Получаем как строку
		Duration    int    `json:"duration"`
		IsOnline    bool   `json:"is_online"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	// Парсим время вручную с учетом локали
	meetingTime, err := time.Parse(time.RFC3339, reqB.MeetingTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid time format")
		return
	}

	// Преобразуем в нужный часовой пояс (например, Europe/Moscow)
	loc, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		loc = time.UTC // Fallback
	}

	localTime := meetingTime.In(loc)

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(id),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var planner mngInterfaces.IPlannerService
	if found {
		planner = h.planners[models.PlannerName(calendarInfo.Type)]
	}

	// Создаем входные данные с корректным временем
	meetingInput := minputdata.AddSlot{
		ProfessorId: uint(id),
		Duration:    reqB.Duration,
		Description: reqB.Description,
		MeetingTime: localTime, // Используем локальное время
		IsOnline:    reqB.IsOnline,
	}

	slotId, err := strconv.Atoi(chi.URLParam(r, "slotID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	err = h.meetingInteractor.UpdateSlot(slotId, meetingInput, planner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *MeetingHandler) AddSlot(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	print("foundUser")

	id, err := strconv.Atoi(user.GetAccId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	print("foundProf")

	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var reqB requestbodies.AddSlot
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&reqB)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(id),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var planner mngInterfaces.IPlannerService
	if found {
		planner = h.planners[models.PlannerName(calendarInfo.Type)]
	}

	meetingInput := minputdata.AddSlot{
		ProfessorId: uint(id),
		Duration:    reqB.Duration,
		Description: reqB.Description,
		MeetingTime: reqB.MeetingTime,
		IsOnline:    reqB.IsOnline,
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	slot_id, err := h.meetingInteractor.AddSlot(meetingInput, planner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(slot_id); err != nil {
		log.Printf("Ошибка при кодировании slot id: %v", err)
	}
}

func (h *MeetingHandler) GetMeetingStatusList(w http.ResponseWriter, r *http.Request) {
	result := responsebodies.MeetingStatuses{
		Statuses: []responsebodies.Status{
			{
				Name:  domainaggregate.MeetingPlanned.String(),
				Value: int(domainaggregate.MeetingPlanned),
			},
			{
				Name:  domainaggregate.MeetingPassed.String(),
				Value: int(domainaggregate.MeetingPassed),
			},
			{
				Name:  domainaggregate.MeetingCancelled.String(),
				Value: int(domainaggregate.MeetingCancelled),
			},
		},
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}
}

func (h *MeetingHandler) GetProfessorMeetings(w http.ResponseWriter, r *http.Request) {
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

	from, err := time.Parse("2006-01-02T15:04:05.000Z", r.URL.Query().Get("from"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	input := minputdata.GetProfessorMeetings{
		ProfessorId: uint(id),
		From:        &from,
	}

	toStr := r.URL.Query().Get("to")
	if toStr != "" {
		to, err := time.Parse("2006-01-02T15:04:05.000Z", toStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		input.To = &to
	}

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(id),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
				log.Printf("Ошибка при кодировании ответа: %v", err)
			}
			return
		}
		found = false
	}

	var planner mngInterfaces.IPlannerService
	if found {
		planner = h.planners[models.PlannerName(calendarInfo.Type)]
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// TODO: pass api key/clone with new key///////////////////////////////////////////////////////////////////////////////
	result, err := h.meetingInteractor.GetProfessorMeetings(input, planner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании результата: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}
}

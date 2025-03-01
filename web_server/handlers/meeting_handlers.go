package handlers

import (
	"encoding/json"
	"errors"
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

	id, err := strconv.Atoi(user.GetProfId())
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

	id, err := strconv.Atoi(user.GetProfId())
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
		From:        from,
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
		input.To = to
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
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Ошибка при кодировании результата: %v", err)
	}
}

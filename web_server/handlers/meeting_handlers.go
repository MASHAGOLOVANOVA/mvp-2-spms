package handlers

import (
	"encoding/json"
	"errors"
	"mvp-2-spms/internal"
	mngInterfaces "mvp-2-spms/services/interfaces"
	ainputdata "mvp-2-spms/services/manage-accounts/inputdata"
	minputdata "mvp-2-spms/services/manage-meetings/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	requestbodies "mvp-2-spms/web_server/handlers/request-bodies"
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
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

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

	integInput := ainputdata.GetPlannerIntegration{
		AccountId: uint(id),
	}

	found := true
	calendarInfo, err := h.accountInteractor.GetPlannerIntegration(integInput)
	if err != nil {
		if !errors.Is(err, models.ErrAccountPlannerDataNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
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
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(meeting_id)
}

func (h *MeetingHandler) GetProfessorMeetings(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	id, err := strconv.Atoi(user.GetProfId())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05.000Z", r.URL.Query().Get("from"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
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
			json.NewEncoder(w).Encode(err.Error())
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
			json.NewEncoder(w).Encode(err.Error())
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
	json.NewEncoder(w).Encode(result)
}

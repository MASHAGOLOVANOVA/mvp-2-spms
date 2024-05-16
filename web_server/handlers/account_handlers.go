package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"mvp-2-spms/internal"
	"mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	accountInteractor interfaces.IAccountInteractor
	cloudDrives       internal.CloudDrives
}

func InitAccountHandler(accountInteractor interfaces.IAccountInteractor, cd internal.CloudDrives) AccountHandler {
	return AccountHandler{
		accountInteractor: accountInteractor,
		cloudDrives:       cd,
	}
}

func (h *AccountHandler) GetAccountIntegrations(w http.ResponseWriter, r *http.Request) {
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

	input := inputdata.GetAccountIntegrations{
		AccountId: uint(id),
	}

	result, err := h.accountInteractor.GetAccountIntegrations(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if result.CloudDrive != nil && result.CloudDrive.BaseFolderId != "" {
		result.CloudDrive.BaseFolderName, err = h.accountInteractor.GetDriveBaseFolderName(
			result.CloudDrive.BaseFolderId, fmt.Sprint(id), h.cloudDrives[models.CloudDriveName(result.CloudDrive.Type.Id)])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *AccountHandler) GetAccountInfo(w http.ResponseWriter, r *http.Request) {
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

	input := inputdata.GetProfessorInfo{
		AccountId: uint(id),
	}

	result, err := h.accountInteractor.GetProfessorInfo(input)
	if err != nil {
		if errors.Is(err, models.ErrProfessorNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

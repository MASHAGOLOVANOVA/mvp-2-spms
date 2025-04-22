package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mvp-2-spms/internal"
	"mvp-2-spms/web_server/handlers/interfaces"
	"net/http"
	"os"
	"strconv"
	"strings"

	"mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/models"
)

type CloudDriveHandler struct {
	drives            internal.CloudDrives
	accountInteractor interfaces.IAccountInteractor
}

func InitCloudDriveHandler(drives internal.CloudDrives, acc interfaces.IAccountInteractor) CloudDriveHandler {
	return CloudDriveHandler{
		drives:            drives,
		accountInteractor: acc,
	}
}

func (h *CloudDriveHandler) GetGoogleDriveLink(w http.ResponseWriter, r *http.Request) {
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

	returnURL := r.URL.Query().Get("redirect")
	redirectURI := os.Getenv("RETURN_URL") + "/api/v1/auth/integration/access/googledrive"

	result, err := h.drives[models.GoogleDrive].GetAuthLink(redirectURI, int(uint(id)), returnURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(result)); err != nil {
		log.Printf("Ошибка при получении интеграции с cloudDrive: %v", err)
	}
}

func (h *CloudDriveHandler) OAuthCallbackGoogleDrive(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	decodedState, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	// needs further update
	params := strings.Split(string(decodedState), ",")
	accountId, err := strconv.Atoi(params[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}
	redirect := params[1]

	input := inputdata.SetDriveIntegration{
		AccountId: uint(accountId),
		AuthCode:  code,
		Type:      int(models.GoogleDrive),
	}

	result, err := h.accountInteractor.SetDriveIntegration(input, h.drives[models.GoogleDrive])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	w.Header().Add("Google-Calendar-Token", result.AccessToken)
	w.Header().Add("Google-Calendar-Token-Exp", result.Expiry.String())
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

// GetYandexDiskLink возвращает URL для авторизации через Яндекс
func (h *CloudDriveHandler) GetYandexDiskLink(w http.ResponseWriter, r *http.Request) {
	user, err := GetSessionUser(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(user.GetAccId())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	returnURL := r.URL.Query().Get("redirect")
	redirectURI := os.Getenv("RETURN_URL") + "api/v1/auth/integration/access/yandexdisk" // Добавьте в .env YANDEX_REDIRECT_URI

	result, err := h.drives[models.YandexDisk].GetAuthLink(redirectURI, int(uint(id)), returnURL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithText(w, http.StatusOK, result)
}

// OAuthCallbackYandexDisk обрабатывает callback от Яндекс OAuth
func (h *CloudDriveHandler) OAuthCallbackYandexDisk(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	decodedState, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	params := strings.Split(string(decodedState), ",")
	if len(params) < 2 {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid state parameter"))
		return
	}

	accountId, err := strconv.Atoi(params[0])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	redirect := params[1]

	input := inputdata.SetDriveIntegration{
		AccountId: uint(accountId),
		AuthCode:  code,
		Type:      int(models.YandexDisk),
	}

	result, err := h.accountInteractor.SetDriveIntegration(input, h.drives[models.YandexDisk])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Yandex-Disk-Token", result.AccessToken)
	w.Header().Add("Yandex-Disk-Token-Exp", result.Expiry.String())
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

// Вспомогательные функции для ответов
func respondWithError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v", err)
	}
}

func respondWithText(w http.ResponseWriter, code int, text string) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(code)
	if _, err := w.Write([]byte(text)); err != nil {
		log.Printf("Ошибка при записи ответа: %v", err)
	}
}

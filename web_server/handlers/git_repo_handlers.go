package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"mvp-2-spms/internal"
	"mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type GitRepoHandler struct {
	repos             internal.GitRepositoryHubs
	accountInteractor interfaces.IAccountInteractor
}

func InitGitRepoHandler(repos internal.GitRepositoryHubs, acc interfaces.IAccountInteractor) GitRepoHandler {
	return GitRepoHandler{
		repos:             repos,
		accountInteractor: acc,
	}
}

func (h *GitRepoHandler) GetGitHubLink(w http.ResponseWriter, r *http.Request) {
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

	returnURL := r.URL.Query().Get("redirect")
	redirectURI := os.Getenv("RETURN_URL") + "/api/v1/auth/integration/access/github"

	result, err := h.repos[models.GitHub].GetAuthLink(redirectURI, int(uint(id)), returnURL)
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
		log.Printf("Ошибка при получении интеграции с gitRepo: %v", err)
	}
}

func (h *GitRepoHandler) OAuthCallbackGitHub(w http.ResponseWriter, r *http.Request) {
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

	input := inputdata.SetRepoHubIntegration{
		AccountId: uint(accountId),
		AuthCode:  code,
		Type:      int(models.GitHub),
	}

	result, err := h.accountInteractor.SetRepoHubIntegration(input, h.repos[models.GitHub])
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

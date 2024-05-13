package handlers

import (
	"encoding/json"
	"mvp-2-spms/web_server/session"
	"net/http"
	"strconv"
)

func BotAuthentificator(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			botToken := r.Header.Get("Bot-Token")
			if botToken != session.BotToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
}

func Authentificator(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			creds, err := GetCredentials(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err.Error())
				return
			}

			userSession, exists := session.Sessions[creds.Session]
			if !exists {
				// If the session token is not present in session map, return an unauthorized error
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if userSession.IsExpired() {
				delete(session.Sessions, creds.Session)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func GetSessionUser(r *http.Request) (session.UserInfo, error) {
	// session is already validated by middleware
	creds, err := GetCredentials(r)
	if err != nil {
		return session.UserInfo{}, err
	}
	userSession := session.Sessions[creds.Session]
	return userSession.GetUser(), nil
}

// ///////////////////////////////////////////////////////////////////////////??
func GetCredentials(r *http.Request) (Credntials, error) {
	professorId, err := strconv.ParseUint(r.Header.Get("Professor-Id"), 10, 32)
	if err != nil {
		return Credntials{}, err
	}

	session := r.Header.Get("Session-Id")
	gcTok := r.Header.Get("Google-Calendar-Token")
	gdTok := r.Header.Get("Google-Drive-Token")
	ghTok := r.Header.Get("GitHub-Token")

	return Credntials{
		ProfessorId:         uint(professorId),
		Session:             session,
		GoogleCalendarToken: gcTok,
		GoogleDriveToken:    gdTok,
		GitHubToken:         ghTok,
	}, err
}

type Credntials struct {
	ProfessorId         uint
	Session             string
	GoogleCalendarToken string
	GoogleDriveToken    string
	GitHubToken         string
}

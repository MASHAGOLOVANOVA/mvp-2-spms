package universityservice

import (
	"encoding/json"
	"log"
	"mvp-2-spms/student-service/session"
	"net/http"
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
				if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
					log.Printf("Ошибка при кодировании результата: %v", err)
				}
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
	session := r.Header.Get("Session-Id")
	gcTok := r.Header.Get("Google-Calendar-Token")
	gdTok := r.Header.Get("Google-Drive-Token")
	ghTok := r.Header.Get("GitHub-Token")

	return Credntials{
		Session:             session,
		GoogleCalendarToken: gcTok,
		GoogleDriveToken:    gdTok,
		GitHubToken:         ghTok,
	}, nil
}

type Credntials struct {
	Session             string
	GoogleCalendarToken string
	GoogleDriveToken    string
	GitHubToken         string
}

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"mvp-2-spms/services/manage-accounts/inputdata"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/handlers/interfaces"
	requestbodies "mvp-2-spms/web_server/handlers/request-bodies"
	responsebodies "mvp-2-spms/web_server/handlers/response-bodies"
	"mvp-2-spms/web_server/session"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuthHandler struct {
	accountInteractor interfaces.IAccountInteractor
}

func InitAuthHandler(acc interfaces.IAccountInteractor) AuthHandler {
	return AuthHandler{
		accountInteractor: acc,
	}
}

func (h *AuthHandler) SignInBot(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var creds requestbodies.CredentialsBot
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inp := inputdata.CheckUsernameExists{
		Login: creds.Phone,
	}

	found, err := h.accountInteractor.CheckUsernameExists(inp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if !found {
		w.WriteHeader(http.StatusConflict)
		if encodeErr := json.NewEncoder(w).Encode("account with phone is not found"); encodeErr != nil {
			log.Printf("Ошибка при кодировании ответа: %v", encodeErr)
		}
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewString() + "/" + creds.Phone
	expiresAt := time.Now().Add(session.SessionDefaultExpTime)

	profId, err := h.accountInteractor.GetAccountProfessorId(creds.Phone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	user := session.InitUserInfo(creds.Phone, profId)
	session.Sessions[sessionToken] = session.InitSession(user, expiresAt)

	resBody := responsebodies.SessionToken{
		Token:  sessionToken,
		Expiry: expiresAt,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var creds requestbodies.Credentials
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := inputdata.CheckCredsValidity{
		Login:    creds.Username,
		Password: creds.Password,
	}

	valid, err := h.accountInteractor.CheckCredsValidity(input)
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewString() + "/" + creds.Username
	expiresAt := time.Now().Add(session.SessionDefaultExpTime)

	profId, err := h.accountInteractor.GetAccountProfessorId(creds.Username)
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	user := session.InitUserInfo(creds.Username, profId)
	session.Sessions[sessionToken] = session.InitSession(user, expiresAt)

	resBody := responsebodies.SessionToken{
		Token:  sessionToken,
		Expiry: expiresAt,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	// проверяем соответсвтвие типа содержимого запроса
	if headerContentTtype != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// декодируем тело запроса
	var creds requestbodies.SignUp
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := inputdata.CheckUsernameExists{
		Login: creds.Username,
	}

	usernameExists, err := h.accountInteractor.CheckUsernameExists(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	if usernameExists {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode("username already exists")
		return
	}

	signupInput := inputdata.SignUp{
		Login:         creds.Username,
		Password:      creds.Password,
		Name:          creds.Name,
		Surname:       creds.Surname,
		Middlename:    creds.Middlename,
		UniId:         creds.UniversityId,
		ScienceDegree: creds.ScienceDegree,
	}

	account, err := h.accountInteractor.SignUp(signupInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
			log.Printf("Ошибка при кодировании ответа: %v", err)
		}
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewString() + "/" + creds.Username
	expiresAt := time.Now().Add(session.SessionDefaultExpTime)

	user := session.InitUserInfo(account.Login, account.Id)
	session.Sessions[sessionToken] = session.InitSession(user, expiresAt)

	resBody := responsebodies.SessionToken{
		Token:  sessionToken,
		Expiry: expiresAt,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// remove the users session from the session map
	delete(session.Sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time

	resBody := responsebodies.SessionToken{
		Token:  sessionToken,
		Expiry: time.Now(),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}

func (h *AuthHandler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	// session is already validated by middleware
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	userSession := session.Sessions[sessionToken]

	newSessionToken := uuid.NewString() + "/" + userSession.GetUser().GetUsername()
	expiresAt := time.Now().Add(session.SessionDefaultExpTime)

	// Set the token in the session map, along with the user whom it represents
	session.Sessions[newSessionToken] = session.InitSession(
		userSession.GetUser(), expiresAt,
	)

	// Delete the older session token
	delete(session.Sessions, sessionToken)

	// Set the new token as the users `session_token` cookie
	resBody := responsebodies.SessionToken{
		Token:  sessionToken,
		Expiry: expiresAt,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}

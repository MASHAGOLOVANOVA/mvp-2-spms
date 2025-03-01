package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mvp-2-spms/services/manage-students/outputdata"
	requestbodies "mvp-2-spms/web_server/handlers/request-bodies"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server = "http://localhost:8080"

func setupAccount() string {
	reqb := requestbodies.Credentials{
		Username: "1",
		Password: "1",
	}
	result := makeRequest("POST", server+"/api/v1/auth/signin", reqb)
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Fatal(err)
	}

	var resp struct {
		Token string `json:"session_token"`
		Exp   string `json:"expires_at"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Token
}

func makeRequest(method, url string, body interface{}) *http.Response {
	requestBody, err := json.Marshal(body) //тело запроса
	if err != nil {
		log.Print(err)
		requestBody = nil
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody)) //создаем запрос
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Print(err)
		return nil
	}
	result, err := http.DefaultClient.Do(request) //выолняем запрос
	if err != nil {
		log.Printf("клиент: произошла ошибка выполнения запроса: %s\n", err)
		return nil
	}

	return result //результат отправляем обратно
}

func TestProject_AddProject(t *testing.T) {

	t.Run("fail, StatusUnsupportedMediaType", func(t *testing.T) {
		s := setupAccount()

		b := requestbodies.AddStudent{
			Name:                   "1123",
			Surname:                "123",
			Middlename:             "13",
			Cource:                 1,
			EducationalProgrammeId: 1,
		}
		requestBody, err := json.Marshal(b) //тело запроса
		if err != nil {
			log.Print(err)
		}
		request, err := http.NewRequest("POST", server+"/api/v1/students/add", bytes.NewBuffer(requestBody)) //создаем запрос
		if err != nil {
			log.Print(err)
		}
		request.Header.Add("Session-Id", s)
		request.Header.Add("Content-Type", "application/json")
		result, err := http.DefaultClient.Do(request) // выполняем запрос
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		defer result.Body.Close()

		body, err := io.ReadAll(result.Body)
		if err != nil {
			log.Fatal(err)
		}

		var resp outputdata.AddStudent
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Fatal(err)
		}

		b1 := requestbodies.AddProject{
			Theme:          "123",
			StudentId:      resp.Id,
			Year:           2024,
			RepoOwner:      "123",
			RepositoryName: "123",
		}

		requestBody, err = json.Marshal(b1) //тело запроса
		if err != nil {
			log.Print(err)
		}
		request, err = http.NewRequest("POST", server+"/api/v1/projects/add", bytes.NewBuffer(requestBody)) //создаем запрос
		if err != nil {
			log.Print(err)
		}
		request.Header.Add("Session-Id", s)
		result, err = http.DefaultClient.Do(request) // выполняем запрос
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		defer result.Body.Close()

		assert.Equal(t, result.StatusCode, http.StatusUnsupportedMediaType)

	})

	t.Run("fail, bad request", func(t *testing.T) {
		s := setupAccount()

		b := requestbodies.AddStudent{
			Name:                   "1123",
			Surname:                "123",
			Middlename:             "13",
			Cource:                 1,
			EducationalProgrammeId: 1,
		}
		requestBody, err := json.Marshal(b) //тело запроса
		if err != nil {
			log.Print(err)
		}
		request, err := http.NewRequest("POST", server+"/api/v1/students/add", bytes.NewBuffer(requestBody)) //создаем запрос
		if err != nil {
			log.Print(err)
		}
		request.Header.Add("Session-Id", s)
		request.Header.Add("Content-Type", "application/json")
		result, err := http.DefaultClient.Do(request) // выполняем запрос
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		defer result.Body.Close()

		body, err := io.ReadAll(result.Body)
		if err != nil {
			log.Fatal(err)
		}

		var resp outputdata.AddStudent
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Fatal(err)
		}

		b1 := requestbodies.AddStudent{
			Name:                   "1123",
			Surname:                "123",
			Middlename:             "13",
			Cource:                 1,
			EducationalProgrammeId: 1,
		}

		requestBody, err = json.Marshal(b1) //тело запроса
		if err != nil {
			log.Print(err)
		}
		request, err = http.NewRequest("POST", server+"/api/v1/projects/add", bytes.NewBuffer(requestBody)) //создаем запрос
		if err != nil {
			log.Print(err)
		}
		request.Header.Add("Session-Id", s)
		request.Header.Add("Content-Type", "application/json")
		result, err = http.DefaultClient.Do(request) // выполняем запрос
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		defer result.Body.Close()

		assert.Equal(t, result.StatusCode, http.StatusBadRequest)

	})

	t.Run("ok", func(t *testing.T) {
		s := setupAccount()

		b := requestbodies.AddStudent{
			Name:                   "1123",
			Surname:                "123",
			Middlename:             "13",
			Cource:                 1,
			EducationalProgrammeId: 1,
		}
		requestBody, err := json.Marshal(b) //тело запроса
		if err != nil {
			log.Print(err)
		}
		request, err := http.NewRequest("POST", server+"/api/v1/students/add", bytes.NewBuffer(requestBody)) //создаем запрос
		if err != nil {
			log.Print(err)
		}
		request.Header.Add("Session-Id", s)
		request.Header.Add("Content-Type", "application/json")
		result, err := http.DefaultClient.Do(request) // выполняем запрос
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		defer result.Body.Close()

		body, err := io.ReadAll(result.Body)
		if err != nil {
			log.Fatal(err)
		}

		var resp outputdata.AddStudent
		err = json.Unmarshal(body, &resp)
		if err != nil {
			log.Fatal(err)
		}

		b1 := requestbodies.AddProject{
			Theme:          "123",
			StudentId:      resp.Id,
			Year:           2024,
			RepoOwner:      "123",
			RepositoryName: "123",
		}

		requestBody, err = json.Marshal(b1) //тело запроса
		if err != nil {
			log.Print(err)
		}
		request, err = http.NewRequest("POST", server+"/api/v1/projects/add", bytes.NewBuffer(requestBody)) //создаем запрос
		if err != nil {
			log.Print(err)
		}
		request.Header.Add("Session-Id", s)
		request.Header.Add("Content-Type", "application/json")
		result, _ = http.DefaultClient.Do(request)   //выолняем запрос
		result, err = http.DefaultClient.Do(request) // выполняем запрос
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		defer result.Body.Close()

		assert.Equal(t, result.StatusCode, http.StatusOK)
	})
}

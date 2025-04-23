// yandexdisc/api.go
package yandexdisc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	entities "mvp-2-spms/domain-aggregate"
	yandexapi "mvp-2-spms/integrations/yandex-api"
	"mvp-2-spms/services/models"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"golang.org/x/oauth2"
)

type YandexDisk struct {
	api *yandexapi.YandexAPI
}

func NewYandexDisk(api *yandexapi.YandexAPI) *YandexDisk {
	return &YandexDisk{api: api}
}

// Authentificate реализует интерфейс ICloudDrive
// Authentificate обновленная версия
func (d *YandexDisk) Authentificate(token *oauth2.Token) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}

	// Убедимся, что клиент создан с актуальным токеном
	err := d.api.SetupClient(token)
	if err != nil {
		return fmt.Errorf("failed to setup client: %w", err)
	}

	// Проверим работоспособность токена
	_, err = d.GetFolderInfo("disk:/")
	if err != nil {
		return fmt.Errorf("authentication check failed: %w", err)
	}

	return nil
}

// CreateFolder создает папку на Яндекс Диске
func (y *YandexDisk) CreateFolder(folderPath string) error {
	url := "https://cloud-api.yandex.net/v1/disk/resources?path=" + url.QueryEscape(folderPath)

	// Создаем новый запрос
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Добавляем заголовок авторизации вручную
	req.Header.Set("Authorization", "OAuth "+y.api.Token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetFolderInfo возвращает информацию о папке
func (y *YandexDisk) GetFolderInfo(folderPath string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://cloud-api.yandex.net/v1/disk/resources?path="+url.QueryEscape(folderPath), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Добавляем заголовок авторизации вручную
	req.Header.Set("Authorization", "OAuth "+y.api.Token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// UploadFile загружает файл на Яндекс Диск
func (y *YandexDisk) UploadFile(filePath string, fileContent []byte, overwrite bool) error {
	// Формируем URL для загрузки файла
	url := "https://cloud-api.yandex.net/v1/disk/resources/upload?path=" + url.QueryEscape(filePath) + "&overwrite=" + strconv.FormatBool(overwrite)

	// Создаем новый запрос для получения ссылки на загрузку
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Добавляем заголовок авторизации
	req.Header.Set("Authorization", y.api.Token.AccessToken)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s, response: %s", resp.Status, string(bodyBytes))
	}

	// Парсим URL для загрузки
	var uploadInfo struct {
		Href string `json:"href"`
	}
	if err := json.Unmarshal(bodyBytes, &uploadInfo); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Загружаем файл
	req, err = http.NewRequest("PUT", uploadInfo.Href, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	// Добавляем заголовок авторизации
	req.Header.Set("Authorization", y.api.Token.AccessToken)

	// Выполняем запрос на загрузку файла
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ = io.ReadAll(resp.Body) // ignore error for simplicity
		return fmt.Errorf("upload error: %s, response: %s", resp.Status, string(bodyBytes))
	}

	// Логирование успешной загрузки
	fmt.Printf("File uploaded successfully to: %s\n", filePath)
	return nil
}

// GetAuthLink возвращает URL для авторизации
func (d *YandexDisk) GetAuthLink(redirectURI string, accountId int, returnURL string) (string, error) {
	statestr := base64.URLEncoding.EncodeToString([]byte(fmt.Sprint(accountId, ",", returnURL)))
	return d.api.GetAuthLink(statestr), nil
}

// AddProjectFolder создает папку проекта
func (d *YandexDisk) AddProjectFolder(project entities.Project, driveInfo models.CloudDriveIntegration) (models.DriveProject, error) {
	folderName := fmt.Sprintf("Project %s_%s", project.Id, project.Theme)
	folderPath := path.Join(driveInfo.BaseFolderId, folderName)

	if err := d.CreateFolder(folderPath); err != nil {
		return models.DriveProject{}, fmt.Errorf("failed to create project folder: %w", err)
	}

	return models.DriveProject{
		Project: project,
		DriveFolder: models.DriveFolder{
			Id:   folderPath,
			Link: fmt.Sprintf("https://disk.yandex.ru/client/disk/%s", url.PathEscape(folderPath)),
		},
	}, nil
}

// AddProfessorBaseFolder создает базовую папку преподавателя
func (d *YandexDisk) AddProfessorBaseFolder() (models.DriveData, error) {
	baseName := "Student Project Management System"
	folderName := baseName

	for i := 1; i <= 10; i++ {
		_, err := d.GetFolderInfo(folderName)
		if err != nil {
			// Папка не существует, создаем
			if err := d.CreateFolder(folderName); err != nil {
				return models.DriveData{}, err
			}
			return models.DriveData{BaseFolderId: folderName}, nil
		}
		folderName = fmt.Sprintf("%s (%d)", baseName, i)
	}

	return models.DriveData{}, fmt.Errorf("failed to create unique folder")
}

// GetFolderNameById возвращает имя папки
func (d *YandexDisk) GetFolderNameById(id string) (string, error) {
	info, err := d.GetFolderInfo(id)
	if err != nil {
		return "", err
	}

	if name, ok := info["name"].(string); ok {
		return name, nil
	}
	return "", fmt.Errorf("invalid folder info")
}

// AddTaskToDrive добавляет задание в папку проекта
func (d *YandexDisk) AddTaskToDrive(task entities.Task, projectFolderId string) (models.DriveTask, error) {
	folderName := fmt.Sprintf("Task %s_%s until %s",
		task.Id, task.Name, task.Deadline.Format("02.01.2006"))
	folderPath := path.Join(projectFolderId, folderName)

	if err := d.CreateFolder(folderPath); err != nil {
		return models.DriveTask{}, err
	}

	fileName := fmt.Sprintf("Task '%s' description.txt", task.Name)
	filePath := path.Join(folderPath, fileName)
	content := fmt.Sprintf("%s\n\n%s", task.Name, task.Description)

	if err := d.UploadFile(filePath, []byte(content), true); err != nil {
		return models.DriveTask{}, err
	}

	return models.DriveTask{
		Task: task,
		DriveFolder: models.DriveFolder{
			Id:   folderPath,
			Link: fmt.Sprintf("https://disk.yandex.ru/client/disk/%s", url.PathEscape(folderPath)),
		},
	}, nil
}

// GetToken реализует получение токена по коду авторизации
func (d *YandexDisk) GetToken(code string) (*oauth2.Token, error) {
	if d.api == nil {
		return nil, fmt.Errorf("YandexAPI not initialized")
	}

	token, err := d.api.ExchangeCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	return token, nil
}

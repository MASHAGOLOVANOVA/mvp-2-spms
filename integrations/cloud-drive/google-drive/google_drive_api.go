package clouddrive

import (
	"fmt"
	"log"
	googleapi "mvp-2-spms/integrations/google-api"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const DAYS_PERIOD = 7
const HOURS_IN_DAY = 24
const EVENT_DURATION_HOURS = 1

type googleDriveApi struct {
	googleapi.Google
	api *drive.Service
}

func InitDriveApi(googleAPI googleapi.GoogleAPI) googleDriveApi {
	d := googleDriveApi{
		Google: googleapi.InintGoogle(googleAPI),
	}
	return d
}

func (d *googleDriveApi) AuthentificateService(token *oauth2.Token) error {
	if err := d.Authentificate(token); err != nil {
		log.Printf("Ошибка при аутентификации CloudDrive: %v", err)
	}

	api, err := drive.NewService(d.GetContext(), option.WithHTTPClient(d.GetClient()))
	if err != nil {
		return err
	}

	d.api = api
	return nil
}

func (d *googleDriveApi) CreateFolder(folderName string, parentFolder ...string) (*drive.File, error) {
	fileMetadata := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  parentFolder,
	}
	file, err := d.api.Files.Create(fileMetadata).Fields("id", "webViewLink").Do()
	if err == nil {
		return file, nil
	}
	return nil, err
}

func (d *googleDriveApi) GetFolderById(folderId string) (*drive.File, error) {
	file, err := d.api.Files.Get(folderId).Fields("id", "webViewLink", "name").Do()
	if err == nil {
		return file, nil
	}
	return nil, err
}

func (d *googleDriveApi) GetFoldersByName(folderName string) (*drive.FileList, error) {
	file, err := d.api.Files.List().Q(fmt.Sprint("name='", folderName, "'")).Do()
	if err == nil {
		return file, nil
	}
	return nil, err
}

func (d *googleDriveApi) AddTextFileToFolder(fileName string, fileText string, parentFolderId string) (*drive.File, error) {
	fileMetadata := &drive.File{
		Name:     fileName,
		MimeType: "application/vnd.google-apps.document", // google document type
		Parents:  []string{parentFolderId},
	}

	r := strings.NewReader(fileText)
	file, err := d.api.Files.Create(fileMetadata).Media(r).Fields("id").Do()
	if err == nil {
		return file, nil
	}
	return nil, err
}

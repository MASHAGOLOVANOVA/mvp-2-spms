package main

import (
	"log"
	"mvp-2-spms/database"
	accountepository "mvp-2-spms/database/account-repository"
	meetingrepository "mvp-2-spms/database/meeting-repository"
	projectrepository "mvp-2-spms/database/project-repository"
	studentrepository "mvp-2-spms/database/student-repository"
	taskrepository "mvp-2-spms/database/task-repository"
	unirepository "mvp-2-spms/database/university-repository"
	googleDrive "mvp-2-spms/integrations/cloud-drive/google-drive"
	"mvp-2-spms/integrations/git-repository-hub/github"
	googleapi "mvp-2-spms/integrations/google-api"
	googleCalendar "mvp-2-spms/integrations/planner-service/google-calendar"
	"mvp-2-spms/internal"
	manageaccounts "mvp-2-spms/services/manage-accounts"
	managemeetings "mvp-2-spms/services/manage-meetings"
	manageprojects "mvp-2-spms/services/manage-projects"
	managestudents "mvp-2-spms/services/manage-students"
	managetasks "mvp-2-spms/services/manage-tasks"
	manageuniversities "mvp-2-spms/services/manage-universities"
	"mvp-2-spms/services/models"
	"mvp-2-spms/web_server/config"
	"mvp-2-spms/web_server/routes"
	"mvp-2-spms/web_server/session"
	"net/http"
	"os"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/drive/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	serverConfig, err := config.ReadConfigFromFile("/mvp-2-spms/web_server/cmd/web_app/server_config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = config.SetConfigEnvVars(serverConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	session.SetBotTokenFromJson("/mvp-2-spms/web_server/cmd/web_app/credentials_bot.json")
	dbConfig, err := database.ReadDBConfigFromFile("/mvp-2-spms/web_server/cmd/web_app/db_config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var gdb *gorm.DB
	gdb = nil
	var err1 error
	err1 = nil
	for i := 0; i < 20; i++ {
		gdb, err1 = gorm.Open(mysql.Open(dbConfig.ConnString), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: dbConfig.SingularTable, // use singular table name, table for `User` would be `user` with this option enabled
			},
		})
	}
	if err1 != nil {
		log.Fatal(err.Error())
	}

	db := database.InitDatabade(gdb)
	repos := internal.Repositories{
		Projects:     projectrepository.InitProjectRepository(*db),
		Students:     studentrepository.InitStudentRepository(*db),
		Universities: unirepository.InitUniversityRepository(*db),
		Meetings:     meetingrepository.InitMeetingRepository(*db),
		Accounts:     accountepository.InitAccountRepository(*db),
		Tasks:        taskrepository.InitTaskRepository(*db),
	}

	repoHub := github.InitGithub(github.InitGithubAPI())

	// scopes := []string{calendar.CalendarScope, drive.DriveScope}

	gCalendarApi := googleCalendar.InitCalendarApi(googleapi.InitGoogleAPI(calendar.CalendarScope))
	gCalendar := googleCalendar.InitGoogleCalendar(gCalendarApi)
	gDriveApi := googleDrive.InitDriveApi(googleapi.InitGoogleAPI(drive.DriveScope))
	gDrive := googleDrive.InitGoogleDrive(gDriveApi)

	interactors := internal.Intercators{
		AccountManager:   manageaccounts.InitAccountInteractor(repos.Accounts, repos.Universities),
		ProjectManager:   manageprojects.InitProjectInteractor(repos.Projects, repos.Students, repos.Universities, repos.Accounts),
		StudentManager:   managestudents.InitStudentInteractor(repos.Students, repos.Projects, repos.Universities),
		MeetingManager:   managemeetings.InitMeetingInteractor(repos.Meetings, repos.Accounts, repos.Students, repos.Projects),
		TaskManager:      managetasks.InitTaskInteractor(repos.Projects, repos.Tasks, repos.Accounts),
		UnversityManager: manageuniversities.InitUniversityInteractor(repos.Universities),
	}

	integrations := internal.Integrations{
		GitRepositoryHubs: make(internal.GitRepositoryHubs),
		CloudDrives:       make(internal.CloudDrives),
		Planners:          make(internal.Planners),
	}

	integrations.Planners[models.GoogleCalendar] = gCalendar
	integrations.CloudDrives[models.GoogleDrive] = gDrive
	integrations.GitRepositoryHubs[models.GitHub] = repoHub

	app := internal.StudentsProjectsManagementApp{
		Intercators:  interactors,
		Integrations: integrations,
	}
	router := routes.SetupRouter(&app)
	if err := http.ListenAndServe(os.Getenv("SERVER_PORT"), router.Router()); err != nil {
		log.Printf("Ошибка при настройке env: %v", err)
	}
}

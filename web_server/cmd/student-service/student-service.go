package main

import (
	"log"
	"mvp-2-spms/database"
	"mvp-2-spms/internal"
	"mvp-2-spms/student-service/config"
	projectrepository "mvp-2-spms/student-service/database/project-repository"
	studentrepository "mvp-2-spms/student-service/database/student-repository"
	unirepository "mvp-2-spms/student-service/database/university-repository"
	managestudents "mvp-2-spms/student-service/manage-students"
	"mvp-2-spms/student-service/session"
	"net/http"
	"os"

	"mvp-2-spms/student-service/routes"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {

	serverConfig, err := config.ReadConfigFromFile("web_server/cmd/student-service/server_config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = config.SetConfigEnvVars(serverConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	session.SetBotTokenFromJson("web_server/cmd/student-service/credentials_bot.json")
	dbConfig, err := database.ReadDBConfigFromFile("web_server/cmd/student-service/db_config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var gdb *gorm.DB

	// Открываем соединение с базой данных
	gdb, err = gorm.Open(mysql.Open(dbConfig.ConnString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: dbConfig.SingularTable, // использовать единственное имя таблицы
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	db := database.InitDatabade(gdb)

	repos := internal.Repositories{
		Projects:     projectrepository.InitProjectRepository(*db),
		Students:     studentrepository.InitStudentRepository(*db),
		Universities: unirepository.InitUniversityRepository(*db),
	}

	interactors := internal.Intercators{
		StudentManager: managestudents.InitStudentInteractor(repos.Students, repos.Projects, repos.Universities),
	}
	integrations := internal.Integrations{}
	app := internal.StudentsProjectsManagementApp{
		Intercators:  interactors,
		Integrations: integrations,
	}

	router := routes.SetupRouter(&app)
	if err := http.ListenAndServe(os.Getenv("SERVER_PORT"), router.Router()); err != nil {
		log.Printf("Ошибка при настройке env: %v", err)
	}

}

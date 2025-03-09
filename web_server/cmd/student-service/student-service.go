package main

import (
	"encoding/json"
	"log"
	"mvp-2-spms/database"
	"mvp-2-spms/internal"
	"mvp-2-spms/services/manage-students/inputdata"
	"mvp-2-spms/student-service/config"
	projectrepository "mvp-2-spms/student-service/database/project-repository"
	studentrepository "mvp-2-spms/student-service/database/student-repository"
	unirepository "mvp-2-spms/student-service/database/university-repository"
	managestudents "mvp-2-spms/student-service/manage-students"
	"mvp-2-spms/student-service/session"
	"net/http"
	"os"

	"mvp-2-spms/student-service/routes"

	"github.com/streadway/amqp"
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

	//rabbitmq

	conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"student_queue", // имя очереди
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // аргументы
	)
	if err != nil {
		log.Fatalf("Ошибка объявления очереди: %s", err)
	}

	// ... остальной код для обработки сообщений ...
	go func() {
		for {
			msgs, err := ch.Consume(
				q.Name, // имя очереди
				"",     // имя потребителя
				true,   // auto-ack
				false,  // exclusive
				false,  // no-local
				false,  // no-wait
				nil,    // аргументы
			)
			if err != nil {
				log.Fatalf("Ошибка потребления: %s", err)
			}

			for d := range msgs {
				log.Printf("Получено сообщение: %s", d.Body) // Отладочное сообщение
				var input inputdata.AddStudent
				if err := json.Unmarshal(d.Body, &input); err != nil {
					log.Printf("Ошибка декодирования сообщения: %s", err)
					continue
				}

				// Обработка добавления студента
				studentID, err := app.Intercators.StudentManager.AddStudent(input)
				if err != nil {
					log.Printf("Ошибка добавления студента: %s", err)
					continue
				}

				log.Printf("Студент добавлен с ID: %d", studentID)
			}
		}
	}()

	router := routes.SetupRouter(&app)
	if err := http.ListenAndServe(os.Getenv("SERVER_PORT"), router.Router()); err != nil {
		log.Printf("Ошибка при настройке env: %v", err)
	}

}

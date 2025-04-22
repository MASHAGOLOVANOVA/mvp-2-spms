package routes

import (
	"mvp-2-spms/internal"
	"mvp-2-spms/web_server/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	router *chi.Mux
	app    *internal.StudentsProjectsManagementApp
}

func (r *Router) Router() *chi.Mux {
	return r.router
}

func InitRouter(app *internal.StudentsProjectsManagementApp) Router {
	r := chi.NewRouter()
	return Router{
		router: r,
		app:    app,
	}
}

func SetupRouter(app *internal.StudentsProjectsManagementApp) Router {
	r := InitRouter(app)
	r.SetupMiddleware()
	r.SetupRoutes()
	return r
}

// middleware for all routes
func (r *Router) SetupMiddleware() {
	r.router.Use(middleware.Logger)
	r.router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://127.0.0.1:3000", "http://127.0.0.1:3000", "https://localhost:3000", "http://localhost:3000", "http://localhost:3000", "https://spams-site.ru/"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Session-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func (r *Router) SetupRoutes() {
	r.router.Get("/api/v1/ping", handlers.Ping)
	r.router.With(handlers.Authentificator).Get("/api/v1/pingauth", handlers.Ping)
	r.setupMeetingRoutes()
	r.setupProjectRoutes()
	r.setupStudentRoutes()
	r.setupTaskRoutes()
	r.setupUniversityRoutes()
	r.setupAccountRoutes()
	r.setupAuthenRoutes()
	r.setupProfessorRoutes()
	r.setupApplicationRoutes()
}

func (r *Router) setupApplicationRoutes() {
	profH := handlers.InitProfessorHandler(r.app.Intercators.ProfessorManager)

	r.router.With(handlers.Authentificator).Route("/api/v1/applications", func(r chi.Router) {
		r.Get("/student", profH.GetStudentApplications)
		r.Get("/professor", profH.GetProfessorApplications)
		r.Put("/{applicationID}", profH.UpdateApplicationStatus)
	})

}

func (r *Router) setupProfessorRoutes() {
	profH := handlers.InitProfessorHandler(r.app.Intercators.ProfessorManager)

	r.router.With(handlers.Authentificator).Route("/api/v1/professors", func(r chi.Router) {
		r.Get("/", profH.GetProfessors)
		r.Route("/{professorID}", func(r chi.Router) {
			r.Post("/apply", profH.Apply)
		})
	})
}

func (r *Router) setupProjectRoutes() {
	projH := handlers.InitProjectHandler(r.app.Intercators.ProjectManager, r.app.Intercators.AccountManager, r.app.Integrations.CloudDrives, r.app.Integrations.GitRepositoryHubs)
	taskH := handlers.InitTaskHandler(r.app.Intercators.TaskManager, r.app.Intercators.AccountManager, r.app.Integrations.CloudDrives)

	// setup middleware for checking if professor is authorized and it's his projects?
	r.router.With(handlers.Authentificator).Route("/api/v1/projects", func(r chi.Router) {
		r.Get("/", projH.GetAllProfProjects)             // GET /projects
		r.Get("/statuslist", projH.GetProjectStatusList) // GET /projects/statuslist
		r.Get("/stagelist", projH.GetProjectStageList)   // GET /projects/stagelist
		r.Get("/filter", dummyHandler)                   // GET /projects/filter?student_id=1 query params are accessed with r.URL.Query().Get("student_id")
		r.Post("/add", projH.AddProject)                 // POST /projects/add
		// Subrouters:
		r.Route("/{projectID}", func(r chi.Router) {
			// r.Use(///) --> context (for handling not found errors for example)
			r.Get("/", projH.GetProject)                           // GET /projects/123
			r.Get("/statistics", projH.GetProjectStatistics)       // GET /projects/123/statistics
			r.Put("/", projH.UpdateProject)                        // PUT /projects/123
			r.Post("/supervisorreport", projH.GetProjectSupReport) // POST /projects/supervisorreport
			r.Put("/grades", projH.UpdateProjectGrading)           // PUT /projects/123/grades
			r.Delete("/", dummyHandler)                            // DELETE /projects/123
			r.Get("/commits", projH.GetProjectCommits)             // GET /projects/123/commits?from=2006-01-02T15:04:05.000Z
			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", taskH.GetAllProjectTasks) // GET /projects/123/tasks
				r.Post("/add", taskH.AddTask)        // POST /projects/123/tasks/add
				r.Route("/{taskID}", func(r chi.Router) {
					r.Put("/", taskH.UpdateTask)
				})
			})
		})
	})
}

func (r *Router) setupTaskRoutes() {
	taskH := handlers.InitTaskHandler(r.app.Intercators.TaskManager, r.app.Intercators.AccountManager, r.app.Integrations.CloudDrives)

	// setup middleware for checking if professor is authorized and it's his projects?
	r.router.With(handlers.Authentificator).Route("/api/v1/tasks", func(r chi.Router) {
		r.Get("/statuslist", taskH.GetTaskStatusList) // GET /tasks/statuslist
	})
}

func (r *Router) setupStudentRoutes() {
	studH := handlers.InitStudentHandler(r.app.Intercators.StudentManager)

	// setup middleware for checking if professor is authorized and it's his projects?
	r.router.With(handlers.Authentificator).Route("/api/v1/students", func(r chi.Router) {
		r.With().Get("/", studH.GetStudents) // GET /students with middleware (currently empty)
		r.Post("/add", studH.AddStudent)     // POST /students/add
	})
}

func (r *Router) setupAccountRoutes() {
	accH := handlers.InitAccountHandler(r.app.Intercators.AccountManager, r.app.Integrations.CloudDrives)

	// setup middleware for checking if professor is authorized and it's his projects?
	r.router.With(handlers.Authentificator).Route("/api/v1/account", func(r chi.Router) {
		r.Get("/", accH.GetAccountInfo)                     // GET /account/
		r.Get("/integrations", accH.GetAccountIntegrations) // GET /account/integrations
		r.Get("/student", accH.GetStudentAccountInfo)
	})
}

func (r *Router) setupUniversityRoutes() {
	uniH := handlers.InitUniversityHandler(r.app.Intercators.UnversityManager)

	// setup middleware for checking if professor is authorized and it's his projects?
	r.router.With(handlers.Authentificator).Route("/api/v1/universities", func(r chi.Router) {
		r.With().Get("/", dummyHandler) // GET /universities with middleware (currently empty)
		r.Route("/{uniID}", func(r chi.Router) {
			r.Get("/", dummyHandler) // GET /universities/123
			r.Route("/edprogrammes", func(r chi.Router) {
				r.Get("/", uniH.GetAllUniEdProgrammes) // GET /universities/123/edprogrammes
				r.Post("/add", dummyHandler)           // POST /universities/123/edprogrammes/add
			})
		})
	})
}

func (r *Router) setupAuthenRoutes() {
	calendarH := handlers.InitPlannerIntegrationHandler(r.app.Integrations.Planners, r.app.Intercators.AccountManager)
	driveH := handlers.InitCloudDriveHandler(r.app.Integrations.CloudDrives, r.app.Intercators.AccountManager)
	repoHubH := handlers.InitGitRepoHandler(r.app.Integrations.GitRepositoryHubs, r.app.Intercators.AccountManager)
	authH := handlers.InitAuthHandler(r.app.Intercators.AccountManager)

	r.router.Route("/api/v1/auth", func(r chi.Router) {
		r.With(handlers.BotAuthentificator).Route("/bot", func(r chi.Router) {
			r.Post("/signinuser", authH.SignInBot)
			r.Post("/signinstudent", authH.SignInStudent)
			r.Post("/signupuser", authH.SignUp)
			r.Post("/signupstudent", authH.StudentSignUp)
		})
		r.Route("/integration", func(r chi.Router) {
			r.With(handlers.Authentificator).Get("/getplanners", calendarH.GetProfessorPlanners)
			r.With(handlers.Authentificator).Post("/setplanner", calendarH.SetProfessorPlanner)
			r.With(handlers.Authentificator).Route("/authlink", func(r chi.Router) {
				r.Get("/googlecalendar", calendarH.GetGoogleCalendarLink) // GET /auth/integration/authlink/googlecalendar
				r.Get("/googledrive", driveH.GetGoogleDriveLink)          // GET /auth/integration/authlink/googledrive
				r.Get("/yandexdisk", driveH.GetYandexDiskLink)
				r.Get("/github", repoHubH.GetGitHubLink) // GET /auth/integration/authlink/github
			})
			r.Route("/access", func(r chi.Router) {
				r.Get("/googlecalendar", calendarH.OAuthCallbackGoogleCalendar) // GET /auth/integration/access/googlecalendar
				r.Get("/googledrive", driveH.OAuthCallbackGoogleDrive)          // GET /auth/integration/access/googledrive
				r.Get("/yandexdisk", driveH.OAuthCallbackYandexDisk)
				r.Get("/github", repoHubH.OAuthCallbackGitHub) // GET /auth/integration/access/github
			})
		})
		r.Post("/signin", authH.SignIn)                                                // POST /auth/signin
		r.Post("/signup", authH.SignUp)                                                // POST /auth/signup
		r.Post("/signout", authH.SignOut)                                              // POST /auth/signout
		r.With(handlers.Authentificator).Post("/refreshsession", authH.RefreshSession) // POST /auth/refreshsession
	})
}

func (r *Router) setupMeetingRoutes() {
	meetH := handlers.InitMeetingHandler(r.app.Intercators.MeetingManager, r.app.Intercators.AccountManager, r.app.Integrations.Planners)
	// RESTy routes for "meetings" resource
	// setup middleware for checking professor?
	r.router.With(handlers.Authentificator).Route("/api/v1/meetings", func(r chi.Router) {
		r.Get("/", meetH.GetProfessorMeetings)           // GET /meetings?from=2006-01-02T15:04:05.000Z
		r.Get("/statuslist", meetH.GetMeetingStatusList) // GET /meetings/statuslist
		r.Get("/filter", dummyHandler)                   // GET /meetings/filter?student_id=1&status=planned query params are accessed with r.URL.Query().Get("student_id")
		r.Post("/add", meetH.AddMeeting)                 // POST /meetings/add
		// Subrouters:
		r.Route("/{meetingID}", func(r chi.Router) {
			// r.Use(///) --> context (for handling not found errors for example)
			r.Get("/", dummyHandler)    // GET /meetings/123
			r.Put("/", dummyHandler)    // PUT /meetings/123
			r.Delete("/", dummyHandler) // DELETE /meetings/123
		})
	})
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

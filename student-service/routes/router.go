package routes

import (
	"mvp-2-spms/internal"
	"mvp-2-spms/web_server/handlers"

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
		AllowedOrigins:   []string{"https://127.0.0.1:3000", "http://127.0.0.1:3000", "https://localhost:3000", "http://localhost:3000", "http://localhost:3000", "https://spams-site.ru/"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Session-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (r *Router) SetupRoutes() {
	r.router.Get("/api/v1/ping", handlers.Ping)
	r.router.With(handlers.Authentificator).Get("/api/v1/pingauth", handlers.Ping)
	r.setupStudentRoutes()
}

func (r *Router) setupStudentRoutes() {
	studH := handlers.InitStudentHandler(r.app.Intercators.StudentManager)

	// setup middleware for checking if professor is authorized and it's his projects?
	r.router.With(handlers.Authentificator).Route("/api/v1/students", func(r chi.Router) {
		r.With().Get("/", studH.GetStudents) // GET /students with middleware (currently empty)
		r.Post("/add", studH.AddStudent)     // POST /students/add
	})
}

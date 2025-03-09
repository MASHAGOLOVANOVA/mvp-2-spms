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
		AllowedOrigins:   []string{"https://127.0.0.1:3000", "http://127.0.0.1:3000", "https://localhost:3000", "http://localhost:3000", "http://localhost:3000", "https://spams-site.ru/"},
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
	r.setupUniversityRoutes()
}

func (r *Router) setupUniversityRoutes() {
	uniH := handlers.InitUniversityHandler(r.app.Intercators.UnversityManager)

	r.router.With(handlers.Authentificator).Route("/api/v1/universities", func(r chi.Router) {
		r.With().Get("/", dummyHandler)
		r.Route("/{uniID}", func(r chi.Router) {
			r.Get("/", dummyHandler)
			r.Route("/edprogrammes", func(r chi.Router) {
				r.Get("/", uniH.GetAllUniEdProgrammes)
				r.Post("/add", dummyHandler)
			})
		})
	})
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

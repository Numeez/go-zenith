package routes

import (
	"github.com/Numeez/go-zenith/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)
		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkOutHandler.HandleGetWorkOutById))
		r.Post("/workouts", app.Middleware.RequireUser(app.WorkOutHandler.HandleCreateWorkOut))
		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkOutHandler.HandlerUpdateWorkoutById))
		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkOutHandler.HandlerDeleteWorkout))

	})
	router.Get("/health", app.HealthCheck)
	router.Post("/users", app.UserHandler.HandlerRegisterUser)
	router.Post("/tokens/authentication", app.TokenHandler.HandlerCreateToken)
	return router
}

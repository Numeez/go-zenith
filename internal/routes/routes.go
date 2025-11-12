package routes

import (
	"github.com/Numeez/go-zenith/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/health", app.HealthCheck)
	router.Get("/workouts/{id}", app.WorkOutHandler.HandleGetWorkOutById)
	router.Post("/workouts", app.WorkOutHandler.HandleCreateWorkOut)
	router.Put("/workouts/{id}", app.WorkOutHandler.HandlerUpdateWorkoutById)
	router.Delete("/workouts/{id}", app.WorkOutHandler.HandlerDeleteWorkout)
	return router
}

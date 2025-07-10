package routes

import (
	"github.com/aliSHAIKH001/api_go/internal/app"
	"github.com/go-chi/chi/v5"
)

// Simple router that links to different handlers which originate from the API folder
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutById)

	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)

	return r
}
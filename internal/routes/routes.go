package routes

import (
	"github.com/aliSHAIKH001/api_go/internal/app"
	"github.com/go-chi/chi/v5"
)

// Simple router that links to different handlers which originate from the API folder
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// This middleware will be applied to all routes in this group.
		// It will set the user , either anonymous or authenticated, in the request context.
		r.Use(app.Middleware.Authenticate)

		// These routes are all wrapped in requireUser middleware to check if the user is authenticated.
		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutById))
		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))
	})

	r.Get("/health", app.HealthCheck)
	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return r
}

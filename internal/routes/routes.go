package routes

import (
	"github.com/makhammatovb/femProject/internal/app"
	"github.com/go-chi/chi/v5"
)

// SetupRoutes sets up the routes for the application using chi router
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts/", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}/", app.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workouts/{id}/", app.WorkoutHandler.HandleDeleteWorkout)

	r.Get("/users/{id}", app.UserHandler.HandleGetUserByID) // checked
	r.Post("/users/", app.UserHandler.HandleRegisterUser) // checked
	r.Put("/users/{id}/", app.UserHandler.HandleUpdateUser) // checked
	r.Delete("/users/{id}/", app.UserHandler.HandleDeleteUser) // checked

	// tokens
	r.Post("/tokens/", app.TokenHandler.HandleCreateToken)
	return r
}

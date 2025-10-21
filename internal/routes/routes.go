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
	return r
}

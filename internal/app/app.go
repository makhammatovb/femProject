package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/makhammatovb/femProject/internal/api"
)

// Application struct includes logger and handler from api package
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

// NewApplication creates a new instance of Application
func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize handlers from api package, creates a new instance of WorkoutHandler and returns pointer to it
	workoutHandler := api.NewWorkoutHandler()
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}
	return app, nil
}

// HealthCheck is a simple handler to check the health of the application
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available")
}

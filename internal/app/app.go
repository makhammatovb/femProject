package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/makhammatovb/femProject/internal/api"
	"github.com/makhammatovb/femProject/internal/store"
	"github.com/makhammatovb/femProject/migrations"
)

// Application struct includes logger and handler from api package
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

// NewApplication creates a new instance of Application
func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(pgDB)

	// Initialize handlers from api package, creates a new instance of WorkoutHandler and returns pointer to it
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}
	return app, nil
}

// HealthCheck is a simple handler to check the health of the application
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available")
}

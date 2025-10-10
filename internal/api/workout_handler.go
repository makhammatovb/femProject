package api

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"strconv"
	"fmt"
)

// WorkoutHandler struct to handle workout-related requests for future use
type WorkoutHandler struct {
}

// NewWorkoutHandler creates a new instance of WorkoutHandler.
func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
}

// HandleGetWorkoutByID handles the GET request to retrieve a workout by its ID.
func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	
	// retrieves the workout ID from the URL parameters
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	// checks the validity of the workout ID
	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Workout ID: %d", workoutID)

}

// HandleCreateWorkout handles the POST request to create a new workout.
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Create Workout")
}

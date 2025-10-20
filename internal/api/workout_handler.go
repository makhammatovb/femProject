package api

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/makhammatovb/femProject/internal/store"
)

// WorkoutHandler struct to handle workout-related requests for future use
type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

// NewWorkoutHandler creates a new instance of WorkoutHandler.
func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
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
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workout)
}

// HandleCreateWorkout handles the POST request to create a new workout.
func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout 
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
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
	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	var updatedWorkoutRequest struct {
		Title           *string               `json:"title"`
		Description     *string 	          `json:"description"`
		DurationMinutes *int            	  `json:"duration_minutes"`
		CaloriesBurned  *int                  `json:"calories_burned"`
		Entries         []store.WorkoutEntry  `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedWorkoutRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if updatedWorkoutRequest.Title != nil {
		existingWorkout.Title = *updatedWorkoutRequest.Title
	}
	if updatedWorkoutRequest.Description != nil {
		existingWorkout.Description = *updatedWorkoutRequest.Description
	}
	if updatedWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkoutRequest.DurationMinutes
	}
	if updatedWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkoutRequest.CaloriesBurned
	}
	if updatedWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updatedWorkoutRequest.Entries
	}
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}

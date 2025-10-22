package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost port=5433 user=postgres password=postgres dbname=test_postgres sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	_, err = db.Exec("TRUNCATE TABLE workouts, workout_entries CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name   string
		workout *Workout
		wantErr bool
	}{
		{
			name: "Valid Workout",
			workout: &Workout{
				Title:           "Morning Routine",
				Description:     "My morning routine",
				DurationMinutes: 30,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName:    "Push Ups",
						Reps:            IntPtr(15),
						Sets:            3,
						Weight:          FloatPtr(10.0),
						DurationSeconds: IntPtr(30),
						Notes:           "Do 3 sets of 15 push ups",
						OrderIndex:      1,
					},
					{
						ExerciseName:    "Running",
						Reps:            nil,
						Sets:            1,
						Weight:          nil,
						DurationSeconds: IntPtr(30),
						Notes:           "Run for 30 seconds",
						OrderIndex:      2,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "Invalid Workout",
			workout: &Workout{
				Title:           "Morning Routine",
				Description:     "My morning routine",
				DurationMinutes: 55,
				CaloriesBurned:  544,
				Entries: []WorkoutEntry{
					{
						ExerciseName:    "Push Ups",
						Reps:            IntPtr(15),
						Sets:            3,
						Notes:           "Do 3 sets of 15 push ups",
						OrderIndex:      1,
					},
					{
						ExerciseName:    "Running",
						Reps:            IntPtr(12),
						Sets:            4,
						DurationSeconds: IntPtr(60),
						Notes:           "Run for 60 seconds",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err, "Expected error but got none")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Reps, retrieved.Entries[i].Reps)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].Weight, retrieved.Entries[i].Weight)
				assert.Equal(t, tt.workout.Entries[i].DurationSeconds, retrieved.Entries[i].DurationSeconds)
				assert.Equal(t, tt.workout.Entries[i].Notes, retrieved.Entries[i].Notes)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}


func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float64) *float64 {
	return &f
}

package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type WorkoutEntry struct {
	ID              int       `json:"id"`
	ExerciseName    string    `json:"exercise_name"`
	Reps            *int      `json:"reps"`
	Sets            int       `json:"sets"`
	Weight          *float64  `json:"weight"`
	DurationSeconds *int      `json:"duration_seconds"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	OrderIndex      int       `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(workout *Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query :=
		`INSERT INTO workouts (title, description, duration_minutes, calories_burned)
	VALUES ($1, $2, $3, $4) RETURNING id;
	`
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}
	for i, entry := range workout.Entries {
		fmt.Printf("DEBUG: Inserting entry %d: %+v\n", i, entry)
		query :=
			`INSERT INTO workout_entries (workout_id, exercise_name, reps, sets, weight, duration_seconds, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;
		`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Reps, entry.Sets, entry.Weight, entry.DurationSeconds, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `
	SELECT id, title, description, duration_minutes, calories_burned, created_at, updated_at
	FROM workouts WHERE id = $1;
	`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned, &workout.CreatedAt, &workout.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workout with ID %d not found", id)
	}
	entriesQuery := `
	SELECT id, exercise_name, reps, sets, weight, duration_seconds, notes, order_index, created_at, updated_at
	FROM workout_entries WHERE workout_id = $1 ORDER BY order_index;
	`
	rows, err := pg.db.Query(entriesQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry WorkoutEntry
		err := rows.Scan(&entry.ID, &entry.ExerciseName, &entry.Reps, &entry.Sets, &entry.Weight, &entry.DurationSeconds, &entry.Notes, &entry.OrderIndex, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}
	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
	UPDATE workouts SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4, updated_at = NOW()
	WHERE id = $5;
	`
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("workout with ID %d not found", workout.ID)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM workout_entries WHERE workout_id = $1;`, workout.ID)
	if err != nil {
		return err
	}
	for _, entry := range workout.Entries {
		query :=
			`INSERT INTO workout_entries (workout_id, exercise_name, reps, sets, weight, duration_seconds, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
		`
		_, err := tx.Exec(query, workout.ID, entry.ExerciseName, entry.Reps, entry.Sets, entry.Weight, entry.DurationSeconds, entry.Notes, entry.OrderIndex)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `
	DELETE FROM workouts WHERE id = $1;
	`
	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("workout with ID %d not found", id)
	}
	return nil
}

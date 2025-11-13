package store

import (
	"database/sql"
)

type Workout struct {
	Id              int            `json:"id"`
	UserId          int            `json:"user_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	Id              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Reps            *int     `json:"reps"`
	Sets            int      `json:"sets"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"`
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkout struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkout {
	return &PostgresWorkout{
		db: db,
	}
}

type WorkoutStore interface {
	CreateWorkout(workout *Workout) (*Workout, error)
	GetWorkOutById(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
	GetWorkoutOwner(id int64) (int, error)
}

func (pg *PostgresWorkout) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	defer func() {
		_ = tx.Rollback()
	}()
	if err != nil {
		return nil, err
	}
	query := `
	INSERT INTO workouts(user_id,title,description,duration_minutes,calories_burned)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id
	`
	err = tx.QueryRow(query, workout.UserId, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.Id)
	if err != nil {
		return nil, err
	}
	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id,exercise_name,sets,reps,duration_seconds,weight,notes,order_index)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id
		`
		err := tx.QueryRow(query, workout.Id, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.Id)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return workout, nil
}

func (pg *PostgresWorkout) GetWorkOutById(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `
	SELECT id,title,description,duration_minutes,calories_burned 
	 from workouts 
	  WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&workout.Id, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	entryQuery := `
  SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
  FROM workout_entries
  WHERE workout_id = $1
  ORDER BY order_index
  `
	entries, err := pg.db.Query(entryQuery, id)
	defer func() {
		_ = entries.Close()
	}()
	if err != nil {
		return nil, err
	}
	for entries.Next() {
		var entry WorkoutEntry
		if err := entries.Scan(
			&entry.Id,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.DurationSeconds,
			&entry.Weight,
			&entry.Notes,
			&entry.OrderIndex,
		); err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}
func (pg *PostgresWorkout) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	defer func() {
		_ = tx.Rollback()
	}()
	if err != nil {
		return err
	}
	query := `
	UPDATE workouts
	SET title=$1,description=$2,duration_minutes=$3,calories_burned=$4
	WHERE id=$5
	`
	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.Id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.Exec("DELETE from workout_entries WHERE id=$1", workout.Id)
	if err != nil {
		return err
	}
	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id,exercise_name,sets,reps,duration_seconds,weight,notes,order_index)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		`
		_, err := tx.Exec(query, workout.Id, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex)
		if err != nil {
			return err
		}

	}

	return nil
}

func (pg *PostgresWorkout) DeleteWorkout(id int64) error {
	query := `DELETE FROM workouts WHERE id=$1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	affectedRow, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if affectedRow == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (pg *PostgresWorkout) GetWorkoutOwner(workoutId int64) (int, error) {
	var userID int
	query := `
	SELECT user_id
	FROM workouts
	WHERE id = $1
	`
	err := pg.db.QueryRow(query, workoutId).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

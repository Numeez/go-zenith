package store

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the DB %v", err)
	}
	fmt.Println("<<<<<<< Connected to TEST Database >>>>>>")
	if err := Migrate(db, "../../migrations/"); err != nil {
		t.Fatalf("Failed to migrate to the DB %v", err)
	}
	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("Failed to truncate  tables in the DB %v", err)
	}
	return db

}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	store := NewPostgresWorkoutStore(db)
	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Push Day",
				Description:     "Upper Body",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						Sets:         3,
						Reps:         IntPtr(10),
						Notes:        "Warm up properly",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title:           "full body",
				Description:     "complete workout",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Notes:        "keep form",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "squats",
						Sets:            4,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(60),
						Weight:          FloatPtr(185.0),
						Notes:           "full depth",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(test.workout)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.Title, test.workout.Title)
			assert.Equal(t, createdWorkout.Description, test.workout.Description)
			assert.Equal(t, createdWorkout.DurationMinutes, test.workout.DurationMinutes)
			assert.Equal(t, createdWorkout.CaloriesBurned, test.workout.CaloriesBurned)
			retrieved, err := store.GetWorkOutById(int64(createdWorkout.Id))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.Id, test.workout.Id)
			assert.Equal(t, len(retrieved.Entries), len(test.workout.Entries))

			for i, entry := range retrieved.Entries {
				assert.Equal(t, test.workout.Entries[i].ExerciseName, entry.ExerciseName)
				assert.Equal(t, test.workout.Entries[i].Sets, entry.Sets)
				assert.Equal(t, test.workout.Entries[i].OrderIndex, entry.OrderIndex)
			}

		})
	}

}

func IntPtr(value int) *int {
	return &value
}

func FloatPtr(value float64) *float64 {
	return &value
}

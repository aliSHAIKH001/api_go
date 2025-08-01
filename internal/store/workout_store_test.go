package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// run migrations for our test db
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating tables error: %v", err)
	}
	
	return db
}

func TestCreateWorkout(t *testing.T) {

	// Setting up a test Db while clearing out the old volumes, we dont want data to persist..
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	// Slice of scenerios we need to test
	tests := []struct {
		name string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title: "Push day",
				Description: "Upper body day",
				DurationMinutes: 60,
				CaloriesBurned: 200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench press",
						Sets: 3, 
						Reps: IntPtr(10),
						Weight: FloatPtr(135.5),
						Notes: "Warm up properly",
						OrderIndex: 1,
					},

				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entries",
			workout: &Workout{
				Title: "full body",
				Description: "Complete Workout",
				DurationMinutes: 90,
				CaloriesBurned: 500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets: 3,
						Reps: IntPtr(60),
						Notes: "keep Form",
						OrderIndex: 1,
					},
					{
						ExerciseName: "Squats",
						Sets: 4,
						Reps: IntPtr(12),
						DurationSeconds: IntPtr(60),
						Weight: FloatPtr(185.0),
						Notes: "full depth",
						OrderIndex: 2,
					},
				},
			},
			wantErr: true,
		},
	}

	// Slice being tested.
	for _, tt := range tests {

		t.Run(tt.name, func (t *testing.T){
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			// Checking the fields are as expected
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)
			
			// The createdWorkout should exist in the database
			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)
			
			// We have to check if the values when this created workout was fetched are consistent
			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			// Now we have to check the actual individual entries
			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}


		})

	}

}

// Because & does not work on constants.
func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
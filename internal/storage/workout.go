package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type WorkoutStorage struct {
	DB *sql.DB
}

func NewWorkoutStorage(db *sql.DB) *WorkoutStorage {
	return &WorkoutStorage{
		DB: db,
	}
}

func (s *WorkoutStorage) InsertWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
	engine := getQueryEngine(ctx, s.DB)
	var id uuid.UUID
	err := engine.QueryRowContext(ctx, `INSERT INTO workouts (total_time, grade_of_training, date_of_training, athlete_id, status)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`, workout.TotalTime, workout.GradeOfTraining, workout.DateOfTraining, workout.AthleteID, workout.Status).Scan(&id)

	return id, err
}

func (s *WorkoutStorage) InsertSetsBulk(ctx context.Context, workoutID uuid.UUID, sets []domain.Set) error {
	if len(sets) == 0 {
		return nil
	}

	engine := getQueryEngine(ctx, s.DB)

	const columnsPerRowSets = 6

	valueStrings := make([]string, 0, len(sets))
	valueArgs := make([]any, 0, len(sets)*columnsPerRowSets)
	plchldrCount := 1
	var sb strings.Builder

	for _, set := range sets {
		sb.WriteString("(")
		for j := 0; j < columnsPerRowSets; j++ {
			fmt.Fprintf(&sb, "$%d", plchldrCount+j)
			if j < columnsPerRowSets-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
		valueStrings = append(valueStrings, sb.String())
		sb.Reset()

		valueArgs = append(valueArgs, workoutID, set.ExerciseID, set.Weight, set.SetOrder, set.Reps, set.Rpe)
		plchldrCount += columnsPerRowSets
	}

	query := fmt.Sprintf(`INSERT INTO sets (workout_id, exercise_id, weight, set_order, reps, rpe) VALUES %s`,
		strings.Join(valueStrings, ", "))

	_, err := engine.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert sets: %w", err)
	}
	return nil
}

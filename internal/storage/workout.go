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

func (s *WorkoutStorage) CreateWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	var id uuid.UUID
	err = tx.QueryRowContext(ctx, `INSERT INTO workouts (total_time, grade_of_training, date_of_training, athlete_id, status)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		workout.TotalTime, workout.GradeOfTraining, workout.DateOfTraining, workout.AthleteID, workout.Status).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	if len(workout.Sets) == 0 {
		return id, nil
	}

	exerciseIDs := make([]int, 0, len(workout.Sets))
	seen := make(map[int]struct{})
	for _, set := range workout.Sets {
		if _, ok := seen[set.ExerciseID]; !ok {
			seen[set.ExerciseID] = struct{}{}
			exerciseIDs = append(exerciseIDs, set.ExerciseID)
		}
	}

	ids := make([]interface{}, len(exerciseIDs))
	for i, id := range exerciseIDs {
		ids[i] = id
	}

	query := `SELECT COUNT(*)
	FROM exercises WHERE id = ANY($1)
	AND (athlete_id IS NULL OR athlete_id = $2)`
	var validCount int
	if err = tx.QueryRowContext(ctx, query, ids, workout.AthleteID).Scan(&validCount); err != nil {
		return uuid.Nil, fmt.Errorf("failed to check exercise permission: %w", err)
	}

	if validCount != len(exerciseIDs) {
		return uuid.Nil, fmt.Errorf("some exercise don't belong to you")
	}

	const columnsPerRowSets = 6

	valueStrings := make([]string, 0, len(workout.Sets))
	valueArgs := make([]any, 0, len(workout.Sets)*columnsPerRowSets)
	//PostgreSQL использует placeholder с $1
	plchldrCount := 1
	var sb strings.Builder
	for _, set := range workout.Sets {

		sb.WriteString("(")
		for j := 0; j < columnsPerRowSets; j++ {
			sb.WriteString(fmt.Sprintf("$%d", plchldrCount+j))
			if j < columnsPerRowSets-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(")")
		valueStrings = append(valueStrings, sb.String())
		sb.Reset()

		valueArgs = append(valueArgs, id, set.ExerciseID, set.Weight, set.SetOrder, set.Reps, set.Rpe)
		plchldrCount += columnsPerRowSets
	}

	query = fmt.Sprintf(`INSERT INTO sets (workout_id, exercise_id, weight, set_order, reps, rpe) VALUES %s`,
		strings.Join(valueStrings, ", "))
	_, err = tx.ExecContext(ctx, query, valueArgs...)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to bulk insert sets: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit tx: %w", err)
	}
	return id, nil
}

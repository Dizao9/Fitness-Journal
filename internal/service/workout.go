package service

import (
	"context"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type WorkoutStorage interface {
	InsertWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error)
	InsertSetsBulk(ctx context.Context, workoutID uuid.UUID, sets []domain.Set) error
}

type ExerciseWorkStorage interface {
	GetExerciseByIDs(ctx context.Context, userID uuid.UUID, IDs []int) ([]domain.Exercise, error)
}

type TxManager interface {
	Do(ctx context.Context, execFn func(ctx context.Context) error) error
}
type WorkoutService struct {
	WorkoutStore  WorkoutStorage
	ExerciseStore ExerciseWorkStorage
	TxManager     TxManager
}

func NewWorkoutService(ws WorkoutStorage, es ExerciseWorkStorage, tm TxManager) *WorkoutService {
	return &WorkoutService{
		WorkoutStore:  ws,
		ExerciseStore: es,
		TxManager:     tm,
	}
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, input CreateWorkoutInput) (uuid.UUID, error) {
	if input.GradeOfTraining != nil && (*input.GradeOfTraining < 0 || *input.GradeOfTraining > 10) {
		return uuid.Nil, domain.ErrInvalidGradeOfTraining
	}

	if input.DateOfTraining.After(time.Now()) {
		return uuid.Nil, domain.DateOfTrainingInFuture
	}

	if input.Status == domain.WorkoutStatusFinished && len(input.Sets) == 0 {
		return uuid.Nil, domain.ErrNoSetsInFinishedWorkout
	}

	workoutDomain := domain.Workout{
		AthleteID:       input.AthleteID,
		TotalTime:       input.TotalTime,
		GradeOfTraining: input.GradeOfTraining,
		DateOfTraining:  input.DateOfTraining,
		Status:          input.Status,
	}

	expectedExercises := len(input.Sets) / 3
	if expectedExercises < 8 {
		expectedExercises = 8
	}

	seen := make(map[int]struct{}, expectedExercises)
	exerciseIDs := make([]int, 0, expectedExercises)

	workoutDomain.Sets = make([]domain.Set, len(input.Sets))
	for i, s := range input.Sets {
		workoutDomain.Sets[i] = domain.Set{
			ExerciseID: s.ExerciseID,
			Weight:     s.Weight,
			SetOrder:   s.SetOrder,
			Reps:       s.Reps,
			Rpe:        s.Rpe,
		}
		if _, exists := seen[s.ExerciseID]; !exists {
			seen[s.ExerciseID] = struct{}{}
			exerciseIDs = append(exerciseIDs, s.ExerciseID)
		}
	}

	var createdWorkoutID uuid.UUID

	err := s.TxManager.Do(ctx, func(txCtx context.Context) error {
		if len(exerciseIDs) > 0 {
			foundExercises, err := s.ExerciseStore.GetExerciseByIDs(txCtx, workoutDomain.AthleteID, exerciseIDs)

			if err != nil {
				return err
			}

			if len(foundExercises) != len(exerciseIDs) {
				return domain.ErrNotEnoughPermission
			}
		}

		var err error
		createdWorkoutID, err = s.WorkoutStore.InsertWorkout(txCtx, workoutDomain)
		if err != nil {
			return err
		}

		if len(workoutDomain.Sets) > 0 {
			err = s.WorkoutStore.InsertSetsBulk(txCtx, createdWorkoutID, workoutDomain.Sets)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return createdWorkoutID, nil
}

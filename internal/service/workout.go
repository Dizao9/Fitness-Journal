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
	GetWorkouts(ctx context.Context, athleteID uuid.UUID, limit, offset int) ([]domain.Workout, error)
	GetWorkoutByID(ctx context.Context, workoutID uuid.UUID) (domain.Workout, error)
	CreateSet(ctx context.Context, workoutID uuid.UUID, set domain.Set) (uuid.UUID, error)
	DeleteWorkout(ctx context.Context, workoutID uuid.UUID) error
	UpdateWorkout(ctx context.Context, workout domain.Workout) error
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
		return uuid.Nil, domain.ErrDateOfTrainingInFuture
	}

	if !input.Status.IsStatusValid() {
		return uuid.Nil, domain.ErrInvalidWorkoutStatus
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

func (s *WorkoutService) GetWorkouts(ctx context.Context, input GetWorkoutsInput) ([]domain.Workout, error) {
	if input.Page < 1 {
		input.Page = 1
	}

	if input.Limit <= 0 {
		input.Limit = 10
	}

	if input.Limit > 50 {
		input.Limit = 50
	}

	offset := ((input.Page - 1) * input.Limit)

	return s.WorkoutStore.GetWorkouts(ctx, input.AthleteID, input.Limit, offset)
}

func (s *WorkoutService) GetWorkout(ctx context.Context, input GetWorkoutInput) (domain.Workout, error) {
	workout, err := s.WorkoutStore.GetWorkoutByID(ctx, input.WorkoutID)
	if err != nil {
		return domain.Workout{}, err
	}

	if workout.AthleteID != input.AthleteID {
		return domain.Workout{}, domain.ErrNotEnoughPermission
	}

	return workout, nil
}

func (s *WorkoutService) CreateSet(ctx context.Context, input CreateSetInputAtom) (uuid.UUID, error) {
	if input.Reps < 0 {
		return uuid.Nil, domain.ErrInvalidRepsCount
	}

	if input.Weight < 0 {
		return uuid.Nil, domain.ErrInvalidWeight
	}

	if input.Rpe != nil && (*input.Rpe <= 0 || *input.Rpe > 10) {
		return uuid.Nil, domain.ErrInvalidRpe
	}

	workout, err := s.WorkoutStore.GetWorkoutByID(ctx, input.WorkoutID)
	if err != nil {
		return uuid.Nil, err
	}

	if workout.AthleteID != input.AthleteID {
		return uuid.Nil, domain.ErrNotEnoughPermission
	}

	if workout.Status != domain.WorkoutStatusInProgress {
		return uuid.Nil, domain.ErrWorkoutAlreadyFinished
	}

	exercise, err := s.ExerciseStore.GetExerciseByIDs(ctx, input.AthleteID, []int{input.ExerciseID})
	if err != nil {
		return uuid.Nil, err
	}

	if len(exercise) != 1 {
		return uuid.Nil, domain.ErrNotEnoughPermission
	}

	set := domain.Set{
		ExerciseID: input.ExerciseID,
		WorkoutID:  input.WorkoutID,
		Weight:     input.Weight,
		SetOrder:   input.SetOrder,
		Rpe:        input.Rpe,
		Reps:       input.Reps,
	}

	return s.WorkoutStore.CreateSet(ctx, input.WorkoutID, set)
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, input DeleteWorkoutInput) error {
	workout, err := s.WorkoutStore.GetWorkoutByID(ctx, input.WorkoutID)
	if err != nil {
		return err
	}

	if workout.AthleteID != input.AthleteID {
		return domain.ErrNotEnoughPermission
	}

	err = s.WorkoutStore.DeleteWorkout(ctx, input.WorkoutID)
	return err
}

func (s *WorkoutService) FinishWorkout(ctx context.Context, input FinishWorkoutInput) error {
	if input.TotalTime != nil && *input.TotalTime < 0 {
		return domain.ErrInvalidTotalTime
	}

	if input.GradeOfTraining != nil && (*input.GradeOfTraining < 0 || *input.GradeOfTraining > 10) {
		return domain.ErrInvalidGradeOfTraining
	}

	workout, err := s.WorkoutStore.GetWorkoutByID(ctx, input.WorkoutID)
	if err != nil {
		return err
	}

	if workout.AthleteID != input.AthleteID {
		return domain.ErrNotEnoughPermission
	}

	if workout.Status == domain.WorkoutStatusFinished {
		return domain.ErrWorkoutAlreadyFinished
	}

	if len(workout.Sets) == 0 {
		return domain.ErrNoSetsInFinishedWorkout
	}

	if input.GradeOfTraining != nil {
		workout.GradeOfTraining = input.GradeOfTraining
	}

	if input.TotalTime != nil {
		workout.TotalTime = input.TotalTime
	}

	workout.Status = domain.WorkoutStatusFinished

	err = s.WorkoutStore.UpdateWorkout(ctx, workout)
	return err
}

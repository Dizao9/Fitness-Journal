package service

import (
	"context"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type WorkoutStorage interface {
	CreateWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error)
}

type WorkoutService struct {
	Store WorkoutStorage
}

func NewWorkoutService(s WorkoutStorage) *WorkoutService {
	return &WorkoutService{
		Store: s,
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

	workoutDomain.Sets = make([]domain.Set, len(input.Sets))
	for i, s := range input.Sets {
		workoutDomain.Sets[i] = domain.Set{
			ExerciseID: s.ExerciseID,
			Weight:     s.Weight,
			SetOrder:   s.SetOrder,
			Reps:       s.Reps,
			Rpe:        s.Rpe,
		}
	}

	return s.Store.CreateWorkout(ctx, workoutDomain)
}

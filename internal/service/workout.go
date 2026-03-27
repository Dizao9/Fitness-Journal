package service

import (
	"context"

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

func (s *WorkoutService) CreateWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
	return s.Store.CreateWorkout(ctx, workout)
}

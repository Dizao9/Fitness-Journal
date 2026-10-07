package service

import (
	"context"

	"github.com/Dizao9/Fitness-Journal/internal/config"
	"github.com/Dizao9/Fitness-Journal/internal/storage"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
type Services struct {
	AuthService     *AuthService
	ExerciseService *ExerciseService
	AthleteService  *AthleteService
	WorkoutService  *WorkoutService
	TxManager       TxManager
}

func NewServices(storage *storage.Storage, config *config.Config, tm TokenManager, ph PasswordHasher, txManager TxManager) *Services {
	return &Services{
		AuthService:     NewAuthService(storage.Athlete, tm, ph),
		ExerciseService: NewExerciseService(storage.Exercise),
		AthleteService:  NewAthleteService(storage.Athlete, config),
		WorkoutService:  NewWorkoutService(storage.Workout, storage.Exercise, txManager),
	}
}

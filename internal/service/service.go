package service

import (
	"github.com/Dizao9/Fitness-Journal/internal/config"
	"github.com/Dizao9/Fitness-Journal/internal/storage"
)

type Services struct {
	AuthService     *AuthService
	ExerciseService *ExerciseService
	AthleteService  *AthleteService
	WorkoutService  *WorkoutService
}

func NewServices(storage *storage.Storage, config *config.Config) *Services {
	return &Services{
		AuthService:     NewAuthService(storage.Athlete, config),
		ExerciseService: NewExerciseService(storage.Exercise),
		AthleteService:  NewAthleteService(storage.Athlete, config),
		WorkoutService:  NewWorkoutService(storage.Workout),
	}
}

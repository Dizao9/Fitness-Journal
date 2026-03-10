package transport

import "github.com/Dizao9/Fitness-Journal/internal/service"

type Handlers struct {
	Auth     *AuthHandler
	Athlete  *AthleteHandler
	Exercise *ExerciseHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:     NewAuthHandler(services.AuthService),
		Athlete:  NewAthleteHandler(services.AthleteService),
		Exercise: NewExerciseHandler(services.ExerciseService),
	}
}

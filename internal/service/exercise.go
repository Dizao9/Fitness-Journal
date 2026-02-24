package service

import (
	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
)

type ExerciseStorage interface {
	CreateExercise(exercise domain.Exercise) (int, error)
}

type ExerciseService struct {
	Store ExerciseStorage
}

func NewExerciseService(s ExerciseStorage) *ExerciseService {
	return &ExerciseService{
		Store: s,
	}
}

func (s *ExerciseService) CreateCustomExercise(exercise dto.CreateExerciseRequestDTO) (int, error) {
	newExercise := domain.Exercise{
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Description: exercise.Description,
		AthleteID:   exercise.AthleteID,
	}
	id, err := s.Store.CreateExercise(newExercise)
	if err != nil {
		return 0, err
	}
	return id, nil
}

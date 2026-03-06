package service

import (
	"context"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/google/uuid"
)

type ExerciseStorage interface {
	CreateExercise(exercise domain.Exercise) (int, error)
	GetExercises(ctx context.Context, userID uuid.UUID, filter string, limit, offset int) ([]domain.Exercise, error)
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

func (s *ExerciseService) GetPageOfExercise(ctx context.Context, userID uuid.UUID, filter string, limit, page int) ([]dto.ExerciseForPageDTO, error) {
	offset := (page - 1) * limit

	ListOfExercises, err := s.Store.GetExercises(ctx, userID, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	DTOExercises := make([]dto.ExerciseForPageDTO, len(ListOfExercises))
	for i, v := range ListOfExercises {
		isOwner := false
		if v.AthleteID != nil {
			isOwner = (*v.AthleteID == userID)
		}
		DTOExercises[i] = dto.ExerciseForPageDTO{
			ID:          v.ID,
			Name:        v.Name,
			MuscleGroup: v.MuscleGroup,
			IsOwner:     isOwner,
			IsSystem:    v.AthleteID == nil,
		}
	}

	return DTOExercises, nil
}

// func (s *ExerciseService) GetExerciseByID(userIDStr string, exerciseID int) (domain.Exercise, error)

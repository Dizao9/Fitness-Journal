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
	GetExerciseByID(userID uuid.UUID, exerciseID int) (domain.Exercise, error)
	DeleteExercise(exerciseID int) error
	UpdateExercise(exercise domain.Exercise) error
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

func (s *ExerciseService) GetExerciseByID(userID uuid.UUID, exerciseID int) (dto.ExerciseDTO, error) {
	exercise, err := s.Store.GetExerciseByID(userID, exerciseID)
	if err != nil {
		return dto.ExerciseDTO{}, err
	}

	isOwner := false
	if exercise.AthleteID != nil {
		isOwner = (*exercise.AthleteID == userID)
	}
	return dto.ExerciseDTO{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Description: exercise.Description,
		IsOwner:     isOwner,
	}, nil
}

func (s *ExerciseService) DeleteExerciseByID(userID uuid.UUID, exerciseID int) error {
	exercise, err := s.Store.GetExerciseByID(userID, exerciseID)
	if err != nil {
		return err
	}

	if exercise.AthleteID == nil {
		return domain.ErrNotEnoughPermission
	}

	return s.Store.DeleteExercise(exerciseID)
}

func (s *ExerciseService) UpdateExercise(userID uuid.UUID, exerciseID int, exerciseUPD dto.ExerciseUpdateReqDTO) error {
	exercise, err := s.Store.GetExerciseByID(userID, exerciseID)
	if err != nil {
		return err
	}

	if exerciseUPD.Name != nil {
		exercise.Name = *exerciseUPD.Name
	}
	if exerciseUPD.MuscleGroup != nil {
		exercise.MuscleGroup = *exerciseUPD.MuscleGroup
	}
	if exerciseUPD.Description != nil {
		exercise.Description = *&exerciseUPD.Description
	}

	return s.Store.UpdateExercise(exercise)
}

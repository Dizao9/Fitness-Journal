package dto

import "github.com/google/uuid"

type CreateExerciseRequestDTO struct {
	Name        string     `json:"name"`
	MuscleGroup string     `json:"muscle_group"`
	Description *string    `json:"description"`
	AthleteID   *uuid.UUID `json:"-"`
}

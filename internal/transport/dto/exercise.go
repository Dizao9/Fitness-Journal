package dto

import "github.com/google/uuid"

type CreateExerciseRequestDTO struct {
	Name        string     `json:"name"`
	MuscleGroup string     `json:"muscle_group"`
	Description *string    `json:"description"`
	AthleteID   *uuid.UUID `json:"-"`
}

type ExerciseForPageDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	MuscleGroup string `json:"muscle_group"`
	IsOwner     bool   `json:"is_owner"`
	IsSystem    bool   `json:"is_system"`
}

type ExerciseDTO struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	MuscleGroup string  `json:"muscle_group"`
	Description *string `json:"description"`
	IsOwner     bool    `json:"is_owner"`
}

type ExerciseUpdateReqDTO struct {
	Name        *string `json:"name"`
	MuscleGroup *string `json:"muscle_group"`
	Description *string `json:"description"`
}

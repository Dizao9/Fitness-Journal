package service

import (
	"time"

	"github.com/google/uuid"
)

type AthleteProfileOutput struct {
	ID           uuid.UUID
	Name         string
	Age          int
	Username     string
	CurrentCycle string
	Gender       string
	Email        string
}

type UpdateProfileInput struct {
	Name         *string
	Age          *int
	Username     *string
	CurrentCycle *string
}

type CreateExerciseInput struct {
	Name        string
	MuscleGroup string
	Description *string
	AthleteID   *uuid.UUID
}

type ExercisePageOutput struct {
	ID          int
	Name        string
	MuscleGroup string
	IsOwner     bool
	IsSystem    bool
}

type ExerciseOutput struct {
	ID          int
	Name        string
	MuscleGroup string
	Description *string
	IsOwner     bool
}

type UpdateExerciseInput struct {
	Name        *string
	MuscleGroup *string
	Description *string
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
	Name     string
	Age      int
}

type LoginInput struct {
	Email    string
	Password string
}

type CreateWorkoutInput struct {
	AthleteID       uuid.UUID
	TotalTime       *int
	GradeOfTraining *int
	DateOfTraining  time.Time
	Status          string
	Sets            []CreateSetInput
}

type CreateSetInput struct {
	ExerciseID int
	Weight     float32
	SetOrder   *int
	Reps       int
	Rpe        *int
}

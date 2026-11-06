package domain

import "github.com/google/uuid"

type Set struct {
	ID           uuid.UUID
	ExerciseName string
	ExerciseID   int
	WorkoutID    uuid.UUID
	Weight       float32
	SetOrder     *int
	Reps         int
	Rpe          *int
}

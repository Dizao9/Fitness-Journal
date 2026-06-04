package domain

import (
	"github.com/google/uuid"
)

type Exercise struct {
	ID          int
	Name        string
	MuscleGroup string
	Description *string
	AthleteID   *uuid.UUID
}

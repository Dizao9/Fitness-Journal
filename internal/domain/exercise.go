package domain

import (
	"github.com/google/uuid"
)

type Exercise struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	MuscleGroup string     `json:"muscle_group"`
	Description *string    `json:"description"`
	AthleteID   *uuid.UUID `json:"athlete_id"`
}

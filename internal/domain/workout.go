package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutStatus string

const WorkoutStatusInProgress WorkoutStatus = "in_progress"
const WorkoutStatusFinished WorkoutStatus = "finished"

func (s WorkoutStatus) IsStatusValid() bool {
	return s == WorkoutStatusInProgress || s == WorkoutStatusFinished
}

type Workout struct {
	ID              uuid.UUID
	TotalTime       *int
	GradeOfTraining *int
	DateOfTraining  time.Time
	AthleteID       uuid.UUID
	Status          WorkoutStatus

	Sets []Set
}

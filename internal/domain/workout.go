package domain

import (
	"time"

	"github.com/google/uuid"
)

const WorkoutStatusInProgress string = "in_progress"
const WorkoutStatusFinished string = "finished"

type Workout struct {
	ID              uuid.UUID
	TotalTime       int
	GradeOfTraining *int
	DateOfTraining  time.Time
	AthleteID       uuid.UUID
	Status          string

	Sets []Set
}

package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateSetReq struct {
	ExerciseID int     `json:"exercise_id"`
	Weight     float32 `json:"weight"`
	SetOrder   *int    `json:"set_order"`
	Reps       int     `json:"reps"`
	Rpe        *int    `json:"rpe"`
}

type CreateWorkoutReq struct {
	TotalTime       *int           `json:"total_time"`
	GradeOfTraining *int           `json:"grade_of_training"`
	DateOfTraining  time.Time      `json:"date_of_training"`
	Status          string         `json:"status"`
	Sets            []CreateSetReq `json:"sets"`
}

type WorkoutSummaryResponse struct {
	ID        uuid.UUID `json:"id"`
	Date      time.Time `json:"date_of_training"`
	Grade     *int      `json:"grade_of_training"`
	Status    string    `json:"status"`
	TotalTime *int      `json:"total_time"`
}

type WorkoutSetResponse struct {
	ID           uuid.UUID `json:"id"`
	ExerciseID   int       `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	Weight       float32   `json:"weight"`
	SetOrder     *int      `json:"set_order"`
	Reps         int       `json:"reps"`
	Rpe          *int      `json:"rpe"`
}
type WorkoutExpandedResponse struct {
	ID        uuid.UUID            `json:"id"`
	Date      time.Time            `json:"date_of_training"`
	Grade     *int                 `json:"grade_of_training"`
	Status    string               `json:"status"`
	TotalTime *int                 `json:"total_time"`
	Sets      []WorkoutSetResponse `json:"sets"`
}

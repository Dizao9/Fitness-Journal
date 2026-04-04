package dto

import "time"

type Set struct {
	ExerciseID int     `json:"exercise_id"`
	Weight     float32 `json:"weight"`
	SetOrder   *int    `json:"set_order"`
	Reps       int     `json:"reps"`
	Rpe        *int    `json:"rpe"`
}

type CreateWorkoutReq struct {
	TotalTime       *int      `json:"total_time"`
	GradeOfTraining *int      `json:"grade_of_training"`
	DateOfTraining  time.Time `json:"date_of_training"`
	Status          string    `json:"status"`
	Sets            []Set     `json:"sets"`
}

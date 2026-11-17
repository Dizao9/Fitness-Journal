package domain

import (
	"errors"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrExerciseNotFound = errors.New("exercise not found")
var ErrNotEnoughPermission = errors.New("not enough permission")

// workout
var ErrInvalidGradeOfTraining = errors.New("invalid grade of training")
var ErrNoSetsInFinishedWorkout = errors.New("finished workout must have at least one set")
var ErrDateOfTrainingInFuture = errors.New("date of training must be in the past")
var ErrInvalidWorkoutStatus = errors.New("invalid status for training")
var ErrWorkoutNotFound = errors.New("workout not found")
var ErrWorkoutAlreadyFinished = errors.New("can't add set in finished workout")
var ErrInvalidRepsCount = errors.New("reps count cannot be negative")
var ErrInvalidRpe = errors.New("rpe should be from 1 to 10")
var ErrInvalidWeight = errors.New("weight cannot be negative")
var ErrInvalidTotalTime = errors.New("total_time cannot be negative")

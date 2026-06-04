package domain

import (
	"errors"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrExerciseNotFound = errors.New("exercise not found")
var ErrNotEnoughPermission = errors.New("exercise not found")

// workout
var ErrInvalidGradeOfTraining = errors.New("invalid grade of training")
var ErrNoSetsInFinishedWorkout = errors.New("finished workout must have at least one set")
var DateOfTrainingInFuture = errors.New("date of training must be in the past")

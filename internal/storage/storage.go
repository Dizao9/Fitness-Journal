package storage

import "database/sql"

type Storage struct {
	Athlete  *AthleteStorage
	Exercise *ExerciseStorage
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Athlete:  NewAthleteStorage(db),
		Exercise: NewExerciseStorage(db),
	}
}

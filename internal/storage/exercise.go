package storage

import (
	"database/sql"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
)

type ExercStorage struct {
	DB *sql.DB
}

func (s *ExercStorage) CreateExercise(exercise domain.Exercise) (int, error) {
	var id int
	err := s.DB.QueryRow(`INSERT INTO exercises 
	(name, muscle_group, description, athlete_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id`, exercise.Name, exercise.MuscleGroup, exercise.Description, exercise.AthleteID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

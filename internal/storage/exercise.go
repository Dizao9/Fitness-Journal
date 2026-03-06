package storage

import (
	"context"
	"database/sql"
	"log"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type ExerciseStorage struct {
	DB *sql.DB
}

func NewExerciseStorage(db *sql.DB) *ExerciseStorage {
	return &ExerciseStorage{DB: db}
}

func (s *ExerciseStorage) CreateExercise(exercise domain.Exercise) (int, error) {
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

func (s *ExerciseStorage) GetExercises(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Exercise, error) {
	rows, err := s.DB.Query(`SELECT * FROM exercises WHERE
	athlete_id IS NULL
	OR
	athlete_id = $1
	ORDER BY name ASC
	LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	exercises := make([]domain.Exercise, 0, 50)
	for rows.Next() {
		var exercise domain.Exercise
		err = rows.Scan(&exercise.ID, &exercise.Name, &exercise.MuscleGroup, &exercise.Description, &exercise.AthleteID)
		exercises = append(exercises, exercise)
	}
	if err != nil {
		log.Printf("[GET_EXERCISES_STORAGE] error: %v", err)
	}

	return exercises, nil
}

package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type ExerciseStorage struct {
	DB *sql.DB
}

func NewExerciseStorage(db *sql.DB) *ExerciseStorage {
	return &ExerciseStorage{DB: db}
}

func ValidateExerciseNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrExerciseNotFound
	}
	return err
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

func (s *ExerciseStorage) GetExercises(ctx context.Context, userID uuid.UUID, filter string, limit, offset int) ([]domain.Exercise, error) {
	var query string
	var args []interface{}

	baseQuery := "SELECT id, name, muscle_group, description, athlete_id FROM exercises WHERE "

	switch filter {
	case "my":
		query = baseQuery + "athlete_id = $1"
		args = append(args, userID)
	case "system":
		query = baseQuery + "athlete_id IS NULL"
	default:
		query = baseQuery + "athlete_id IS NULL OR athlete_id = $1"
		args = append(args, userID)
	}

	query += " ORDER BY name ASC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)
	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exercises := make([]domain.Exercise, 0, 50)
	for rows.Next() {
		var exercise domain.Exercise
		err = rows.Scan(&exercise.ID, &exercise.Name, &exercise.MuscleGroup, &exercise.Description, &exercise.AthleteID)
		if err != nil {
			log.Printf("[GET_EXERCISES_STORAGE] error: %v", err)
		}

		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (s *ExerciseStorage) GetExerciseByID(userID uuid.UUID, exerciseID int) (domain.Exercise, error) {
	var exercise domain.Exercise

	err := s.DB.QueryRow(`SELECT id, name, muscle_group, description, athlete_id FROM exercises WHERE ID = $1
	AND (athlete_id IS NULL OR athlete_id = $2)`, exerciseID, userID).
		Scan(&exercise.ID, &exercise.Name, &exercise.MuscleGroup, &exercise.Description, &exercise.AthleteID)
	if err != nil {
		return exercise, ValidateExerciseNotFound(err)
	}

	return exercise, nil
}

func (s *ExerciseStorage) DeleteExercise(exerciseID int) error {
	res, err := s.DB.Exec(`DELETE FROM exercises WHERE id = $1`, exerciseID)
	if err != nil {
		return ValidateExerciseNotFound(err)
	}

	countAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if countAff == 0 {
		return domain.ErrExerciseNotFound
	}

	return nil
}

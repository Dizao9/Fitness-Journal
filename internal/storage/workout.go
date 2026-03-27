package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type WorkoutStorage struct {
	DB *sql.DB
}

func NewWorkoutStorage(db *sql.DB) *WorkoutStorage {
	return &WorkoutStorage{
		DB: db,
	}
}

func (s *WorkoutStorage) CreateWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()
}

package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	DB *sql.DB
}

func (s *Storage) CreateAthlete(athlete domain.Athlete) (string, error) {
	var id string
	err := s.DB.QueryRow(`INSERT INTO athletes (username, email, name, age, password_hash, created_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, athlete.Username, athlete.Email, athlete.Name,
		athlete.Age, athlete.PasswordHash, athlete.CreatedAt).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", domain.ErrUserAlreadyExists
		}
		return "", fmt.Errorf("failed to create athlet: %w", err)
	}

	return id, nil
}

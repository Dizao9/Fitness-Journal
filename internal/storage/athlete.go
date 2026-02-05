package storage

import (
	"database/sql"
	"fmt"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
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
		return "", fmt.Errorf("failed to create athlet: %w", err)
	}

	return id, nil
}

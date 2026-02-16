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

func ValidateErrorUserAlreadyExists(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}

func (s *Storage) CreateAthlete(athlete domain.Athlete) (string, error) {
	var id string
	err := s.DB.QueryRow(`INSERT INTO athletes (username, email, name, age, password_hash, created_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, athlete.Username, athlete.Email, athlete.Name,
		athlete.Age, athlete.PasswordHash, athlete.CreatedAt).Scan(&id)
	if err != nil {
		if ValidateErrorUserAlreadyExists(err) {
			return "", domain.ErrUserAlreadyExists
		}
		return "", fmt.Errorf("failed to create athlet: %w", err)
	}

	return id, nil
}

func (s *Storage) GetByEmail(email string) (domain.Athlete, error) {
	var a domain.Athlete
	err := s.DB.QueryRow(`SELECT id, age, name, username, password_hash, current_cycle, created_at, email, gender, role FROM athletes WHERE email = $1`, email).
		Scan(&a.ID, &a.Age, &a.Name, &a.Username, &a.PasswordHash, &a.CurrentCycle, &a.CreatedAt, &a.Email, &a.Gender, &a.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return a, domain.ErrUserNotFound
		}
		return a, err
	}

	return a, nil
}

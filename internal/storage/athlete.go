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

func ValidateUserNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrUserNotFound
	}
	return err
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
		return a, ValidateUserNotFound(err)
	}

	return a, nil
}

func (s *Storage) GetByUserID(userID string) (domain.Athlete, error) {
	var a domain.Athlete
	err := s.DB.QueryRow("SELECT id, age, name, username, current_cycle, email, gender, role FROM athletes WHERE id = $1", userID).
		Scan(&a.ID, &a.Age, &a.Name, &a.Username, &a.CurrentCycle, &a.Email, &a.Gender, &a.Role)
	if err != nil {
		return a, ValidateUserNotFound(err)
	}

	return a, nil
}

func (s *Storage) UpdateUser(id string, a domain.Athlete) error {
	res, err := s.DB.Exec("UPDATE athletes SET name = $1, age = $2, username = $3, current_cycle = $4 WHERE id = $5", a.Name, a.Age, a.Username, a.CurrentCycle, id)
	if err != nil {
		return err
	}

	countAff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if countAff == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (s *Storage) DeleteUser(id string) error {
	res, err := s.DB.Exec("DELETE FROM athletes WHERE id = $1", id)
	if err != nil {
		return err
	}

	countAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if countAff == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type AthleteStorage struct {
	DB *sql.DB
}

func NewAthleteStorage(db *sql.DB) *AthleteStorage {
	return &AthleteStorage{DB: db}
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

func (s *AthleteStorage) CreateAthlete(athlete domain.Athlete) (uuid.UUID, error) {
	var idStr string
	err := s.DB.QueryRow(`INSERT INTO athletes (username, email, name, age, password_hash, created_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, athlete.Username, athlete.Email, athlete.Name,
		athlete.Age, athlete.PasswordHash, athlete.CreatedAt).Scan(&idStr)
	if err != nil {
		if ValidateErrorUserAlreadyExists(err) {
			return uuid.Nil, domain.ErrUserAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("failed to create athlet: %w", err)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return id, nil
	}

	return id, nil
}

func (s *AthleteStorage) GetByEmail(email string) (domain.Athlete, error) {
	var a domain.Athlete
	err := s.DB.QueryRow(`SELECT id, age, name, username, password_hash, current_cycle, created_at, email, gender, role FROM athletes WHERE email = $1`, email).
		Scan(&a.ID, &a.Age, &a.Name, &a.Username, &a.PasswordHash, &a.CurrentCycle, &a.CreatedAt, &a.Email, &a.Gender, &a.Role)
	if err != nil {
		return a, ValidateUserNotFound(err)
	}

	return a, nil
}

func (s *AthleteStorage) GetByUserID(userID uuid.UUID) (domain.Athlete, error) {
	var a domain.Athlete
	err := s.DB.QueryRow("SELECT id, age, name, username, current_cycle, email, gender, role FROM athletes WHERE id = $1", userID).
		Scan(&a.ID, &a.Age, &a.Name, &a.Username, &a.CurrentCycle, &a.Email, &a.Gender, &a.Role)
	if err != nil {
		return a, ValidateUserNotFound(err)
	}

	return a, nil
}

func (s *AthleteStorage) UpdateUser(id uuid.UUID, a domain.Athlete) error {
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

func (s *AthleteStorage) DeleteUser(id uuid.UUID) error {
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

func (s *AthleteStorage) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM athletes WHERE ID = $1)`
	err := s.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("db check failed: %v", err)
	}
	return exists, err
}

func (s *AthleteStorage) SaveRefreshToken(athleteID uuid.UUID, jti uuid.UUID, expiresAt time.Time, refreshToken string) error {
	query := `INSERT INTO refresh_tokens (athlete_id, jti, token, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := s.DB.Exec(query, athleteID, jti, refreshToken, expiresAt)
	return err
}

func (s *AthleteStorage) DeleteRefreshToken(jti uuid.UUID) (bool, error) {
	res, err := s.DB.Exec(`DELETE FROM refresh_tokens WHERE jti = $1`, jti)
	if err != nil {
		return false, err
	}
	countAff, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return countAff > 0, err
}


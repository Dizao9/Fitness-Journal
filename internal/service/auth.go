package service

import (
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"golang.org/x/crypto/bcrypt"
)

type AthleteStorage interface {
	CreateAthlete(athlete domain.Athlete) (string, error)
}

type AuthService struct {
	Store AthleteStorage
}

func (s *AuthService) Register(req dto.RegisterUser) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return "", err
	}

	u := domain.Athlete{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(bytes),
		Age:          req.Age,
		CreatedAt:    time.Now(),
		Name:         req.Name,
	}

	return s.Store.CreateAthlete(u)
}

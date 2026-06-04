package service

import (
	"github.com/Dizao9/Fitness-Journal/internal/config"
	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type AthleteStorage interface {
	GetByUserID(userID uuid.UUID) (domain.Athlete, error)
	UpdateUser(userID uuid.UUID, a domain.Athlete) error
	DeleteUser(userID uuid.UUID) error
}

type AthleteService struct {
	Store AthleteStorage
	Conf  *config.Config
}

func NewAthleteService(Str AthleteStorage, Cnf *config.Config) *AthleteService {
	return &AthleteService{
		Store: Str,
		Conf:  Cnf,
	}
}

func (a *AthleteService) GetProfile(userID uuid.UUID) (AthleteProfileOutput, error) {
	athlete, err := a.Store.GetByUserID(userID)
	if err != nil {
		return AthleteProfileOutput{}, err
	}

	resp := AthleteProfileOutput{
		ID:       athlete.ID,
		Username: athlete.Username,
		Email:    athlete.Email,
	}

	if athlete.Age != nil {
		resp.Age = *athlete.Age
	}
	if athlete.Name != nil {
		resp.Name = *athlete.Name
	}
	if athlete.Gender != nil {
		resp.Gender = *athlete.Gender
	}
	if athlete.CurrentCycle != nil {
		resp.CurrentCycle = *athlete.CurrentCycle
	}

	return resp, nil
}

func (s *AthleteService) UpdateProfile(userID uuid.UUID, upd UpdateProfileInput) error {
	athlete, err := s.Store.GetByUserID(userID)
	if err != nil {
		return err
	}

	if upd.Age != nil {
		athlete.Age = upd.Age
	}

	if upd.Name != nil {
		athlete.Name = upd.Name
	}

	if upd.CurrentCycle != nil {
		athlete.CurrentCycle = upd.CurrentCycle
	}

	if upd.Username != nil {
		athlete.Username = *upd.Username
	}

	return s.Store.UpdateUser(userID, athlete)
}

func (s *AthleteService) DeleteProfile(userID uuid.UUID) error {
	return s.Store.DeleteUser(userID)
}

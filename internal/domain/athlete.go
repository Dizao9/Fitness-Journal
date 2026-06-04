package domain

import (
	"time"

	"github.com/google/uuid"
)

type Athlete struct {
	ID           uuid.UUID
	Age          *int
	Name         *string
	Username     string
	PasswordHash string
	CurrentCycle *string
	CreatedAt    time.Time
	Email        string
	Gender       *string
	Role         *string
}

func (a *Athlete) GetRole() string {
	if a.Role == nil {
		return ""
	}
	return *a.Role
}

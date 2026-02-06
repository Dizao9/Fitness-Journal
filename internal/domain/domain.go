package domain

import (
	"time"
)

type Athlete struct {
	ID           string    `json:"id"`
	Age          int       `json:"age"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CurrentCycle string    `json:"current_cycle"`
	CreatedAt    time.Time `json:"created_at"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"`
	Role         string    `json:"role"`
}

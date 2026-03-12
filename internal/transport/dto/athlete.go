package dto

import "github.com/google/uuid"

type UserProfileResponseDTO struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Age          int       `json:"age"`
	Username     string    `json:"username"`
	CurrentCycle string    `json:"current_cycle,omitempty"`
	Gender       string    `json:"gender"`
	Email        string    `json:"email"`
}

type UpdateProfileRequest struct {
	Name         *string `json:"name"`
	Age          *int    `json:"age"`
	Username     *string `json:"username"`
	CurrentCycle *string `json:"current_cycle"`
}

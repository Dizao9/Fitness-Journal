package dto

type UserProfileResponseDTO struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Age          int    `json:"age"`
	Username     string `json:"username"`
	CurrentCycle string `json:"current_cycle,omitempty"`
	Gender       string `json:"gender"`
	Email        string `json:"email"`
}

type UpdateProfileRequest struct {
	Name         *string `json:"name"`
	Age          *int    `json:"age"`
	Username     *string `json:"username"`
	CurrentCycle *string `json:"current_cycle"`
}

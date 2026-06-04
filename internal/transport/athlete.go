package transport

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/service"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/google/uuid"
)

type AthleteService interface {
	GetProfile(userID uuid.UUID) (service.AthleteProfileOutput, error)
	UpdateProfile(userID uuid.UUID, upd service.UpdateProfileInput) error
	DeleteProfile(userID uuid.UUID) error
}

type AthleteHandler struct {
	AthlServ AthleteService
}

func NewAthleteHandler(AthlSvc AthleteService) *AthleteHandler {
	return &AthleteHandler{
		AthlServ: AthlSvc,
	}
}

func (h *AthleteHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Print("[PROFILE] false from context")
		http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
		return
	}

	profile, err := h.AthlServ.GetProfile(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("[PROFILE] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	responseDTO := dto.UserProfileResponseDTO{
		ID:           profile.ID,
		Name:         profile.Name,
		Age:          profile.Age,
		Username:     profile.Username,
		CurrentCycle: profile.CurrentCycle,
		Gender:       profile.Gender,
		Email:        profile.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responseDTO); err != nil {
		log.Printf("[PROFILE] encoder failed: %v", err)
	}
}

func (h *AthleteHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Print("[PROFILE] false from context")
		http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
		return
	}

	var req dto.UpdateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request data", http.StatusBadRequest)
		return
	}

	input := service.UpdateProfileInput{
		Name:         req.Name,
		Age:          req.Age,
		Username:     req.Username,
		CurrentCycle: req.CurrentCycle,
	}

	if err := h.AthlServ.UpdateProfile(userID, input); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("[UPDATE] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "successful"}); err != nil {
		log.Printf("[UPDATE] internal server error: %v", err)
	}
}

func (h *AthleteHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Print("[PROFILE] false from context")
		http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
		return
	}

	err := h.AthlServ.DeleteProfile(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

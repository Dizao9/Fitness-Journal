package transport

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
)

type AthleteService interface {
	GetByUserID(userID string) (dto.UserProfileResponseDTO, error)
	UpdateUser(userID string, upd dto.UpdateProfileRequest) error
	DeleteUser(userID string) error
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

	response, err := h.AthlServ.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		log.Printf("[PROFILE] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
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

	var upd dto.UpdateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		http.Error(w, "invalid request data", http.StatusBadRequest)
		return
	}

	if err := h.AthlServ.UpdateUser(userID, upd); err != nil {
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

	err := h.AthlServ.DeleteUser(userID)
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
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "ok"}); err != nil {
		log.Printf("[DELETE] internal server error: %v", err)
	}
}

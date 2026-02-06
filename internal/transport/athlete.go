package transport

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
)

type AuthService interface {
	Register(req dto.RegisterUser) (string, error)
}

type Handler struct {
	AuthSvc AuthService
}

func (h *Handler) ValidateUser(u dto.RegisterUser) []string {
	var errs []string

	switch {
	case utf8.RuneCountInString(u.Username) < 3:
		errs = append(errs, "username: слишком короткий (мин. 3)")
	case utf8.RuneCountInString(u.Username) > 20:
		errs = append(errs, "username: слишком длинный (max. 20)")
	}

	if !strings.Contains(u.Email, "@") {
		errs = append(errs, "email: некорректный формат")
	}

	if len(u.Password) < 8 {
		errs = append(errs, "password: слишком короткий (min. 8)")
	}
	return errs
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var u dto.RegisterUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "problem to decode your data", http.StatusBadRequest)
		return
	}

	if errs := h.ValidateUser(u); len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"errors": errs,
		})
		return
	}

	idStr, err := h.AuthSvc.Register(u)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}
		log.Printf("[AUTH] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{
		"id": idStr,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[AUTH] encode error: %v\n", err)
		return
	}
}

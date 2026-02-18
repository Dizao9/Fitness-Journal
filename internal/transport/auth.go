package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/service"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
)

type ctxKey int

const userIDKey ctxKey = iota

func ContextWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

type AuthService interface {
	Register(req dto.RegisterUser) (string, error)
	Login(email string, password string) (string, error)
	ParseToken(token string) (*service.CustomClaims, error)
}

type AuthHandler struct {
	AuthSvc AuthService
}

func NewAuthHandler(ASvc AuthService) AuthHandler {
	return AuthHandler{
		AuthSvc: ASvc,
	}
}

func (h *AuthHandler) ValidateUser(u dto.RegisterUser) []string {
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

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u dto.LoginUser
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !strings.Contains(u.Email, "@") {
		http.Error(w, "invalid request, email required", http.StatusBadRequest)
		return
	}

	if u.Password == "" {
		http.Error(w, "invalid request, password required", http.StatusBadRequest)
		return
	}

	token, err := h.AuthSvc.Login(u.Email, u.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			http.Error(w, "login error", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //on prodaction = true
		SameSite: http.SameSiteStrictMode,
		MaxAge:   15 * 60 * 60,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "login successful"}); err != nil {
		log.Printf("[LOGIN] encoder failed:%v", err)
	}
}

func (h *AuthHandler) AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "token cookie is required", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value
		if tokenStr == "" {
			http.Error(w, "token cookie is required", http.StatusBadRequest)
			return
		}
		claims, err := h.AuthSvc.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "some tokens problem", http.StatusUnauthorized)
			return
		}
		ctx := ContextWithUserID(r.Context(), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

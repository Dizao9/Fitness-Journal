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
	"github.com/google/uuid"
)

type ctxKey int

const userIDKey ctxKey = iota

func ContextWithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

type AuthService interface {
	Register(req dto.RegisterUser) (uuid.UUID, error)
	Login(email string, password string) (service.TokenPair, error)
	ParseAccessToken(tokenStr string) (*service.CustomClaims, error)
	ParseRefreshToken(tokenStr string) (*service.RefreshClaims, error)
	ExistsByID(id uuid.UUID) (bool, error)
	Refresh(refreshToken string) (service.TokenPair, error)
	LogOut(jti uuid.UUID) (bool, error)
}

type AuthHandler struct {
	AuthSvc AuthService
}

func NewAuthHandler(ASvc AuthService) *AuthHandler {
	return &AuthHandler{
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

func (h *AuthHandler) SetTokenCookie(w http.ResponseWriter, tokenPair service.TokenPair) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //on prodaction = true
		SameSite: http.SameSiteStrictMode,
		MaxAge:   30 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //on prodaction = true
		SameSite: http.SameSiteStrictMode,
		MaxAge:   30 * 24 * 60 * 60,
	})
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

	tokenPair, err := h.AuthSvc.Login(u.Email, u.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			http.Error(w, "login error", http.StatusUnauthorized)
			return
		}
		log.Printf("[LOGIN] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.SetTokenCookie(w, tokenPair)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "login successful"}); err != nil {
		log.Printf("[LOGIN] encoder failed:%v", err)
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh_token cookie is required", http.StatusUnauthorized)
		return
	}
	refreshTokenStr := refreshCookie.Value
	if refreshTokenStr == "" {
		http.Error(w, "refresh_token cookie is required", http.StatusUnauthorized)
		return
	}

	tokenPair, err := h.AuthSvc.Refresh(refreshTokenStr)
	if err != nil {
		http.Error(w, "authorization problem", http.StatusUnauthorized)
		return
	}

	h.SetTokenCookie(w, tokenPair)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "ok"}); err != nil {
		log.Printf("[REFRESH] encoder failed: %v", err)
	}
}

func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	refreshCokie, err := r.Cookie("refresh_token")
	if err == nil && refreshCokie.Value != "" {
		claims, err := h.AuthSvc.ParseRefreshToken(refreshCokie.Value)
		if err == nil {
			_, err := h.AuthSvc.LogOut(claims.JTI)
			if err != nil {
				log.Printf("[LOGOUT] db failed: %v", err)
			}
		}
	}

	h.clearAuthCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) clearAuthCookies(w http.ResponseWriter) {
	cookieNames := []string{"access_token", "refresh_token"}
	for _, name := range cookieNames {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})
	}
}

func (h *AuthHandler) AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "access_token cookie is required", http.StatusUnauthorized)
			return
		}
		tokenStr := cookie.Value
		if tokenStr == "" {
			http.Error(w, "access_token cookie is required", http.StatusUnauthorized)
			return
		}
		claims, err := h.AuthSvc.ParseAccessToken(tokenStr)
		if err != nil {
			http.Error(w, "some tokens problem", http.StatusUnauthorized)
			return
		}

		ctx := ContextWithUserID(r.Context(), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

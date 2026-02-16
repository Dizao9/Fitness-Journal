package service

import (
	"fmt"
	"log"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/config"
	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type AthleteStorage interface {
	CreateAthlete(athlete domain.Athlete) (string, error)
	GetByEmail(email string) (domain.Athlete, error)
}

type AuthService struct {
	Store AthleteStorage
	Conf  *config.Config
}

func NewAuthService(s AthleteStorage, c *config.Config) *AuthService {
	if s == nil {
		log.Fatalf("[AUTH] Storage is required")
	}
	return &AuthService{
		Store: s,
		Conf:  c,
	}
}

func (a *AuthService) Register(req dto.RegisterUser) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return "", err
	}

	u := domain.Athlete{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(bytes),
		Age:          &req.Age,
		CreatedAt:    time.Now(),
		Name:         domain.PtrString(req.Name),
	}

	return a.Store.CreateAthlete(u)
}

func (a *AuthService) generateToken(id, role, email string) (string, error) {
	claims := CustomClaims{
		UserID: id,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(a.Conf.JWTSecret)
	return token.SignedString(secretKey)
}

func (a *AuthService) Login(email string, password string) (string, error) {
	user, err := a.Store.GetByEmail(email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return "", domain.ErrInvalidCredentials
		}
		log.Printf("[LOGIN] database process was failed: %v", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return a.generateToken(user.ID, user.GetRole(), user.Email)
}

func (a *AuthService) ParseToken(tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method", domain.ErrInvalidCredentials)
		}
		return []byte(a.Conf.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, domain.ErrInvalidCredentials
	}
	return claims, nil
}

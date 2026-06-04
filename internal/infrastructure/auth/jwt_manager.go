package auth

import (
	"fmt"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	accessSecret  string
	refreshSecret string
}

type accessClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	UserID string `json:"user_id"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func NewJWTManager(accessSecret, refreshSecret string) *JWTManager {
	return &JWTManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (m *JWTManager) GenerateAccessToken(userID uuid.UUID, role, email string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"email":   email,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.accessSecret))
}

func (m *JWTManager) GenerateRefreshToken(userID uuid.UUID, jti uuid.UUID, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"jti":     jti.String(),
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.refreshSecret))
}

func (m *JWTManager) ParseAccessToken(tokenStr string) (*service.AccessTokenClaims, error) {
	claims := &accessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method", domain.ErrInvalidCredentials)
		}
		return []byte(m.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, domain.ErrInvalidCredentials
	}

	uID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}
	return &service.AccessTokenClaims{
		UserID: uID,
		Role:   claims.Role,
		Email:  claims.Email,
	}, err
}

func (m *JWTManager) ParseRefreshToken(tokenStr string) (*service.RefreshTokenClaims, error) {
	claims := &refreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method", domain.ErrInvalidCredentials)
		}
		return []byte(m.refreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, domain.ErrInvalidCredentials
	}

	uID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}
	jti, err := uuid.Parse(claims.JTI)
	if err != nil {
		return nil, err
	}

	return &service.RefreshTokenClaims{
		UserID: uID,
		JTI:    jti,
	}, nil
}

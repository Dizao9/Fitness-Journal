package service

import (
	"fmt"
	"log"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/config"
	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshClaims struct {
	UserID uuid.UUID `json:"user_id"`
	JTI    uuid.UUID `json:"jti"`
	jwt.RegisteredClaims
}
type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type AuthStorage interface {
	CreateAthlete(athlete domain.Athlete) (uuid.UUID, error)
	GetByEmail(email string) (domain.Athlete, error)
	GetByUserID(userID uuid.UUID) (domain.Athlete, error)
	ExistsByID(userID uuid.UUID) (bool, error)
	SaveRefreshToken(userID, jti uuid.UUID, expiresAt time.Time, refreshToken string) error
	DeleteRefreshToken(jti uuid.UUID) (bool, error)
}

type AuthService struct {
	Store AuthStorage
	Conf  *config.Config
}

func NewAuthService(s AuthStorage, c *config.Config) *AuthService {
	if s == nil {
		log.Fatalf("[AUTH] Storage is required")
	}
	return &AuthService{
		Store: s,
		Conf:  c,
	}
}

func (a *AuthService) Register(req dto.RegisterUser) (uuid.UUID, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return uuid.Nil, err
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

func (a *AuthService) generateToken(id uuid.UUID, role, email string) (TokenPair, error) {
	accessClaims := CustomClaims{
		UserID: id,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(a.Conf.JWTSecret))
	if err != nil {
		return TokenPair{}, err
	}

	expiresDate := time.Now().Add(30 * 24 * time.Hour)
	jti := uuid.New()
	refreshClaims := RefreshClaims{
		UserID: id,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresDate),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(a.Conf.RefreshJWTSecret))
	if err != nil {
		return TokenPair{}, err
	}

	err = a.Store.SaveRefreshToken(id, jti, expiresDate, refreshToken)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) Login(email string, password string) (TokenPair, error) {
	user, err := a.Store.GetByEmail(email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return TokenPair{}, domain.ErrInvalidCredentials
		}
		log.Printf("[LOGIN] database process was failed: %v", err)
		return TokenPair{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return TokenPair{}, domain.ErrInvalidCredentials
	}

	return a.generateToken(user.ID, user.GetRole(), user.Email)
}

func (a *AuthService) ParseAccessToken(tokenStr string) (*CustomClaims, error) {
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

func (a *AuthService) ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method", domain.ErrInvalidCredentials)
		}
		return []byte(a.Conf.RefreshJWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, domain.ErrInvalidCredentials
	}
	return claims, err
}

func (a *AuthService) ExistsByID(id uuid.UUID) (bool, error) {
	return a.Store.ExistsByID(id)
}

func (a *AuthService) Refresh(refreshToken string) (TokenPair, error) {
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidCredentials
		}
		return ([]byte(a.Conf.RefreshJWTSecret)), nil
	})
	if err != nil || !token.Valid {
		if err != nil {
			log.Printf("[REFRESH SERVICE] error: %v", err)
		}
		return TokenPair{}, domain.ErrInvalidCredentials
	}

	exists, err := a.Store.DeleteRefreshToken(claims.JTI)
	if err != nil {
		return TokenPair{}, err
	}

	if !exists {
		return TokenPair{}, domain.ErrInvalidCredentials
	}

	user, err := a.Store.GetByUserID(claims.UserID)
	if err != nil {
		return TokenPair{}, err
	}

	return a.generateToken(user.ID, user.GetRole(), user.Email)
}

func (a *AuthService) LogOut(jti uuid.UUID) (bool, error) {
	return a.Store.DeleteRefreshToken(jti)
}

package service

import (
	"log"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
type AccessTokenClaims struct {
	UserID uuid.UUID
	Role   string
	Email  string
}

type RefreshTokenClaims struct {
	UserID uuid.UUID
	JTI    uuid.UUID
}

type AuthStorage interface {
	CreateAthlete(athlete domain.Athlete) (uuid.UUID, error)
	GetByEmail(email string) (domain.Athlete, error)
	GetByUserID(userID uuid.UUID) (domain.Athlete, error)
	ExistsByID(userID uuid.UUID) (bool, error)
	SaveRefreshToken(userID, jti uuid.UUID, expiresAt time.Time, refreshToken string) error
	DeleteRefreshToken(jti uuid.UUID) (bool, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	ComparePassword(hash, password string) error
}

type TokenManager interface {
	GenerateAccessToken(userID uuid.UUID, role, email string, ttl time.Duration) (string, error)
	GenerateRefreshToken(userID uuid.UUID, jti uuid.UUID, ttl time.Duration) (string, error)

	ParseAccessToken(tokenStr string) (*AccessTokenClaims, error)
	ParseRefreshToken(tokenStr string) (*RefreshTokenClaims, error)
}

type AuthService struct {
	Store          AuthStorage
	TokenManager   TokenManager
	PasswordHasher PasswordHasher
}

func NewAuthService(s AuthStorage, tm TokenManager, ph PasswordHasher) *AuthService {
	if s == nil {
		log.Fatalf("[AUTH] Storage is required")
	}
	return &AuthService{
		Store:          s,
		TokenManager:   tm,
		PasswordHasher: ph,
	}
}

func (a *AuthService) Register(req RegisterInput) (uuid.UUID, error) {
	bytes, err := a.PasswordHasher.Hash(req.Password)
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

	accessToken, err := a.TokenManager.GenerateAccessToken(id, role, email, 30*time.Minute)
	if err != nil {
		return TokenPair{}, err
	}

	expiresDate := time.Now().Add(30 * 24 * time.Hour)
	jti := uuid.New()

	refreshToken, err := a.TokenManager.GenerateRefreshToken(id, jti, 30*24*time.Hour)
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

func (a *AuthService) Login(loginData LoginInput) (TokenPair, error) {
	user, err := a.Store.GetByEmail(loginData.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return TokenPair{}, domain.ErrInvalidCredentials
		}
		log.Printf("[LOGIN] database process was failed: %v", err)
		return TokenPair{}, err
	}

	err = a.PasswordHasher.ComparePassword(user.PasswordHash, loginData.Password)
	if err != nil {
		return TokenPair{}, domain.ErrInvalidCredentials
	}

	return a.generateToken(user.ID, user.GetRole(), user.Email)
}

func (a *AuthService) ParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	return a.TokenManager.ParseAccessToken(tokenStr)
}

func (a *AuthService) ExistsByID(id uuid.UUID) (bool, error) {
	return a.Store.ExistsByID(id)
}

func (a *AuthService) Refresh(refreshToken string) (TokenPair, error) {
	claims, err := a.TokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
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

func (a *AuthService) LogOut(refreshToken string) error {
	claims, err := a.TokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return domain.ErrInvalidCredentials
	}

	_, err = a.Store.DeleteRefreshToken(claims.JTI)
	return err
}

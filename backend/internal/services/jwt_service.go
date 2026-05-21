package services

import (
	"fmt"
	"time"

	"github.com/gabrzb/auth-go-gin/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims is shared by both access and refresh tokens. Email is omitted in refresh tokens.
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewJWTService(cfg *config.Config) (*JWTService, error) {
	accessExp, err := time.ParseDuration(cfg.JWTAccessExpiration)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRATION: %w", err)
	}
	refreshExp, err := time.ParseDuration(cfg.JWTRefreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRATION: %w", err)
	}
	return &JWTService{
		secret:            cfg.JWTSecret,
		accessExpiration:  accessExp,
		refreshExpiration: refreshExp,
	}, nil
}

func (s *JWTService) GenerateAccessToken(userID uint, email string) (string, error) {
	return s.generate(userID, email, s.accessExpiration)
}

func (s *JWTService) GenerateRefreshToken(userID uint) (string, error) {
	return s.generate(userID, "", s.refreshExpiration)
}

func (s *JWTService) AccessExpiresIn() int {
	return int(s.accessExpiration.Seconds())
}

// generate is the single place that creates and signs a JWT (DRY).
func (s *JWTService) generate(userID uint, email string, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

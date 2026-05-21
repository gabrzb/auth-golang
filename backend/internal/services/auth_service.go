package services

import (
	"errors"
	"time"

	"github.com/gabrzb/auth-go-gin/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

// jwtService is the narrow interface AuthService needs (SOLID-I, SOLID-D).
type jwtService interface {
	GenerateAccessToken(userID uint, email string) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	AccessExpiresIn() int
}

// tokenBlacklist is the narrow interface for revoking tokens (SOLID-I, SOLID-D).
type tokenBlacklist interface {
	Add(token string, ttl time.Duration) error
	Contains(token string) (bool, error)
}

type AuthService struct {
	db        *gorm.DB
	jwt       jwtService
	blacklist tokenBlacklist
}

func NewAuthService(db *gorm.DB, jwt jwtService, blacklist tokenBlacklist) *AuthService {
	return &AuthService{db: db, jwt: jwt, blacklist: blacklist}
}

func (s *AuthService) Register(email, password string) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{Email: email, Password: string(hashed)}

	if err := s.db.Create(user).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (accessToken, refreshToken string, expiresIn int, err error) {
	var user models.User
	if err = s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", 0, ErrInvalidCredentials
		}
		return "", "", 0, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", 0, ErrInvalidCredentials
	}

	accessToken, err = s.jwt.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err = s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, s.jwt.AccessExpiresIn(), nil
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) Refresh(refreshToken string) (string, int, error) {
	claims, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return "", 0, err // ErrExpiredToken or ErrInvalidToken propagate as-is
	}

	user, err := s.GetUserByID(claims.UserID)
	if err != nil {
		return "", 0, err
	}

	access, err := s.jwt.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", 0, err
	}

	return access, s.jwt.AccessExpiresIn(), nil
}

func (s *AuthService) Logout(accessToken, refreshToken string) error {
	if err := s.blacklistToken(accessToken); err != nil {
		return err
	}
	return s.blacklistToken(refreshToken)
}

// blacklistToken adds a token to the blacklist with a TTL equal to its remaining lifetime (DRY).
func (s *AuthService) blacklistToken(tokenString string) error {
	claims, err := s.jwt.ValidateToken(tokenString)
	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			return nil // already expired — no need to store
		}
		return ErrInvalidToken
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	return s.blacklist.Add(tokenString, ttl)
}

// isUniqueViolation checks for Postgres error code 23505 (unique_violation).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

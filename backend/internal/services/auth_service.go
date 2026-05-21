package services

import (
	"errors"

	"github.com/gabrzb/auth-go-gin/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// jwtService is the narrow interface AuthService needs — defined here, at the consumption site (SOLID-I, SOLID-D).
type jwtService interface {
	GenerateAccessToken(userID uint, email string) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	AccessExpiresIn() int
}

type AuthService struct {
	db  *gorm.DB
	jwt jwtService
}

func NewAuthService(db *gorm.DB, jwt jwtService) *AuthService {
	return &AuthService{db: db, jwt: jwt}
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

// isUniqueViolation checks for Postgres error code 23505 (unique_violation).
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

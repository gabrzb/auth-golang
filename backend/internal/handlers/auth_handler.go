package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gabrzb/auth-go-gin/internal/middleware"
	"github.com/gabrzb/auth-go-gin/internal/models"
	"github.com/gabrzb/auth-go-gin/internal/services"
	"github.com/gin-gonic/gin"
)

// authService is the narrow interface this handler needs — defined here, not in the service package (SOLID-I, SOLID-D).
type authService interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (accessToken, refreshToken string, accessExpiresIn, refreshExpiresIn int, err error)
	Rotate(refreshToken string) (accessToken, newRefreshToken string, accessExpiresIn, refreshExpiresIn int, err error)
	Logout(accessToken, refreshToken string) error
}

type AuthHandler struct {
	svc          authService
	cookieSecure bool
}

func NewAuthHandler(svc authService, cookieSecure bool) *AuthHandler {
	return &AuthHandler{svc: svc, cookieSecure: cookieSecure}
}

const refreshCookieName = "refresh_token"
const refreshCookiePath = "/auth"

// setRefreshCookie writes the refresh token as an httpOnly + SameSite=Strict cookie
// scoped to /auth so it isn't sent on every API call. Secure is toggled per env.
func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string, maxAgeSeconds int) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(refreshCookieName, token, maxAgeSeconds, refreshCookiePath, "", h.cookieSecure, true)
}

// clearRefreshCookie expires the refresh cookie immediately.
func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(refreshCookieName, "", -1, refreshCookiePath, "", h.cookieSecure, true)
}

type registerRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Logout(c *gin.Context) {
	val, exists := c.Get(middleware.TokenKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	accessToken, ok := val.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		return
	}

	// Refresh cookie is best-effort: if it's missing or malformed we still blacklist
	// the access token. AuthService.Logout tolerates an empty refresh string.
	refreshToken, _ := c.Cookie(refreshCookieName)

	// Always clear the cookie so the browser stops sending a dead token, regardless
	// of whether blacklisting succeeds below.
	h.clearRefreshCookie(c)

	if err := h.svc.Logout(accessToken, refreshToken); err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		log.Printf("logout error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	access, newRefresh, accessExpiresIn, refreshExpiresIn, err := h.svc.Rotate(refreshToken)
	if err != nil {
		// Old refresh is no longer valid — clear the cookie so the browser doesn't keep replaying it.
		h.clearRefreshCookie(c)
		if errors.Is(err, services.ErrExpiredToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	h.setRefreshCookie(c, newRefresh, refreshExpiresIn)

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		"expires_in":   accessExpiresIn,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, accessExpiresIn, refreshExpiresIn, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.setRefreshCookie(c, refresh, refreshExpiresIn)

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		"expires_in":   accessExpiresIn,
	})
}

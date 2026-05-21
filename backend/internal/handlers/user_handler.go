package handlers

import (
	"net/http"

	"github.com/gabrzb/auth-go-gin/internal/middleware"
	"github.com/gabrzb/auth-go-gin/internal/models"
	"github.com/gin-gonic/gin"
)

// userService is the narrow interface this handler needs (SOLID-I, SOLID-D).
type userService interface {
	GetUserByID(id uint) (*models.User, error)
}

type UserHandler struct {
	svc userService
}

func NewUserHandler(svc userService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user, err := h.svc.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

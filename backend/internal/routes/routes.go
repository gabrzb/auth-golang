package routes

import (
	"net/http"

	"github.com/gabrzb/auth-go-gin/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, authMiddleware gin.HandlerFunc) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public routes — no authentication required
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Protected routes — JWT required
	protected := r.Group("/")
	protected.Use(authMiddleware)
	{
		protected.GET("/me", userHandler.Me)
	}
}

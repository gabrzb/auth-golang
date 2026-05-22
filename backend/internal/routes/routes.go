package routes

import (
	"net/http"

	"github.com/gabrzb/auth-go-gin/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, authMiddleware, originCheck gin.HandlerFunc) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public auth routes — origin-checked as CSRF defense-in-depth on top of SameSite=Strict.
	auth := r.Group("/auth")
	auth.Use(originCheck)
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Protected routes — JWT required (origin-check too, since logout is state-changing).
	protected := r.Group("/")
	protected.Use(authMiddleware)
	{
		protected.GET("/me", userHandler.Me)
		protected.POST("/auth/logout", originCheck, authHandler.Logout)
	}
}

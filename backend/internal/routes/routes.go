package routes

import (
	"net/http"

	"github.com/gabrzb/auth-go-gin/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, authHandler *handlers.AuthHandler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}

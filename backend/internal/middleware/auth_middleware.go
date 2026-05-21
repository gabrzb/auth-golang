package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gabrzb/auth-go-gin/internal/services"
	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

// tokenValidator is the narrow interface the middleware needs (SOLID-I, SOLID-D).
type tokenValidator interface {
	ValidateToken(tokenString string) (*services.Claims, error)
}

func Auth(validator tokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		claims, err := validator.ValidateToken(tokenString)
		if err != nil {
			if errors.Is(err, services.ErrExpiredToken) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

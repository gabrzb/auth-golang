package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gabrzb/auth-go-gin/internal/services"
	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"
const TokenKey = "token"

// tokenValidator is the narrow interface for JWT validation (SOLID-I, SOLID-D).
type tokenValidator interface {
	ValidateToken(tokenString string) (*services.Claims, error)
}

// blacklistChecker is the narrow interface for blacklist lookup (SOLID-I, SOLID-D).
type blacklistChecker interface {
	Contains(ctx context.Context, token string) (bool, error)
}

func Auth(validator tokenValidator, blacklist blacklistChecker) gin.HandlerFunc {
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

		revoked, err := blacklist.Contains(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if revoked {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(TokenKey, tokenString)
		c.Next()
	}
}

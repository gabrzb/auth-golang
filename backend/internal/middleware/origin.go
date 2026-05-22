package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireOrigin rejects browser requests whose Origin header isn't in the allow-list.
// Requests without an Origin header are allowed through — those are non-browser callers
// (curl, Postman, server-to-server) which a CSRF mitigation doesn't apply to. This is
// belt-and-suspenders on top of SameSite=Strict on the refresh cookie.
func RequireOrigin(allowed []string) gin.HandlerFunc {
	set := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		set[a] = struct{}{}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}
		if _, ok := set[origin]; ok {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "origin not allowed"})
	}
}

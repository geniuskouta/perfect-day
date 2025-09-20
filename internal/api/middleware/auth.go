package middleware

import (
	"net/http"
	"perfect-day/pkg/auth"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthRequired is middleware that validates authentication for protected routes
func AuthRequired(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Not authenticated",
				},
				"meta": gin.H{
					"timestamp": time.Now().UTC().Format(time.RFC3339),
					"version":   "0.1.0",
				},
			})
			c.Abort()
			return
		}

		user, err := authService.ValidateSession(sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Not authenticated",
				},
				"meta": gin.H{
					"timestamp": time.Now().UTC().Format(time.RFC3339),
					"version":   "0.1.0",
				},
			})
			c.Abort()
			return
		}

		// Store user information in context for use by handlers
		c.Set("user", user)
		c.Set("username", user.Username)
		c.Next()
	}
}
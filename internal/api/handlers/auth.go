package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
}

func (h *Handlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	user, session, err := h.AuthService.Login(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "User not found",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	// Set session cookie
	c.SetCookie("session_id", session.ID, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": gin.H{
				"username":   user.Username,
				"timezone":   user.Timezone,
				"created_at": user.CreatedAt.Format(time.RFC3339),
			},
			"session": gin.H{
				"id":         session.ID,
				"expires_at": session.ExpiresAt.Format(time.RFC3339),
			},
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) GetCurrentUser(c *gin.Context) {
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
		return
	}

	user, err := h.AuthService.ValidateSession(sessionID)
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
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"username":   user.Username,
			"timezone":   user.Timezone,
			"created_at": user.CreatedAt.Format(time.RFC3339),
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}
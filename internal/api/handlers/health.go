package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status":         "healthy",
			"uptime_seconds": 3600, // TODO: Track actual uptime
			"version":        "0.1.0",
			"checks": gin.H{
				"storage":        "ok",
				"google_places":  "ok",
			},
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"version":    "0.1.0",
			"build":      "dev",
			"go_version": "1.21.0",
			"built_at":   time.Now().UTC().Format(time.RFC3339),
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}
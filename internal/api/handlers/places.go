package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) SearchPlaces(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_QUERY",
				"message": "Search query is required",
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
			"places": []interface{}{},
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) GetAreas(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"areas": []interface{}{},
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}
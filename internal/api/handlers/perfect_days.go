package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) ListPerfectDays(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"perfect_days": []interface{}{},
			"pagination": gin.H{
				"total":    0,
				"offset":   0,
				"limit":    10,
				"has_more": false,
			},
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) CreatePerfectDay(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Not implemented yet",
		},
	})
}

func (h *Handlers) GetPerfectDay(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Not implemented yet",
		},
	})
}

func (h *Handlers) UpdatePerfectDay(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Not implemented yet",
		},
	})
}

func (h *Handlers) DeletePerfectDay(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Not implemented yet",
		},
	})
}
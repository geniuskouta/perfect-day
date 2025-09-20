package handlers

import (
	"context"
	"net/http"
	"strconv"
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

	// Get optional parameters
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	// Search places using Google Places API
	ctx := context.Background()
	places, err := h.PlacesService.SearchPlaces(ctx, query)
	if err != nil {
		// If Places API fails, return empty results (graceful degradation)
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"places": []interface{}{},
				"query":  query,
				"limit":  limit,
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
				"notice":    "Places API unavailable, showing fallback results",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"places": places,
			"query":  query,
			"limit":  limit,
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) GetAreas(c *gin.Context) {
	// Load all perfect days to extract unique areas
	allPerfectDays, err := h.Storage.PerfectDayStorage.LoadAll(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to load areas",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	// Extract unique areas using search service
	areas := h.SearchService.GetUniqueAreas(allPerfectDays)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"areas": areas,
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}
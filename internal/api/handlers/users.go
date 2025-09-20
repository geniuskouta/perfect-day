package handlers

import (
	"net/http"
	"perfect-day/pkg/models"
	"perfect-day/pkg/search"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetUserProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_USERNAME",
				"message": "Username is required",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	// Check if user exists
	if !h.Storage.UserStorage.Exists(username) {
		c.JSON(http.StatusNotFound, gin.H{
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

	// Load user profile
	user, err := h.Storage.UserStorage.Load(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to load user profile",
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

func (h *Handlers) GetUserPerfectDays(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_USERNAME",
				"message": "Username is required",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	// Check if user exists
	if !h.Storage.UserStorage.Exists(username) {
		c.JSON(http.StatusNotFound, gin.H{
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

	// Get query parameters for pagination
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	includeDeleted := c.DefaultQuery("include_deleted", "false") == "true"

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Load user's perfect days
	allUserPerfectDays, err := h.Storage.PerfectDayStorage.LoadAllByUser(username, includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to load user's perfect days",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	// Apply pagination
	total := len(allUserPerfectDays)
	end := offset + limit
	if end > total {
		end = total
	}

	var paginatedResults []*models.PerfectDay
	if offset < total {
		paginatedResults = allUserPerfectDays[offset:end]
	}

	// Use search service to format results consistently
	searchResults := search.SearchResult{
		PerfectDays: paginatedResults,
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": searchResults,
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}
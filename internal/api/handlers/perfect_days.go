package handlers

import (
	"net/http"
	"perfect-day/pkg/models"
	"perfect-day/pkg/search"
	"perfect-day/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreatePerfectDayRequest struct {
	Title       string                   `json:"title" binding:"required"`
	Description string                   `json:"description"`
	Date        string                   `json:"date" binding:"required"`
	Activities  []CreateActivityRequest  `json:"activities"`
}

type CreateActivityRequest struct {
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description"`
	Location    CreateLocationRequest `json:"location" binding:"required"`
	StartTime   string              `json:"start_time" binding:"required"`
	Duration    int                 `json:"duration" binding:"required"`
	Commentary  string              `json:"commentary"`
}

type CreateLocationRequest struct {
	Type    string  `json:"type" binding:"required"` // "google_place" or "custom_text"
	PlaceID string  `json:"place_id"`
	Name    string  `json:"name" binding:"required"`
	Area    string  `json:"area" binding:"required"`
	Address string  `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (h *Handlers) ListPerfectDays(c *gin.Context) {
	// Get query parameters
	userFilter := c.Query("user")
	areas := c.Query("areas")
	query := c.Query("q")
	from := c.Query("from")
	to := c.Query("to")
	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	// Load all perfect days
	allPerfectDays, err := h.Storage.PerfectDayStorage.LoadAll(false) // exclude deleted
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to load perfect days",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	// Apply search filters
	searchCriteria := search.SearchCriteria{
		Query:     query,
		Username:  userFilter,
		DateFrom:  from,
		DateTo:    to,
		SortBy:    sortBy,
		SortOrder: order,
		Limit:     limit,
		Offset:    offset,
	}

	if areas != "" {
		searchCriteria.Areas = []string{areas}
	}

	searchResult := h.SearchService.Search(allPerfectDays, searchCriteria)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"perfect_days": searchResult.PerfectDays,
			"pagination": gin.H{
				"total":    searchResult.Total,
				"offset":   searchResult.Offset,
				"limit":    searchResult.Limit,
				"has_more": searchResult.Offset+searchResult.Limit < searchResult.Total,
			},
		},
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) CreatePerfectDay(c *gin.Context) {
	var req CreatePerfectDayRequest
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

	// For now, use a default username - in a real app, this would come from authentication
	username := "testuser"

	// Create perfect day
	perfectDay, err := models.NewPerfectDay(utils.GenerateID(), req.Title, req.Description, username, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	// Add activities
	for _, actReq := range req.Activities {
		location := createLocationFromRequest(actReq.Location)
		activity, err := models.NewActivity(
			utils.GenerateID(),
			actReq.Name,
			*location,
			actReq.StartTime,
			actReq.Duration,
			actReq.Description,
			actReq.Commentary,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "Invalid activity: " + err.Error(),
				},
				"meta": gin.H{
					"timestamp": time.Now().UTC().Format(time.RFC3339),
					"version":   "0.1.0",
				},
			})
			return
		}
		perfectDay.AddActivity(*activity)
	}

	// Save to storage
	if err := h.Storage.PerfectDayStorage.Save(perfectDay); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "STORAGE_ERROR",
				"message": "Failed to save perfect day",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": perfectDay,
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) GetPerfectDay(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Perfect day ID is required",
			},
		})
		return
	}

	// Try to find in all users' data
	allPerfectDays, err := h.Storage.PerfectDayStorage.LoadAll(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to load perfect days",
			},
		})
		return
	}

	var foundPerfectDay *models.PerfectDay
	for _, pd := range allPerfectDays {
		if pd.ID == id {
			foundPerfectDay = pd
			break
		}
	}

	if foundPerfectDay == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Perfect day not found",
			},
			"meta": gin.H{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "0.1.0",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": foundPerfectDay,
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) UpdatePerfectDay(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Perfect day ID is required",
			},
		})
		return
	}

	var req CreatePerfectDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	// Find existing perfect day
	allPerfectDays, err := h.Storage.PerfectDayStorage.LoadAll(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to load perfect days",
			},
		})
		return
	}

	var existingPerfectDay *models.PerfectDay
	for _, pd := range allPerfectDays {
		if pd.ID == id {
			existingPerfectDay = pd
			break
		}
	}

	if existingPerfectDay == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Perfect day not found",
			},
		})
		return
	}

	// Update the perfect day
	updatedPerfectDay, err := models.NewPerfectDay(id, req.Title, req.Description, existingPerfectDay.Username, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Copy creation time
	updatedPerfectDay.CreatedAt = existingPerfectDay.CreatedAt

	// Add activities
	for _, actReq := range req.Activities {
		location := createLocationFromRequest(actReq.Location)
		activity, err := models.NewActivity(
			utils.GenerateID(),
			actReq.Name,
			*location,
			actReq.StartTime,
			actReq.Duration,
			actReq.Description,
			actReq.Commentary,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "Invalid activity: " + err.Error(),
				},
			})
			return
		}
		updatedPerfectDay.AddActivity(*activity)
	}

	// Save to storage
	if err := h.Storage.PerfectDayStorage.Save(updatedPerfectDay); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "STORAGE_ERROR",
				"message": "Failed to save perfect day",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedPerfectDay,
		"meta": gin.H{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "0.1.0",
		},
	})
}

func (h *Handlers) DeletePerfectDay(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "MISSING_ID",
				"message": "Perfect day ID is required",
			},
		})
		return
	}

	// Find existing perfect day
	allPerfectDays, err := h.Storage.PerfectDayStorage.LoadAll(false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to load perfect days",
			},
		})
		return
	}

	var existingPerfectDay *models.PerfectDay
	for _, pd := range allPerfectDays {
		if pd.ID == id {
			existingPerfectDay = pd
			break
		}
	}

	if existingPerfectDay == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Perfect day not found",
			},
		})
		return
	}

	// Soft delete
	existingPerfectDay.SoftDelete()

	// Save to storage
	if err := h.Storage.PerfectDayStorage.Save(existingPerfectDay); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "STORAGE_ERROR",
				"message": "Failed to delete perfect day",
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func createLocationFromRequest(req CreateLocationRequest) *models.Location {
	if req.Type == "google_place" && req.PlaceID != "" {
		var coords *models.Coordinates
		if req.Latitude != nil && req.Longitude != nil {
			coords = &models.Coordinates{
				Latitude:  *req.Latitude,
				Longitude: *req.Longitude,
			}
		}
		return models.NewGooglePlaceLocation(req.PlaceID, req.Name, req.Address, req.Area, coords)
	}
	return models.NewCustomTextLocation(req.Name, req.Area)
}
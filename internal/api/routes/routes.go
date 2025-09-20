package routes

import (
	"perfect-day/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, h *handlers.Handlers) {
	// API v1 routes
	v1 := router.Group("/api/v1")

	// Health check
	v1.GET("/health", h.HealthCheck)
	v1.GET("/version", h.Version)

	// Authentication
	auth := v1.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.GET("/me", h.GetCurrentUser)
	}

	// Perfect days
	perfectDays := v1.Group("/perfect-days")
	{
		perfectDays.GET("", h.ListPerfectDays)
		perfectDays.POST("", h.CreatePerfectDay)
		perfectDays.GET("/:id", h.GetPerfectDay)
		perfectDays.PUT("/:id", h.UpdatePerfectDay)
		perfectDays.DELETE("/:id", h.DeletePerfectDay)
	}

	// Users
	users := v1.Group("/users")
	{
		users.GET("/:username", h.GetUserProfile)
		users.GET("/:username/perfect-days", h.GetUserPerfectDays)
	}

	// Places
	places := v1.Group("/places")
	{
		places.GET("/search", h.SearchPlaces)
	}

	// Areas
	v1.GET("/areas", h.GetAreas)
}
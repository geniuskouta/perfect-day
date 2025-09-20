package routes

import (
	"perfect-day/internal/api/handlers"
	"perfect-day/internal/api/middleware"
	"perfect-day/pkg/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, h *handlers.Handlers, authService *auth.AuthService) {
	// API v1 routes
	v1 := router.Group("/api/v1")

	// Health check
	v1.GET("/health", h.HealthCheck)
	v1.GET("/version", h.Version)

	// Authentication
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
		authGroup.GET("/me", middleware.AuthRequired(authService), h.GetCurrentUser)
	}

	// Perfect days
	perfectDays := v1.Group("/perfect-days")
	{
		perfectDays.GET("", h.ListPerfectDays) // Public read access
		perfectDays.POST("", middleware.AuthRequired(authService), h.CreatePerfectDay)
		perfectDays.GET("/:id", h.GetPerfectDay) // Public read access
		perfectDays.PUT("/:id", middleware.AuthRequired(authService), h.UpdatePerfectDay)
		perfectDays.DELETE("/:id", middleware.AuthRequired(authService), h.DeletePerfectDay)
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
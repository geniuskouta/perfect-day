package handlers

import (
	"perfect-day/pkg/auth"
	"perfect-day/pkg/places"
	"perfect-day/pkg/search"
	"perfect-day/pkg/storage"
)

type Handlers struct {
	AuthService   *auth.AuthService
	Storage       *storage.Storage
	PlacesService *places.PlacesService
	SearchService *search.SearchService
}
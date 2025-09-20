package places

import (
	"context"
	"fmt"
	"perfect-day/src/models"
	"strings"

	"googlemaps.github.io/maps"
)

type PlacesService struct {
	client *maps.Client
}

func NewPlacesService(apiKey string) (*PlacesService, error) {
	if apiKey == "" {
		return &PlacesService{client: nil}, nil
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Maps client: %v", err)
	}

	return &PlacesService{client: client}, nil
}

func (ps *PlacesService) IsEnabled() bool {
	return ps.client != nil
}

func (ps *PlacesService) SearchPlaces(ctx context.Context, query string) ([]PlaceResult, error) {
	if ps.client == nil {
		return nil, fmt.Errorf("Google Places API is not enabled")
	}

	request := &maps.TextSearchRequest{
		Query: query,
	}

	response, err := ps.client.TextSearch(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to search places: %v", err)
	}

	var results []PlaceResult
	for _, place := range response.Results {
		result := PlaceResult{
			PlaceID:   place.PlaceID,
			Name:      place.Name,
			Address:   place.FormattedAddress,
			Latitude:  place.Geometry.Location.Lat,
			Longitude: place.Geometry.Location.Lng,
		}
		results = append(results, result)
	}

	return results, nil
}

func (ps *PlacesService) GetPlaceDetails(ctx context.Context, placeID string) (*PlaceResult, error) {
	if ps.client == nil {
		return nil, fmt.Errorf("Google Places API is not enabled")
	}

	request := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
	}

	response, err := ps.client.PlaceDetails(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get place details: %v", err)
	}

	result := &PlaceResult{
		PlaceID:   response.PlaceID,
		Name:      response.Name,
		Address:   response.FormattedAddress,
		Latitude:  response.Geometry.Location.Lat,
		Longitude: response.Geometry.Location.Lng,
	}

	return result, nil
}

func (ps *PlacesService) CreateLocationFromPlace(place PlaceResult, area string) *models.Location {
	coords := &models.Coordinates{
		Latitude:  place.Latitude,
		Longitude: place.Longitude,
	}

	return models.NewGooglePlaceLocation(
		place.PlaceID,
		place.Name,
		place.Address,
		area,
		coords,
	)
}

func (ps *PlacesService) SuggestAreaFromAddress(address string) string {
	if address == "" {
		return ""
	}

	parts := strings.Split(address, ",")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[len(parts)-2])
	}

	return strings.TrimSpace(parts[0])
}

type PlaceResult struct {
	PlaceID   string  `json:"place_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
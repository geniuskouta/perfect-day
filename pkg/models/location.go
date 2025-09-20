package models

type LocationType string

const (
	GooglePlaceLocation LocationType = "google_place"
	CustomTextLocation  LocationType = "custom_text"
)

type Location struct {
	Type        LocationType `json:"type"`
	PlaceID     string       `json:"place_id,omitempty"`
	Name        string       `json:"name"`
	Address     string       `json:"address,omitempty"`
	Area        string       `json:"area"`
	Coordinates *Coordinates `json:"coordinates,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewGooglePlaceLocation(placeID, name, address, area string, coords *Coordinates) *Location {
	return &Location{
		Type:        GooglePlaceLocation,
		PlaceID:     placeID,
		Name:        name,
		Address:     address,
		Area:        area,
		Coordinates: coords,
	}
}

func NewCustomTextLocation(name, area string) *Location {
	return &Location{
		Type: CustomTextLocation,
		Name: name,
		Area: area,
	}
}
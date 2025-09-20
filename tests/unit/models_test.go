package unit

import (
	"perfect-day/src/models"
	"strings"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		timezone string
		wantErr  bool
	}{
		{"valid user", "testuser", "America/New_York", false},
		{"valid user tokyo", "kouta", "Asia/Tokyo", false},
		{"empty username", "", "UTC", true},
		{"short username", "ab", "UTC", true},
		{"long username", "this_is_a_very_long_username_that_exceeds_20_chars", "UTC", true},
		{"invalid chars", "user@name", "UTC", true},
		{"invalid timezone", "testuser", "Invalid/Timezone", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := models.NewUser(tt.username, tt.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if user.Username != tt.username {
					t.Errorf("NewUser().Username = %v, want %v", user.Username, tt.username)
				}
				if user.Timezone != tt.timezone {
					t.Errorf("NewUser().Timezone = %v, want %v", user.Timezone, tt.timezone)
				}
				if user.CreatedAt.IsZero() {
					t.Error("NewUser().CreatedAt should not be zero")
				}
			}
		})
	}
}

func TestNewLocation(t *testing.T) {
	coords := &models.Coordinates{
		Latitude:  35.6762,
		Longitude: 139.6503,
	}

	googlePlace := models.NewGooglePlaceLocation("place123", "Tokyo Station", "Tokyo, Japan", "Tokyo", coords)
	if googlePlace.Type != models.GooglePlaceLocation {
		t.Errorf("Expected GooglePlaceLocation type, got %v", googlePlace.Type)
	}
	if googlePlace.PlaceID != "place123" {
		t.Errorf("Expected PlaceID 'place123', got %v", googlePlace.PlaceID)
	}

	customLocation := models.NewCustomTextLocation("My Cafe", "Shibuya")
	if customLocation.Type != models.CustomTextLocation {
		t.Errorf("Expected CustomTextLocation type, got %v", customLocation.Type)
	}
	if customLocation.PlaceID != "" {
		t.Errorf("Expected empty PlaceID for custom location, got %v", customLocation.PlaceID)
	}
}

func TestNewActivity(t *testing.T) {
	location := models.NewCustomTextLocation("Test Location", "Test Area")

	tests := []struct {
		name        string
		activityName string
		startTime   string
		duration    int
		wantErr     bool
	}{
		{"valid activity", "Coffee", "09:00", 60, false},
		{"valid activity 2", "Lunch", "12:30", 90, false},
		{"empty name", "", "10:00", 30, true},
		{"invalid time", "Walk", "25:00", 30, true},
		{"invalid time format", "Walk", "9:0", 30, true},
		{"zero duration", "Break", "14:00", 0, true},
		{"negative duration", "Rest", "15:00", -30, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			activity, err := models.NewActivity("test-id", tt.activityName, *location, tt.startTime, tt.duration, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("NewActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if activity.Name != tt.activityName {
					t.Errorf("NewActivity().Name = %v, want %v", activity.Name, tt.activityName)
				}
				endTime, err := activity.EndTime()
				if err != nil {
					t.Errorf("EndTime() error = %v", err)
				}
				if endTime == "" {
					t.Error("EndTime() should not be empty")
				}
			}
		})
	}
}

func TestNewPerfectDay(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		username    string
		date        string
		wantErr     bool
	}{
		{"valid perfect day", "Great Day", "testuser", "2023-12-01", false},
		{"empty title", "", "testuser", "2023-12-01", true},
		{"empty username", "Great Day", "", "2023-12-01", true},
		{"invalid date", "Great Day", "testuser", "2023-13-01", true},
		{"invalid date format", "Great Day", "testuser", "12/01/2023", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pd, err := models.NewPerfectDay("test-id", tt.title, "description", tt.username, tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPerfectDay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if pd.Title != tt.title {
					t.Errorf("NewPerfectDay().Title = %v, want %v", pd.Title, tt.title)
				}
				if pd.IsDeleted {
					t.Error("NewPerfectDay().IsDeleted should be false initially")
				}
				if len(pd.Areas) != 0 {
					t.Error("NewPerfectDay().Areas should be empty initially")
				}
				if len(pd.Activities) != 0 {
					t.Error("NewPerfectDay().Activities should be empty initially")
				}
			}
		})
	}
}

func TestPerfectDayAddActivity(t *testing.T) {
	pd, _ := models.NewPerfectDay("test-id", "Test Day", "description", "testuser", "2023-12-01")
	location1 := models.NewCustomTextLocation("Cafe", "Shibuya")
	location2 := models.NewCustomTextLocation("Restaurant", "Shinjuku")

	activity1, _ := models.NewActivity("act1", "Coffee", *location1, "09:00", 60, "", "")
	activity2, _ := models.NewActivity("act2", "Lunch", *location2, "12:00", 90, "", "")

	initialUpdatedAt := pd.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	pd.AddActivity(*activity1)
	if len(pd.Activities) != 1 {
		t.Errorf("Expected 1 activity, got %d", len(pd.Activities))
	}
	if len(pd.Areas) != 1 || pd.Areas[0] != "Shibuya" {
		t.Errorf("Expected areas [Shibuya], got %v", pd.Areas)
	}
	if !pd.UpdatedAt.After(initialUpdatedAt) {
		t.Error("UpdatedAt should be updated after adding activity")
	}

	pd.AddActivity(*activity2)
	if len(pd.Activities) != 2 {
		t.Errorf("Expected 2 activities, got %d", len(pd.Activities))
	}
	if len(pd.Areas) != 2 {
		t.Errorf("Expected 2 areas, got %d", len(pd.Areas))
	}

	pd.SortActivitiesByTime()
	if pd.Activities[0].StartTime != "09:00" {
		t.Errorf("Expected first activity at 09:00, got %s", pd.Activities[0].StartTime)
	}
	if pd.Activities[1].StartTime != "12:00" {
		t.Errorf("Expected second activity at 12:00, got %s", pd.Activities[1].StartTime)
	}
}

func TestPerfectDaySoftDelete(t *testing.T) {
	pd, _ := models.NewPerfectDay("test-id", "Test Day", "description", "testuser", "2023-12-01")

	if pd.IsDeleted {
		t.Error("PerfectDay should not be deleted initially")
	}

	initialUpdatedAt := pd.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	pd.SoftDelete()
	if !pd.IsDeleted {
		t.Error("PerfectDay should be deleted after SoftDelete()")
	}
	if !pd.UpdatedAt.After(initialUpdatedAt) {
		t.Error("UpdatedAt should be updated after soft delete")
	}
}

func TestPerfectDaySearchableContent(t *testing.T) {
	pd, _ := models.NewPerfectDay("test-id", "Great Coffee Day", "A day about coffee", "testuser", "2023-12-01")
	location := models.NewCustomTextLocation("Blue Bottle Coffee", "Shibuya")
	activity, _ := models.NewActivity("act1", "Coffee Tasting", *location, "10:00", 60, "Trying new beans", "Amazing experience")

	pd.AddActivity(*activity)

	content := pd.SearchableContent()
	expectedTerms := []string{"great coffee day", "coffee", "shibuya", "coffee tasting", "blue bottle coffee", "trying new beans", "amazing experience"}

	for _, term := range expectedTerms {
		if !containsIgnoreCase(content, term) {
			t.Errorf("Searchable content should contain '%s', got: %s", term, content)
		}
	}
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(substr) > 0 && indexIgnoreCase(s, substr) >= 0))
}

func indexIgnoreCase(s, substr string) int {
	return strings.Index(strings.ToLower(s), strings.ToLower(substr))
}
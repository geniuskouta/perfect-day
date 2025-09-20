package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"perfect-day/internal/api/server"
	"perfect-day/pkg/config"
	"perfect-day/pkg/models"
	"strings"
	"testing"
	"time"
)

func TestPlacesSearch(t *testing.T) {
	// Create test server with unique data directory
	timestamp := time.Now().UnixNano()
	testDataDir := fmt.Sprintf("/tmp/test-places-%d", timestamp)

	cfg := &config.Config{
		DataDir: testDataDir,
	}

	testServer := server.NewServer(cfg)

	// Test 1: Search places without query (should fail)
	t.Run("search_places_no_query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/places/search", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test 2: Search places with query (should return empty results since no API key)
	t.Run("search_places_with_query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/places/search?q=coffee", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		places := data["places"].([]interface{})
		if len(places) != 0 {
			t.Errorf("Expected empty places array without API key, got %d places", len(places))
		}

		// Should have notice about API unavailability
		meta := response["meta"].(map[string]interface{})
		if notice, exists := meta["notice"]; exists {
			if !containsIgnoreCase(notice.(string), "unavailable") {
				t.Errorf("Expected notice about API unavailability, got: %v", notice)
			}
		}
	})

	// Test 3: Search places with limit parameter
	t.Run("search_places_with_limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/places/search?q=restaurant&limit=5", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		if data["limit"].(float64) != 5 {
			t.Errorf("Expected limit 5, got %v", data["limit"])
		}
		if data["query"] != "restaurant" {
			t.Errorf("Expected query 'restaurant', got %v", data["query"])
		}
	})
}

func TestAreas(t *testing.T) {
	// Create test server with unique data directory
	timestamp := time.Now().UnixNano()
	testDataDir := fmt.Sprintf("/tmp/test-areas-%d", timestamp)

	cfg := &config.Config{
		DataDir: testDataDir,
	}

	testServer := server.NewServer(cfg)

	// Create test data for areas endpoint
	user, _ := models.NewUser("testuser", "UTC")
	testServer.Storage.UserStorage.Save(user)

	pd1, _ := models.NewPerfectDay("pd1", "Day in Shibuya", "Great day", "testuser", "2023-12-01")
	pd2, _ := models.NewPerfectDay("pd2", "Day in Shinjuku", "Another day", "testuser", "2023-12-02")

	location1 := models.NewCustomTextLocation("Cafe", "Shibuya")
	location2 := models.NewCustomTextLocation("Restaurant", "Shinjuku")
	location3 := models.NewCustomTextLocation("Park", "Shibuya")

	activity1, _ := models.NewActivity("act1", "Coffee", *location1, "10:00", 60, "", "")
	activity2, _ := models.NewActivity("act2", "Lunch", *location2, "12:00", 90, "", "")
	activity3, _ := models.NewActivity("act3", "Walk", *location3, "15:00", 30, "", "")

	pd1.AddActivity(*activity1)
	pd1.AddActivity(*activity3)
	pd2.AddActivity(*activity2)

	testServer.Storage.PerfectDayStorage.Save(pd1)
	testServer.Storage.PerfectDayStorage.Save(pd2)

	// Test: Get areas (should return unique areas from perfect days)
	t.Run("get_areas", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/areas", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		areas := data["areas"].([]interface{})

		// Should have Shibuya and Shinjuku
		if len(areas) != 2 {
			t.Errorf("Expected 2 unique areas, got %d", len(areas))
		}

		areaNames := make([]string, len(areas))
		for i, area := range areas {
			areaNames[i] = area.(string)
		}

		if !contains(areaNames, "Shibuya") {
			t.Errorf("Expected Shibuya in areas, got %v", areaNames)
		}
		if !contains(areaNames, "Shinjuku") {
			t.Errorf("Expected Shinjuku in areas, got %v", areaNames)
		}
	})

	// Test: Get areas when no data exists
	t.Run("get_areas_empty", func(t *testing.T) {
		// Create fresh server with no data
		timestamp := time.Now().UnixNano()
		emptyDataDir := fmt.Sprintf("/tmp/test-areas-empty-%d", timestamp)

		emptyCfg := &config.Config{
			DataDir: emptyDataDir,
		}

		emptyServer := server.NewServer(emptyCfg)

		req := httptest.NewRequest("GET", "/api/v1/areas", nil)
		w := httptest.NewRecorder()

		emptyServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		areas := data["areas"].([]interface{})

		if len(areas) != 0 {
			t.Errorf("Expected 0 areas when no data exists, got %d", len(areas))
		}
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(substr) > 0 && indexIgnoreCase(s, substr) >= 0))
}

func indexIgnoreCase(s, substr string) int {
	return strings.Index(strings.ToLower(s), strings.ToLower(substr))
}
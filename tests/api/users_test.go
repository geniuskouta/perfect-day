package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"perfect-day/internal/api/server"
	"perfect-day/pkg/config"
	"perfect-day/pkg/models"
	"testing"
	"time"
)

func TestUserProfile(t *testing.T) {
	// Create test server with unique data directory
	timestamp := time.Now().UnixNano()
	testDataDir := fmt.Sprintf("/tmp/test-users-%d", timestamp)

	cfg := &config.Config{
		DataDir: testDataDir,
	}

	testServer := server.NewServer(cfg)

	// Create test users
	user1, _ := models.NewUser("alice", "UTC")
	user2, _ := models.NewUser("bob", "Asia/Tokyo")

	testServer.Storage.UserStorage.Save(user1)
	testServer.Storage.UserStorage.Save(user2)

	// Test 1: Get existing user profile
	t.Run("get_existing_user_profile", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/alice", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		if data["username"] != "alice" {
			t.Errorf("Expected username alice, got %v", data["username"])
		}
		if data["timezone"] != "UTC" {
			t.Errorf("Expected timezone UTC, got %v", data["timezone"])
		}
	})

	// Test 2: Get non-existent user profile
	t.Run("get_nonexistent_user_profile", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/charlie", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestUserPerfectDays(t *testing.T) {
	// Create test server with unique data directory
	timestamp := time.Now().UnixNano()
	testDataDir := fmt.Sprintf("/tmp/test-user-perfect-days-%d", timestamp)

	cfg := &config.Config{
		DataDir: testDataDir,
	}

	testServer := server.NewServer(cfg)

	// Create test user and perfect days
	user, _ := models.NewUser("alice", "UTC")
	testServer.Storage.UserStorage.Save(user)

	// Create perfect days for alice
	pd1, _ := models.NewPerfectDay("pd1", "Alice's Day 1", "Great day", "alice", "2023-12-01")
	pd2, _ := models.NewPerfectDay("pd2", "Alice's Day 2", "Another great day", "alice", "2023-12-02")

	location := models.NewCustomTextLocation("Test Location", "Shibuya")
	activity, _ := models.NewActivity("act1", "Test Activity", *location, "10:00", 60, "desc", "commentary")
	pd1.AddActivity(*activity)

	testServer.Storage.PerfectDayStorage.Save(pd1)
	testServer.Storage.PerfectDayStorage.Save(pd2)

	// Test 1: Get user's perfect days
	t.Run("get_user_perfect_days", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/alice/perfect-days", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		if data["total"].(float64) != 2 {
			t.Errorf("Expected 2 perfect days, got %v", data["total"])
		}
	})

	// Test 2: Get perfect days for non-existent user
	t.Run("get_nonexistent_user_perfect_days", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/charlie/perfect-days", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	// Test 3: Get user's perfect days with pagination
	t.Run("get_user_perfect_days_paginated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/alice/perfect-days?limit=1&offset=0", nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		if data["total"].(float64) != 2 {
			t.Errorf("Expected total 2, got %v", data["total"])
		}
		if data["limit"].(float64) != 1 {
			t.Errorf("Expected limit 1, got %v", data["limit"])
		}
		if data["offset"].(float64) != 0 {
			t.Errorf("Expected offset 0, got %v", data["offset"])
		}

		perfectDays := data["perfect_days"].([]interface{})
		if len(perfectDays) != 1 {
			t.Errorf("Expected 1 perfect day in page, got %d", len(perfectDays))
		}
	})
}
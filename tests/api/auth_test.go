package api

import (
	"bytes"
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

func TestAuthenticationWorkflow(t *testing.T) {
	// Create test server with unique data directory
	timestamp := time.Now().UnixNano()
	testDataDir := fmt.Sprintf("/tmp/test-auth-%d", timestamp)

	cfg := &config.Config{
		DataDir: testDataDir,
	}

	testServer := server.NewServer(cfg)

	// Create a test user
	user, err := models.NewUser("testuser", "UTC")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	err = testServer.Storage.UserStorage.Save(user)
	if err != nil {
		t.Fatalf("Failed to save test user: %v", err)
	}

	// Test 1: Try to create perfect day without authentication (should fail)
	t.Run("create_without_auth", func(t *testing.T) {
		perfectDayReq := map[string]interface{}{
			"title":       "Test Day",
			"description": "Test description",
			"date":        "2023-12-01",
			"activities":  []interface{}{},
		}

		reqBody, _ := json.Marshal(perfectDayReq)
		req := httptest.NewRequest("POST", "/api/v1/perfect-days", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	// Test 2: Login to get session
	var sessionCookie string
	t.Run("login", func(t *testing.T) {
		loginReq := map[string]string{
			"username": "testuser",
		}

		reqBody, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Login failed with status %d: %s", w.Code, w.Body.String())
		}

		// Extract session cookie
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie.Value
				break
			}
		}

		if sessionCookie == "" {
			t.Fatal("No session cookie received")
		}
	})

	// Test 3: Create perfect day with authentication (should succeed)
	var createdPerfectDayID string
	t.Run("create_with_auth", func(t *testing.T) {
		perfectDayReq := map[string]interface{}{
			"title":       "Test Day",
			"description": "Test description",
			"date":        "2023-12-01",
			"activities":  []interface{}{},
		}

		reqBody, _ := json.Marshal(perfectDayReq)
		req := httptest.NewRequest("POST", "/api/v1/perfect-days", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionCookie})
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		createdPerfectDayID = data["id"].(string)
	})

	// Test 4: Update perfect day with authentication (should succeed)
	t.Run("update_with_auth", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"title":       "Updated Test Day",
			"description": "Updated description",
			"date":        "2023-12-01",
			"activities":  []interface{}{},
		}

		reqBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", "/api/v1/perfect-days/"+createdPerfectDayID, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionCookie})
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	// Test 5: Update perfect day without authentication (should fail)
	t.Run("update_without_auth", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"title":       "Unauthorized Update",
			"description": "Should fail",
			"date":        "2023-12-01",
			"activities":  []interface{}{},
		}

		reqBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", "/api/v1/perfect-days/"+createdPerfectDayID, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	// Test 6: Delete perfect day with authentication (should succeed)
	t.Run("delete_with_auth", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/perfect-days/"+createdPerfectDayID, nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: sessionCookie})
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d: %s", http.StatusNoContent, w.Code, w.Body.String())
		}
	})

	// Test 7: Delete perfect day without authentication (should fail)
	t.Run("delete_without_auth", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/perfect-days/"+createdPerfectDayID, nil)
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestOwnershipVerification(t *testing.T) {
	// Create test server with unique data directory
	timestamp := time.Now().UnixNano()
	testDataDir := fmt.Sprintf("/tmp/test-ownership-%d", timestamp)

	cfg := &config.Config{
		DataDir: testDataDir,
	}

	testServer := server.NewServer(cfg)

	// Create two test users
	user1, _ := models.NewUser("user1", "UTC")
	user2, _ := models.NewUser("user2", "UTC")

	testServer.Storage.UserStorage.Save(user1)
	testServer.Storage.UserStorage.Save(user2)

	// Login as user1 and create a perfect day
	var user1Session, user2Session string
	var perfectDayID string

	// Login user1
	loginReq := map[string]string{"username": "user1"}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testServer.ServeHTTP(w, req)

	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			user1Session = cookie.Value
			break
		}
	}

	// Login user2
	loginReq = map[string]string{"username": "user2"}
	reqBody, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	testServer.ServeHTTP(w, req)

	cookies = w.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			user2Session = cookie.Value
			break
		}
	}

	// User1 creates a perfect day
	perfectDayReq := map[string]interface{}{
		"title":       "User1's Day",
		"description": "Created by user1",
		"date":        "2023-12-01",
		"activities":  []interface{}{},
	}

	reqBody, _ = json.Marshal(perfectDayReq)
	req = httptest.NewRequest("POST", "/api/v1/perfect-days", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: user1Session})
	w = httptest.NewRecorder()
	testServer.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]interface{})
	perfectDayID = data["id"].(string)

	// Test: User2 tries to update user1's perfect day (should fail)
	t.Run("cross_user_update_forbidden", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"title":       "Hijacked by User2",
			"description": "Should not work",
			"date":        "2023-12-01",
			"activities":  []interface{}{},
		}

		reqBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PUT", "/api/v1/perfect-days/"+perfectDayID, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "session_id", Value: user2Session})
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})

	// Test: User2 tries to delete user1's perfect day (should fail)
	t.Run("cross_user_delete_forbidden", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/perfect-days/"+perfectDayID, nil)
		req.AddCookie(&http.Cookie{Name: "session_id", Value: user2Session})
		w := httptest.NewRecorder()

		testServer.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d: %s", http.StatusForbidden, w.Code, w.Body.String())
		}
	})
}
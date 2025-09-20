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
	"strings"
	"testing"
	"time"
)

func setupTestServer() *server.Server {
	cfg := &config.Config{
		DataDir: "/tmp/perfect-day-test-" + time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%d", time.Now().UnixNano()),
	}
	return server.NewServer(cfg)
}

func createTestUser(srv *server.Server, username string) {
	user, _ := models.NewUser(username, "Asia/Tokyo")
	srv.Storage.UserStorage.Save(user)
}

func TestCreatePerfectDay(t *testing.T) {
	srv := setupTestServer()
	createTestUser(srv, "testuser")

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid perfect day",
			requestBody: map[string]interface{}{
				"title":       "Great Tokyo Day",
				"description": "Amazing exploration",
				"date":        "2025-01-15",
				"activities": []map[string]interface{}{
					{
						"name": "Visit Temple",
						"location": map[string]interface{}{
							"type": "custom_text",
							"name": "Senso-ji Temple",
							"area": "Asakusa",
						},
						"start_time":  "09:00",
						"duration":    120,
						"description": "Morning visit",
						"commentary":  "Beautiful experience",
					},
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing title",
			requestBody: map[string]interface{}{
				"description": "Day without title",
				"date":        "2025-01-15",
				"activities":  []map[string]interface{}{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name: "invalid date",
			requestBody: map[string]interface{}{
				"title":      "Bad Date Day",
				"date":       "invalid-date",
				"activities": []map[string]interface{}{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/perfect-days", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})

			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)

			if tt.expectedError != "" {
				if errorData, ok := response["error"].(map[string]interface{}); ok {
					if code, ok := errorData["code"].(string); !ok || code != tt.expectedError {
						t.Errorf("Expected error code %s, got %v", tt.expectedError, errorData["code"])
					}
				} else {
					t.Errorf("Expected error response, got %v", response)
				}
			} else {
				if data, ok := response["data"].(map[string]interface{}); ok {
					if title, ok := data["title"].(string); !ok || title != tt.requestBody["title"] {
						t.Errorf("Expected title %v, got %v", tt.requestBody["title"], data["title"])
					}
					if id, ok := data["id"].(string); !ok || id == "" {
						t.Error("Expected non-empty ID in response")
					}
				} else {
					t.Errorf("Expected data in response, got %v", response)
				}
			}
		})
	}
}

func TestGetPerfectDay(t *testing.T) {
	srv := setupTestServer()
	createTestUser(srv, "testuser")

	// Create a test perfect day
	pd, _ := models.NewPerfectDay("test-id-123", "Test Day", "Test description", "testuser", "2025-01-15")
	location := models.NewCustomTextLocation("Test Location", "Test Area")
	activity, _ := models.NewActivity("act-1", "Test Activity", *location, "10:00", 60, "Test desc", "Great time")
	pd.AddActivity(*activity)
	srv.Storage.PerfectDayStorage.Save(pd)

	tests := []struct {
		name           string
		perfectDayID   string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "existing perfect day",
			perfectDayID:   "test-id-123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existent perfect day",
			perfectDayID:   "non-existent-id",
			expectedStatus: http.StatusNotFound,
			expectedError:  "NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/perfect-days/"+tt.perfectDayID, nil)
			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)

			if tt.expectedError != "" {
				if errorData, ok := response["error"].(map[string]interface{}); ok {
					if code, ok := errorData["code"].(string); !ok || code != tt.expectedError {
						t.Errorf("Expected error code %s, got %v", tt.expectedError, errorData["code"])
					}
				} else {
					t.Errorf("Expected error response, got %v", response)
				}
			} else {
				if data, ok := response["data"].(map[string]interface{}); ok {
					if id, ok := data["id"].(string); !ok || id != tt.perfectDayID {
						t.Errorf("Expected ID %s, got %v", tt.perfectDayID, data["id"])
					}
					if title, ok := data["title"].(string); !ok || title != "Test Day" {
						t.Errorf("Expected title 'Test Day', got %v", data["title"])
					}
				} else {
					t.Errorf("Expected data in response, got %v", response)
				}
			}
		})
	}
}

func TestUpdatePerfectDay(t *testing.T) {
	srv := setupTestServer()
	createTestUser(srv, "testuser")

	// Create a test perfect day
	pd, _ := models.NewPerfectDay("test-id-update", "Original Title", "Original description", "testuser", "2025-01-15")
	srv.Storage.PerfectDayStorage.Save(pd)

	tests := []struct {
		name           string
		perfectDayID   string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:         "valid update",
			perfectDayID: "test-id-update",
			requestBody: map[string]interface{}{
				"title":       "Updated Title",
				"description": "Updated description",
				"date":        "2025-01-16",
				"activities":  []map[string]interface{}{},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "non-existent perfect day",
			perfectDayID: "non-existent-id",
			requestBody: map[string]interface{}{
				"title":      "Updated Title",
				"date":       "2025-01-16",
				"activities": []map[string]interface{}{},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/v1/perfect-days/"+tt.perfectDayID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})

			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Response: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)

			if tt.expectedError != "" {
				if errorData, ok := response["error"].(map[string]interface{}); ok {
					if code, ok := errorData["code"].(string); !ok || code != tt.expectedError {
						t.Errorf("Expected error code %s, got %v", tt.expectedError, errorData["code"])
					}
				} else {
					t.Errorf("Expected error response, got %v", response)
				}
			} else {
				if data, ok := response["data"].(map[string]interface{}); ok {
					if title, ok := data["title"].(string); !ok || title != tt.requestBody["title"] {
						t.Errorf("Expected title %v, got %v", tt.requestBody["title"], data["title"])
					}
				} else {
					t.Errorf("Expected data in response, got %v", response)
				}
			}
		})
	}
}

func TestDeletePerfectDay(t *testing.T) {
	srv := setupTestServer()
	createTestUser(srv, "testuser")

	// Create a test perfect day
	pd, _ := models.NewPerfectDay("test-id-delete", "To Delete", "Will be deleted", "testuser", "2025-01-15")
	srv.Storage.PerfectDayStorage.Save(pd)

	tests := []struct {
		name           string
		perfectDayID   string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "existing perfect day",
			perfectDayID:   "test-id-delete",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "non-existent perfect day",
			perfectDayID:   "non-existent-id",
			expectedStatus: http.StatusNotFound,
			expectedError:  "NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/v1/perfect-days/"+tt.perfectDayID, nil)
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})

			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(rr.Body.Bytes(), &response)
				if errorData, ok := response["error"].(map[string]interface{}); ok {
					if code, ok := errorData["code"].(string); !ok || code != tt.expectedError {
						t.Errorf("Expected error code %s, got %v", tt.expectedError, errorData["code"])
					}
				}
			}
		})
	}
}

func TestListPerfectDays(t *testing.T) {
	srv := setupTestServer()
	createTestUser(srv, "testuser")
	createTestUser(srv, "otheruser")

	// Create test data
	pd1, _ := models.NewPerfectDay("list-test-1", "First Day", "First description", "testuser", "2025-01-15")
	pd2, _ := models.NewPerfectDay("list-test-2", "Second Day", "Second description", "testuser", "2025-01-16")
	pd3, _ := models.NewPerfectDay("list-test-3", "Other User Day", "Other description", "otheruser", "2025-01-17")

	srv.Storage.PerfectDayStorage.Save(pd1)
	srv.Storage.PerfectDayStorage.Save(pd2)
	srv.Storage.PerfectDayStorage.Save(pd3)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "list all perfect days",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "list with limit",
			queryParams:    "?limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "list by user",
			queryParams:    "?user=testuser",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "list by non-existent user",
			queryParams:    "?user=nonexistent",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/v1/perfect-days"+tt.queryParams, nil)
			rr := httptest.NewRecorder()
			srv.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			var response map[string]interface{}
			json.Unmarshal(rr.Body.Bytes(), &response)

			if data, ok := response["data"].(map[string]interface{}); ok {
				if perfectDays, ok := data["perfect_days"].([]interface{}); ok {
					if len(perfectDays) != tt.expectedCount {
						t.Errorf("Expected %d perfect days, got %d", tt.expectedCount, len(perfectDays))
					}
				} else {
					t.Error("Expected perfect_days array in response")
				}

				if pagination, ok := data["pagination"].(map[string]interface{}); ok {
					if total, ok := pagination["total"].(float64); ok {
						if strings.Contains(tt.queryParams, "user=testuser") && int(total) != 2 {
							t.Errorf("Expected total 2 for testuser, got %v", total)
						}
					}
				} else {
					t.Error("Expected pagination in response")
				}
			} else {
				t.Errorf("Expected data in response, got %v", response)
			}
		})
	}
}
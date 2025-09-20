package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"perfect-day/internal/api/server"
	"perfect-day/pkg/config"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		DataDir: "/tmp/perfect-day-test",
	}

	// Create server
	srv := server.NewServer(cfg)

	// Create test request
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	srv.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health endpoint returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response content type
	expected := "application/json; charset=utf-8"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("Health endpoint returned wrong content type: got %v want %v",
			ct, expected)
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse JSON response: %v", err)
	}

	// Check response structure
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("Response should have 'data' field")
	}

	if status, ok := data["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", data["status"])
	}

	if version, ok := data["version"].(string); !ok || version != "0.1.0" {
		t.Errorf("Expected version '0.1.0', got %v", data["version"])
	}

	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Error("Response should have 'meta' field")
	}

	if timestamp, ok := meta["timestamp"].(string); !ok || timestamp == "" {
		t.Error("Meta should have non-empty timestamp")
	}
}

func TestVersionEndpoint(t *testing.T) {
	cfg := &config.Config{
		DataDir: "/tmp/perfect-day-test",
	}

	srv := server.NewServer(cfg)

	req, err := http.NewRequest("GET", "/api/v1/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Version endpoint returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not parse JSON response: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("Response should have 'data' field")
	}

	if version, ok := data["version"].(string); !ok || version != "0.1.0" {
		t.Errorf("Expected version '0.1.0', got %v", data["version"])
	}
}
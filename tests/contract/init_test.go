package contract

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("init", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Set up Perfect Day configuration")
	result.AssertStdoutContains(t, "--api-key")
	result.AssertStdoutContains(t, "--data-dir")
	result.AssertStdoutContains(t, "--interactive")
}

func TestInitWithAPIKey(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("init", "--api-key", "test-key-123")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Configuration saved")
	result.AssertStdoutContains(t, "Google Places API: Configured")
	result.AssertStdoutContains(t, "Perfect Day is ready to use")
}

func TestInitWithDataDir(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	customDir := filepath.Join(helper.tempDir, "custom-data")
	result := helper.ExecuteCommand("init", "--data-dir", customDir)
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Configuration saved")
	result.AssertStdoutContains(t, "Data Directory: "+customDir)

	// Verify directory was created
	if _, err := os.Stat(customDir); os.IsNotExist(err) {
		t.Errorf("Custom data directory was not created: %s", customDir)
	}
}

func TestInitInteractiveEmpty(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommandWithInput("\n\n", "init", "--interactive")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Perfect Day Configuration Setup")
	result.AssertStdoutContains(t, "Configuration saved")
	result.AssertStdoutContains(t, "Google Places API: Not configured")
}

func TestInitInteractiveWithValues(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	customDir := filepath.Join(helper.tempDir, "interactive-data")
	input := "interactive-api-key\n" + customDir + "\n"

	result := helper.ExecuteCommandWithInput(input, "init", "--interactive")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Configuration saved")
	result.AssertStdoutContains(t, "Data Directory: "+customDir)
	result.AssertStdoutContains(t, "Google Places API: Configured")
}

func TestInitCreatesValidConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("init", "--api-key", "test-key", "--data-dir", helper.tempDir)
	result.AssertExitCode(t, 0)

	// Verify config file was created and is valid
	configFile := filepath.Join(helper.tempDir, ".perfect-day", "config.json")
	configData, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(configData, &config); err != nil {
		t.Fatalf("Config file is not valid JSON: %v", err)
	}

	if config["google_places_api_key"] != "test-key" {
		t.Errorf("Expected API key 'test-key', got %v", config["google_places_api_key"])
	}

	if config["data_directory"] != helper.tempDir {
		t.Errorf("Expected data directory %s, got %v", helper.tempDir, config["data_directory"])
	}
}
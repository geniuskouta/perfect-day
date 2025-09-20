package contract

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestHelper provides utilities for contract testing CLI commands
type TestHelper struct {
	t           *testing.T
	tempDir     string
	binaryPath  string
	originalDir string
}

// NewTestHelper creates a new test helper with temporary directory and binary
func NewTestHelper(t *testing.T) *TestHelper {
	tempDir, err := os.MkdirTemp("", "perfectday-contract-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Get original directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Build binary for testing
	binaryPath := filepath.Join(tempDir, "perfectday")
	if err := buildBinary(binaryPath); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to build binary: %v", err)
	}

	return &TestHelper{
		t:           t,
		tempDir:     tempDir,
		binaryPath:  binaryPath,
		originalDir: originalDir,
	}
}

// Cleanup removes temporary directory and files
func (h *TestHelper) Cleanup() {
	os.RemoveAll(h.tempDir)
}

// CommandResult holds the result of executing a CLI command
type CommandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
}

// ExecuteCommand runs the perfectday binary with given arguments
func (h *TestHelper) ExecuteCommand(args ...string) *CommandResult {
	cmd := exec.Command(h.binaryPath, args...)

	// Set environment variables for test
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PERFECT_DAY_DATA_DIR=%s", h.tempDir),
		fmt.Sprintf("HOME=%s", h.tempDir), // Override home directory for config
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return &CommandResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Error:    err,
	}
}

// ExecuteCommandWithInput runs command with stdin input
func (h *TestHelper) ExecuteCommandWithInput(input string, args ...string) *CommandResult {
	cmd := exec.Command(h.binaryPath, args...)

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PERFECT_DAY_DATA_DIR=%s", h.tempDir),
		fmt.Sprintf("HOME=%s", h.tempDir), // Override home directory for config
	)

	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return &CommandResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Error:    err,
	}
}

// AssertExitCode checks the command exit code
func (r *CommandResult) AssertExitCode(t *testing.T, expected int) {
	if r.ExitCode != expected {
		t.Errorf("Expected exit code %d, got %d\nStdout: %s\nStderr: %s",
			expected, r.ExitCode, r.Stdout, r.Stderr)
	}
}

// AssertContains checks if output contains expected string
func (r *CommandResult) AssertContains(t *testing.T, expected string) {
	combined := r.Stdout + r.Stderr
	if !strings.Contains(combined, expected) {
		t.Errorf("Expected output to contain %q\nActual output:\nStdout: %s\nStderr: %s",
			expected, r.Stdout, r.Stderr)
	}
}

// AssertNotContains checks if output does not contain string
func (r *CommandResult) AssertNotContains(t *testing.T, notExpected string) {
	combined := r.Stdout + r.Stderr
	if strings.Contains(combined, notExpected) {
		t.Errorf("Expected output to NOT contain %q\nActual output:\nStdout: %s\nStderr: %s",
			notExpected, r.Stdout, r.Stderr)
	}
}

// AssertStdoutContains checks if stdout contains expected string
func (r *CommandResult) AssertStdoutContains(t *testing.T, expected string) {
	if !strings.Contains(r.Stdout, expected) {
		t.Errorf("Expected stdout to contain %q\nActual stdout: %s", expected, r.Stdout)
	}
}

// AssertStderrContains checks if stderr contains expected string
func (r *CommandResult) AssertStderrContains(t *testing.T, expected string) {
	if !strings.Contains(r.Stderr, expected) {
		t.Errorf("Expected stderr to contain %q\nActual stderr: %s", expected, r.Stderr)
	}
}

// buildBinary builds the perfectday binary for testing
func buildBinary(outputPath string) error {
	// Get current working directory and find project root
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	// Navigate to project root (from wherever we are)
	projectRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return fmt.Errorf("could not find project root with go.mod")
		}
		projectRoot = parent
	}

	cmd := exec.Command("go", "build", "-o", outputPath, ".")
	cmd.Dir = projectRoot

	// Capture build output for debugging
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %v\nBuild output: %s", err, stderr.String())
	}

	return nil
}

// SetupConfigFile creates a test config file in the temp directory
func (h *TestHelper) SetupConfigFile(apiKey, dataDir string) error {
	configDir := filepath.Join(h.tempDir, ".perfect-day")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configContent := fmt.Sprintf(`{
  "google_places_api_key": "%s",
  "data_directory": "%s"
}`, apiKey, dataDir)

	configFile := filepath.Join(configDir, "config.json")
	return os.WriteFile(configFile, []byte(configContent), 0600)
}

// CreateTestDataFiles creates sample data files for testing
func (h *TestHelper) CreateTestDataFiles() error {
	// First set up config to use our temp directory as data directory
	if err := h.SetupConfigFile("", h.tempDir); err != nil {
		return err
	}

	dataDir := h.tempDir

	// Create users directory with test user
	usersDir := filepath.Join(dataDir, "users")
	if err := os.MkdirAll(usersDir, 0755); err != nil {
		return err
	}

	testUser := `{
  "username": "testuser",
  "timezone": "UTC",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}`

	userFile := filepath.Join(usersDir, "testuser.json")
	if err := os.WriteFile(userFile, []byte(testUser), 0644); err != nil {
		return err
	}

	// Create current_user file
	currentUserFile := filepath.Join(dataDir, "current_user")
	if err := os.WriteFile(currentUserFile, []byte("testuser"), 0644); err != nil {
		return err
	}

	// Create perfect_days directory with test data
	perfectDaysDir := filepath.Join(dataDir, "perfect_days", "testuser")
	if err := os.MkdirAll(perfectDaysDir, 0755); err != nil {
		return err
	}

	testPerfectDay := `{
  "id": "test-perfect-day-123",
  "title": "Test Perfect Day",
  "description": "A test perfect day for contract testing",
  "username": "testuser",
  "date": "2024-01-01",
  "areas": ["Test Area"],
  "activities": [
    {
      "name": "Test Activity",
      "description": "A test activity",
      "location": {
        "type": "custom_text",
        "name": "Test Location",
        "area": "Test Area"
      },
      "start_time": "09:00",
      "duration": 60,
      "commentary": "This was great!"
    }
  ],
  "is_deleted": false,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}`

	perfectDayFile := filepath.Join(perfectDaysDir, "test-perfect-day-123.json")
	return os.WriteFile(perfectDayFile, []byte(testPerfectDay), 0644)
}

// TestCase represents a single test case for contract testing
type TestCase struct {
	Name        string
	Args        []string
	Input       string
	WantExitCode int
	WantStdout  []string
	WantStderr  []string
	NotWantOut  []string
	Setup       func(*TestHelper) error
}

// RunTestCases executes a slice of test cases
func (h *TestHelper) RunTestCases(testCases []TestCase) {
	for _, tc := range testCases {
		h.t.Run(tc.Name, func(t *testing.T) {
			// Setup for this test case
			if tc.Setup != nil {
				if err := tc.Setup(h); err != nil {
					t.Fatalf("Test setup failed: %v", err)
				}
			}

			// Execute command
			var result *CommandResult
			if tc.Input != "" {
				result = h.ExecuteCommandWithInput(tc.Input, tc.Args...)
			} else {
				result = h.ExecuteCommand(tc.Args...)
			}

			// Check exit code
			result.AssertExitCode(t, tc.WantExitCode)

			// Check stdout contains
			for _, want := range tc.WantStdout {
				result.AssertStdoutContains(t, want)
			}

			// Check stderr contains
			for _, want := range tc.WantStderr {
				result.AssertStderrContains(t, want)
			}

			// Check not contains
			for _, notWant := range tc.NotWantOut {
				result.AssertNotContains(t, notWant)
			}
		})
	}
}
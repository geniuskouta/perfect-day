package integration

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestEditCommandHelp(t *testing.T) {
	output, err := runCLI("edit", "--help")
	if err != nil {
		t.Fatalf("Edit help command failed: %v", err)
	}

	if !strings.Contains(output, "Edit an existing perfect day") {
		t.Errorf("Expected edit help description, got: %s", output)
	}

	if !strings.Contains(output, "edit <ID>") {
		t.Errorf("Expected edit usage, got: %s", output)
	}
}

func TestEditCommandWithoutArgs(t *testing.T) {
	_, err := runCLI("edit")
	if err == nil {
		t.Error("Expected error when running edit without ID")
	}
}

func TestEditCommandWithNonExistentID(t *testing.T) {
	tempDir := t.TempDir()
	buildBinary(t, tempDir)
	binaryPath := filepath.Join(tempDir, "perfect-day")

	output, err := runCLIWithEnv(binaryPath, tempDir, "edit", "nonexistent")
	if err == nil {
		t.Error("Expected error when editing non-existent perfect day")
	}

	if !strings.Contains(output, "Please login first") {
		// If no login, should get login error first
		if !strings.Contains(output, "not found") {
			t.Errorf("Expected error message about not found or login, got: %s", output)
		}
	}
}
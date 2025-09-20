package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIVersion(t *testing.T) {
	output, err := runCLI("version")
	if err != nil {
		t.Fatalf("CLI version command failed: %v", err)
	}

	if !strings.Contains(output, "perfect-day version 0.1.0") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

func TestCLIHelp(t *testing.T) {
	output, err := runCLI("--help")
	if err != nil {
		t.Fatalf("CLI help command failed: %v", err)
	}

	expectedCommands := []string{"create", "delete", "list", "login", "search", "show", "version"}
	for _, cmd := range expectedCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("Expected command '%s' in help output", cmd)
		}
	}
}

func TestCLIWorkflow(t *testing.T) {
	tempDir := t.TempDir()

	buildBinary(t, tempDir)

	binaryPath := filepath.Join(tempDir, "perfect-day")

	t.Run("list without login", func(t *testing.T) {
		output, err := runCLIWithEnv(binaryPath, tempDir, "list")
		if err == nil {
			t.Error("Expected error when listing without login")
		}
		if !strings.Contains(output, "Please login first") {
			t.Errorf("Expected login prompt, got: %s", output)
		}
	})

	t.Run("show non-existent", func(t *testing.T) {
		output, err := runCLIWithEnv(binaryPath, tempDir, "show", "nonexistent")
		if err == nil {
			t.Error("Expected error when showing non-existent perfect day")
		}
		if !strings.Contains(output, "not found") {
			t.Errorf("Expected not found message, got: %s", output)
		}
	})
}

func buildBinary(t *testing.T, tempDir string) {
	cmd := exec.Command("go", "build", "-o", filepath.Join(tempDir, "perfect-day"), ".")
	cmd.Dir = getProjectRoot()
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
}

func runCLI(args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = getProjectRoot()
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func runCLIWithEnv(binaryPath, dataDir string, args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Env = append(os.Environ(), "PERFECT_DAY_DATA_DIR="+dataDir)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func getProjectRoot() string {
	wd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			panic("Could not find project root")
		}
		wd = parent
	}
}
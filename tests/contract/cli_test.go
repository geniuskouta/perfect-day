package contract

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIHelp(t *testing.T) {
	// Build CLI if not exists
	cliPath := "./../../bin/perfectday-cli"
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		cmd := exec.Command("go", "build", "-o", cliPath, "./../../cmd/perfectday-cli")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to build CLI: %v", err)
		}
	}

	// Test help command
	cmd := exec.Command(cliPath, "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI help command failed: %v", err)
	}

	outputStr := string(output)

	// Check expected content
	expectedStrings := []string{
		"Perfect Day",
		"Usage:",
		"Available Commands:",
		"create",
		"list",
		"show",
		"version",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Help output should contain '%s', got: %s", expected, outputStr)
		}
	}
}

func TestCLIVersion(t *testing.T) {
	cliPath := "./../../bin/perfectday-cli"

	cmd := exec.Command(cliPath, "version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI version command failed: %v", err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "perfect-day version") {
		t.Errorf("Version output should contain 'perfect-day version', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "0.1.0") {
		t.Errorf("Version output should contain '0.1.0', got: %s", outputStr)
	}
}

func TestCLIList(t *testing.T) {
	cliPath := "./../../bin/perfectday-cli"

	cmd := exec.Command(cliPath, "list")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI list command failed: %v", err)
	}

	outputStr := string(output)

	// Should contain either "Found X perfect days" or "No perfect days found"
	if !strings.Contains(outputStr, "Found") && !strings.Contains(outputStr, "No perfect days found") {
		t.Errorf("List output should show results or no results message, got: %s", outputStr)
	}
}

func TestCLIListWithUser(t *testing.T) {
	cliPath := "./../../bin/perfectday-cli"

	cmd := exec.Command(cliPath, "list", "--user", "kouta")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI list --user command failed: %v", err)
	}

	outputStr := string(output)

	// Should contain either "Found X perfect days" or "No perfect days found"
	if !strings.Contains(outputStr, "Found") && !strings.Contains(outputStr, "No perfect days found") {
		t.Errorf("List output should show results or no results message, got: %s", outputStr)
	}
}
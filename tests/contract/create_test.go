package contract

import (
	"testing"
)

func TestCreateHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("create", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Create a new perfect day")
	result.AssertStdoutContains(t, "interactive prompts")
}

func TestCreateRequiresLogin(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// No current user set up
	result := helper.ExecuteCommand("create")
	result.AssertExitCode(t, 1)
	result.AssertStdoutContains(t, "Please login first")
}

// Note: Interactive create flows are tested in integration tests
// Contract tests focus on CLI interface validation only
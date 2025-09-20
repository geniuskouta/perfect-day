package contract

import (
	"testing"
)

func TestEditHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("edit", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Edit an existing perfect day")
	result.AssertStdoutContains(t, "Edit an existing perfect day")
}

func TestEditRequiresLogin(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("edit", "some-id")
	result.AssertExitCode(t, 1)
	result.AssertStdoutContains(t, "Please login first")
}

func TestEditWithoutID(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("edit")
	result.AssertExitCode(t, 1)
	result.AssertStderrContains(t, "accepts 1 arg(s), received 0")
}

// Note: Interactive edit flows and data validation are tested in integration tests
// Contract tests focus on CLI interface validation only
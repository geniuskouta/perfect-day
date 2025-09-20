package contract

import (
	"testing"
)

func TestDeleteHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("delete", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Soft delete a perfect day")
	result.AssertStdoutContains(t, "Soft delete")
}

func TestDeleteRequiresLogin(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// No current user set up
	result := helper.ExecuteCommand("delete", "some-id")
	result.AssertExitCode(t, 1)
	result.AssertStdoutContains(t, "Please login first")
}

// Note: Data validation and actual delete operations are tested in integration tests
// Contract tests focus on CLI interface validation only


func TestDeleteWithoutID(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("delete")
	result.AssertExitCode(t, 1)
	result.AssertStderrContains(t, "accepts 1 arg(s), received 0")
}
package contract

import (
	"testing"
)

func TestShowHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("show", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Display detailed information")
	result.AssertStdoutContains(t, "Display detailed information")
}

// Note: Data retrieval and display logic are tested in integration tests
// Contract tests focus on CLI interface validation only


func TestShowWithoutID(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("show")
	result.AssertExitCode(t, 1)
	result.AssertStderrContains(t, "accepts 1 arg(s), received 0")
}
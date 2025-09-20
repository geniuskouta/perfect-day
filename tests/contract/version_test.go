package contract

import (
	"testing"
)

func TestVersionHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("version", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "version")
}

func TestVersionCommand(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("version")
	result.AssertExitCode(t, 0)
	// Should show version information
}
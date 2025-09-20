package contract

import (
	"testing"
)

func TestListHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("list", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "List perfect days")
	result.AssertStdoutContains(t, "--user")
	result.AssertStdoutContains(t, "--all")
	result.AssertStdoutContains(t, "--deleted")
}

func TestListWithoutData(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// No test data setup - should show no results
	result := helper.ExecuteCommand("list")
	result.AssertExitCode(t, 1)
	result.AssertStdoutContains(t, "Please login first")
}

func TestListAllUsers(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("list", "--all")
	result.AssertExitCode(t, 0)
	// Should exit successfully even with no data
}

func TestListSpecificUser(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("list", "--user", "someuser")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "No perfect days found")
}

func TestListNonExistentUser(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("list", "--user", "nonexistent")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "No perfect days found")
}
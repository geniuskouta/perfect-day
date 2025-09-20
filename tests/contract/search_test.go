package contract

import (
	"testing"
)

func TestSearchHelp(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("search", "--help")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "Search perfect days")
	result.AssertStdoutContains(t, "--query")
	result.AssertStdoutContains(t, "--areas")
	result.AssertStdoutContains(t, "--user")
	result.AssertStdoutContains(t, "--from")
	result.AssertStdoutContains(t, "--to")
	result.AssertStdoutContains(t, "--sort")
	result.AssertStdoutContains(t, "--limit")
}

func TestSearchByQuery(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("search", "--query", "test")
	result.AssertExitCode(t, 0)
	// Should succeed even with no data
}

func TestSearchByArea(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("search", "--areas", "Test Area")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "No perfect days found")
}

func TestSearchByUser(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("search", "--user", "testuser")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "No perfect days found")
}

func TestSearchNoResults(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("search", "--query", "nonexistent")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "No perfect days found")
}

func TestSearchWithLimit(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	result := helper.ExecuteCommand("search", "--query", "test", "--limit", "1")
	result.AssertExitCode(t, 0)
	result.AssertStdoutContains(t, "No perfect days found")
}
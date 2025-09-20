package unit

import (
	"perfect-day/src/models"
	"perfect-day/src/storage"
	"testing"
	"time"
)

func TestPerfectDayUpdateAreas(t *testing.T) {
	pd, err := models.NewPerfectDay("test-id", "Test Day", "description", "testuser", "2023-12-01")
	if err != nil {
		t.Fatalf("Failed to create perfect day: %v", err)
	}

	location1 := models.NewCustomTextLocation("Cafe", "Shibuya")
	location2 := models.NewCustomTextLocation("Restaurant", "Shinjuku")

	activity1, _ := models.NewActivity("act1", "Coffee", *location1, "09:00", 60, "", "")
	activity2, _ := models.NewActivity("act2", "Lunch", *location2, "12:00", 90, "", "")

	pd.Activities = append(pd.Activities, *activity1, *activity2)

	initialUpdatedAt := pd.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	pd.UpdateAreas()

	if len(pd.Areas) != 2 {
		t.Errorf("Expected 2 areas, got %d", len(pd.Areas))
	}

	expectedAreas := []string{"Shibuya", "Shinjuku"}
	for i, expected := range expectedAreas {
		if i >= len(pd.Areas) || pd.Areas[i] != expected {
			t.Errorf("Expected area %s at position %d, got %v", expected, i, pd.Areas)
		}
	}

	if !pd.UpdatedAt.After(initialUpdatedAt) {
		t.Error("UpdatedAt should be updated after UpdateAreas()")
	}
}

func TestEditCommandLoadPerfectDay(t *testing.T) {
	tempDir := t.TempDir()
	storage := storage.NewStorage(tempDir)

	user, _ := models.NewUser("testuser", "UTC")
	storage.UserStorage.Save(user)

	pd, _ := models.NewPerfectDay("test-id", "Test Day", "description", "testuser", "2023-12-01")
	location := models.NewCustomTextLocation("Test Location", "Test Area")
	activity, _ := models.NewActivity("act1", "Test Activity", *location, "10:00", 60, "", "")
	pd.AddActivity(*activity)

	storage.PerfectDayStorage.Save(pd)

	// Test loading with full ID
	loadedPD, err := storage.PerfectDayStorage.Load("testuser", "test-id")
	if err != nil {
		t.Fatalf("Failed to load perfect day: %v", err)
	}

	if loadedPD.ID != pd.ID {
		t.Errorf("Expected ID %s, got %s", pd.ID, loadedPD.ID)
	}
	if loadedPD.Title != pd.Title {
		t.Errorf("Expected title %s, got %s", pd.Title, loadedPD.Title)
	}
}

func TestEditValidationFunctions(t *testing.T) {
	// Test date validation
	validDates := []string{"2023-12-01", "2024-01-15", "2025-06-30"}
	for _, date := range validDates {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			t.Errorf("Valid date %s should not return error: %v", date, err)
		}
	}

	invalidDates := []string{"2023-13-01", "12/01/2023", "2023-12-32", "invalid"}
	for _, date := range invalidDates {
		_, err := time.Parse("2006-01-02", date)
		if err == nil {
			t.Errorf("Invalid date %s should return error", date)
		}
	}

	// Test time validation
	validTimes := []string{"09:00", "12:30", "23:59", "00:00"}
	for _, timeStr := range validTimes {
		_, err := time.Parse("15:04", timeStr)
		if err != nil {
			t.Errorf("Valid time %s should not return error: %v", timeStr, err)
		}
	}

	invalidTimes := []string{"25:00", "12:60", "invalid", "9:0"}
	for _, timeStr := range invalidTimes {
		_, err := time.Parse("15:04", timeStr)
		if err == nil {
			t.Errorf("Invalid time %s should return error", timeStr)
		}
	}
}

func TestPerfectDayActivityManipulation(t *testing.T) {
	pd, _ := models.NewPerfectDay("test-id", "Test Day", "description", "testuser", "2023-12-01")

	location1 := models.NewCustomTextLocation("Cafe", "Shibuya")
	location2 := models.NewCustomTextLocation("Restaurant", "Shinjuku")
	location3 := models.NewCustomTextLocation("Park", "Harajuku")

	activity1, _ := models.NewActivity("act1", "Coffee", *location1, "09:00", 60, "", "")
	activity2, _ := models.NewActivity("act2", "Lunch", *location2, "12:00", 90, "", "")
	activity3, _ := models.NewActivity("act3", "Walk", *location3, "14:00", 30, "", "")

	// Test adding activities
	pd.AddActivity(*activity1)
	pd.AddActivity(*activity3)
	pd.AddActivity(*activity2)

	if len(pd.Activities) != 3 {
		t.Errorf("Expected 3 activities, got %d", len(pd.Activities))
	}

	// Test sorting by time
	pd.SortActivitiesByTime()
	expectedOrder := []string{"09:00", "12:00", "14:00"}
	for i, expected := range expectedOrder {
		if pd.Activities[i].StartTime != expected {
			t.Errorf("Expected activity at position %d to start at %s, got %s",
				i, expected, pd.Activities[i].StartTime)
		}
	}

	// Test removing activity (simulate removing middle activity)
	pd.Activities = append(pd.Activities[:1], pd.Activities[2:]...)
	if len(pd.Activities) != 2 {
		t.Errorf("Expected 2 activities after removal, got %d", len(pd.Activities))
	}

	// Verify remaining activities
	if pd.Activities[0].StartTime != "09:00" || pd.Activities[1].StartTime != "14:00" {
		t.Error("Wrong activities remain after removal")
	}

	// Test updating areas after activity removal
	pd.UpdateAreas()
	expectedAreas := []string{"Harajuku", "Shibuya"}
	if len(pd.Areas) != 2 {
		t.Errorf("Expected 2 areas after update, got %d", len(pd.Areas))
	}
	for i, expected := range expectedAreas {
		if i >= len(pd.Areas) || pd.Areas[i] != expected {
			t.Errorf("Expected area %s at position %d, got %v", expected, i, pd.Areas)
		}
	}
}

func TestEditDeletedPerfectDay(t *testing.T) {
	pd, _ := models.NewPerfectDay("test-id", "Test Day", "description", "testuser", "2023-12-01")
	pd.SoftDelete()

	if !pd.IsDeleted {
		t.Error("Perfect day should be marked as deleted")
	}

	// In a real edit command, we would check pd.IsDeleted and refuse to edit
	// This test verifies the soft delete functionality works as expected
}

func TestActivityModification(t *testing.T) {
	location := models.NewCustomTextLocation("Test Location", "Test Area")
	activity, _ := models.NewActivity("act1", "Original Activity", *location, "10:00", 60, "original desc", "original commentary")

	// Test modifying activity properties (simulating edit functionality)
	activity.Name = "Modified Activity"
	activity.StartTime = "11:00"
	activity.Duration = 90
	activity.Description = "modified desc"
	activity.Commentary = "modified commentary"

	if activity.Name != "Modified Activity" {
		t.Errorf("Expected modified name, got %s", activity.Name)
	}
	if activity.StartTime != "11:00" {
		t.Errorf("Expected modified start time, got %s", activity.StartTime)
	}
	if activity.Duration != 90 {
		t.Errorf("Expected modified duration, got %d", activity.Duration)
	}

	// Test end time calculation with new duration
	endTime, err := activity.EndTime()
	if err != nil {
		t.Fatalf("Failed to calculate end time: %v", err)
	}
	if endTime != "12:30" {
		t.Errorf("Expected end time 12:30, got %s", endTime)
	}
}
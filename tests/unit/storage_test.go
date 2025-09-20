package unit

import (
	"os"
	"path/filepath"
	"perfect-day/src/models"
	"perfect-day/src/storage"
	"testing"
)

func TestUserStorageSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	userStorage := storage.NewUserStorage(tempDir)

	user, err := models.NewUser("testuser", "America/New_York")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := userStorage.Save(user); err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	if !userStorage.Exists("testuser") {
		t.Error("User should exist after saving")
	}

	loadedUser, err := userStorage.Load("testuser")
	if err != nil {
		t.Fatalf("Failed to load user: %v", err)
	}

	if loadedUser.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, loadedUser.Username)
	}
	if loadedUser.Timezone != user.Timezone {
		t.Errorf("Expected timezone %s, got %s", user.Timezone, loadedUser.Timezone)
	}
	if !loadedUser.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("Expected created_at %v, got %v", user.CreatedAt, loadedUser.CreatedAt)
	}
}

func TestUserStorageNotFound(t *testing.T) {
	tempDir := t.TempDir()
	userStorage := storage.NewUserStorage(tempDir)

	if userStorage.Exists("nonexistent") {
		t.Error("Non-existent user should not exist")
	}

	_, err := userStorage.Load("nonexistent")
	if err == nil {
		t.Error("Loading non-existent user should return error")
	}
}

func TestPerfectDayStorageSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	pdStorage := storage.NewPerfectDayStorage(tempDir)

	pd, err := models.NewPerfectDay("test-id", "Test Day", "Test description", "testuser", "2023-12-01")
	if err != nil {
		t.Fatalf("Failed to create perfect day: %v", err)
	}

	location := models.NewCustomTextLocation("Test Location", "Test Area")
	activity, _ := models.NewActivity("act1", "Test Activity", *location, "10:00", 60, "description", "commentary")
	pd.AddActivity(*activity)

	if err := pdStorage.Save(pd); err != nil {
		t.Fatalf("Failed to save perfect day: %v", err)
	}

	loadedPD, err := pdStorage.Load("testuser", "test-id")
	if err != nil {
		t.Fatalf("Failed to load perfect day: %v", err)
	}

	if loadedPD.ID != pd.ID {
		t.Errorf("Expected ID %s, got %s", pd.ID, loadedPD.ID)
	}
	if loadedPD.Title != pd.Title {
		t.Errorf("Expected title %s, got %s", pd.Title, loadedPD.Title)
	}
	if len(loadedPD.Activities) != len(pd.Activities) {
		t.Errorf("Expected %d activities, got %d", len(pd.Activities), len(loadedPD.Activities))
	}
	if len(loadedPD.Areas) != len(pd.Areas) {
		t.Errorf("Expected %d areas, got %d", len(pd.Areas), len(loadedPD.Areas))
	}
}

func TestPerfectDayStorageLoadAllByUser(t *testing.T) {
	tempDir := t.TempDir()
	pdStorage := storage.NewPerfectDayStorage(tempDir)

	pd1, _ := models.NewPerfectDay("id1", "Day 1", "", "testuser", "2023-12-01")
	pd2, _ := models.NewPerfectDay("id2", "Day 2", "", "testuser", "2023-12-02")
	pd3, _ := models.NewPerfectDay("id3", "Day 3", "", "otheruser", "2023-12-03")
	pd2.SoftDelete()

	pdStorage.Save(pd1)
	pdStorage.Save(pd2)
	pdStorage.Save(pd3)

	userPDs, err := pdStorage.LoadAllByUser("testuser", false)
	if err != nil {
		t.Fatalf("Failed to load user perfect days: %v", err)
	}
	if len(userPDs) != 1 {
		t.Errorf("Expected 1 non-deleted perfect day, got %d", len(userPDs))
	}

	userPDsWithDeleted, err := pdStorage.LoadAllByUser("testuser", true)
	if err != nil {
		t.Fatalf("Failed to load user perfect days with deleted: %v", err)
	}
	if len(userPDsWithDeleted) != 2 {
		t.Errorf("Expected 2 perfect days including deleted, got %d", len(userPDsWithDeleted))
	}
}

func TestPerfectDayStorageLoadAll(t *testing.T) {
	tempDir := t.TempDir()
	pdStorage := storage.NewPerfectDayStorage(tempDir)

	pd1, _ := models.NewPerfectDay("id1", "Day 1", "", "user1", "2023-12-01")
	pd2, _ := models.NewPerfectDay("id2", "Day 2", "", "user2", "2023-12-02")
	pd2.SoftDelete()

	pdStorage.Save(pd1)
	pdStorage.Save(pd2)

	allPDs, err := pdStorage.LoadAll(false)
	if err != nil {
		t.Fatalf("Failed to load all perfect days: %v", err)
	}
	if len(allPDs) != 1 {
		t.Errorf("Expected 1 non-deleted perfect day, got %d", len(allPDs))
	}

	allPDsWithDeleted, err := pdStorage.LoadAll(true)
	if err != nil {
		t.Fatalf("Failed to load all perfect days with deleted: %v", err)
	}
	if len(allPDsWithDeleted) != 2 {
		t.Errorf("Expected 2 perfect days including deleted, got %d", len(allPDsWithDeleted))
	}
}

func TestPerfectDayStorageDelete(t *testing.T) {
	tempDir := t.TempDir()
	pdStorage := storage.NewPerfectDayStorage(tempDir)

	pd, _ := models.NewPerfectDay("test-id", "Test Day", "", "testuser", "2023-12-01")
	pdStorage.Save(pd)

	filePath := filepath.Join(tempDir, "perfect-days", "testuser", "test-id.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Perfect day file should exist before deletion")
	}

	if err := pdStorage.Delete("testuser", "test-id"); err != nil {
		t.Fatalf("Failed to delete perfect day: %v", err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Perfect day file should not exist after deletion")
	}

	if err := pdStorage.Delete("testuser", "nonexistent"); err != nil {
		t.Errorf("Deleting non-existent file should not return error, got: %v", err)
	}
}

func TestStorageInitialization(t *testing.T) {
	tempDir := t.TempDir()
	storage := storage.NewStorage(tempDir)

	if storage.GetDataDir() != tempDir {
		t.Errorf("Expected data dir %s, got %s", tempDir, storage.GetDataDir())
	}

	if err := storage.Initialize(); err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Data directory should be created after initialization")
	}
}

func TestStorageDefaultDataDir(t *testing.T) {
	storage := storage.NewStorage("")
	dataDir := storage.GetDataDir()

	if dataDir == "" {
		t.Error("Default data directory should not be empty")
	}

	expectedSuffix := ".perfect-day"
	if filepath.Base(dataDir) != expectedSuffix {
		t.Errorf("Expected data directory to end with %s, got %s", expectedSuffix, dataDir)
	}
}
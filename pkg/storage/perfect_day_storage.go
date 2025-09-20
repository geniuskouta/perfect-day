package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"perfect-day/pkg/models"
	"strings"
)

type PerfectDayStorage struct {
	dataDir string
}

func NewPerfectDayStorage(dataDir string) *PerfectDayStorage {
	return &PerfectDayStorage{dataDir: dataDir}
}

func (pds *PerfectDayStorage) Save(perfectDay *models.PerfectDay) error {
	if err := pds.ensureDataDir(); err != nil {
		return err
	}

	userDir := filepath.Join(pds.dataDir, "perfect-days", perfectDay.Username)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return fmt.Errorf("failed to create user directory: %v", err)
	}

	filePath := filepath.Join(userDir, perfectDay.ID+".json")
	data, err := json.MarshalIndent(perfectDay, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal perfect day: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write perfect day file: %v", err)
	}

	return nil
}

func (pds *PerfectDayStorage) Load(username, id string) (*models.PerfectDay, error) {
	filePath := filepath.Join(pds.dataDir, "perfect-days", username, id+".json")

	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("perfect day not found: %s/%s", username, id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read perfect day file: %v", err)
	}

	var perfectDay models.PerfectDay
	if err := json.Unmarshal(data, &perfectDay); err != nil {
		return nil, fmt.Errorf("failed to unmarshal perfect day: %v", err)
	}

	return &perfectDay, nil
}

func (pds *PerfectDayStorage) LoadAllByUser(username string, includeDeleted bool) ([]*models.PerfectDay, error) {
	userDir := filepath.Join(pds.dataDir, "perfect-days", username)

	entries, err := os.ReadDir(userDir)
	if os.IsNotExist(err) {
		return []*models.PerfectDay{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read user directory: %v", err)
	}

	var perfectDays []*models.PerfectDay
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")
		perfectDay, err := pds.Load(username, id)
		if err != nil {
			continue
		}

		if !includeDeleted && perfectDay.IsDeleted {
			continue
		}

		perfectDays = append(perfectDays, perfectDay)
	}

	return perfectDays, nil
}

func (pds *PerfectDayStorage) LoadAll(includeDeleted bool) ([]*models.PerfectDay, error) {
	perfectDaysDir := filepath.Join(pds.dataDir, "perfect-days")

	userEntries, err := os.ReadDir(perfectDaysDir)
	if os.IsNotExist(err) {
		return []*models.PerfectDay{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read perfect-days directory: %v", err)
	}

	var allPerfectDays []*models.PerfectDay
	for _, userEntry := range userEntries {
		if !userEntry.IsDir() {
			continue
		}

		username := userEntry.Name()
		userPerfectDays, err := pds.LoadAllByUser(username, includeDeleted)
		if err != nil {
			continue
		}

		allPerfectDays = append(allPerfectDays, userPerfectDays...)
	}

	return allPerfectDays, nil
}

func (pds *PerfectDayStorage) Delete(username, id string) error {
	filePath := filepath.Join(pds.dataDir, "perfect-days", username, id+".json")

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete perfect day file: %v", err)
	}

	return nil
}

func (pds *PerfectDayStorage) ensureDataDir() error {
	return os.MkdirAll(pds.dataDir, 0755)
}
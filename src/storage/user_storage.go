package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"perfect-day/src/models"
)

type UserStorage struct {
	dataDir string
}

func NewUserStorage(dataDir string) *UserStorage {
	return &UserStorage{dataDir: dataDir}
}

func (us *UserStorage) Save(user *models.User) error {
	if err := us.ensureDataDir(); err != nil {
		return err
	}

	filePath := filepath.Join(us.dataDir, "users", user.Username+".json")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create user directory: %v", err)
	}

	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write user file: %v", err)
	}

	return nil
}

func (us *UserStorage) Load(username string) (*models.User, error) {
	filePath := filepath.Join(us.dataDir, "users", username+".json")

	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("user not found: %s", username)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read user file: %v", err)
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil
}

func (us *UserStorage) Exists(username string) bool {
	filePath := filepath.Join(us.dataDir, "users", username+".json")
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (us *UserStorage) ensureDataDir() error {
	return os.MkdirAll(us.dataDir, 0755)
}
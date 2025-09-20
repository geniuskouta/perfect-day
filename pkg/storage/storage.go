package storage

import (
	"os"
	"path/filepath"
)

type Storage struct {
	UserStorage       *UserStorage
	PerfectDayStorage *PerfectDayStorage
	dataDir           string
}

func NewStorage(dataDir string) *Storage {
	if dataDir == "" {
		homeDir, _ := os.UserHomeDir()
		dataDir = filepath.Join(homeDir, ".perfect-day")
	}

	return &Storage{
		UserStorage:       NewUserStorage(dataDir),
		PerfectDayStorage: NewPerfectDayStorage(dataDir),
		dataDir:           dataDir,
	}
}

func (s *Storage) GetDataDir() string {
	return s.dataDir
}

func (s *Storage) Initialize() error {
	return os.MkdirAll(s.dataDir, 0755)
}
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Username         string `json:"username"`
	Timezone         string `json:"timezone"`
	DataDir          string `json:"data_dir"`
	GooglePlacesAPIKey string `json:"google_places_api_key,omitempty"`
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		homeDir, _ := os.UserHomeDir()
		configPath = filepath.Join(homeDir, ".perfect-day", "config.json")
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		homeDir, _ := os.UserHomeDir()
		configPath = filepath.Join(homeDir, ".perfect-day", "config.json")
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func NewConfig(username, timezone string) *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".perfect-day")

	return &Config{
		Username: username,
		Timezone: timezone,
		DataDir:  dataDir,
	}
}
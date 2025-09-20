package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	GooglePlacesAPIKey string `json:"google_places_api_key,omitempty"`
	DataDirectory      string `json:"data_directory,omitempty"`
}

var (
	initAPIKey    string
	initDataDir   string
	initInteractive bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Perfect Day configuration",
	Long:  "Set up Perfect Day configuration including Google Places API key and data directory.",
	Run:   runInit,
}

func init() {
	initCmd.Flags().StringVar(&initAPIKey, "api-key", "", "Google Places API key")
	initCmd.Flags().StringVar(&initDataDir, "data-dir", "", "Data directory path")
	initCmd.Flags().BoolVarP(&initInteractive, "interactive", "i", false, "Interactive setup")
}

func runInit(cmd *cobra.Command, args []string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".perfect-day")
	configFile := filepath.Join(configDir, "config.json")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}

	// Load existing config or create new one
	config := &Config{}
	if existingData, err := os.ReadFile(configFile); err == nil {
		if err := json.Unmarshal(existingData, config); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading existing config: %v\n", err)
			os.Exit(1)
		}
	}

	// Interactive mode or flag-based setup
	if initInteractive || (initAPIKey == "" && initDataDir == "") {
		runInteractiveSetup(config)
	} else {
		if initAPIKey != "" {
			config.GooglePlacesAPIKey = initAPIKey
		}
		if initDataDir != "" {
			config.DataDirectory = initDataDir
		}
	}

	// Set default data directory if not specified
	if config.DataDirectory == "" {
		config.DataDirectory = configDir
	}

	// Save config
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling config: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(configFile, configData, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration saved to: %s\n", configFile)
	fmt.Println("\nConfiguration:")
	fmt.Printf("  Data Directory: %s\n", config.DataDirectory)
	if config.GooglePlacesAPIKey != "" {
		fmt.Printf("  Google Places API: Configured\n")
	} else {
		fmt.Printf("  Google Places API: Not configured (custom locations only)\n")
	}

	// Create data directory
	if err := os.MkdirAll(config.DataDirectory, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nPerfect Day is ready to use! Run 'perfectday create' to get started.\n")
}

func runInteractiveSetup(config *Config) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🌟 Perfect Day Configuration Setup")
	fmt.Println("Press Enter to keep existing values or leave blank for defaults.\n")

	// Google Places API Key
	currentAPI := config.GooglePlacesAPIKey
	if currentAPI != "" {
		currentAPI = "***configured***"
	}
	fmt.Printf("Google Places API Key [%s]: ", currentAPI)
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey != "" {
		config.GooglePlacesAPIKey = apiKey
	}

	// Data Directory
	defaultDataDir := filepath.Join(getUserHomeDir(), ".perfect-day")
	currentDataDir := config.DataDirectory
	if currentDataDir == "" {
		currentDataDir = defaultDataDir
	}
	fmt.Printf("Data Directory [%s]: ", currentDataDir)
	dataDir, _ := reader.ReadString('\n')
	dataDir = strings.TrimSpace(dataDir)
	if dataDir != "" {
		config.DataDirectory = dataDir
	} else if config.DataDirectory == "" {
		config.DataDirectory = defaultDataDir
	}

	fmt.Println()
}

func getUserHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~"
	}
	return homeDir
}

func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".perfect-day", "config.json")
}

func LoadConfig() (*Config, error) {
	configFile := GetConfigPath()

	config := &Config{}
	if data, err := os.ReadFile(configFile); err == nil {
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("error parsing config file: %v", err)
		}
	}

	// Apply environment variable overrides
	if envAPIKey := os.Getenv("GOOGLE_PLACES_API_KEY"); envAPIKey != "" {
		config.GooglePlacesAPIKey = envAPIKey
	}
	if envDataDir := os.Getenv("PERFECT_DAY_DATA_DIR"); envDataDir != "" {
		config.DataDirectory = envDataDir
	}

	// Set default data directory if still empty
	if config.DataDirectory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		config.DataDirectory = filepath.Join(homeDir, ".perfect-day")
	}

	return config, nil
}
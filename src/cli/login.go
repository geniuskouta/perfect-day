package cli

import (
	"fmt"
	"os"
	"perfect-day/src/lib"
	"perfect-day/src/models"
	"perfect-day/src/storage"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login or create a new user",
	Long:  "Login with an existing username or create a new user with timezone information.",
	Run:   runLogin,
}

func runLogin(cmd *cobra.Command, args []string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	storage := storage.NewStorage(config.DataDirectory)
	if err := storage.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	username := lib.PromptInput("Username: ")
	if username == "" {
		fmt.Println("Username cannot be empty")
		os.Exit(1)
	}

	if storage.UserStorage.Exists(username) {
		user, err := storage.UserStorage.Load(username)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading user: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Welcome back, %s! (Timezone: %s)\n", user.Username, user.Timezone)
		saveCurrentUser(username)
	} else {
		fmt.Printf("User '%s' not found. Let's create a new account.\n", username)
		timezone := lib.PromptInput("Timezone (e.g., Asia/Tokyo, America/New_York): ")
		if timezone == "" {
			timezone = "UTC"
		}

		user, err := models.NewUser(username, timezone)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating user: %v\n", err)
			os.Exit(1)
		}

		if err := storage.UserStorage.Save(user); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving user: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Welcome, %s! Your account has been created with timezone: %s\n", user.Username, user.Timezone)
		saveCurrentUser(username)
	}
}

func saveCurrentUser(username string) {
	config, _ := LoadConfig()
	storage := storage.NewStorage(config.DataDirectory)
	currentUserFile := fmt.Sprintf("%s/current_user", storage.GetDataDir())

	if err := os.WriteFile(currentUserFile, []byte(username), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not save current user: %v\n", err)
	}
}

func getCurrentUser() string {
	config, _ := LoadConfig()
	storage := storage.NewStorage(config.DataDirectory)
	currentUserFile := fmt.Sprintf("%s/current_user", storage.GetDataDir())

	data, err := os.ReadFile(currentUserFile)
	if err != nil {
		return ""
	}
	return string(data)
}
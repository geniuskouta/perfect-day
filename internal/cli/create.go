package cli

import (
	"context"
	"fmt"
	"os"
	"perfect-day/pkg/utils"
	"perfect-day/pkg/models"
	"perfect-day/pkg/places"
	"perfect-day/pkg/storage"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new perfect day",
	Long:  "Create a new perfect day with interactive prompts for activities and locations.",
	Run:   runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	username := getCurrentUser()
	if username == "" {
		fmt.Println("Please login first using 'perfect-day login'")
		os.Exit(1)
	}

	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	storage := storage.NewStorage(config.DataDirectory)
	placesService, _ := places.NewPlacesService(config.GooglePlacesAPIKey)

	fmt.Println("Creating a new Perfect Day...")

	title := utils.PromptInput("Title: ")
	if title == "" {
		fmt.Println("Title is required")
		os.Exit(1)
	}

	description := utils.PromptInput("Description (optional): ")

	dateStr := utils.PromptInput("Date (YYYY-MM-DD, or press Enter for today): ")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	perfectDay, err := models.NewPerfectDay(utils.GenerateID(), title, description, username, dateStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating perfect day: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nNow let's add activities to your perfect day...")

	for {
		fmt.Println("\n--- Adding Activity ---")

		activityName := utils.PromptInput("Activity name: ")
		if activityName == "" {
			break
		}

		location := promptForLocation(placesService)
		if location == nil {
			continue
		}

		startTime := utils.PromptInput("Start time (HH:MM): ")
		durationStr := utils.PromptInput("Duration in minutes: ")
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			fmt.Println("Invalid duration, skipping activity")
			continue
		}

		activityDescription := utils.PromptInput("Activity description (optional): ")
		commentary := utils.PromptInput("Personal commentary (optional): ")

		activity, err := models.NewActivity(
			utils.GenerateID(),
			activityName,
			*location,
			startTime,
			duration,
			activityDescription,
			commentary,
		)
		if err != nil {
			fmt.Printf("Error creating activity: %v\n", err)
			continue
		}

		perfectDay.AddActivity(*activity)
		fmt.Printf("Added activity: %s at %s\n", activityName, location.Name)

		if !utils.PromptConfirm("Add another activity?") {
			break
		}
	}

	perfectDay.SortActivitiesByTime()

	if err := storage.PerfectDayStorage.Save(perfectDay); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving perfect day: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nPerfect Day '%s' created successfully!\n", perfectDay.Title)
	fmt.Printf("ID: %s\n", perfectDay.ID)
}

func promptForLocation(placesService *places.PlacesService) *models.Location {
	fmt.Println("Location options:")
	fmt.Println("1. Search Google Places")
	fmt.Println("2. Enter custom location")

	choice := utils.PromptInput("Choose option (1 or 2): ")

	switch choice {
	case "1":
		if !placesService.IsEnabled() {
			fmt.Println("Google Places API is not configured. Please set GOOGLE_PLACES_API_KEY environment variable.")
			return promptForCustomLocation()
		}
		return promptForGooglePlace(placesService)
	case "2":
		return promptForCustomLocation()
	default:
		fmt.Println("Invalid choice, using custom location")
		return promptForCustomLocation()
	}
}

func promptForGooglePlace(placesService *places.PlacesService) *models.Location {
	query := utils.PromptInput("Search for place: ")
	if query == "" {
		return nil
	}

	ctx := context.Background()
	results, err := placesService.SearchPlaces(ctx, query)
	if err != nil {
		fmt.Printf("Error searching places: %v\n", err)
		return promptForCustomLocation()
	}

	if len(results) == 0 {
		fmt.Println("No places found")
		return promptForCustomLocation()
	}

	fmt.Println("\nFound places:")
	for i, result := range results {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s - %s\n", i+1, result.Name, result.Address)
	}

	choiceStr := utils.PromptInput("Select place (1-5) or 0 for custom: ")
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 0 || choice > len(results) {
		return promptForCustomLocation()
	}

	if choice == 0 {
		return promptForCustomLocation()
	}

	selectedPlace := results[choice-1]
	area := utils.PromptInput(fmt.Sprintf("Area (suggested: %s): ", placesService.SuggestAreaFromAddress(selectedPlace.Address)))
	if area == "" {
		area = placesService.SuggestAreaFromAddress(selectedPlace.Address)
	}

	return placesService.CreateLocationFromPlace(selectedPlace, area)
}

func promptForCustomLocation() *models.Location {
	name := utils.PromptInput("Location name: ")
	if name == "" {
		return nil
	}

	area := utils.PromptInput("Area: ")
	return models.NewCustomTextLocation(name, area)
}
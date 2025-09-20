package cli

import (
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

var editCmd = &cobra.Command{
	Use:   "edit <ID>",
	Short: "Edit a perfect day",
	Long:  "Edit an existing perfect day with interactive prompts for all properties.",
	Args:  cobra.ExactArgs(1),
	Run:   runEdit,
}

func runEdit(cmd *cobra.Command, args []string) {
	perfectDayID := args[0]
	currentUser := getCurrentUser()

	if currentUser == "" {
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

	perfectDay, err := loadPerfectDayForEdit(storage, currentUser, perfectDayID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading perfect day: %v\n", err)
		os.Exit(1)
	}

	if perfectDay.IsDeleted {
		fmt.Println("Cannot edit deleted perfect day. Use 'perfect-day list --deleted' to see deleted items.")
		os.Exit(1)
	}

	fmt.Printf("Editing Perfect Day: %s\n", perfectDay.Title)
	fmt.Printf("Current date: %s\n", perfectDay.Date)
	fmt.Printf("Current activities: %d\n", len(perfectDay.Activities))
	fmt.Println()

	for {
		choice := showEditMenu()
		if !handleEditChoice(choice, perfectDay, placesService, storage) {
			break
		}
	}
}

func loadPerfectDayForEdit(storage *storage.Storage, username, perfectDayID string) (*models.PerfectDay, error) {
	perfectDay, err := storage.PerfectDayStorage.Load(username, perfectDayID)
	if err != nil {
		userPerfectDays, err := storage.PerfectDayStorage.LoadAllByUser(username, true)
		if err != nil {
			return nil, fmt.Errorf("failed to load perfect days: %v", err)
		}

		for _, pd := range userPerfectDays {
			if pd.ID == perfectDayID || pd.ID[:8] == perfectDayID {
				return pd, nil
			}
		}
		return nil, fmt.Errorf("perfect day with ID '%s' not found", perfectDayID)
	}
	return perfectDay, nil
}

func showEditMenu() string {
	fmt.Println("=== Edit Menu ===")
	fmt.Println("1. Edit basic info (title, description, date)")
	fmt.Println("2. Manage activities")
	fmt.Println("3. Preview current perfect day")
	fmt.Println("4. Save and exit")
	fmt.Println("5. Exit without saving")
	fmt.Println()

	return utils.PromptInput("Choose an option (1-5): ")
}

func handleEditChoice(choice string, perfectDay *models.PerfectDay, placesService *places.PlacesService, storage *storage.Storage) bool {
	switch choice {
	case "1":
		editBasicInfo(perfectDay)
	case "2":
		manageActivities(perfectDay, placesService)
	case "3":
		previewPerfectDay(perfectDay)
	case "4":
		return saveAndExit(perfectDay, storage)
	case "5":
		fmt.Println("Exiting without saving changes.")
		return false
	default:
		fmt.Println("Invalid choice. Please enter 1-5.")
	}
	return true
}

func editBasicInfo(perfectDay *models.PerfectDay) {
	fmt.Println("\n=== Edit Basic Info ===")

	fmt.Printf("Current title: %s\n", perfectDay.Title)
	newTitle := utils.PromptInput("New title (press Enter to keep current): ")
	if newTitle != "" {
		perfectDay.Title = newTitle
		fmt.Println("Title updated.")
	}

	fmt.Printf("Current description: %s\n", perfectDay.Description)
	newDescription := utils.PromptInput("New description (press Enter to keep current): ")
	if newDescription != "" {
		perfectDay.Description = newDescription
		fmt.Println("Description updated.")
	}

	fmt.Printf("Current date: %s\n", perfectDay.Date)
	newDate := utils.PromptInput("New date (YYYY-MM-DD, press Enter to keep current): ")
	if newDate != "" {
		if err := validateDate(newDate); err != nil {
			fmt.Printf("Invalid date format: %v\n", err)
		} else {
			perfectDay.Date = newDate
			fmt.Println("Date updated.")
		}
	}

	perfectDay.UpdatedAt = time.Now()
	fmt.Println()
}

func manageActivities(perfectDay *models.PerfectDay, placesService *places.PlacesService) {
	for {
		choice := showActivitiesMenu(perfectDay)
		if !handleActivitiesChoice(choice, perfectDay, placesService) {
			break
		}
	}
}

func showActivitiesMenu(perfectDay *models.PerfectDay) string {
	fmt.Println("\n=== Manage Activities ===")
	fmt.Printf("Current activities: %d\n", len(perfectDay.Activities))

	if len(perfectDay.Activities) > 0 {
		fmt.Println("\nCurrent activities:")
		for i, activity := range perfectDay.Activities {
			fmt.Printf("%d. %s at %s (%s)\n",
				i+1, activity.Name, activity.Location.Name,
				utils.FormatTimeRange(activity.StartTime, activity.Duration))
		}
	}

	fmt.Println("\nOptions:")
	fmt.Println("1. Add new activity")
	if len(perfectDay.Activities) > 0 {
		fmt.Println("2. Edit activity")
		fmt.Println("3. Remove activity")
		fmt.Println("4. Reorder activities")
	}
	fmt.Println("9. Back to main menu")

	return utils.PromptInput("Choose an option: ")
}

func handleActivitiesChoice(choice string, perfectDay *models.PerfectDay, placesService *places.PlacesService) bool {
	switch choice {
	case "1":
		addNewActivity(perfectDay, placesService)
	case "2":
		if len(perfectDay.Activities) > 0 {
			editActivity(perfectDay, placesService)
		} else {
			fmt.Println("No activities to edit.")
		}
	case "3":
		if len(perfectDay.Activities) > 0 {
			removeActivity(perfectDay)
		} else {
			fmt.Println("No activities to remove.")
		}
	case "4":
		if len(perfectDay.Activities) > 1 {
			reorderActivities(perfectDay)
		} else {
			fmt.Println("Need at least 2 activities to reorder.")
		}
	case "9":
		return false
	default:
		fmt.Println("Invalid choice.")
	}
	return true
}

func addNewActivity(perfectDay *models.PerfectDay, placesService *places.PlacesService) {
	fmt.Println("\n=== Add New Activity ===")

	activityName := utils.PromptInput("Activity name: ")
	if activityName == "" {
		fmt.Println("Activity name is required.")
		return
	}

	location := promptForLocation(placesService)
	if location == nil {
		fmt.Println("Location is required.")
		return
	}

	startTime := utils.PromptInput("Start time (HH:MM): ")
	durationStr := utils.PromptInput("Duration in minutes: ")
	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		fmt.Println("Invalid duration.")
		return
	}

	description := utils.PromptInput("Description (optional): ")
	commentary := utils.PromptInput("Commentary (optional): ")

	activity, err := models.NewActivity(
		utils.GenerateID(),
		activityName,
		*location,
		startTime,
		duration,
		description,
		commentary,
	)
	if err != nil {
		fmt.Printf("Error creating activity: %v\n", err)
		return
	}

	perfectDay.AddActivity(*activity)
	perfectDay.SortActivitiesByTime()
	fmt.Printf("Activity '%s' added successfully!\n", activityName)
}

func editActivity(perfectDay *models.PerfectDay, placesService *places.PlacesService) {
	fmt.Println("\n=== Edit Activity ===")

	indexStr := utils.PromptInput("Enter activity number to edit: ")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 1 || index > len(perfectDay.Activities) {
		fmt.Println("Invalid activity number.")
		return
	}

	activity := &perfectDay.Activities[index-1]
	fmt.Printf("Editing: %s\n", activity.Name)

	fmt.Printf("Current name: %s\n", activity.Name)
	newName := utils.PromptInput("New name (press Enter to keep): ")
	if newName != "" {
		activity.Name = newName
	}

	if utils.PromptConfirm("Edit location?") {
		newLocation := promptForLocation(placesService)
		if newLocation != nil {
			activity.Location = *newLocation
		}
	}

	fmt.Printf("Current start time: %s\n", activity.StartTime)
	newStartTime := utils.PromptInput("New start time (HH:MM, press Enter to keep): ")
	if newStartTime != "" {
		if err := validateActivityTime(newStartTime); err != nil {
			fmt.Printf("Invalid time format: %v\n", err)
		} else {
			activity.StartTime = newStartTime
		}
	}

	fmt.Printf("Current duration: %d minutes\n", activity.Duration)
	newDurationStr := utils.PromptInput("New duration (minutes, press Enter to keep): ")
	if newDurationStr != "" {
		newDuration, err := strconv.Atoi(newDurationStr)
		if err != nil || newDuration <= 0 {
			fmt.Println("Invalid duration.")
		} else {
			activity.Duration = newDuration
		}
	}

	fmt.Printf("Current description: %s\n", activity.Description)
	newDescription := utils.PromptInput("New description (press Enter to keep): ")
	if newDescription != "" {
		activity.Description = newDescription
	}

	fmt.Printf("Current commentary: %s\n", activity.Commentary)
	newCommentary := utils.PromptInput("New commentary (press Enter to keep): ")
	if newCommentary != "" {
		activity.Commentary = newCommentary
	}

	perfectDay.SortActivitiesByTime()
	perfectDay.UpdatedAt = time.Now()
	fmt.Println("Activity updated successfully!")
}

func removeActivity(perfectDay *models.PerfectDay) {
	fmt.Println("\n=== Remove Activity ===")

	indexStr := utils.PromptInput("Enter activity number to remove: ")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 1 || index > len(perfectDay.Activities) {
		fmt.Println("Invalid activity number.")
		return
	}

	activity := perfectDay.Activities[index-1]
	fmt.Printf("Remove: %s at %s?\n", activity.Name, activity.Location.Name)

	if utils.PromptConfirm("Are you sure?") {
		perfectDay.Activities = append(perfectDay.Activities[:index-1], perfectDay.Activities[index:]...)
		perfectDay.UpdateAreas()
		perfectDay.UpdatedAt = time.Now()
		fmt.Println("Activity removed successfully!")
	}
}

func reorderActivities(perfectDay *models.PerfectDay) {
	fmt.Println("\n=== Reorder Activities ===")
	fmt.Println("Current order:")
	for i, activity := range perfectDay.Activities {
		fmt.Printf("%d. %s (%s)\n", i+1, activity.Name, activity.StartTime)
	}

	if utils.PromptConfirm("Sort by time automatically?") {
		perfectDay.SortActivitiesByTime()
		fmt.Println("Activities sorted by time!")
		return
	}

	// Manual reordering could be implemented here if needed
	fmt.Println("Manual reordering not implemented yet. Activities sorted by time.")
	perfectDay.SortActivitiesByTime()
}

func previewPerfectDay(perfectDay *models.PerfectDay) {
	fmt.Println("\n=== Preview ===")
	printPerfectDayDetails(perfectDay)
	fmt.Println()
}

func saveAndExit(perfectDay *models.PerfectDay, storage *storage.Storage) bool {
	if utils.PromptConfirm("Save changes?") {
		perfectDay.UpdatedAt = time.Now()
		if err := storage.PerfectDayStorage.Save(perfectDay); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving perfect day: %v\n", err)
			return true // Stay in edit mode
		}
		fmt.Println("Perfect day saved successfully!")
	} else {
		fmt.Println("Changes discarded.")
	}
	return false
}

func validateDate(dateStr string) error {
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %v", err)
	}
	return nil
}

func validateActivityTime(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:MM: %v", err)
	}
	return nil
}
package cli

import (
	"fmt"
	"os"
	"perfect-day/src/lib"
	"perfect-day/src/storage"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <ID>",
	Short: "Delete a perfect day",
	Long:  "Soft delete a perfect day (marks as deleted but preserves data).",
	Args:  cobra.ExactArgs(1),
	Run:   runDelete,
}

func runDelete(cmd *cobra.Command, args []string) {
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

	perfectDay, err := storage.PerfectDayStorage.Load(currentUser, perfectDayID)
	if err != nil {
		allPerfectDays, err := storage.PerfectDayStorage.LoadAllByUser(currentUser, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading perfect days: %v\n", err)
			os.Exit(1)
		}

		for _, pd := range allPerfectDays {
			if pd.ID[:8] == perfectDayID {
				perfectDay = pd
				break
			}
		}

		if perfectDay == nil {
			fmt.Printf("Perfect day with ID '%s' not found or you don't have permission to delete it\n", perfectDayID)
			os.Exit(1)
		}
	}

	if perfectDay.IsDeleted {
		fmt.Println("Perfect day is already deleted")
		return
	}

	fmt.Printf("Perfect Day: %s\n", perfectDay.Title)
	fmt.Printf("Date: %s\n", perfectDay.Date)
	fmt.Printf("Activities: %d\n", len(perfectDay.Activities))

	if !lib.PromptConfirm("Are you sure you want to delete this perfect day?") {
		fmt.Println("Delete cancelled")
		return
	}

	perfectDay.SoftDelete()

	if err := storage.PerfectDayStorage.Save(perfectDay); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting perfect day: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Perfect day '%s' has been deleted\n", perfectDay.Title)
}
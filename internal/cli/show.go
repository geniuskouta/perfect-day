package cli

import (
	"fmt"
	"os"
	"perfect-day/pkg/utils"
	"perfect-day/pkg/models"
	"perfect-day/pkg/storage"
	"strings"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <ID>",
	Short: "Show perfect day details",
	Long:  "Display detailed information about a specific perfect day.",
	Args:  cobra.ExactArgs(1),
	Run:   runShow,
}

func runShow(cmd *cobra.Command, args []string) {
	perfectDayID := args[0]
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	storage := storage.NewStorage(config.DataDirectory)

	var perfectDay *models.PerfectDay

	currentUser := getCurrentUser()
	if currentUser != "" {
		perfectDay, err = storage.PerfectDayStorage.Load(currentUser, perfectDayID)
		if err == nil {
			printPerfectDayDetails(perfectDay)
			return
		}
	}

	allPerfectDays, err := storage.PerfectDayStorage.LoadAll(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading perfect days: %v\n", err)
		os.Exit(1)
	}

	for _, pd := range allPerfectDays {
		if pd.ID == perfectDayID || pd.ID[:8] == perfectDayID {
			perfectDay = pd
			break
		}
	}

	if perfectDay == nil {
		fmt.Printf("Perfect day with ID '%s' not found\n", perfectDayID)
		os.Exit(1)
	}

	printPerfectDayDetails(perfectDay)
}

func printPerfectDayDetails(pd *models.PerfectDay) {
	fmt.Printf("Perfect Day: %s\n", pd.Title)
	fmt.Printf("ID: %s\n", pd.ID)
	fmt.Printf("Username: %s\n", pd.Username)
	fmt.Printf("Date: %s\n", pd.Date)

	if pd.Description != "" {
		fmt.Printf("Description: %s\n", pd.Description)
	}

	if len(pd.Areas) > 0 {
		fmt.Printf("Areas: %v\n", pd.Areas)
	}

	if pd.IsDeleted {
		fmt.Println("Status: DELETED")
	}

	fmt.Printf("Created: %s\n", pd.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", pd.UpdatedAt.Format("2006-01-02 15:04:05"))

	if len(pd.Activities) == 0 {
		fmt.Println("\nNo activities found")
		return
	}

	fmt.Printf("\nActivities (%d):\n", len(pd.Activities))
	fmt.Println(strings.Repeat("=", 80))

	for i, activity := range pd.Activities {
		fmt.Printf("\n%d. %s\n", i+1, activity.Name)
		fmt.Printf("   Time: %s (%s)\n",
			utils.FormatTimeRange(activity.StartTime, activity.Duration),
			utils.FormatDuration(activity.Duration))
		fmt.Printf("   Location: %s", activity.Location.Name)
		if activity.Location.Area != "" {
			fmt.Printf(" (%s)", activity.Location.Area)
		}
		if activity.Location.Address != "" {
			fmt.Printf("\n   Address: %s", activity.Location.Address)
		}
		fmt.Println()

		if activity.Description != "" {
			fmt.Printf("   Description: %s\n", activity.Description)
		}
		if activity.Commentary != "" {
			fmt.Printf("   Commentary: %s\n", activity.Commentary)
		}
	}
}
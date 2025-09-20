package cli

import (
	"fmt"
	"os"
	"perfect-day/src/lib"
	"perfect-day/src/models"
	"perfect-day/src/storage"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listUser    string
	listAll     bool
	listDeleted bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List perfect days",
	Long:  "List perfect days for current user or all users.",
	Run:   runList,
}

func init() {
	listCmd.Flags().StringVarP(&listUser, "user", "u", "", "List perfect days for specific user")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "List perfect days from all users")
	listCmd.Flags().BoolVar(&listDeleted, "deleted", false, "Include deleted perfect days")
}

func runList(cmd *cobra.Command, args []string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	storage := storage.NewStorage(config.DataDirectory)

	var perfectDays []*models.PerfectDay

	if listAll {
		perfectDays, err = storage.PerfectDayStorage.LoadAll(listDeleted)
	} else if listUser != "" {
		perfectDays, err = storage.PerfectDayStorage.LoadAllByUser(listUser, listDeleted)
	} else {
		currentUser := getCurrentUser()
		if currentUser == "" {
			fmt.Println("Please login first using 'perfect-day login' or use --all flag")
			os.Exit(1)
		}
		perfectDays, err = storage.PerfectDayStorage.LoadAllByUser(currentUser, listDeleted)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading perfect days: %v\n", err)
		os.Exit(1)
	}

	if len(perfectDays) == 0 {
		fmt.Println("No perfect days found")
		return
	}

	printPerfectDaysList(perfectDays)
}

func printPerfectDaysList(perfectDays []*models.PerfectDay) {
	fmt.Printf("Found %d perfect days:\n\n", len(perfectDays))

	fmt.Printf("%-8s %-20s %-12s %-15s %-20s %s\n",
		"ID", "Title", "Username", "Date", "Areas", "Activities")
	fmt.Println(strings.Repeat("-", 80))

	for _, pd := range perfectDays {
		idShort := pd.ID[:8]
		title := lib.TruncateString(pd.Title, 20)
		username := lib.TruncateString(pd.Username, 12)
		areas := lib.TruncateString(fmt.Sprintf("%v", pd.Areas), 20)
		activityCount := fmt.Sprintf("%d activities", len(pd.Activities))

		if pd.IsDeleted {
			fmt.Printf("%-8s %-20s %-12s %-15s %-20s %s [DELETED]\n",
				idShort, title, username, pd.Date, areas, activityCount)
		} else {
			fmt.Printf("%-8s %-20s %-12s %-15s %-20s %s\n",
				idShort, title, username, pd.Date, areas, activityCount)
		}
	}

	fmt.Println("\nUse 'perfect-day show <ID>' to view details")
}
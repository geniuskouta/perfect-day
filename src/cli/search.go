package cli

import (
	"fmt"
	"os"
	"perfect-day/src/lib"
	"perfect-day/src/models"
	"perfect-day/src/search"
	"perfect-day/src/storage"
	"strings"

	"github.com/spf13/cobra"
)

var (
	searchQuery    string
	searchAreas    []string
	searchUser     string
	searchDateFrom string
	searchDateTo   string
	searchSortBy   string
	searchSortOrder string
	searchLimit    int
	searchOffset   int
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search perfect days",
	Long:  "Search perfect days by query, area, user, or date range.",
	Run:   runSearch,
}

func init() {
	searchCmd.Flags().StringVarP(&searchQuery, "query", "q", "", "Search query")
	searchCmd.Flags().StringSliceVarP(&searchAreas, "areas", "a", []string{}, "Filter by areas (comma-separated)")
	searchCmd.Flags().StringVarP(&searchUser, "user", "u", "", "Filter by username")
	searchCmd.Flags().StringVar(&searchDateFrom, "from", "", "Filter from date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&searchDateTo, "to", "", "Filter to date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&searchSortBy, "sort", "created_at", "Sort by: date, created_at, title")
	searchCmd.Flags().StringVar(&searchSortOrder, "order", "desc", "Sort order: asc, desc")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 10, "Number of results to show")
	searchCmd.Flags().IntVar(&searchOffset, "offset", 0, "Number of results to skip")
}

func runSearch(cmd *cobra.Command, args []string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	storage := storage.NewStorage(config.DataDirectory)
	searchService := search.NewSearchService()

	allPerfectDays, err := storage.PerfectDayStorage.LoadAll(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading perfect days: %v\n", err)
		os.Exit(1)
	}

	criteria := search.SearchCriteria{
		Query:     searchQuery,
		Areas:     searchAreas,
		Username:  searchUser,
		DateFrom:  searchDateFrom,
		DateTo:    searchDateTo,
		SortBy:    searchSortBy,
		SortOrder: searchSortOrder,
		Limit:     searchLimit,
		Offset:    searchOffset,
	}

	results := searchService.Search(allPerfectDays, criteria)

	if results.Total == 0 {
		fmt.Println("No perfect days found matching your criteria")
		return
	}

	fmt.Printf("Found %d perfect days", results.Total)
	if searchLimit > 0 {
		start := results.Offset + 1
		end := results.Offset + len(results.PerfectDays)
		fmt.Printf(" (showing %d-%d)", start, end)
	}
	fmt.Println(":")
	fmt.Println()

	printSearchResults(results.PerfectDays)

	if results.Total > searchLimit && searchLimit > 0 {
		fmt.Printf("\nShowing %d of %d results. Use --offset and --limit to see more.\n",
			len(results.PerfectDays), results.Total)
	}
}

func printSearchResults(perfectDays []*models.PerfectDay) {
	for i, pd := range perfectDays {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 60))
		}

		fmt.Printf("Title: %s\n", pd.Title)
		fmt.Printf("ID: %s\n", pd.ID[:8])
		fmt.Printf("User: %s | Date: %s\n", pd.Username, pd.Date)

		if len(pd.Areas) > 0 {
			fmt.Printf("Areas: %s\n", strings.Join(pd.Areas, ", "))
		}

		if pd.Description != "" {
			fmt.Printf("Description: %s\n", pd.Description)
		}

		fmt.Printf("Activities: %d\n", len(pd.Activities))

		if len(pd.Activities) > 0 {
			fmt.Println("Timeline:")
			for _, activity := range pd.Activities {
				fmt.Printf("  • %s at %s (%s)\n",
					activity.Name,
					activity.Location.Name,
					lib.FormatTimeRange(activity.StartTime, activity.Duration))
			}
		}

		fmt.Println()
	}
}
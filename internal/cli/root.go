package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "perfect-day",
	Short: "Perfect Day - Share your perfect day experiences",
	Long: `Perfect Day is a terminal application for creating, documenting, and sharing
detailed day plans including activities, locations, time spent, and personal commentary.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(versionCmd)
}
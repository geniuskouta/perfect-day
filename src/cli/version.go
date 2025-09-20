package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
	Build   = "1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display the current version and build information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("perfect-day version %s (build %s)\n", Version, Build)
	},
}
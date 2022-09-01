package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of sda",
	Long: `Writes the current version to stdout and terminates.

We use Semantic Versioning (see semver.org). In a nutshell this means
* first digit changed -> you have to adjust some scripts/CI job if present
* second digit changed -> yeah, new features
* third digit changed -> nothing new but less bugs
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.2.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

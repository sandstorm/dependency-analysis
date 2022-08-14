package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helloWorldCmd represents the helloWorld command
var helloWorldCmd = &cobra.Command{
	Use:   "helloWorld",
	Short: "Prints HelloWorld!",
	Long:  `Prints HelloWorld! to stdout and terminates`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HelloWorld!")
	},
}

func init() {
	rootCmd.AddCommand(helloWorldCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helloWorldCmd represents the helloWorld command
var helloWorldCmd = &cobra.Command{
	Use:   "helloWorld",
	Short: "Prints HelloWorld!",
	Long: `Prints HelloWorld! to stdout and terminates`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HelloWorld!")
	},
}

func init() {
	rootCmd.AddCommand(helloWorldCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloWorldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloWorldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

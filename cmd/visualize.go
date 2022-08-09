package cmd

import (
	"github.com/sandstorm/dependency-analysis/analysis"

	"github.com/spf13/cobra"
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize [path to sources]",
	Short: "Renders the dependencies into an image.",
	Long: `In order to get an overview over the project structure all root packages and their dependencies are rendered.

We ignore dependencies to packages outside the project.

We try to stick to the following overall layout if possible:
* packages without incoming dependencies come top
* packages without dependencies are at the bottom

Currently we display unused as well as commented imports appear. Please clean up your code to remove them.

To zoom into a sub-package you can set the input path accordingly, e.g.
$ sda visualize src/main/java/de/sandstorm/sso/services

File extensions determine the languages. Currently supported are:
* Java (*.java)
`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := "."
		if len(args) > 0 {
			sourcePath = args[0]
		}  
		analysis.Analyse(sourcePath)
	},
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// visualizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// visualizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

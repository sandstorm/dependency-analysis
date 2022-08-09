package cmd

import (
	"fmt"
	"github.com/sandstorm/dependency-analysis/rendering"

	"github.com/spf13/cobra"
)

// listSupportedOutputTypesCmd represents the listSupportedOutputTypes command
var listSupportedOutputTypesCmd = &cobra.Command{
	Use:   "listSupportedOutputTypes",
	Short: "List all supported images types for the output file",
	Long: `When visualizing the dependencies we use GraphViz to render the image file.
Thus we support types supported by GraphViz. See https://graphviz.org/docs/outputs/

Types are printed to stdout.
`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, imageType := range rendering.SupportedOutputFormats {
			fmt.Printf("%-10s%s %s\n", imageType.FlagValue, imageType.Label, imageType.Description)
		}
	},
}

func init() {
	visualizeCmd.AddCommand(listSupportedOutputTypesCmd)
}

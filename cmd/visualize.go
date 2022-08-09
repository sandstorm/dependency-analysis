package cmd

import (
	"log"
	"os"
	"os/exec"
	"github.com/sandstorm/dependency-analysis/analysis"
	"github.com/sandstorm/dependency-analysis/rendering"

	"github.com/spf13/cobra"
)

// variables for CLI flags
const defaultOutput = "output.svg"
var output string = ""
const detaultTargetType = "svg"
var targetType string = ""
var openImage = true

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
		graph, err := analysis.Analyse(sourcePath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if output == "stdout" {
			rendering.RenderDotStdout(graph)
		} else {
			outputFormat := rendering.GetOutputFormatByFlagValue(targetType)
			if outputFormat == nil {
				log.Fatalf("unknown type '%s', for available types see listSupportedOutputTypes", targetType)
				os.Exit(2)
			}
			if output == defaultOutput && targetType != detaultTargetType {
				// replace .svg with correct file ending
				output = output[0:len(output) - 3] + outputFormat.FileEnding
			}
			dotFilePath := output + ".dot"
			if err := rendering.RenderDotFile(graph, dotFilePath); err != nil {
				log.Fatal(err)
				os.Exit(3)
			}
			if err := rendering.GraphVizCmd(dotFilePath, output, outputFormat).Run(); err != nil {
				log.Fatal(err)
				os.Exit(4)
			}
			if err := exec.Command("open", output).Run(); err != nil {
				log.Fatal(err)
				os.Exit(5)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	visualizeCmd.Flags().StringVarP(&output, "output", "o", defaultOutput, "path to the image file to generate, use 'stdout' to output DOT graph without image rendering")
	visualizeCmd.Flags().StringVarP(&targetType, "type", "T", detaultTargetType, "type of the image file, for available formats see listSupportedOutputTypes")
	visualizeCmd.Flags().BoolVarP(&openImage, "show-image", "s", true, "Automatically open the image after rendering")
}

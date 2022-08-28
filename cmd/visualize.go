package cmd

import (
	"github.com/sandstorm/dependency-analysis/analysis"
	"github.com/sandstorm/dependency-analysis/rendering"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
)

// variables for CLI flags
var visualizeCmdFlags = struct {
	defaultOutput     string
	output            string
	defaultTargetType string
	targetType        string
	openImage         bool
	depthString       string
	renderingSettings *rendering.RenderingSettings
}{
	defaultOutput:     "output.svg",
	output:            "",
	defaultTargetType: "svg",
	targetType:        "",
	openImage:         true,
	depthString:       "",
	renderingSettings: rendering.NewRenderingSettings(),
}

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize [path to sources]",
	Short: "Renders the dependencies into an image.",
	Long: `In order to get an overview over the project structure all root packages and their dependencies are rendered.

We ignore dependencies to packages outside the project's root packge.

We try to stick to the following overall layout if possible:
* packages without incoming dependencies come top
* packages without outgoing dependencies are at the bottom

Currently we display unused as well as commented imports appear. Please clean up your code to remove them.

To zoom into a sub-package you can set the input path accordingly, e.g.
$ sda visualize src/main/java/de/sandstorm/sso/services

File extensions determine the languages. Currently supported are:
* Golang (*.go)
* Groovy (*.groovy)
* Java (*.java)
* PHP (*.php)
`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := "."
		if len(args) > 0 {
			sourcePath = args[0]
		}
		depthString := visualizeCmdFlags.depthString
		depth, err := strconv.Atoi(depthString)
		if err != nil {
			log.Fatal("failed to parse parameter 'depth'")
			log.Fatal(err)
			os.Exit(6)
		}
		graph, err := analysis.BuildDependencyGraph(sourcePath, depth)
		if err != nil {
			log.Fatal("failed to read source files and build dependency graph")
			log.Fatal(err)
			os.Exit(1)
		}
		wGraph := analysis.WeightByNumberOfDescendant(graph)
		cycles := analysis.FindCycles(graph)
		output := visualizeCmdFlags.output
		if output == "stdout" {
			rendering.RenderDotStdout(visualizeCmdFlags.renderingSettings, wGraph, cycles)
		} else {
			targetType := visualizeCmdFlags.targetType
			outputFormat := rendering.GetOutputFormatByFlagValue(targetType)
			if outputFormat == nil {
				log.Fatalf("unknown type '%s', for available types see listSupportedOutputTypes", targetType)
				os.Exit(2)
			}
			if output == visualizeCmdFlags.defaultOutput && targetType != visualizeCmdFlags.defaultTargetType {
				// replace .svg with correct file ending
				output = output[0:len(output)-3] + outputFormat.FileEnding
			}
			dotFilePath := output + ".dot"
			if err := rendering.RenderDotFile(visualizeCmdFlags.renderingSettings, wGraph, cycles, dotFilePath); err != nil {
				log.Fatal("failed to render graph into a DOT file")
				log.Fatal(err)
				os.Exit(3)
			}
			if err := rendering.GraphVizCmd(dotFilePath, output, outputFormat).Run(); err != nil {
				log.Fatal("failed to render graph into an image file")
				log.Fatal(err)
				os.Exit(4)
			}
			if isOSX() && visualizeCmdFlags.openImage {
				if err := exec.Command("open", output).Run(); err != nil {
					log.Fatal("failed to open image file")
					log.Fatal(err)
					os.Exit(5)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	renderingDefaults := rendering.NewRenderingSettings()
	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.output, "output", "o", visualizeCmdFlags.defaultOutput, "path to the image file to generate, use 'stdout' to output DOT graph without image rendering")
	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.targetType, "type", "T", visualizeCmdFlags.defaultTargetType, "type of the image file, for available formats see listSupportedOutputTypes")
	if isOSX() {
		visualizeCmd.Flags().BoolVarP(&visualizeCmdFlags.openImage, "show-image", "s", true, "automatically open the image after rendering")
	}
	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.depthString, "depth", "d", "1", "number of steps to go further down into the package hierarchy starting at the root package")
	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.renderingSettings.GraphLabel, "graphLabel", "l", renderingDefaults.GraphLabel, "the graph label is located at the bottom center of the resulting image")
	visualizeCmd.Flags().BoolVarP(&visualizeCmdFlags.renderingSettings.ShowNodeLabels, "show-node-labels", "", renderingDefaults.ShowNodeLabels, "render graph with node labels")
}

func isOSX() bool {
	return runtime.GOOS == "darwin"
}

package cmd

import (
	"github.com/sandstorm/dependency-analysis/analysis"
	"github.com/sandstorm/dependency-analysis/rendering"
	"log"
	"os"
	"os/exec"
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
	label             string
}{
	defaultOutput:     "output.svg",
	output:            "",
	defaultTargetType: "svg",
	targetType:        "",
	openImage:         true,
	depthString:       "",
	label:             "",
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
or manually overwrite the root package, e.g.
$ sda visualize src/main/java --root-package de.sandstorm.sso.services

File extensions determine the languages. Currently supported are:
* Java (*.java)
`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := "."
		if len(args) > 0 {
			sourcePath = args[0]
		}
		defaultOutput := visualizeCmdFlags.defaultOutput
		output := visualizeCmdFlags.output
		defaultTargetType := visualizeCmdFlags.defaultTargetType
		targetType := visualizeCmdFlags.targetType
		openImage := visualizeCmdFlags.openImage
		depthString := visualizeCmdFlags.depthString
		depth, err := strconv.Atoi(depthString)
		label := visualizeCmdFlags.label
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
		if output == "stdout" {
			rendering.RenderDotStdout(label, wGraph, cycles)
		} else {
			outputFormat := rendering.GetOutputFormatByFlagValue(targetType)
			if outputFormat == nil {
				log.Fatalf("unknown type '%s', for available types see listSupportedOutputTypes", targetType)
				os.Exit(2)
			}
			if output == defaultOutput && targetType != defaultTargetType {
				// replace .svg with correct file ending
				output = output[0:len(output)-3] + outputFormat.FileEnding
			}
			dotFilePath := output + ".dot"
			if err := rendering.RenderDotFile(label, wGraph, cycles, dotFilePath); err != nil {
				log.Fatal("failed to render graph into a DOT file")
				log.Fatal(err)
				os.Exit(3)
			}
			if err := rendering.GraphVizCmd(dotFilePath, output, outputFormat).Run(); err != nil {
				log.Fatal("failed to render graph into an image file")
				log.Fatal(err)
				os.Exit(4)
			}
			if openImage {
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

	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.output, "output", "o", visualizeCmdFlags.defaultOutput, "path to the image file to generate, use 'stdout' to output DOT graph without image rendering")
	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.targetType, "type", "T", visualizeCmdFlags.defaultTargetType, "type of the image file, for available formats see listSupportedOutputTypes")
	visualizeCmd.Flags().BoolVarP(&visualizeCmdFlags.openImage, "show-image", "s", true, "automatically open the image after rendering")
	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.depthString, "depth", "d", "1", "number of steps to go further down into the package hierarchy starting at the root package")
	visualizeCmd.Flags().StringVarP(&visualizeCmdFlags.label, "label", "l", "rendered by github.com/sandstorm/dependency-analysis", "the graph is located at the bottom of the resulting image")
	// TODO visualizeCmd.Flags().StringVarP(&rootPackage, "root-package", "r", "", "root package of the project")
}

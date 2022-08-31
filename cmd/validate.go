package cmd

import (
	"fmt"
	"github.com/sandstorm/dependency-analysis/analysis"
	"github.com/sandstorm/dependency-analysis/dataStructures"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

// variables for CLI flags
var validateCmdFlags = struct {
	analyzerSettings *analyzerSettingsFlags
	maxCycles        string
}{
	analyzerSettings: NewAnalyzerSettingsFlags(),
	maxCycles:        "0",
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates the dependencies within the project",
	Long: `All dependencies within the project are analyzed and searched for cylces.

We do not want cyclic dependencies anywhere within our projects:
* not between classes
* not between packages

Since this happens easily by accident this command exits with an error when
there cycles exist between
* the root packages of the project
* (more coming later)

The parameter '--max-cycles' is intended as follows:
* remove cycles step-by-step from legacy projects with the goal to set --max-cycles to zero eventually
* rare corner-cases where you consider cycles a valid option
	`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := "."
		if len(args) > 0 {
			sourcePath = args[0]
		}
		analyzerSettings, err := validateCmdFlags.analyzerSettings.toAnalyzerSettings(sourcePath)
		if err != nil {
			os.Exit(6)
		}
		maxCycles, err := strconv.Atoi(validateCmdFlags.maxCycles)
		if err != nil {
			log.Fatal("failed to parse parameter '--max-cycles'")
			log.Fatal(err)
			os.Exit(7)
		}
		graph, err := analysis.BuildDependencyGraph(analyzerSettings)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		cycles := analysis.FindCycles(graph)
		hasCylces := len(cycles) > 0
		if hasCylces {
			fmt.Printf("found %d cycles:\n", len(cycles))
			prefix := " |  "
			for _, cycle := range cycles {
				fmt.Println("")
				var head = getSmallestNode(cycle)
				var current = ""
				for current != head {
					if current == "" {
						current = head
					}
					var firstPrefix = prefix
					if current == head {
						firstPrefix = " ┌▶ "
					}
					fmt.Printf("%s%s\n", firstPrefix, current)

					current = cycle[current]
					isLastNode := current == head
					if isLastNode {
						fmt.Printf(" └───┘\n")
					} else {
						fmt.Printf("%s ▼\n", prefix)
					}
				}
			}
			if len(cycles) > maxCycles {
				os.Exit(len(cycles))
			} else {
				fmt.Printf("Found %d cycles are below threshold of %d\n", len(cycles), maxCycles)
			}
		} else {
			fmt.Println("No cycles found. Good work :)")
		}
	},
}

func getSmallestNode(cycle dataStructures.Cycle) string {
	nodes := make([]string, len(cycle))
	var i = 0
	for _, node := range cycle {
		nodes[i] = node
		i++
	}
	sort.Strings(nodes)
	return nodes[0]
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&validateCmdFlags.maxCycles, "max-cycles", "", "0", "Maximum number of cycles to attribute with exit-code '0'")
	addAnalyzerSettingsFlags(validateCmd, validateCmdFlags.analyzerSettings)
}

/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"fmt"
	"os"
	"sort"
	"github.com/sandstorm/dependency-analysis/analysis"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validates the dependencies within the project",
	Long: `All dependencies within the project are analyzed and searched for cylces.

We do not want cyclic dependencies anywhere within our projects:
* not between classes
* not between packages

Since this happens easily by accident this command exits with an error when
there cycles exist betwee
* the root packages of the project
* (more coming later)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := "."
		if len(args) > 0 {
			sourcePath = args[0]
		}  
		graph, err := analysis.BuildDependencyGraph(sourcePath)
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
				os.Exit(len(cycles))
			}
		}
	},
}

func getSmallestNode(cycle analysis.Cycle) string {
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
}

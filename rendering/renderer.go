package rendering

import (
	"fmt"
	"os"
	"github.com/sandstorm/dependency-analysis/dataStructures"
	"github.com/sandstorm/dependency-analysis/utils"

	"github.com/mazznoer/colorgrad"
)

func RenderDotFile(sourceGraph *dataStructures.WeightedStringGraph, cycles []dataStructures.Cycle, targetPath string) error {
    file, err := os.Create(targetPath)
    if err != nil {
        return err
    }
    defer file.Close()
	err2 := renderDot(sourceGraph, cycles, func(value string) error {
		_, err := file.WriteString(value)
		return err
	})
	return err2
}

func RenderDotStdout(sourceGraph *dataStructures.WeightedStringGraph, cycles []dataStructures.Cycle) error {
	err := renderDot(sourceGraph, cycles, func(value string) error {
		fmt.Print(value)
		return nil
	})
	return err
}

type writeFunc func(string) error
func renderDot(sourceGraph *dataStructures.WeightedStringGraph, cycles []dataStructures.Cycle, write writeFunc) error {
		if err := write(fmt.Sprintln("digraph {")); err != nil {
			return err
		}
		// global settings
		if err := write(fmt.Sprintf("label = \"%s\"\n", "TODO: root package")); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("labelloc = \"t\";\n\n")); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("node [shape = box];\n\n")); err != nil {
			return err
		}

		// print nodes
		nodesByWeight, maxWeight := sourceGraph.GetNodesGroupedByWeight()
		nodeColorScale := colorgrad.Cool()
		for weight, nodes := range nodesByWeight {
			var rank = "same"
			if weight == 0 {
				rank = "max"
			} else if weight == maxWeight {
				rank = "min"
			}
			color := nodeColorScale.At(0.2 + 0.8*(1.0 - float64(weight)/float64(maxWeight))).Hex()
			if err := write(fmt.Sprintf("{\nrank=%s;\n", rank)); err != nil {
				return err
			}
			for _, node := range nodes {
				if err := write(fmt.Sprintf(
					"n_%s [label=\"%s\",color=\"%s\",style=\"filled\",fillcolor=\"%s\"]\n",
					utils.MD5String(node),
					node,
					color,
					color,
				)); err != nil {
					return err
				}
			}
			if err := write(fmt.Sprintln("}")); err != nil {
				return err
			}
		}

		// print edges of cycles
		printedEdges := dataStructures.NewStringSet()
		cycleColorScale := colorgrad.Warm()
		for i, cycle := range cycles {
			color := cycleColorScale.At(0.6 * (1.0 - float64(i)/float64(len(cycles)))).Hex()
			for caller, callee := range cycle {
				if err := write(fmt.Sprintf(
					"n_%s -> n_%s [penwidth=2,color=\"%s\"]",
					utils.MD5String(caller),
					utils.MD5String(callee),
					color,
				)); err != nil {
					return err
				}
				printedEdges.Add(edgeToString(caller, callee))
			}
		}

		// print remaining edges
		for caller, callees := range sourceGraph.GetEdges() {
			if len(callees) > 0 {
				if err := write(fmt.Sprintf("n_%s -> {", utils.MD5String(caller))); err != nil {
					return err
				}
				for _, callee := range callees {
					if !printedEdges.Contains(edgeToString(caller, callee)) {
						if err := write(fmt.Sprintf(" n_%s", utils.MD5String(callee))); err != nil {
							return err
						}
					}
				}
				if err := write(fmt.Sprintln(" }")); err != nil {
					return err
				}
			}
		}
		if err := write(fmt.Sprintln("}")); err != nil {
			return err
		}
		return nil
}

func edgeToString(start string, end string) string {
	return start + " ➡️ " + end
}

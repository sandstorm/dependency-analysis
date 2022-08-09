package rendering

import (
	"fmt"
	"os"
	"github.com/sandstorm/dependency-analysis/dataStructures"
	"github.com/sandstorm/dependency-analysis/utils"
)

func RenderDotFile(sourceGraph *dataStructures.WeightedStringGraph, targetPath string) error {
    file, err := os.Create(targetPath)
    if err != nil {
        return err
    }
    defer file.Close()
	err2 := renderDot(sourceGraph, func(value string) error {
		_, err := file.WriteString(value)
		return err
	})
	return err2
}

func RenderDotStdout(sourceGraph *dataStructures.WeightedStringGraph) error {
	err := renderDot(sourceGraph, func(value string) error {
		fmt.Print(value)
		return nil
	})
	return err
}

type writeFunc func(string) error
func renderDot(sourceGraph *dataStructures.WeightedStringGraph, write writeFunc) error {
		// rendering
		nodesByWeight, maxWeight := sourceGraph.GetNodesGroupedByWeight()
		if err := write(fmt.Sprintln("digraph {")); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("label = \"%s\"\n", "TODO: root package")); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("labelloc = \"t\";\n\n")); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("node [shape = box];\n\n")); err != nil {
			return err
		}
		for caller, callees := range sourceGraph.GetEdges() {
			if len(callees) > 0 {
				if err := write(fmt.Sprintf("n_%s -> {", utils.MD5String(caller))); err != nil {
					return err
				}
				for _, callee := range callees {
					if err := write(fmt.Sprintf(" n_%s", utils.MD5String(callee))); err != nil {
						return err
					}
				}
				if err := write(fmt.Sprintln(" }")); err != nil {
					return err
				}
			}
		}
		for weight, nodes := range nodesByWeight {
			var rank = "same"
			if weight == 0 {
				rank = "max"
			} else if weight == maxWeight {
				rank = "min"
			}
			if err := write(fmt.Sprintf("{\nrank=%s;\n", rank)); err != nil {
				return err
			}
			for _, node := range nodes {
				if err := write(fmt.Sprintf("n_%s [label=\"%s\"]\n", utils.MD5String(node), node)); err != nil {
					return err
				}
			}
			if err := write(fmt.Sprintln("}")); err != nil {
				return err
			}
		}
		if err := write(fmt.Sprintln("}")); err != nil {
			return err
		}
		return nil
}

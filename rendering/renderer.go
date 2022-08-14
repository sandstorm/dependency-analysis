package rendering

import (
	"fmt"
	"github.com/sandstorm/dependency-analysis/dataStructures"
	"github.com/sandstorm/dependency-analysis/utils"
	"math"
	"os"

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
	/*
		TODO: generate label form root package
		if err := write(fmt.Sprintf("label = \"%s\"\n", "TODO: root package")); err != nil {
			return err
		}
		if err := write(fmt.Sprintf("labelloc = \"t\";\n\n")); err != nil {
			return err
		}
	*/
	if err := write(fmt.Sprintf("node [shape = box];\n\n")); err != nil {
		return err
	}

	// print colored edges of cycles
	coloredEdges := make(map[string]*coloredEdgeStruct)
	cycleColorScale := colorgrad.Warm()
	// group by edge
	for i, cycle := range cycles {
		color := cycleColorScale.At(0.6 * (1.0 - float64(i)/float64(len(cycles)))).Hex()
		for caller, callee := range cycle {
			edgeString := edgeToString(caller, callee)
			edge, isPresent := coloredEdges[edgeString]
			if isPresent {
				edge.color += fmt.Sprintf(":%s", color)
				edge.colorCount++
			} else {
				coloredEdges[edgeString] = &coloredEdgeStruct{
					caller:     caller,
					callee:     callee,
					color:      color,
					colorCount: 0,
				}
			}
		}
	}
	// output
	for _, coloredEdge := range coloredEdges {
		if err := write(fmt.Sprintf(
			"n_%s -> n_%s [color=\"%s\",arrowsize=\"%f\"]\n",
			utils.MD5String(coloredEdge.caller),
			utils.MD5String(coloredEdge.callee),
			coloredEdge.color,
			math.Min(5, math.Max(1, float64(coloredEdge.colorCount)/10)),
		)); err != nil {
			return err
		}
	}

	// print remaining edges
	for caller, callees := range sourceGraph.GetEdges() {
		if len(callees) > 0 {
			if err := write(fmt.Sprintf("n_%s -> {", utils.MD5String(caller))); err != nil {
				return err
			}
			for _, callee := range callees {
				if _, isPresent := coloredEdges[edgeToString(caller, callee)]; !isPresent {
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
		color := nodeColorScale.At(0.2 + 0.8*(1.0-float64(weight)/float64(maxWeight))).Hex()
		if err := write(fmt.Sprintf("{rank=%s;\n", rank)); err != nil {
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

	// close graph
	if err := write(fmt.Sprintln("}")); err != nil {
		return err
	}
	return nil
}

type coloredEdgeStruct struct {
	caller     string
	callee     string
	color      string
	colorCount int
}

func edgeToString(caller string, callee string) string {
	return caller + " ➡️ " + callee
}

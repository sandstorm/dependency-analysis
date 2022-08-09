package analysis

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/sandstorm/dependency-analysis/parsing"
	"github.com/sandstorm/dependency-analysis/utils"
)

// mapping from file path to source-unit
type sourceUnitByFile = map[string][]string
// mapping from source-unit to all its imports
type dependenciesBySourceUnit = map[string]*utils.StringSet

func Analyse(sourcePath string) {
	// Step 1)
	sourceUnits := make(sourceUnitByFile)
	err := filepath.Walk(sourcePath, findSourceUnits(sourceUnits))
	if err != nil {
        log.Fatal(err)
    }
	var rootPackage []string = nil
	for _, sourceUnit := range sourceUnits {
		if rootPackage == nil {
			rootPackage = sourceUnit
		} else {
			commonPrefixLength := getCommonPrefixLength(rootPackage, sourceUnit)
			rootPackage = rootPackage[0:commonPrefixLength]
		}
	}

	// Step 2
	unweightedDependencyGraph := findDependencies(rootPackage, sourceUnits)

	// Step 3
	dependencyGraph := WeightByNumberOfDescendant(unweightedDependencyGraph)

	// rendering
	nodesByWeight, maxWeight := dependencyGraph.GetNodesGroupedByWeight()
	fmt.Println("digraph {")
	fmt.Printf("label = \"%s\"\n", parsing.ParseJavaJoinPathSegments(rootPackage));
	fmt.Printf("labelloc = \"t\";\n\n")
	fmt.Printf("node [shape = box];\n\n")
	for caller, callees := range dependencyGraph.edges {
		calleesArray := callees.ToArray()
		if len(calleesArray) > 0 {
			fmt.Printf("n_%s -> {", utils.MD5String(caller))
			for _, callee := range calleesArray {
				fmt.Printf(" n_%s", utils.MD5String(callee))
			}
			fmt.Println(" }");
		}
	}
	for weight, nodes := range nodesByWeight {
		var rank = "same"
		if weight == 0 {
			rank = "max"
		} else if weight == maxWeight {
			rank = "min"
		}
		fmt.Printf("{\nrank=%s;\n", rank)
		for _, node := range nodes {
			fmt.Printf("n_%s [label=\"%s\"]\n", utils.MD5String(node), node)
		}
		fmt.Println("}")
	}
	fmt.Println("}")
}

func findSourceUnits(result sourceUnitByFile) filepath.WalkFunc {
	return func (path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if !info.IsDir() {
			fileReader, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fileReader.Close()
			sourceUnit := parsing.ParseJavaSourceUnit(fileReader)
			if len(sourceUnit) > 0 {
				result[path] = sourceUnit
			}
		}
		return nil
	}
}

func getCommonPrefixLength(left []string, right []string) int {
	limit := min(len(left), len(right))
	for i := 0; i < limit; i++ {
		if left[i] != right[i] {
			return i
		}
	}
	return limit
}

func findDependencies(rootPackage []string, sourceUnits sourceUnitByFile) *DirectedStringGraph {
	dependencyGraph := NewDirectedStringGraph()
	// TODO: make configurable
	prefixLength := len(rootPackage)
	segmentLimit := len(rootPackage) + 1
	for path, sourceUnit := range sourceUnits {
    	fileReader, err := os.Open(path)
    	if err != nil {
			log.Fatal(err)
    	}
		allDependencies, err := parsing.ParseJavaImports(fileReader)
		fileReader.Close()
		if err != nil {
			log.Fatal(err)
    	}
		sourceUnitString := parsing.ParseJavaJoinPathSegments(
			sourceUnit[prefixLength:min(segmentLimit, len(sourceUnit))])
		dependencyGraph.AddNode(sourceUnitString)
		for _, dependency := range(allDependencies) {
			if arrayStartsWith(dependency, rootPackage) {
				target := parsing.ParseJavaJoinPathSegments(
					dependency[prefixLength:min(segmentLimit, len(dependency))])
				if sourceUnitString != target {
					dependencyGraph.AddEdge(sourceUnitString, target)
				}
			}
		}
	}
	return dependencyGraph
}

func arrayStartsWith(value []string, prefix []string) bool {
	if len(prefix) > len(value) {
		return false
	}
	for i, v := range(prefix) {
		if v != value[i] {
			return false
		}
	}
	return true
}

func min(left int, right int) int {
	if left < right {
		return left
	} else {
		return right
	}
}


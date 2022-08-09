package analysis

import (
	"os"
	"path/filepath"
	"github.com/sandstorm/dependency-analysis/parsing"
	"github.com/sandstorm/dependency-analysis/dataStructures"
)

// mapping from file path to source-unit
type sourceUnitByFile = map[string][]string
// mapping from source-unit to all its imports
type dependenciesBySourceUnit = map[string]*dataStructures.StringSet

func Analyse(sourcePath string) (*dataStructures.WeightedStringGraph, error) {
	// Step 1)
	sourceUnits := make(sourceUnitByFile)
	err := filepath.Walk(sourcePath, findSourceUnits(sourceUnits))
	if err != nil {
		return nil, err
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
	unweightedDependencyGraph, err := findDependencies(rootPackage, sourceUnits)
	if err != nil {
		return nil, err
    }

	// Step 3
	dependencyGraph := dataStructures.WeightByNumberOfDescendant(unweightedDependencyGraph)

	return dependencyGraph, nil
}

func findSourceUnits(result sourceUnitByFile) filepath.WalkFunc {
	return func (path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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

func findDependencies(rootPackage []string, sourceUnits sourceUnitByFile) (*dataStructures.DirectedStringGraph, error) {
	dependencyGraph := dataStructures.NewDirectedStringGraph()
	// TODO: make configurable
	prefixLength := len(rootPackage)
	segmentLimit := len(rootPackage) + 1
	for path, sourceUnit := range sourceUnits {
    	fileReader, err := os.Open(path)
    	if err != nil {
			return nil, err
    	}
		allDependencies, err := parsing.ParseJavaImports(fileReader)
		fileReader.Close()
		if err != nil {
			return nil, err
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
	return dependencyGraph, nil
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

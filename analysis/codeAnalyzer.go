package analysis

import (
	"github.com/sandstorm/dependency-analysis/dataStructures"
	"github.com/sandstorm/dependency-analysis/parsing"
	"os"
	"path/filepath"
	"regexp"
)

// mapping from file path to source-unit
type sourceUnitByFile = map[string][]string

// mapping from source-unit to all its imports
type dependenciesBySourceUnit = map[string]*dataStructures.StringSet

func BuildDependencyGraph(settings *AnalyzerSettings) (*dataStructures.DirectedStringGraph, error) {
	sourceUnits := make(sourceUnitByFile)
	if err := filepath.Walk(settings.SourcePath, initializeParsers(settings.IncludePattern)); err != nil {
		return nil, err
	}
	if err := filepath.Walk(settings.SourcePath, findSourceUnits(settings.IncludePattern, sourceUnits)); err != nil {
		return nil, err
	}
	var rootPackage []string = nil
	for _, sourceUnit := range sourceUnits {
		if rootPackage == nil {
			rootPackage = sourceUnit
		} else {
			commonPrefixLength := getCommonPrefixLength(rootPackage, sourceUnit)
			rootPackage = rootPackage[:commonPrefixLength]
		}
	}

	return findDependencies(rootPackage, sourceUnits, settings.Depth)
}

func initializeParsers(includePattern *regexp.Regexp) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && includePattern.MatchString(path) {
			return parsing.InitializeParsers(path)
		}
		return nil
	}
}

func findSourceUnits(includePattern *regexp.Regexp, result sourceUnitByFile) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && includePattern.MatchString(path) {
			fileReader, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fileReader.Close()
			sourceUnit := parsing.ParseSourceUnit(path, fileReader)
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

func findDependencies(rootPackage []string, sourceUnits sourceUnitByFile, depth int) (*dataStructures.DirectedStringGraph, error) {
	dependencyGraph := dataStructures.NewDirectedStringGraph()
	prefixLength := len(rootPackage)
	segmentLimit := len(rootPackage) + depth
	for path, sourceUnit := range sourceUnits {
		fileReader, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		allDependencies, err := parsing.ParseImports(path, fileReader)
		fileReader.Close()
		if err != nil {
			return nil, err
		}
		sourceUnitString := parsing.JoinPathSegments(
			path,
			sourceUnit[prefixLength:min(segmentLimit, len(sourceUnit))])
		dependencyGraph.AddNode(sourceUnitString)
		for _, dependency := range allDependencies {
			if arrayStartsWith(dependency, rootPackage) {
				dependencyString := parsing.JoinPathSegments(
					path,
					dependency[prefixLength:min(segmentLimit, len(dependency))])
				if sourceUnitString != dependencyString {
					dependencyGraph.AddEdge(sourceUnitString, dependencyString)
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
	for i, v := range prefix {
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

package analysis

import (
	"regexp"
	"testing"

	"github.com/sandstorm/dependency-analysis/dataStructures"
)

// this test runs the code analysis against this project
// if you change the package structure, you need to adjust the expectations as well
func TestBuildDependencyGraph(t *testing.T) {
	t.Run("analyze sda", func(t *testing.T) {
		settings := &AnalyzerSettings{
			SourcePath: "..",
			Depth: 1,
			IncludePattern: regexp.MustCompile(".*"),
		}
		actual, err := BuildDependencyGraph(settings)
		AssertNil(t, err)
		expected := dataStructures.NewDirectedStringGraph().
			AddEdge("main", "cmd").
			AddEdge("cmd", "analysis").
			AddEdge("cmd", "dataStructures").
			AddEdge("cmd", "rendering").
			AddEdge("analysis", "parsing").
			AddEdge("analysis", "dataStructures").
			AddEdge("parsing", "dataStructures").
			AddEdge("rendering", "dataStructuress")
		AssertEquals(t, "incorrect graph", expected, actual)
	})
}

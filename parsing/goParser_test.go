package parsing

import (
	"bytes"
	"testing"
)

// TODO: func TestParseGoMod(t *testing.T)

func TestParseGoSourceUnit(t *testing.T) {
	modulePath := []string{
		"github.com",
		"sandstom",
		"dependency-analysis"}
	golangParser.modulePath = modulePath
	testCases := []struct {
		name        string
		fileName    string
		fileContent string
		expected    []string
	}{
		{
			name:     "simple file without imports",
			fileName: "cycle.go",
			fileContent: `package dataStructures

			// all edges of a cylce in a mapping from source to destination
			type Cycle map[string]string`,
			expected: append(modulePath, []string{
				"dataStructures",
				"cycle.go"}...),
		},
		{
			name:     "simple file with imports and sub-package",
			fileName: "codeAnalyzer.go",
			fileContent: `package analysis/details

			import (
				"os"
				"path/filepath"
				"github.com/sandstorm/dependency-analysis/parsing"
				"github.com/sandstorm/dependency-analysis/dataStructures"
			)
			
			// mapping from file path to source-unit
			type sourceUnitByFile = map[string][]string
			…`,
			expected: append(modulePath, []string{
				"analysis",
				"details",
				"codeAnalyzer.go"}...),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			AssertEquals(t,
				testCase.expected,
				ParseGoSourceUnit(testCase.fileName, file),
			)
		})
	}
}

func TestParseGoImports(t *testing.T) {
	testCases := []struct {
		name        string
		fileContent string
		expected    [][]string
	}{
		{
			name: "simple file without imports",
			fileContent: `package dataStructures

			// all edges of a cylce in a mapping from source to destination
			type Cycle map[string]string`,
			expected: [][]string{},
		},
		{
			name: "simple file with single imports",
			fileContent: `package analysis

			import "os"
			
			// mapping from file path to source-unit
			type sourceUnitByFile = map[string][]string
			…`,
			expected: [][]string{
				[]string{"os"},
			},
		},
		{
			name: "simple file with grouped imports",
			fileContent: `package analysis

			import (
				_ "os"
				file "path/filepath"
				"github.com/sandstorm/dependency-analysis/parsing"
				"github.com/sandstorm/dependency-analysis/dataStructures"
			)
			
			// mapping from file path to source-unit
			type sourceUnitByFile = map[string][]string
			…`,
			expected: [][]string{
				[]string{"os"},
				[]string{"path", "filepath"},
				[]string{"github.com", "sandstorm", "dependency-analysis", "parsing"},
				[]string{"github.com", "sandstorm", "dependency-analysis", "dataStructures"},
			},
		},
		{
			name: "simple file with individual imports",
			fileContent: `package analysis

			import _ "os"
			import file "path/filepath"
			import "github.com/sandstorm/dependency-analysis/parsing"
			import "github.com/sandstorm/dependency-analysis/dataStructures"
			
			// mapping from file path to source-unit
			type sourceUnitByFile = map[string][]string
			…`,
			expected: [][]string{
				[]string{"os"},
				[]string{"path", "filepath"},
				[]string{"github.com", "sandstorm", "dependency-analysis", "parsing"},
				[]string{"github.com", "sandstorm", "dependency-analysis", "dataStructures"},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			file := bytes.NewBufferString(testCase.fileContent)
			actual, _ := ParseGoImports(
				file)
			AssertEquals(t,
				testCase.expected,
				actual)
		})
	}
}

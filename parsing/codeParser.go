package parsing

import (
	"io"
	"strings"
	"path/filepath"
)

// 1st step during the code analysis, called for each file in the source directory
// During this step you can search for global information of the project
// to initiailze your parser.
func InitializeParsers(filePath string) error {
	switch {
	case filePath == "go.mod" || strings.HasSuffix(filePath, "/go.mod"):
		return ParseGoMod(filePath)
	}
	return nil
}

// 2nd step during code analysis, called for each file in the source directory
// Determines the full package path of the source unit within the
// given source file. If the file contains no source unit an emtpy slice/array
// is returned.
// The package path is already split by the language's delimiter,
// e.g. in Java de.sandstorm.test.helpers.ListHelpers results in
// [de sandstorm test helpers ListHelpers]
func ParseSourceUnit(sourcePath string, fileReader io.Reader) []string {
	switch {
	case strings.HasSuffix(sourcePath, ".go"):
		filePathSplit := strings.Split(sourcePath, "/")
		fileName := filePathSplit[len(filePathSplit)-1]
		if fileName != "doc.go" {
			return ParseGoSourceUnit(fileName, fileReader)
		}
	case strings.HasSuffix(sourcePath, ".java"):
		return ParseJavaSourceUnit(fileReader)
	case strings.HasSuffix(sourcePath, ".php"):
		return ParsePhpSourceUnit(fileReader)
	case strings.HasSuffix(sourcePath, ".groovy"):
		return ParseGroovySourceUnit(fileReader)
	case strings.HasSuffix(sourcePath, ".ts"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".tsx"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".js"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".jsx"):
		return ParseJavaScriptSourceUnit(sourcePath)
	}
	return []string{}
}

// 3rd step during code analysis, called for each source unit (see 2nc step)
// Provides the full package path of all dependencies of the given source unit within or outside
// the project. The dependencies are split as in the 2nd step.
// Eg a Java source unit could provide
// - [java util List]
// - [de sandstorm test helpers ListHelpers]
func ParseImports(sourcePath string, fileReader io.Reader) ([][]string, error) {
	switch {
	case strings.HasSuffix(sourcePath, ".go"):
		return ParseGoImports(fileReader)
	case strings.HasSuffix(sourcePath, ".java"):
		return ParseJavaImports(fileReader)
	case strings.HasSuffix(sourcePath, ".php"):
		return ParsePhpImports(fileReader)
	case strings.HasSuffix(sourcePath, ".groovy"):
		return ParseGroovyImports(fileReader)
	case strings.HasSuffix(sourcePath, ".ts"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".tsx"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".js"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".jsx"):
		return ParseJavaScriptImports(sourcePath, fileReader)
	}
	return [][]string{}, nil
}

// Utility function for creating nice labels
// Joins the full package path with the according delimiter of the language.
func JoinPathSegments(sourcePath string, segments []string) string {
	switch {
	case strings.HasSuffix(sourcePath, ".go"):
		return strings.Join(segments, "/")
	case strings.HasSuffix(sourcePath, ".java"):
		return strings.Join(segments, ".")
	case strings.HasSuffix(sourcePath, ".php"):
		return strings.Join(segments, "\\")
	case strings.HasSuffix(sourcePath, ".ts"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".tsx"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".js"):
		fallthrough
	case strings.HasSuffix(sourcePath, ".jsx"):
		if len(segments) == 0 {
			return "index" + filepath.Ext(sourcePath)
		} else {
			return strings.Join(segments, "/")
		}
	default:
		return strings.Join(segments, ".")
	}
}

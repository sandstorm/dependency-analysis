package parsing

import (
	"io"
	"strings"
)

func InitializeParsers(filePath string) error {
	switch {
	case filePath == "go.mod" || strings.HasSuffix(filePath, "/go.mod"):
		return ParseGoMod(filePath)
	}
	return nil
}

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
	}
	return []string{}
}

func ParseImports(sourcePath string, fileReader io.Reader) ([][]string, error) {
	switch {
	case strings.HasSuffix(sourcePath, ".go"):
		return ParseGoImports(fileReader)
	case strings.HasSuffix(sourcePath, ".java"):
		return ParseJavaImports(fileReader)
	case strings.HasSuffix(sourcePath, ".php"):
		return ParsePhpImports(fileReader)
	}
	return [][]string{}, nil
}

func JoinPathSegments(sourcePath string, segments []string) string {
	switch {
	case strings.HasSuffix(sourcePath, ".go"):
		return strings.Join(segments, "/")
	case strings.HasSuffix(sourcePath, ".java"):
		return strings.Join(segments, ".")
	default:
		return strings.Join(segments, ".")
	}
}

package parsing

import (
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

var javaScriptParser = struct {
	importRegex *regexp.Regexp
}{
	importRegex: regexp.MustCompile(`(?:import\s+.*)?from\s+["']([^'"]+)["'];?`),
}

func ParseJavaScriptSourceUnit(sourcePath string) [][]string {
	if strings.Contains(sourcePath, "node_modules") {
		return [][]string{}
	} else {
		parent := filepath.Dir(sourcePath)
		parentSegments := strings.Split(parent, "/")
		fileName := filepath.Base(sourcePath)
		fileBasename := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		if fileBasename == "index" {
			return [][]string{parentSegments}
		} else {
			return [][]string{append(parentSegments, fileBasename)}
		}
	}
}

func ParseJavaScriptImports(sourcePath string, fileReader io.Reader) ([][]string, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	parent := filepath.Dir(sourcePath)
	relativeImportPaths := getAllMatches(content, javaScriptParser.importRegex)
	fullImportPaths := mapElements(relativeImportPaths, func(e string) string {
		if strings.HasPrefix(e, ".") {
			return filepath.Clean(
				filepath.Join(parent, e))
		} else {
			return e
		}
	})
	return splitAll(fullImportPaths, "/"), nil
}

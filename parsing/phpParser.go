package parsing

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var phpParser = struct {
	namespaceRegex *regexp.Regexp
	classRegex     *regexp.Regexp
	useRegex       *regexp.Regexp
}{
	namespaceRegex: regexp.MustCompile(`namespace\s+([^; ]+)\s*;`),
	classRegex:     regexp.MustCompile(`(?:(?:public|protected|private)\s*)?class\s+([^{ ]+)`),
	useRegex:       regexp.MustCompile(`use\s+([^; ]+)\s*;`),
}

func ParsePhpSourceUnit(fileReader io.Reader) []string {
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)
	namespace := getFirstLineMatchInScanner(scanner, phpParser.namespaceRegex)
	className := getFirstLineMatchInScanner(scanner, phpParser.classRegex)
	if namespace != "" && className != "" {
		return append(strings.Split(namespace, "\\"), className)
	}
	if className != "" {
		return []string{className}
	}
	return []string{}
}

func ParsePhpImports(fileReader io.Reader) ([][]string, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	usages := getAllMatches(content, phpParser.useRegex)
	return splitAll(usages, "\\"), nil
}

func ParsePhpJoinPathSegments(segments []string) string {
	return strings.Join(segments, "\\")
}

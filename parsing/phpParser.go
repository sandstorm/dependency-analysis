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

func ParsePhpSourceUnit(fileReader io.Reader) []fullyQualifiedType {
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)
	namespace := getFirstLineMatchInScanner(scanner, phpParser.namespaceRegex)
	className := getFirstLineMatchInScanner(scanner, phpParser.classRegex)
	if namespace != "" && className != "" {
		return []fullyQualifiedType{append(strings.Split(namespace, "\\"), className)}
	}
	if className != "" {
		return []fullyQualifiedType{{className}}
	}
	return []fullyQualifiedType{}
}

func ParsePhpImports(fileReader io.Reader) ([]fullyQualifiedType, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	usages := getAllMatches(content, phpParser.useRegex)
	return splitAll(usages, "\\"), nil
}

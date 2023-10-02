package parsing

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var groovyParser = struct {
	packageRegex *regexp.Regexp
	classRegex   *regexp.Regexp
	importRegex  *regexp.Regexp
}{
	packageRegex: regexp.MustCompile(`package\s+([^; \n]+)\s*;?`),
	classRegex:   regexp.MustCompile(`(?:public|protected|private)?\s*(?:static\s+|final\s+)*class\s+([^{ ]+)`),
	importRegex:  regexp.MustCompile(`import\s+(?:static\s+)?([^; \n]+)\s*;?`),
}

func ParseGroovySourceUnit(fileReader io.Reader) [][]string {
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)
	packageString := getFirstLineMatchInScanner(scanner, groovyParser.packageRegex)
	className := getFirstLineMatchInScanner(scanner, groovyParser.classRegex)
	if packageString != "" && className != "" {
		return [][]string{append(strings.Split(packageString, "."), className)}
	}
	if className != "" {
		return [][]string{[]string{className}}
	}
	return [][]string{}
}

func ParseGroovyImports(fileReader io.Reader) ([][]string, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	imports := getAllMatches(content, groovyParser.importRegex)
	return splitAll(imports, "."), nil
}

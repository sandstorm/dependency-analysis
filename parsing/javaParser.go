package parsing

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var javaParser = struct {
	packageRegex *regexp.Regexp
	classRegex   *regexp.Regexp
	importRegex  *regexp.Regexp
}{
	packageRegex: regexp.MustCompile(`package\s+([^; ]+)\s*;`),
	classRegex:   regexp.MustCompile(`(?:public|protected|private)?\s*(?:static\s+|final\s+)*class\s+([^{ ]+)`),
	importRegex:  regexp.MustCompile(`import\s+(?:static\s+)?([^; ]+)\s*;`),
}

func ParseJavaSourceUnit(fileReader io.Reader) [][]string {
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)
	packageString := getFirstLineMatchInScanner(scanner, javaParser.packageRegex)
	className := getFirstLineMatchInScanner(scanner, javaParser.classRegex)
	if packageString != "" && className != "" {
		return [][]string{append(strings.Split(packageString, "."), className)}
	}
	if className != "" {
		return [][]string{[]string{className}}
	}
	return [][]string{}
}

func ParseJavaImports(fileReader io.Reader) ([][]string, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	imports := getAllMatches(content, javaParser.importRegex)
	return splitAll(imports, "."), nil
}

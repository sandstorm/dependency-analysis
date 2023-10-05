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

func ParseJavaSourceUnit(fileReader io.Reader) []fullyQualifiedType {
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)
	packageString := getFirstLineMatchInScanner(scanner, javaParser.packageRegex)
	className := getFirstLineMatchInScanner(scanner, javaParser.classRegex)
	if packageString != "" && className != "" {
		return []fullyQualifiedType{append(strings.Split(packageString, "."), className)}
	}
	if className != "" {
		return []fullyQualifiedType{{className}}
	}
	return []fullyQualifiedType{}
}

func ParseJavaImports(fileReader io.Reader) ([]fullyQualifiedType, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	imports := getAllMatches(content, javaParser.importRegex)
	return splitAll(imports, "."), nil
}

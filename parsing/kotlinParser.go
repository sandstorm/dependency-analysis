package parsing

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

var kotlinParser = struct {
	packageRegex *regexp.Regexp
	classRegex   *regexp.Regexp
	valRegex     *regexp.Regexp
	importRegex  *regexp.Regexp
}{
	packageRegex: regexp.MustCompile(`package\s+([^; ]+)`),
	classRegex:   regexp.MustCompile(`(?:public|protected|private)?\s*(?:open\s+)?class\s+([^({ ]+)`),
	importRegex:  regexp.MustCompile(`import\s+([^; \n]+)\s*;?`),
}

func ParseKotlinSourceUnit(fileReader io.Reader) []fullyQualifiedType {
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)
	packageString := getFirstLineMatchInScanner(scanner, kotlinParser.packageRegex)
	unitNames := getAllFollowingMatches(scanner, kotlinParser.classRegex)
	result := make([]fullyQualifiedType, 0, len(unitNames))
	for _, unitName := range unitNames {
		if unitName != "" {
			if packageString != "" {
				result = append(result, append(strings.Split(packageString, "."), unitName))
			} else {
				result = append(result, fullyQualifiedType{unitName})
			}
		}
	}
	if len(result) == 0 && packageString != "" {
		return []fullyQualifiedType{strings.Split(packageString, ".")}
	}
	return result
}

func ParseKotlinImports(fileReader io.Reader) ([]fullyQualifiedType, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	imports := getAllMatches(content, kotlinParser.importRegex)
	return splitAll(imports, "."), nil
}

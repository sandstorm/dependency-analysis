package parsing

import (
	"io"
	"regexp"
	"strings"
)

var golangParser = struct {
	modulePath         fullyQualifiedType
	moduleRegex        *regexp.Regexp
	packageRegex       *regexp.Regexp
	singleImportRegex  *regexp.Regexp
	importGroupRegex   *regexp.Regexp
	quotedPackageRegex *regexp.Regexp
}{
	modulePath:         make(fullyQualifiedType, 0),
	moduleRegex:        regexp.MustCompile(`module\s+([^ ]+)`),
	packageRegex:       regexp.MustCompile(`package\s+([^ ]+)`),
	singleImportRegex:  regexp.MustCompile(`import\s+(?:[^"()]+\s+)?"([^"]+)"`),
	importGroupRegex:   regexp.MustCompile(`import\s+\(([^()]+)\)`),
	quotedPackageRegex: regexp.MustCompile(`(?:[^"()]+\s+)?"([^"]+)"`),
}

func ParseGoMod(filePath string) error {
	module, err := getFirstLineMatchInPath(filePath, golangParser.moduleRegex)
	if err != nil {
		return err
	}
	golangParser.modulePath = strings.Split(module, "/")
	return nil
}

func ParseGoSourceUnit(fileName string, fileReader io.Reader) []fullyQualifiedType {
	packageString := getFirstLineMatchInReader(fileReader, golangParser.packageRegex)
	if packageString != "" {
		packagePath := strings.Split(packageString, "/")
		return []fullyQualifiedType{append(golangParser.modulePath, append(packagePath, fileName)...)}
	} else {
		return []fullyQualifiedType{}
	}
}

func ParseGoImports(fileReader io.Reader) ([]fullyQualifiedType, error) {
	content, err := readerToString(fileReader)
	if err != nil {
		return nil, err
	}
	singleImports := getAllMatches(content, golangParser.singleImportRegex)
	importGroups := getAllMatches(content, golangParser.importGroupRegex)
	importsInGroups := make(fullyQualifiedType, 100 /* hard-coded upper bound of imports in one import (…) we can handle */)
	importsInGroupsIndex := 0
	for _, group := range importGroups {
		packageMatches := golangParser.quotedPackageRegex.FindAllStringSubmatch(group, -1)
		for _, v := range packageMatches {
			importsInGroups[importsInGroupsIndex] = v[1]
			importsInGroupsIndex++
		}
	}
	allImports := append(singleImports, importsInGroups[:importsInGroupsIndex]...)
	return splitAll(allImports, "/"), nil
}

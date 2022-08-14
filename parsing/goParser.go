package parsing

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

var golangParser = struct {
	modulePath []string
}{
	modulePath: make([]string, 0),
}

func ParseGoMod(filePath string) error {
	fileReader, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)

	moduleRegex := regexp.MustCompile(`module\s+([^ ]+)`)
	for scanner.Scan() {
		line := scanner.Text()
		moduleMatch := moduleRegex.FindStringSubmatch(line)
		if moduleMatch != nil {
			moduleString := moduleMatch[1]
			golangParser.modulePath = strings.Split(moduleString, "/")
			return nil
		}
	}
	return nil
}

func ParseGoSourceUnit(fileName string, fileReader io.Reader) []string {
	packageRegex := regexp.MustCompile(`package\s+([^ ]+)`)

	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		packageMatch := packageRegex.FindStringSubmatch(line)
		if packageMatch != nil {
			packageString := packageMatch[1]
			packagePath := strings.Split(packageString, "/")
			return append(golangParser.modulePath, append(packagePath, fileName)...)
		}
	}
	return []string{}
}

func ParseGoImports(fileReader io.Reader) ([][]string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, fileReader)
	if err != nil {
		return nil, err
	}
	content := buffer.String()
	packagePattern := `(?:[^"()]+\s+)?"([^"]+)"`

	singleImportRegex := regexp.MustCompile(`import\s+` + packagePattern)
	singleImportMatches := singleImportRegex.FindAllStringSubmatch(content, -1)
	singleImportResults := make([][]string, len(singleImportMatches))
	for i, v := range singleImportMatches {
		singleImportResults[i] = strings.Split(v[1], "/")
	}

	importGroupRegex := regexp.MustCompile(`import\s+\(([^()]+)\)`)
	importGroupMatches := importGroupRegex.FindAllStringSubmatch(content, -1)
	importGroupResults := make([][]string, 100 /* hard-coded upper bound of imports in one import (…) we can handle */)
	packageRegex := regexp.MustCompile(packagePattern)
	index := 0
	for _, group := range importGroupMatches {
		packageMatches := packageRegex.FindAllStringSubmatch(group[1], -1)
		for _, v := range packageMatches {
			importGroupResults[index] = strings.Split(v[1], "/")
			index++
		}
	}
	return append(singleImportResults, importGroupResults[:index]...), nil
}

func ParseGoJoinPathSegments(segments []string) string {
	return strings.Join(segments, "/")
}

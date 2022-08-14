package parsing

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

func ParseJavaSourceUnit(fileReader io.Reader) []string {
	packageRegex := regexp.MustCompile(`package\s+([^; ]+)\s*;`)
	classRegex := regexp.MustCompile(`(?:public|protected|private)\s+(?:static\s+|final\s+)*class\s+([^{ ]+)`)

	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)

	var packagePath []string = nil
	var class string = ""
	for scanner.Scan() {
		line := scanner.Text()

		packageMatch := packageRegex.FindStringSubmatch(line)
		if packageMatch != nil {
			packageString := packageMatch[1]
			packagePath = strings.Split(packageString, ".")
		}

		classMatch := classRegex.FindStringSubmatch(line)
		if classMatch != nil {
			class = classMatch[1]
		}

		if packagePath != nil && class != "" {
			return append(packagePath, class)
		}
	}
	return []string{}
}

func ParseJavaImports(fileReader io.Reader) ([][]string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, fileReader)
	if err != nil {
		return nil, err
	}
	content := buffer.String()
	importRegex := regexp.MustCompile(`import\s+(?:static\s+)?([^; ]+)\s*;`)
	matches := importRegex.FindAllStringSubmatch(content, -1)
	result := make([][]string, len(matches))
	for i, v := range matches {
		result[i] = strings.Split(v[1], ".")
	}
	return result, nil
}

func ParseJavaJoinPathSegments(segments []string) string {
	return strings.Join(segments, ".")
}

package parsing

import (
	"regexp"
	"io"
	"bufio"
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


func ParseJavaImports(rootPackage string, segmentLimit int, fileReader io.Reader) ([]string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, fileReader)
	if err != nil {
		return nil, err
	}
	content := buffer.String()
	importRegex := regexp.MustCompile(`import\s+(?:static\s+)?([^; ]+)\s*;`)
	matches := importRegex.FindAllStringSubmatch(content, -1)
	result :=  make([]string, len(matches))
	var resultCount = 0
	for _, v := range matches {
		fullPackage := v[1]
		if strings.HasPrefix(fullPackage, rootPackage) {
			packgageSegments := strings.Split(fullPackage, ".")
			croppedSegments := packgageSegments[0:segmentLimit]
			result[resultCount] = strings.Join(croppedSegments, ".")
			resultCount++
		}
	}
	return result[0:resultCount], nil
}



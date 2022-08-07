package parsing

import (
	"fmt"
	"regexp"
	"io"
	"bufio"
	"strings"
)

func ParseJavaPackage(file io.Reader) []string {
	packageRegex := regexp.MustCompile(`package\s+([^; ]+)\s*;`)
	classRegex := regexp.MustCompile(`(?:public|protected|private)\s+(?:static\s+|final\s+)*class\s+([^{ ]+)`)
	
	scanner := bufio.NewScanner(file)
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


func ParseJavaImports(content string) []string {
	importRegex := regexp.MustCompile(`import\s+(?:static\s+)?([^; ]+)\s*;`)
	matches := importRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		fullPackage := v[1]
		fmt.Printf("%s\n", fullPackage)
	}
	return []string{}
}



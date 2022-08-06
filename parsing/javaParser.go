package parsing

import (
	"fmt"
	"regexp"
)

func ParseJavaFile(content string) []string {
	importRegex := regexp.MustCompile(`import\s+(?:static\s+)?([^; ]+)\s*;\s*`)
	matches := importRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		fullPackage := v[1]
		fmt.Printf("%s\n", fullPackage)
	}
	return []string{}
}

package parsing

import (
	"fmt"
	"regexp"
)

func ParseJavaFile(content string) {
	importRegex := regexp.MustCompile(`import\s+(static\s+)?([^; ]+)\s*;\s*`)
	matches := importRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		isStaticImport := len(v[1]) > 0
		fullPackage := v[2]
		fmt.Printf("%s %s\n", isStaticImport, fullPackage)
	}
}

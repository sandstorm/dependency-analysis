package analyser

import (
	"fmt"
	"log"
	"os"
    "path/filepath"
	"github.com/sandstorm/dependency-analysis/parsing"
)

var basePackage []string = nil

func Analyse(sourcePath string) {
	err := filepath.Walk(sourcePath, analyzeFile)
	if err != nil {
        log.Fatal(err)
    }
	fmt.Printf("%s\n", basePackage)
}

func analyzeFile(path string, info os.FileInfo, err error) error {
    if err != nil {
        log.Print(err)
        return nil
	}
	if !info.IsDir() {
    	fileReader, err := os.Open(path)
    	if err != nil {
			return err
    	}
		defer fileReader.Close()
		sourceUnit := parsing.ParseJavaSourceUnit(fileReader)
		if len(sourceUnit) > 0 {
			if basePackage == nil {
				basePackage = sourceUnit
			} else {
				commonPrefixLength := getCommonPrefixLength(basePackage, sourceUnit)
				basePackage = basePackage[0:commonPrefixLength]
			}
		}
	}
	return nil
}

func getCommonPrefixLength(left []string, right []string) int {
	limit := min(len(left), len(right))
	for i := 0; i < limit; i++ {
		if left[i] != right[i] {
			return i
		}
	}
	return limit
}

func min(left int, right int) int {
	if left < right {
		return left
	} else {
		return right
	}
}

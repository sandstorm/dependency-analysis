package analyser

import (
	"fmt"
	"log"
	"os"
    "path/filepath"
	"github.com/sandstorm/dependency-analysis/parsing"
)



func Analyse(sourcePath string) {
	err := filepath.Walk(sourcePath, analyzeFile)
	if err != nil {
        log.Fatal(err)
    }
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
		fmt.Printf("%s\n", sourceUnit)
	}
	return nil
}

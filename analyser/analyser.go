package analyser

import (
	"fmt"
	"log"
	"os"
    "path/filepath"
)

func Analyse(language string, sourcePath string) {
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
    	fmt.Println(path)
	}
	return nil
}

package analyser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/sandstorm/dependency-analysis/parsing"
	"github.com/sandstorm/dependency-analysis/utils"
)

// base package/namespace of the project, eg [de sandstorm sso]
var rootPackage []string = nil
// mapping from file path to content unit identifier
// eg src/main/java/de/sandstorm/sso/Main.java -> [de sandstorm sso Main]
var sourceUnits map[string][]string = nil

func Analyse(sourcePath string) {
	sourceUnits = make(map[string][]string)
	err := filepath.Walk(sourcePath, findSourceUnits)
	if err != nil {
        log.Fatal(err)
    }
	dependencies := findDependencies()
	fmt.Println("digraph {")
	fmt.Printf("label = \"%s\"\n", parsing.ParseJavaJoinPathSegments(rootPackage));
	fmt.Println("labelloc = \"t\";")
	for caller, _ := range dependencies {
		fmt.Printf("n_%s [label=\"%s\"]\n", utils.MD5String(caller), caller)
	}
	for caller, callees := range dependencies {
		for _, callee := range callees.ToArray() {
			fmt.Printf("n_%s -> n_%s;\n", utils.MD5String(caller), utils.MD5String(callee))
		}
	}
	fmt.Println("}")
}

func findSourceUnits(path string, info os.FileInfo, err error) error {
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
			if rootPackage == nil {
				rootPackage = sourceUnit
			} else {
				commonPrefixLength := getCommonPrefixLength(rootPackage, sourceUnit)
				rootPackage = rootPackage[0:commonPrefixLength]
			}
			sourceUnits[path] = sourceUnit
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

func findDependencies() map[string]*utils.StringSet {
	dependencies := make(map[string]*utils.StringSet)
	// TODO: make configurable
	prefixLength := len(rootPackage)
	segmentLimit := len(rootPackage) + 1
	for path, sourceUnit := range sourceUnits {
    	fileReader, err := os.Open(path)
    	if err != nil {
			log.Fatal(err)
    	}
		allDependencies, err := parsing.ParseJavaImports(fileReader)
		fileReader.Close()
		if err != nil {
			log.Fatal(err)
    	}
		sanitizedDependencies := utils.NewStringSet()
		for _, dependency := range(allDependencies) {
			if arrayStartsWith(dependency, rootPackage) {
				sanitizedDependencies.Add(parsing.ParseJavaJoinPathSegments(
					dependency[prefixLength:min(segmentLimit, len(dependency))]))
			}
		}
		sourceUnitString := parsing.ParseJavaJoinPathSegments(
			sourceUnit[prefixLength:min(segmentLimit, len(sourceUnit))])
		sanitizedDependencies.Remove(sourceUnitString)
		if dependencies[sourceUnitString] != nil {
			dependencies[sourceUnitString].AddSet(sanitizedDependencies)
		} else {
			dependencies[sourceUnitString] = sanitizedDependencies
		}
	}
	return dependencies
}

func arrayStartsWith(value []string, prefix []string) bool {
	if len(prefix) > len(value) {
		return false
	}
	for i, v := range(prefix) {
		if v != value[i] {
			return false
		}
	}
	return true
}


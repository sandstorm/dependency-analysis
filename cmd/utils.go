package cmd

import (
	"github.com/sandstorm/dependency-analysis/analysis"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strconv"
)

type analyzerSettingsFlags struct {
	depth          string
	includePattern string
}

func NewAnalyzerSettingsFlags() *analyzerSettingsFlags {
	return &analyzerSettingsFlags{
		depth:          "1",
		includePattern: ".*",
	}
}

func addAnalyzerSettingsFlags(cmd *cobra.Command, target *analyzerSettingsFlags) {
	defaultValues := NewAnalyzerSettingsFlags()

	cmd.Flags().StringVarP(&target.depth, "depth", "d", defaultValues.depth, "number of steps to go further down into the package hierarchy starting at the root package")
	cmd.Flags().StringVarP(&target.includePattern, "include", "", defaultValues.includePattern, "regular expression to filter files by their full path before analysis")
}

func (this *analyzerSettingsFlags) toAnalyzerSettings(sourcePath string) (*analysis.AnalyzerSettings, error) {
	result := &analysis.AnalyzerSettings{}
	var err error
	result.Depth, err = strconv.Atoi(this.depth)
	if err != nil {
		log.Fatal("failed to parse parameter 'depth'")
		log.Fatal(err)
		return nil, err
	}
	result.SourcePath = sourcePath
	result.IncludePattern = regexp.MustCompile(this.includePattern)
	return result, nil
}

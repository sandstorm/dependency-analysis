package analysis

import (
	"regexp"
)

type AnalyzerSettings struct {
	SourcePath     string
	Depth          int
	IncludePattern *regexp.Regexp
}

package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "customLinter",
	Doc:  "first linter demo",
	Run:  run,
}

package analyzer

import (
	"go/ast"

	"customlinter/internal/analyzer/extractor"
	"customlinter/internal/analyzer/rules/english"
	"customlinter/internal/analyzer/rules/lowercase"
	"customlinter/internal/analyzer/rules/sensitive"
	"customlinter/internal/analyzer/rules/specials"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "customLinter",
	Doc:  "first linter demo",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			msg, ok := extractor.Extract(pass, call)
			if !ok {
				return true
			}

			if lowercase.Check(pass, msg) {
				pass.Reportf(msg.MsgExpr.Pos(), lowercase.Message)
			}
			if english.Check(pass, msg) {
				pass.Reportf(msg.MsgExpr.Pos(), english.Message)
			}
			if specials.Check(pass, msg) {
				pass.Reportf(msg.MsgExpr.Pos(), specials.Message)
			}
			if sensitive.Check(pass, msg) {
				pass.Reportf(msg.MsgExpr.Pos(), sensitive.Message)
			}

			return true
		})
	}

	return nil, nil
}
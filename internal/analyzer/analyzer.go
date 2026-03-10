package analyzer

import (
	"go/ast"

	"github.com/myau69/customLinter/internal/analyzer/autofix"
	"github.com/myau69/customLinter/internal/analyzer/config"
	"github.com/myau69/customLinter/internal/analyzer/extractor"
	"github.com/myau69/customLinter/internal/analyzer/model"
	"github.com/myau69/customLinter/internal/analyzer/rules/english"
	"github.com/myau69/customLinter/internal/analyzer/rules/lowercase"
	"github.com/myau69/customLinter/internal/analyzer/rules/sensitive"
	"github.com/myau69/customLinter/internal/analyzer/rules/specials"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "customLinter",
	Doc:  "first linter demo",
	Run:  run,
}

var configPath string

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", "", "path to customlinter JSON config file")
}

func run(pass *analysis.Pass) (any, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}

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

			if cfg.Rules.Lowercase && lowercase.Check(pass, msg) {
				fixes := []analysis.SuggestedFix(nil)
				if cfg.AutoFix.Enabled && cfg.AutoFix.Lowercase {
					if suggested, ok := autofix.Lowercase(msg); ok {
						fixes = suggested
					}
				}
				report(pass, msg, lowercase.Message, fixes)
			}

			if cfg.Rules.English && english.Check(pass, msg) {
				report(pass, msg, english.Message, nil)
			}

			if cfg.Rules.Specials && specials.Check(pass, msg) {
				fixes := []analysis.SuggestedFix(nil)
				if cfg.AutoFix.Enabled && cfg.AutoFix.Specials {
					if suggested, ok := autofix.Specials(msg); ok {
						fixes = suggested
					}
				}
				report(pass, msg, specials.Message, fixes)
			}

			if cfg.Rules.Sensitive && sensitive.CheckWithPatterns(pass, msg, cfg.Sensitive.Patterns) {
				report(pass, msg, sensitive.Message, nil)
			}

			return true
		})
	}

	return nil, nil
}

func report(pass *analysis.Pass, msg model.LogMessage, message string, fixes []analysis.SuggestedFix) {
	diag := analysis.Diagnostic{
		Pos:     msg.MsgExpr.Pos(),
		End:     msg.MsgExpr.End(),
		Message: message,
	}

	if len(fixes) > 0 {
		diag.SuggestedFixes = fixes
	}

	pass.Report(diag)
}

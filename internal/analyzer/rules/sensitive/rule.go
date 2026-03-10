package sensitive

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"github.com/myau69/customLinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

const Message = "log message contains potentially sensitive data"

var defaultSensitiveWords = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"apikey",
}

func Check(pass *analysis.Pass, msg model.LogMessage) bool {
	return CheckWithPatterns(pass, msg, nil)
}

func CheckWithPatterns(pass *analysis.Pass, msg model.LogMessage, customPatterns []string) bool {
	words := buildSensitiveWords(customPatterns)
	if exprContainsSensitive(pass, msg.MsgExpr, words) {
		return true
	}

	if !msg.HasText || msg.Text == "" {
		return false
	}

	return containsSensitive(msg.Text, words)
}

func buildSensitiveWords(customPatterns []string) []string {
	seen := make(map[string]struct{}, len(defaultSensitiveWords)+len(customPatterns))
	words := make([]string, 0, len(defaultSensitiveWords)+len(customPatterns))

	add := func(raw string) {
		normalized := normalize(raw)
		if normalized == "" {
			return
		}
		if _, ok := seen[normalized]; ok {
			return
		}
		seen[normalized] = struct{}{}
		words = append(words, normalized)
	}

	for _, word := range defaultSensitiveWords {
		add(word)
	}
	for _, word := range customPatterns {
		add(word)
	}

	return words
}

func exprContainsSensitive(pass *analysis.Pass, expr ast.Expr, words []string) bool {
	if expr == nil {
		return false
	}

	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		if found {
			return false
		}

		switch n := node.(type) {
		case *ast.Ident:
			if isPkgName(pass, n) {
				return true
			}
			if containsSensitive(n.Name, words) {
				found = true
				return false
			}
		case *ast.SelectorExpr:
			if containsSensitive(n.Sel.Name, words) {
				found = true
				return false
			}
		case *ast.BasicLit:
			if n.Kind == token.STRING && containsSensitive(n.Value, words) {
				found = true
				return false
			}
		}

		return true
	})

	return found
}

func isPkgName(pass *analysis.Pass, ident *ast.Ident) bool {
	if pass == nil || pass.TypesInfo == nil || pass.TypesInfo.Uses == nil {
		return false
	}

	_, ok := pass.TypesInfo.Uses[ident].(*types.PkgName)
	return ok
}

func containsSensitive(s string, words []string) bool {
	normalized := normalize(s)
	for _, word := range words {
		if strings.Contains(normalized, word) {
			return true
		}
	}

	return false
}

func normalize(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range strings.ToLower(s) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}

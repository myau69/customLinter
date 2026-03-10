package sensitive

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"

	"customlinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

const Message = "log message contains potentially sensitive data"

var sensitiveWords = []string{
	"password",
	"passwd",
	"pwd",
	"secret",
	"token",
	"apikey",
}

func Check(pass *analysis.Pass, msg model.LogMessage) bool {
	if exprContainsSensitive(pass, msg.MsgExpr) {
		return true
	}

	if !msg.HasText || msg.Text == "" {
		return false
	}

	return containsSensitive(msg.Text)
}

func exprContainsSensitive(pass *analysis.Pass, expr ast.Expr) bool {
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
			if containsSensitive(n.Name) {
				found = true
				return false
			}
		case *ast.SelectorExpr:
			if containsSensitive(n.Sel.Name) {
				found = true
				return false
			}
		case *ast.BasicLit:
			if n.Kind == token.STRING && containsSensitive(n.Value) {
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

func containsSensitive(s string) bool {
	normalized := normalize(s)
	for _, word := range sensitiveWords {
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

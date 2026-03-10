package autofix

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/myau69/customLinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

func Lowercase(msg model.LogMessage) ([]analysis.SuggestedFix, bool) {
	lit, text, ok := stringLiteral(msg)
	if !ok || text == "" {
		return nil, false
	}

	firstRune, size := utf8.DecodeRuneInString(text)
	if firstRune == utf8.RuneError || unicode.IsLower(firstRune) {
		return nil, false
	}

	fixed := string(unicode.ToLower(firstRune)) + text[size:]
	if fixed == text {
		return nil, false
	}

	return singleReplace(lit, fixed, "make log message start with lowercase"), true
}

func Specials(msg model.LogMessage) ([]analysis.SuggestedFix, bool) {
	lit, text, ok := stringLiteral(msg)
	if !ok || text == "" {
		return nil, false
	}

	var b strings.Builder
	changed := false

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
			continue
		}
		changed = true
	}

	fixed := strings.Join(strings.Fields(b.String()), " ")
	if !changed || fixed == text {
		return nil, false
	}

	return singleReplace(lit, fixed, "remove special symbols from log message"), true
}

func stringLiteral(msg model.LogMessage) (*ast.BasicLit, string, bool) {
	lit, ok := msg.MsgExpr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil, "", false
	}

	text, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nil, "", false
	}

	return lit, text, true
}

func singleReplace(lit *ast.BasicLit, text, title string) []analysis.SuggestedFix {
	replacement := strconv.Quote(text)
	if strings.HasPrefix(lit.Value, "`") && !strings.ContainsRune(text, '`') {
		replacement = "`" + text + "`"
	}

	return []analysis.SuggestedFix{
		{
			Message: title,
			TextEdits: []analysis.TextEdit{
				{
					Pos:     lit.Pos(),
					End:     lit.End(),
					NewText: []byte(replacement),
				},
			},
		},
	}
}

package lowercase

import (
	"unicode"
	"unicode/utf8"

	"github.com/myau69/customLinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

const Message = "log message should start with a lowercase letter"

func Check(_ *analysis.Pass, msg model.LogMessage) bool {
	if !msg.HasText || msg.Text == "" {
		return false
	}

	firstRune, _ := utf8.DecodeRuneInString(msg.Text)
	if firstRune == utf8.RuneError {
		return false
	}

	return !unicode.IsLower(firstRune)
}

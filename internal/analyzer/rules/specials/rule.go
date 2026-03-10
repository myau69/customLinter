package specials

import (
	"unicode"

	"customlinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

const Message = "log message should not contain special symbols or emoji"

func Check(_ *analysis.Pass, msg model.LogMessage) bool {
	if !msg.HasText || msg.Text == "" {
		return false
	}

	for _, r := range msg.Text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		return true
	}

	return false
}

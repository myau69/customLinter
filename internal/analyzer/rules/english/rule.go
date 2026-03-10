package english

import (
	"unicode"

	"github.com/myau69/customLinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

const Message = "log message should contain only English letters"

func Check(_ *analysis.Pass, msg model.LogMessage) bool {
	if !msg.HasText || msg.Text == "" {
		return false
	}

	for _, r := range msg.Text {
		if unicode.IsLetter(r) && !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return true
		}
	}

	return false
}

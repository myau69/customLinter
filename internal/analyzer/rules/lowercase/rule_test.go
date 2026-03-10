package lowercase

import (
	"testing"

	"github.com/myau69/customLinter/internal/analyzer/model"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name string
		msg  model.LogMessage
		want bool
	}{
		{name: "no text", msg: model.LogMessage{HasText: false}, want: false},
		{name: "empty", msg: model.LogMessage{HasText: true, Text: ""}, want: false},
		{name: "lowercase start", msg: model.LogMessage{HasText: true, Text: "started"}, want: false},
		{name: "uppercase start", msg: model.LogMessage{HasText: true, Text: "Started"}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(nil, tt.msg); got != tt.want {
				t.Fatalf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

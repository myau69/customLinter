package english

import (
	"testing"

	"customlinter/internal/analyzer/model"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name string
		msg  model.LogMessage
		want bool
	}{
		{name: "no text", msg: model.LogMessage{HasText: false}, want: false},
		{name: "english only", msg: model.LogMessage{HasText: true, Text: "started server"}, want: false},
		{name: "cyrillic", msg: model.LogMessage{HasText: true, Text: "запуск"}, want: true},
		{name: "mixed", msg: model.LogMessage{HasText: true, Text: "server запущен"}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(nil, tt.msg); got != tt.want {
				t.Fatalf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

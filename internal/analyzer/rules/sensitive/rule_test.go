package sensitive

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"customlinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

func TestCheck(t *testing.T) {
	makePass := func() *analysis.Pass {
		return &analysis.Pass{TypesInfo: &types.Info{Uses: map[*ast.Ident]types.Object{}}}
	}

	makeString := func(value string) ast.Expr {
		return &ast.BasicLit{Kind: token.STRING, Value: value}
	}

	tests := []struct {
		name string
		msg  model.LogMessage
		want bool
	}{
		{
			name: "sensitive identifier",
			msg:  model.LogMessage{Call: &ast.CallExpr{}, MsgExpr: ast.NewIdent("password"), HasText: false},
			want: true,
		},
		{
			name: "delimiter pattern",
			msg:  model.LogMessage{Call: &ast.CallExpr{}, MsgExpr: makeString("\"token: abc\""), HasText: true, Text: "token: abc"},
			want: true,
		},
		{
			name: "format string with sensitive word",
			msg: model.LogMessage{
				Call:     &ast.CallExpr{Args: []ast.Expr{makeString("\"token %s\""), ast.NewIdent("value")}},
				MsgExpr:  makeString("\"token %s\""),
				MsgIndex: 0,
				HasText:  true,
				Text:     "token %s",
			},
			want: true,
		},
		{
			name: "safe message",
			msg:  model.LogMessage{Call: &ast.CallExpr{}, MsgExpr: makeString("\"started\""), HasText: true, Text: "started"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(makePass(), tt.msg); got != tt.want {
				t.Fatalf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckWithPatterns(t *testing.T) {
	pass := &analysis.Pass{TypesInfo: &types.Info{Uses: map[*ast.Ident]types.Object{}}}
	msg := model.LogMessage{
		Call:    &ast.CallExpr{},
		MsgExpr: &ast.BasicLit{Kind: token.STRING, Value: "\"session_id present\""},
		HasText: true,
		Text:    "session_id present",
	}

	if got := CheckWithPatterns(pass, msg, nil); got {
		t.Fatalf("CheckWithPatterns() = %v, want false without custom patterns", got)
	}
	if got := CheckWithPatterns(pass, msg, []string{"session_id"}); !got {
		t.Fatalf("CheckWithPatterns() = %v, want true with custom patterns", got)
	}
}

package autofix

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/myau69/customLinter/internal/analyzer/model"
)

func TestLowercase(t *testing.T) {
	msg := model.LogMessage{
		MsgExpr: &ast.BasicLit{Kind: token.STRING, Value: "\"Bad message\""},
	}

	fixes, ok := Lowercase(msg)
	if !ok {
		t.Fatal("Lowercase() expected suggested fix")
	}
	if len(fixes) != 1 || len(fixes[0].TextEdits) != 1 {
		t.Fatalf("unexpected suggested fixes shape: %#v", fixes)
	}
	if got, want := string(fixes[0].TextEdits[0].NewText), "\"bad message\""; got != want {
		t.Fatalf("new text = %q, want %q", got, want)
	}
}

func TestSpecials(t *testing.T) {
	msg := model.LogMessage{
		MsgExpr: &ast.BasicLit{Kind: token.STRING, Value: "\"started!!! now\""},
	}

	fixes, ok := Specials(msg)
	if !ok {
		t.Fatal("Specials() expected suggested fix")
	}
	if got, want := string(fixes[0].TextEdits[0].NewText), "\"started now\""; got != want {
		t.Fatalf("new text = %q, want %q", got, want)
	}
}

func TestNoFixForNonLiteral(t *testing.T) {
	msg := model.LogMessage{MsgExpr: ast.NewIdent("msg")}

	if _, ok := Lowercase(msg); ok {
		t.Fatal("Lowercase() should not suggest fix for non-literal")
	}
	if _, ok := Specials(msg); ok {
		t.Fatal("Specials() should not suggest fix for non-literal")
	}
}

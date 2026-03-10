package model

import "go/ast"

type LogMessage struct {
	Call     *ast.CallExpr
	MsgExpr  ast.Expr
	MsgIndex int
	Text     string
	HasText  bool
}

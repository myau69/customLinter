package extractor

import (
	"go/ast"
	"go/constant"
	"go/types"
	"strings"

	"customlinter/internal/analyzer/model"

	"golang.org/x/tools/go/analysis"
)

const (
	slogPkg = "log/slog"
	zapPkg  = "go.uber.org/zap"
)

var slogMsgIdxByFunc = map[string]int{
	"Debug":        0,
	"Info":         0,
	"Warn":         0,
	"Error":        0,
	"DebugContext": 1,
	"InfoContext":  1,
	"WarnContext":  1,
	"ErrorContext": 1,
	"Log":          2,
	"LogAttrs":     2,
}

func Extract(pass *analysis.Pass, call *ast.CallExpr) (model.LogMessage, bool) {
	msgExpr, msgIndex, ok := findMsg(pass, call)
	if !ok {
		return model.LogMessage{}, false
	}

	tv, ok := pass.TypesInfo.Types[msgExpr]
	if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
		return model.LogMessage{
			Call:     call,
			MsgExpr:  msgExpr,
			MsgIndex: msgIndex,
			HasText:  false,
		}, true
	}

	return model.LogMessage{
		Call:     call,
		MsgExpr:  msgExpr,
		MsgIndex: msgIndex,
		Text:     strings.TrimSpace(constant.StringVal(tv.Value)),
		HasText:  true,
	}, true
}

func findMsg(pass *analysis.Pass, call *ast.CallExpr) (ast.Expr, int, bool) {
	if len(call.Args) == 0 {
		return nil, -1, false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, -1, false
	}

	if idx, ok := slogMsgIdx(pass, sel, len(call.Args)); ok {
		return call.Args[idx], idx, true
	}

	selection := pass.TypesInfo.Selections[sel]
	if selection == nil {
		return nil, -1, false
	}

	fn, ok := selection.Obj().(*types.Func)
	if !ok || fn.Pkg() == nil {
		return nil, -1, false
	}

	path := fn.Pkg().Path()
	if path != slogPkg && path != zapPkg {
		return nil, -1, false
	}

	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig == nil {
		return nil, -1, false
	}

	idx := stringParamIdx(sig)
	if idx < 0 || idx >= len(call.Args) {
		return nil, -1, false
	}

	return call.Args[idx], idx, true
}

func slogMsgIdx(pass *analysis.Pass, sel *ast.SelectorExpr, argc int) (int, bool) {
	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return -1, false
	}

	pkgName, ok := pass.TypesInfo.Uses[pkgIdent].(*types.PkgName)
	if !ok || pkgName.Imported().Path() != slogPkg {
		return -1, false
	}

	idx, ok := slogMsgIdxByFunc[sel.Sel.Name]
	if !ok || idx < 0 || idx >= argc {
		return -1, false
	}

	return idx, true
}

func stringParamIdx(sig *types.Signature) int {
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		basic, ok := params.At(i).Type().Underlying().(*types.Basic)
		if ok && basic.Kind() == types.String {
			return i
		}
	}

	return -1
}
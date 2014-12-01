package tree

type defUseInfo struct {
	Defs int
	Uses int
}

type defUse struct {
	decl          *FuncDecl
	localToStruct map[LocalInfo_Ref]*defUseInfo
}

func (du *defUse) GetLocalInfo(index LocalInfo_Ref) *defUseInfo {
	info, ok := du.localToStruct[index]
	if !ok {
		info = &defUseInfo{}
		du.localToStruct[index] = info
	}
	return info
}

func makeApproxDefUse(decl *FuncDecl) *defUse {
	return &defUse{
		decl:          decl,
		localToStruct: map[LocalInfo_Ref]*defUseInfo{},
	}
}

func defUseExpr(expr Expr, du *defUse) {
	if expr == nil {
		return
	}
	switch expr := expr.(type) {
	case *GetLocal:
		du.GetLocalInfo(expr.Info).Uses += 1
	case *IntLiteral, *RuneLiteral, *BoolLiteral, *StringLiteral, *NilLiteral, *GetGlobal:
		// Leaf
	case *UnaryExpr:
		defUseExpr(expr.Expr, du)
	case *BinaryExpr:
		defUseExpr(expr.Left, du)
		defUseExpr(expr.Right, du)
	case *Call:
		defUseExpr(expr.Expr, du)
		for _, arg := range expr.Args {
			defUseExpr(arg, du)
		}
	case *StructLiteral:
		for _, arg := range expr.Args {
			defUseExpr(arg.Expr, du)
		}
	case *ListLiteral:
		for _, arg := range expr.Args {
			defUseExpr(arg, du)
		}
	case *Selector:
		defUseExpr(expr.Expr, du)
	case *Index:
		defUseExpr(expr.Expr, du)
		defUseExpr(expr.Index, du)
	case *TypeCoerce:
		defUseExpr(expr.Expr, du)
	case *TypeAssert:
		defUseExpr(expr.Expr, du)
	default:
		panic(expr)
	}
}

func defUseTarget(expr Target, du *defUse) {
	switch expr := expr.(type) {
	case *SetLocal:
		du.GetLocalInfo(expr.Info).Defs += 1
	default:
		panic(du.decl.Name)
	}
}

func defUseStmt(stmt Stmt, du *defUse) {
	switch stmt := stmt.(type) {
	case *Goto, *Label:
		// Leaf
	case *Assign:
		defUseExprList(stmt.Sources, du)
		for _, tgt := range stmt.Targets {
			defUseTarget(tgt, du)
		}
	case *Var:
		if stmt.Expr != nil {
			defUseExpr(stmt.Expr, du)
			du.GetLocalInfo(stmt.Info).Defs += 1
		}
	case *If:
		defUseExpr(stmt.Cond, du)
		defUseBlock(stmt.Body, du)
		if stmt.Else != nil {
			defUseStmt(stmt.Else, du)
		}
	case *BlockStmt:
		defUseBlock(stmt.Body, du)
	case *Return:
		defUseExprList(stmt.Args, du)
	default:
		expr, ok := stmt.(Expr)
		if ok {
			defUseExpr(expr, du)
		} else {
			panic(stmt)
		}
	}
}

func defUseExprList(exprs []Expr, du *defUse) {
	for _, expr := range exprs {
		defUseExpr(expr, du)
	}
}

func defUseBlock(stmts []Stmt, du *defUse) {
	for _, stmt := range stmts {
		defUseStmt(stmt, du)
	}
}

func defUseParam(param *Param, input bool, du *defUse) {
	// Outputs are implicitly zeroed.
	du.GetLocalInfo(param.Info).Defs += 1
	if !input {
		du.GetLocalInfo(param.Info).Uses += 1
	}
}

func defUseFunc(decl *FuncDecl, du *defUse) {
	if decl.Recv != nil {
		defUseParam(decl.Recv, true, du)
	}
	for _, p := range decl.Type.Params {
		defUseParam(p, true, du)
	}
	for _, p := range decl.Type.Results {
		defUseParam(p, false, du)
	}
	defUseBlock(decl.Body, du)
}

func pullLocal(expr *GetLocal, du *defUse, out []Stmt) (Expr, []Stmt) {
	info := du.GetLocalInfo(expr.Info)
	if info.Uses != 1 || info.Defs != 1 {
		return expr, out
	}

	n := len(out)
	if n <= 0 {
		return expr, out
	}
	last := out[n-1]
	lastAssign, ok := last.(*Assign)
	if !ok {
		return expr, out
	}
	if len(lastAssign.Targets) != 1 || len(lastAssign.Sources) != 1 {
		return expr, out
	}
	target, ok := lastAssign.Targets[0].(*SetLocal)
	if !ok || target.Info != expr.Info {
		return expr, out
	}
	return lastAssign.Sources[0], out[:n-1]
}

func consolidateExprList(exprs []Expr, du *defUse, out []Stmt) []Stmt {
	for i := len(exprs) - 1; i >= 0; i-- {
		exprs[i], out = consolidateExpr(exprs[i], du, out)
	}
	return out
}

func consolidateExpr(expr Expr, du *defUse, out []Stmt) (Expr, []Stmt) {
	if expr == nil {
		return nil, out
	}
	switch expr := expr.(type) {
	case *GetLocal:
		return pullLocal(expr, du, out)
	case *IntLiteral, *RuneLiteral, *BoolLiteral, *StringLiteral, *NilLiteral, *GetGlobal:
		// Leaf
	case *UnaryExpr:
		expr.Expr, out = consolidateExpr(expr.Expr, du, out)
	case *BinaryExpr:
		expr.Right, out = consolidateExpr(expr.Right, du, out)
		expr.Left, out = consolidateExpr(expr.Left, du, out)
	case *Call:
		out = consolidateExprList(expr.Args, du, out)
		expr.Expr, out = consolidateExpr(expr.Expr, du, out)
	case *Selector:
		expr.Expr, out = consolidateExpr(expr.Expr, du, out)
	case *Index:
		expr.Index, out = consolidateExpr(expr.Index, du, out)
		expr.Expr, out = consolidateExpr(expr.Expr, du, out)
	case *TypeCoerce:
		expr.Expr, out = consolidateExpr(expr.Expr, du, out)
	case *TypeAssert:
		expr.Expr, out = consolidateExpr(expr.Expr, du, out)
	case *StructLiteral:
		for i := len(expr.Args) - 1; i >= 0; i-- {
			expr.Args[i].Expr, out = consolidateExpr(expr.Args[i].Expr, du, out)
		}
	case *ListLiteral:
		out = consolidateExprList(expr.Args, du, out)
	default:
		panic(expr)
	}
	return expr, out
}

func consolidateStmt(stmt Stmt, du *defUse, out []Stmt) []Stmt {
	switch stmt := stmt.(type) {
	case *Goto, *Label:
		// Leaf
	case *Assign:
		for i := len(stmt.Sources) - 1; i >= 0; i-- {
			stmt.Sources[i], out = consolidateExpr(stmt.Sources[i], du, out)
		}
	case *Var:
		stmt.Expr, out = consolidateExpr(stmt.Expr, du, out)
	case *If:
		stmt.Cond, out = consolidateExpr(stmt.Cond, du, out)
		stmt.Body = consolidateBlock(stmt.Body, du)
		if stmt.Else != nil {
			consolidateStmt(stmt.Else, du, nil)
		}
	case *BlockStmt:
		stmt.Body = consolidateBlock(stmt.Body, du)
	case *Return:
		out = consolidateExprList(stmt.Args, du, out)
	default:
		expr, ok := stmt.(Expr)
		if ok {
			expr, out = consolidateExpr(expr, du, out)
			stmt, _ = expr.(Stmt)
		} else {
			panic(stmt)
		}
	}
	return out
}

func consolidateBlock(stmts []Stmt, du *defUse) []Stmt {
	out := []Stmt{}
	for _, stmt := range stmts {
		out = consolidateStmt(stmt, du, out)
		out = append(out, stmt)
	}
	return out
}

func ConsolidateFunc(decl *FuncDecl) {
	if decl.Body != nil {
		du := makeApproxDefUse(decl)
		defUseFunc(decl, du)
		decl.Body = consolidateBlock(decl.Body, du)
	}
}

func consolidateDecl(decl Decl) {
	switch decl := decl.(type) {
	case *InterfaceDecl, *StructDecl, *TypeDefDecl, *VarDecl:
		// Leaf
	case *FuncDecl:
		ConsolidateFunc(decl)
	default:
		panic(decl)
	}
}

func Consolidate(prog *ProgramAST) {
	for _, pkg := range prog.Packages {
		if pkg.Extern {
			continue
		}
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				consolidateDecl(decl)
			}
		}
	}

}

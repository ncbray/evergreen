package tree

type approxDefUseInfo struct {
	Defs int
	Uses int
}

type approxDefUse struct {
	nameToStruct map[string]*approxDefUseInfo
	removed      map[string]bool
}

func (du *approxDefUse) GetInfo(name string) *approxDefUseInfo {
	info, ok := du.nameToStruct[name]
	if !ok {
		info = &approxDefUseInfo{}
		du.nameToStruct[name] = info
	}
	return info
}

func makeApproxDefUse() *approxDefUse {
	return &approxDefUse{
		nameToStruct: map[string]*approxDefUseInfo{},
		removed:      map[string]bool{},
	}
}

func approxDefUseExpr(expr Expr, du *approxDefUse) {
	if expr == nil {
		return
	}
	switch expr := expr.(type) {
	case *NameRef:
		du.GetInfo(expr.Text).Uses += 1
	case *IntLiteral, *RuneLiteral, *BoolLiteral, *StringLiteral, *NilLiteral:
		// Leaf
	case *UnaryExpr:
		approxDefUseExpr(expr.Expr, du)
	case *BinaryExpr:
		approxDefUseExpr(expr.Left, du)
		approxDefUseExpr(expr.Right, du)
	case *Call:
		approxDefUseExpr(expr.Expr, du)
		for _, arg := range expr.Args {
			approxDefUseExpr(arg, du)
		}
	case *StructLiteral:
		for _, arg := range expr.Args {
			approxDefUseExpr(arg.Expr, du)
		}
	case *ListLiteral:
		for _, arg := range expr.Args {
			approxDefUseExpr(arg, du)
		}
	case *Selector:
		approxDefUseExpr(expr.Expr, du)
	case *Index:
		approxDefUseExpr(expr.Expr, du)
		approxDefUseExpr(expr.Index, du)
	case *TypeCoerce:
		approxDefUseExpr(expr.Expr, du)
	case *TypeAssert:
		approxDefUseExpr(expr.Expr, du)
	default:
		panic(expr)
	}
}

func approxDefUseTarget(expr Expr, du *approxDefUse) {
	switch expr := expr.(type) {
	case *NameRef:
		du.GetInfo(expr.Text).Defs += 1
	default:
		panic(expr)
	}
}

func approxDefUseStmt(stmt Stmt, du *approxDefUse) {
	switch stmt := stmt.(type) {
	case *Goto, *Label:
		// Leaf
	case *Assign:
		approxDefUseExprList(stmt.Sources, du)
		for _, tgt := range stmt.Targets {
			approxDefUseTarget(tgt, du)
		}
	case *Var:
		if stmt.Expr != nil {
			approxDefUseExpr(stmt.Expr, du)
			du.GetInfo(stmt.Name).Defs += 1
		}
	case *If:
		approxDefUseExpr(stmt.Cond, du)
		approxDefUseBlock(stmt.Body, du)
		if stmt.Else != nil {
			approxDefUseStmt(stmt.Else, du)
		}
	case *BlockStmt:
		approxDefUseBlock(stmt.Body, du)
	case *Return:
		approxDefUseExprList(stmt.Args, du)
	default:
		expr, ok := stmt.(Expr)
		if ok {
			approxDefUseExpr(expr, du)
		} else {
			panic(stmt)
		}
	}
}

func approxDefUseExprList(exprs []Expr, du *approxDefUse) {
	for _, expr := range exprs {
		approxDefUseExpr(expr, du)
	}
}

func approxDefUseBlock(stmts []Stmt, du *approxDefUse) {
	for _, stmt := range stmts {
		approxDefUseStmt(stmt, du)
	}
}

func approxDefUseParam(param *Param, input bool, du *approxDefUse) {
	// Outputs are implicitly zeroed.
	du.GetInfo(param.Name).Defs += 1
	if !input {
		du.GetInfo(param.Name).Uses += 1
	}
}

func approxDefUseFunc(decl *FuncDecl, du *approxDefUse) {
	if decl.Recv != nil {
		approxDefUseParam(decl.Recv, true, du)
	}
	for _, p := range decl.Type.Params {
		approxDefUseParam(p, true, du)
	}
	for _, p := range decl.Type.Results {
		approxDefUseParam(p, false, du)
	}
	approxDefUseBlock(decl.Body, du)
}

func pullName(expr *NameRef, du *approxDefUse, out []Stmt) (Expr, []Stmt) {
	info := du.GetInfo(expr.Text)
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
	target, ok := lastAssign.Targets[0].(*NameRef)
	if !ok || target.Text != expr.Text {
		return expr, out
	}
	du.removed[expr.Text] = true
	return lastAssign.Sources[0], out[:n-1]
}

func retreeExprList(exprs []Expr, du *approxDefUse, out []Stmt) []Stmt {
	for i := len(exprs) - 1; i >= 0; i-- {
		exprs[i], out = retreeExpr(exprs[i], du, out)
	}
	return out
}

func retreeExpr(expr Expr, du *approxDefUse, out []Stmt) (Expr, []Stmt) {
	if expr == nil {
		return nil, out
	}
	switch expr := expr.(type) {
	case *NameRef:
		return pullName(expr, du, out)
	case *IntLiteral, *RuneLiteral, *BoolLiteral, *StringLiteral, *NilLiteral:
		// Leaf
	case *UnaryExpr:
		expr.Expr, out = retreeExpr(expr.Expr, du, out)
	case *BinaryExpr:
		expr.Right, out = retreeExpr(expr.Right, du, out)
		expr.Left, out = retreeExpr(expr.Left, du, out)
	case *Call:
		out = retreeExprList(expr.Args, du, out)
		expr.Expr, out = retreeExpr(expr.Expr, du, out)
	case *Selector:
		expr.Expr, out = retreeExpr(expr.Expr, du, out)
	case *Index:
		expr.Index, out = retreeExpr(expr.Index, du, out)
		expr.Expr, out = retreeExpr(expr.Expr, du, out)
	case *TypeCoerce:
		expr.Expr, out = retreeExpr(expr.Expr, du, out)
	case *TypeAssert:
		expr.Expr, out = retreeExpr(expr.Expr, du, out)
	case *StructLiteral:
		for i := len(expr.Args) - 1; i >= 0; i-- {
			expr.Args[i].Expr, out = retreeExpr(expr.Args[i].Expr, du, out)
		}
	case *ListLiteral:
		out = retreeExprList(expr.Args, du, out)
	default:
		panic(expr)
	}
	return expr, out
}

func retreeStmt(stmt Stmt, du *approxDefUse, out []Stmt) []Stmt {
	switch stmt := stmt.(type) {
	case *Goto, *Label:
		// Leaf
	case *Assign:
		for i := len(stmt.Sources) - 1; i >= 0; i-- {
			stmt.Sources[i], out = retreeExpr(stmt.Sources[i], du, out)
		}
	case *Var:
		stmt.Expr, out = retreeExpr(stmt.Expr, du, out)
	case *If:
		stmt.Cond, out = retreeExpr(stmt.Cond, du, out)
		stmt.Body = retreeBlock(stmt.Body, du)
		if stmt.Else != nil {
			retreeStmt(stmt.Else, du, nil)
		}
	case *BlockStmt:
		stmt.Body = retreeBlock(stmt.Body, du)
	case *Return:
		out = retreeExprList(stmt.Args, du, out)
	default:
		expr, ok := stmt.(Expr)
		if ok {
			expr, out = retreeExpr(expr, du, out)
			stmt, _ = expr.(Stmt)
		} else {
			panic(stmt)
		}
	}
	return out
}

func retreeBlock(stmts []Stmt, du *approxDefUse) []Stmt {
	out := []Stmt{}
	for _, stmt := range stmts {
		out = retreeStmt(stmt, du, out)
		out = append(out, stmt)
	}
	return out
}

func removeVarHackStmt(stmt Stmt, du *approxDefUse, out []Stmt) []Stmt {
	switch stmt := stmt.(type) {
	case *Goto, *Label, *Assign, *Call, *Return:
		return append(out, stmt)
	case *If:
		stmt.Body = removeVarHackBlock(stmt.Body, du)
	case *BlockStmt:
		stmt.Body = removeVarHackBlock(stmt.Body, du)
	case *Var:
		removed, _ := du.removed[stmt.Name]
		if removed {
			return out
		}
	default:
		panic(stmt)
	}
	return append(out, stmt)
}

func removeVarHackBlock(stmts []Stmt, du *approxDefUse) []Stmt {
	out := []Stmt{}
	for _, stmt := range stmts {
		out = removeVarHackStmt(stmt, du, out)
	}
	return out
}

func retreeDecl(decl Decl) {
	switch decl := decl.(type) {
	case *InterfaceDecl, *StructDecl:
		// Leaf
	case *FuncDecl:
		du := makeApproxDefUse()
		if decl.Body != nil {
			approxDefUseFunc(decl, du)
			decl.Body = retreeBlock(decl.Body, du)
			decl.Body = removeVarHackBlock(decl.Body, du)
		}
	default:
		panic(decl)
	}
}

func retreeFile(file *File) {
	for _, decl := range file.Decls {
		retreeDecl(decl)
	}
}

func Retree(prog *Program) {
	for _, pkg := range prog.Packages {
		if pkg.Extern {
			continue
		}
		for _, file := range pkg.Files {
			retreeFile(file)
		}
	}

}

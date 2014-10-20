package tree

type refRewriter interface {
	rewriteLocalInfo(index int) int
}

// Find all LocalInfos
type funcGC struct {
	live []bool
}

func (rewriter *funcGC) rewriteLocalInfo(index int) int {
	if index < 0 {
		panic(index)
	}
	rewriter.live[index] = true
	return index
}

// Rewrite all LocalInfos
type funcRemap struct {
	remap []int
}

func (rewriter *funcRemap) rewriteLocalInfo(index int) int {
	if index < 0 {
		panic(index)
	}
	return rewriter.remap[index]
}

func sweepExprList(exprs []Expr, rewriter refRewriter) {
	for _, expr := range exprs {
		sweepExpr(expr, rewriter)
	}
}

func sweepExpr(expr Expr, rewriter refRewriter) {
	switch expr := expr.(type) {
	case *IntLiteral, *BoolLiteral, *StringLiteral, *RuneLiteral, *NilLiteral, *GetGlobal, *GetName:
		// Leaf
	case *GetLocal:
		expr.Info = rewriter.rewriteLocalInfo(expr.Info)
	case *TypeAssert:
		sweepExpr(expr.Expr, rewriter)
	case *TypeCoerce:
		sweepExpr(expr.Expr, rewriter)
	case *Selector:
		sweepExpr(expr.Expr, rewriter)
	case *Index:
		sweepExpr(expr.Expr, rewriter)
		sweepExpr(expr.Index, rewriter)
	case *Call:
		sweepExpr(expr.Expr, rewriter)
		sweepExprList(expr.Args, rewriter)
	case *UnaryExpr:
		sweepExpr(expr.Expr, rewriter)
	case *BinaryExpr:
		sweepExpr(expr.Left, rewriter)
		sweepExpr(expr.Right, rewriter)
	case *StructLiteral:
		for _, arg := range expr.Args {
			sweepExpr(arg.Expr, rewriter)
		}
	case *ListLiteral:
		sweepExprList(expr.Args, rewriter)
	default:
		panic(expr)
	}
}

func sweepTarget(expr Target, rewriter refRewriter) {
	switch expr := expr.(type) {
	case *SetLocal:
		expr.Info = rewriter.rewriteLocalInfo(expr.Info)
	case *SetName:
		// Leaf
	default:
		panic(expr)
	}
}

func sweepStmt(stmt Stmt, rewriter refRewriter) {
	switch stmt := stmt.(type) {
	case *Goto, *Label:
		// Leaf
	case *Assign:
		for _, src := range stmt.Sources {
			sweepExpr(src, rewriter)
		}
		for _, tgt := range stmt.Targets {
			sweepTarget(tgt, rewriter)
		}
	case *If:
		sweepExpr(stmt.Cond, rewriter)
		sweepBlock(stmt.Body, rewriter)
		if stmt.Else != nil {
			sweepStmt(stmt.Else, rewriter)
		}
	case *BlockStmt:
		sweepBlock(stmt.Body, rewriter)
	case *Return:
		sweepExprList(stmt.Args, rewriter)
	default:
		expr, ok := stmt.(Expr)
		if ok {
			sweepExpr(expr, rewriter)
		} else {
			panic(stmt)
		}
	}
}

func sweepBlock(stmts []Stmt, rewriter refRewriter) {
	for _, stmt := range stmts {
		sweepStmt(stmt, rewriter)
	}
}

func sweepParam(param *Param, rewriter refRewriter) {
	param.Info = rewriter.rewriteLocalInfo(param.Info)
}

func sweepFunc(decl *FuncDecl, rewriter refRewriter) {
	// TODO self
	if decl.Recv != nil {
		sweepParam(decl.Recv, rewriter)
	}
	for _, p := range decl.Type.Params {
		sweepParam(p, rewriter)
	}
	for _, p := range decl.Type.Results {
		sweepParam(p, rewriter)
	}
	sweepBlock(decl.Body, rewriter)
}

func MakeRemap(live []bool) ([]int, int) {
	remap := make([]int, len(live))
	count := 0
	for i, isLive := range live {
		var idx int
		if isLive {
			idx = count
			count += 1
		} else {
			idx = -1
		}
		remap[i] = idx
	}
	return remap, count
}

func CompactFunc(decl *FuncDecl) {
	sweep := &funcGC{live: make([]bool, len(decl.Locals))}
	sweepFunc(decl, sweep)

	remap, count := MakeRemap(sweep.live)
	sweepFunc(decl, &funcRemap{remap: remap})

	locals := make([]*LocalInfo, count)
	for i, info := range decl.Locals {
		idx := remap[i]
		if idx >= 0 {
			locals[idx] = info
		}
	}
	decl.Locals = locals
}

func isParam(decl *FuncDecl, ref int) bool {
	if ref < 0 {
		return false
	}
	if decl.Recv != nil {
		if decl.Recv.Info == ref {
			return true
		}
	}
	for _, p := range decl.Type.Params {
		if p.Info == ref {
			return true
		}
	}
	for _, p := range decl.Type.Results {
		if p.Info == ref {
			return true
		}
	}
	return false
}

// Assumes the function has been compacted.
func InsertVarDecls(decl *FuncDecl) {
	stmts := []Stmt{}

	// Declare the variables up front.
	// It is easier to do this than precisely calculate where they need to be defined.
	for i, info := range decl.Locals {
		if isParam(decl, i) {
			continue
		}
		stmts = append(stmts, &Var{
			Name: info.Name,
			// HACK should copy this instead of sharing it
			Type: info.T,
			Info: i,
		})
	}

	decl.Body = append(stmts, decl.Body...)
}

func (decl *FuncDecl) CreateLocalInfo(name string, T Type) int {
	idx := len(decl.Locals)
	decl.Locals = append(decl.Locals, &LocalInfo{
		Name: name,
		T:    T,
	})
	return idx
}

func (decl *FuncDecl) GetLocalInfo(idx int) *LocalInfo {
	if idx >= len(decl.Locals) || idx < 0 {
		panic(idx)
	}
	return decl.Locals[idx]
}

func (decl *FuncDecl) MakeParam(idx int) *Param {
	localInfo := decl.GetLocalInfo(idx)
	return &Param{
		Name: localInfo.Name,
		Type: localInfo.T,
		Info: idx,
	}
}

func (decl *FuncDecl) MakeGetLocal(idx int) Expr {
	if idx >= len(decl.Locals) || idx < 0 {
		panic(idx)
	}
	return &GetLocal{
		Info: idx,
	}
}

func (decl *FuncDecl) MakeSetLocal(idx int) Target {
	if idx >= len(decl.Locals) || idx < 0 {
		panic(idx)
	}
	return &SetLocal{
		Info: idx,
	}
}

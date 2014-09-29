package tree

type refRewriter interface {
	rewriteLocalInfo(index int) int
}

// Find all LocalInfos
type funcGC struct {
	live []bool
}

func (rewriter *funcGC) rewriteLocalInfo(index int) int {
	if index >= 0 {
		rewriter.live[index] = true
	}
	return index
}

// Rewrite all LocalInfos
type funcRemap struct {
	remap []int
}

func (rewriter *funcRemap) rewriteLocalInfo(index int) int {
	if index >= 0 {
		return rewriter.remap[index]
	}
	return index
}

func sweepExprList(exprs []Expr, rewriter refRewriter) {
	for _, expr := range exprs {
		sweepExpr(expr, rewriter)
	}
}

func sweepExpr(expr Expr, rewriter refRewriter) {
	switch expr := expr.(type) {
	case *IntLiteral, *BoolLiteral, *StringLiteral, *RuneLiteral, *NilLiteral:
		// Leaf
	case *NameRef:
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

func sweepStmt(stmt Stmt, rewriter refRewriter) {
	switch stmt := stmt.(type) {
	case *Goto, *Label:
		// Leaf
	case *Assign:
		for _, src := range stmt.Sources {
			sweepExpr(src, rewriter)
		}
		for _, tgt := range stmt.Targets {
			sweepExpr(tgt, rewriter)
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
	sweepBlock(decl.Body, sweep)

	remap, count := MakeRemap(sweep.live)
	sweepBlock(decl.Body, &funcRemap{remap: remap})

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

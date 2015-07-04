package tree

import (
	"evergreen/go/core"
)

type refRewriter interface {
	rewriteLocalInfo(index *LocalInfo) *LocalInfo
}

// Find all LocalInfos
type funcGC struct {
	live []bool
}

func (rewriter *funcGC) rewriteLocalInfo(lcl *LocalInfo) *LocalInfo {
	rewriter.live[lcl.Index] = true
	return lcl
}

func sweepExprList(exprs []Expr, rewriter refRewriter) {
	for _, expr := range exprs {
		sweepExpr(expr, rewriter)
	}
}

func sweepExpr(expr Expr, rewriter refRewriter) {
	switch expr := expr.(type) {
	case *IntLiteral, *Float32Literal, *BoolLiteral, *StringLiteral, *RuneLiteral, *NilLiteral, *GetGlobal, *GetFunction, *GetName:
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
		sweepBlock(stmt.T, rewriter)
		if stmt.F != nil {
			sweepBlock(stmt.F, rewriter)
		}
	case *For:
		sweepBlock(stmt.Block, rewriter)
	case *BlockStmt:
		sweepBlock(stmt.Block, rewriter)
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

func sweepBlock(block *Block, rewriter refRewriter) {
	for _, stmt := range block.Body {
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
	if decl.Block == nil {
		panic(decl.Name)
	}
	sweepBlock(decl.Block, rewriter)
}

func MakeRemap(live []bool) ([]LocalInfo_Ref, int) {
	remap := make([]LocalInfo_Ref, len(live))
	count := 0
	for i, isLive := range live {
		var idx int
		if isLive {
			idx = count
			count += 1
		} else {
			idx = -1
		}
		remap[i] = LocalInfo_Ref(idx)
	}
	return remap, count
}

func CompactFunc(decl *FuncDecl) {
	sweep := &funcGC{live: make([]bool, decl.LocalInfo_Scope.Len())}
	sweepFunc(decl, sweep)

	remap, count := MakeRemap(sweep.live)
	decl.LocalInfo_Scope.Remap(remap, count)
}

func isParam(decl *FuncDecl, lcl *LocalInfo) bool {
	if lcl == nil {
		return false
	}
	if decl.Recv != nil && decl.Recv.Info == lcl {
		return true
	}
	for _, p := range decl.Type.Params {
		if p.Info == lcl {
			return true
		}
	}
	for _, p := range decl.Type.Results {
		if p.Info == lcl {
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
	iter := decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		info := iter.Value()
		if isParam(decl, info) {
			continue
		}
		stmts = append(stmts, &Var{
			Name: info.Name,
			// HACK should copy this instead of sharing it
			Type: info.T,
			Info: info,
		})
	}

	decl.Block.Body = append(stmts, decl.Block.Body...)
}

func (scope *LocalInfo_Scope) Get(ref LocalInfo_Ref) *LocalInfo {
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *LocalInfo_Scope) Register(info *LocalInfo) *LocalInfo {
	info.Index = LocalInfo_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info
}

func (scope *LocalInfo_Scope) Len() int {
	return len(scope.objects)
}

func (scope *LocalInfo_Scope) Iter() *localInfoIterator {
	return &localInfoIterator{scope: scope, current: -1}
}

func (scope *LocalInfo_Scope) Remap(remap []LocalInfo_Ref, count int) {
	objects := make([]*LocalInfo, count)
	for i, info := range scope.objects {
		idx := remap[i]
		info.Index = idx
		if idx != ^LocalInfo_Ref(0) {
			objects[idx] = info
		}
	}
	scope.objects = objects
}

type localInfoIterator struct {
	scope   *LocalInfo_Scope
	current int
}

func (iter *localInfoIterator) Next() bool {
	iter.current += 1
	return iter.current < len(iter.scope.objects)
}

func (iter *localInfoIterator) Index() LocalInfo_Ref {
	return LocalInfo_Ref(iter.current)
}

func (iter *localInfoIterator) Value() *LocalInfo {
	return iter.scope.objects[iter.current]
}

func (decl *FuncDecl) CreateLocalInfo(name string, T TypeRef) *LocalInfo {
	return decl.LocalInfo_Scope.Register(&LocalInfo{
		Name: name,
		T:    T,
	})
}

func (decl *FuncDecl) GetLocalInfo(idx LocalInfo_Ref) *LocalInfo {
	return decl.LocalInfo_Scope.Get(idx)
}

func RefForType(t core.GoType) TypeRef {
	switch t := t.(type) {
	case *core.ExternalType:
		return &NameRef{
			Name: t.Name,
			T:    t,
		}
	case *core.StructType:
		return &NameRef{
			Name: t.Name,
			T:    t,
		}
	case *core.InterfaceType:
		return &NameRef{
			Name: t.Name,
			T:    t,
		}
	case *core.TypeDefType:
		return &NameRef{
			Name: t.Name,
			T:    t,
		}
	case *core.FuncType:
		params := make([]*Param, len(t.Params))
		for i, pt := range t.Params {
			params[i] = &Param{
				Type: RefForType(pt),
			}
		}
		results := make([]*Param, len(t.Results))
		for i, pt := range t.Results {
			results[i] = &Param{
				Type: RefForType(pt),
			}
		}
		return &FuncTypeRef{
			Params:  params,
			Results: results,
		}
	case *core.PointerType:
		return &PointerRef{
			Element: RefForType(t.Element),
			T:       t,
		}
	case *core.SliceType:
		return &SliceRef{
			Element: RefForType(t.Element),
			T:       t,
		}
	default:
		panic(t)
	}
}

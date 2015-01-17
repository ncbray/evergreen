package tree

import (
	"evergreen/go/core"
)

type refRewriter interface {
	rewriteLocalInfo(index LocalInfo_Ref) LocalInfo_Ref
}

// Find all LocalInfos
type funcGC struct {
	live []bool
}

func (rewriter *funcGC) rewriteLocalInfo(index LocalInfo_Ref) LocalInfo_Ref {
	if index < 0 {
		panic(index)
	}
	rewriter.live[index] = true
	return index
}

// Rewrite all LocalInfos
type funcRemap struct {
	remap []LocalInfo_Ref
}

func (rewriter *funcRemap) rewriteLocalInfo(index LocalInfo_Ref) LocalInfo_Ref {
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
	case *IntLiteral, *Float32Literal, *BoolLiteral, *StringLiteral, *RuneLiteral, *NilLiteral, *GetGlobal, *GetName:
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
	sweepFunc(decl, &funcRemap{remap: remap})

	decl.LocalInfo_Scope.Remap(remap, count)
}

func isParam(decl *FuncDecl, ref LocalInfo_Ref) bool {
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
	iter := decl.LocalInfo_Scope.Iter()
	for iter.Next() {
		i := iter.Index()
		info := iter.Value()
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

func (scope *LocalInfo_Scope) Get(ref LocalInfo_Ref) *LocalInfo {
	return scope.objects[ref]
}

func (scope *LocalInfo_Scope) Register(info *LocalInfo) LocalInfo_Ref {
	index := LocalInfo_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return index
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
		if idx != NoLocalInfo {
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

func (decl *FuncDecl) CreateLocalInfo(name string, T TypeRef) LocalInfo_Ref {
	return decl.LocalInfo_Scope.Register(&LocalInfo{
		Name: name,
		T:    T,
	})
}

func (decl *FuncDecl) GetLocalInfo(idx LocalInfo_Ref) *LocalInfo {
	return decl.LocalInfo_Scope.Get(idx)
}

func (decl *FuncDecl) MakeParam(idx LocalInfo_Ref) *Param {
	localInfo := decl.GetLocalInfo(idx)
	return &Param{
		Name: localInfo.Name,
		Type: localInfo.T,
		Info: idx,
	}
}

func (decl *FuncDecl) MakeGetLocal(idx LocalInfo_Ref) Expr {
	return &GetLocal{
		Info: idx,
	}
}

func (decl *FuncDecl) MakeSetLocal(idx LocalInfo_Ref) Target {
	return &SetLocal{
		Info: idx,
	}
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

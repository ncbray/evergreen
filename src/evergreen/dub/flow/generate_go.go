package flow

import (
	"evergreen/base"
	ast "evergreen/go/tree"
	"fmt"
)

type DubToGoLinker interface {
	TypeRef(r *ast.TypeRef, s *LLStruct)
	ForwardType(s *LLStruct, impl ast.TypeImpl)
	Finish()
}

type linkElement struct {
	impl ast.TypeImpl
	refs []*ast.TypeRef
}

type linkerImpl struct {
	types map[*LLStruct]*linkElement
}

func (l *linkerImpl) get(s *LLStruct) *linkElement {
	e, ok := l.types[s]
	if !ok {
		e = &linkElement{}
		l.types[s] = e
	}
	return e
}

func (l *linkerImpl) TypeRef(r *ast.TypeRef, s *LLStruct) {
	e := l.get(s)
	e.refs = append(e.refs, r)
}

func (l *linkerImpl) ForwardType(s *LLStruct, impl ast.TypeImpl) {
	e := l.get(s)
	e.impl = impl
}

func (l *linkerImpl) Finish() {
	for s, e := range l.types {
		if e.impl == nil {
			panic(s)
		}
		for _, r := range e.refs {
			r.Impl = e.impl
		}
		e.refs = nil
	}
}

func MakeLinker() DubToGoLinker {
	return &linkerImpl{
		types: map[*LLStruct]*linkElement{},
	}
}

// Begin AST construction wrappers

func id(name string) ast.Expr {
	return &ast.NameRef{Text: name}
}

func attr(expr ast.Expr, name string) ast.Expr {
	return &ast.Selector{Expr: expr, Text: name}
}

func strLiteral(value string) ast.Expr {
	return &ast.StringLiteral{Value: value}
}

func intLiteral(value int) ast.Expr {
	return &ast.IntLiteral{Value: value}
}

// End AST construction wrappers

type DubToGoContext struct {
	index *BuiltinIndex
	state *ast.StructDecl
	link  DubToGoLinker
}

type blockInfo struct {
	label string
}

func blockName(i int) string {
	return fmt.Sprintf("block%d", i)
}

func builtinStmt(name string, args ...ast.Expr) ast.Stmt {
	return &ast.Call{Expr: attr(id("frame"), name), Args: args}
}

func builtinExpr(name string, args ...ast.Expr) ast.Expr {
	return &ast.Call{Expr: attr(id("frame"), name), Args: args}
}

func constInt(v int64) ast.Expr {
	// Type converstion hack
	return intLiteral(int(v))
}

func reg(r DubRegister) ast.Expr {
	return id(RegisterName(r))
}

func returnVarName(i int) string {
	return fmt.Sprintf("ret%d", i)
}

func opAssign(expr ast.Expr, dst DubRegister) ast.Stmt {
	if dst != NoRegister {
		return &ast.Assign{
			Targets: []ast.Expr{reg(dst)},
			Op:      "=",
			Sources: []ast.Expr{expr},
		}
	} else {
		// TODO fix expr / stmt duality.
		return expr.(ast.Stmt)
	}
}

func opMultiAssign(expr ast.Expr, dsts []DubRegister) ast.Stmt {
	if len(dsts) != 0 {
		lhs := make([]ast.Expr, len(dsts))
		for i, dst := range dsts {
			if dst != NoRegister {
				lhs[i] = reg(dst)
			} else {
				lhs[i] = id("_")
			}
		}
		return &ast.Assign{
			Targets: lhs,
			Op:      "=",
			Sources: []ast.Expr{expr},
		}
	} else {
		// TODO unhack
		return expr.(ast.Stmt)
	}
}

func goTypeName(t DubType, ctx *DubToGoContext) ast.Type {
	switch t := t.(type) {
	case *BoolType:
		return &ast.TypeRef{Name: "bool", Impl: ctx.index.BoolType}
	case *IntType:
		return &ast.TypeRef{Name: "int", Impl: ctx.index.IntType}
	case *RuneType:
		return &ast.TypeRef{Name: "rune", Impl: ctx.index.RuneType}
	case *StringType:
		return &ast.TypeRef{Name: "string", Impl: ctx.index.StringType}
	case *ListType:
		return &ast.SliceType{Element: goTypeName(t.Type, ctx)}
	case *LLStruct:
		out := goStructType(t, ctx)
		if t.Abstract {
			return out
		} else {
			return &ast.PointerType{Element: out}
		}
	default:
		panic(t)
	}
}

func goStructType(t *LLStruct, ctx *DubToGoContext) *ast.TypeRef {
	out := &ast.TypeRef{Name: t.Name}
	ctx.link.TypeRef(out, t)
	return out
}

func GenerateOp(f *LLFunc, op DubOp, ctx *DubToGoContext, block []ast.Stmt) []ast.Stmt {
	if IsNop(op) {
		return block
	}

	switch op := op.(type) {
	case *BinaryOp:
		// TODO validate Op?
		block = append(block, opAssign(
			&ast.BinaryExpr{
				Left:  reg(op.Left),
				Op:    op.Op,
				Right: reg(op.Right),
			},
			op.Dst,
		))
	case *CallOp:
		args := []ast.Expr{
			id("frame"),
		}
		for _, arg := range op.Args {
			args = append(args, reg(arg))
		}
		block = append(block, opMultiAssign(
			&ast.Call{
				Expr: id(op.Name),
				Args: args,
			},
			op.Dsts,
		))
	case *ConstructOp:
		elts := make([]*ast.KeywordExpr, len(op.Args))
		for i, arg := range op.Args {
			elts[i] = &ast.KeywordExpr{
				Name: arg.Key,
				Expr: reg(arg.Value),
			}
		}
		block = append(block, opAssign(
			&ast.UnaryExpr{
				Op: "&",
				Expr: &ast.StructLiteral{
					Type: goStructType(op.Type, ctx),
					Args: elts,
				},
			},
			op.Dst,
		))
	case *ConstructListOp:
		elts := make([]ast.Expr, len(op.Args))
		for i, arg := range op.Args {
			elts[i] = reg(arg)
		}
		block = append(block, opAssign(
			&ast.ListLiteral{
				// TODO unhack
				Type: goTypeName(op.Type, ctx).(*ast.SliceType),
				Args: elts,
			},
			op.Dst,
		))
	case *CoerceOp:
		block = append(block, opAssign(
			&ast.TypeCoerce{
				Type: goTypeName(op.T, ctx),
				Expr: reg(op.Src),
			},
			op.Dst,
		))
	case *ConstantNilOp:
		block = append(block, opAssign(
			&ast.NilLiteral{},
			op.Dst,
		))
	case *ConstantBoolOp:
		block = append(block, opAssign(
			&ast.BoolLiteral{Value: op.Value},
			op.Dst,
		))
	case *ConstantIntOp:
		block = append(block, opAssign(
			// TODO unhack
			&ast.IntLiteral{Value: int(op.Value)},
			op.Dst,
		))
	case *ConstantRuneOp:
		block = append(block, opAssign(
			&ast.RuneLiteral{Value: op.Value},
			op.Dst,
		))
	case *ConstantStringOp:
		block = append(block, opAssign(
			&ast.StringLiteral{Value: op.Value},
			op.Dst,
		))
	case *Peek:
		block = append(block, opAssign(
			builtinExpr("Peek"),
			op.Dst,
		))
	case *Consume:
		block = append(block,
			builtinStmt("Consume"),
		)
	case *AppendOp:
		block = append(block, opAssign(
			&ast.Call{
				Expr: id("append"),
				Args: []ast.Expr{
					reg(op.List),
					reg(op.Value),
				},
			},
			op.Dst,
		))
	case *ReturnOp:
		if len(op.Exprs) != len(f.ReturnTypes) {
			panic(fmt.Sprintf("Wrong number of return values.  Expected %d, got %d.", len(f.ReturnTypes), len(op.Exprs)))
		}
		for i, e := range op.Exprs {
			block = append(block, &ast.Assign{
				Targets: []ast.Expr{id(returnVarName(i))},
				Op:      "=",
				Sources: []ast.Expr{reg(e)},
			})
		}
	case *Fail:
		block = append(block, builtinStmt("Fail"))
	case *Checkpoint:
		block = append(block, opAssign(
			builtinExpr("Checkpoint"),
			op.Dst,
		))
	case *Recover:
		block = append(block, builtinStmt("Recover", reg(op.Src)))
	case *LookaheadBegin:
		block = append(block, opAssign(
			builtinExpr("LookaheadBegin"),
			op.Dst,
		))
	case *LookaheadEnd:
		if op.Failed {
			block = append(block, builtinStmt("LookaheadFail", reg(op.Src)))
		} else {
			block = append(block, builtinStmt("LookaheadNormal", reg(op.Src)))
		}
	case *Slice:
		block = append(block, opAssign(
			builtinExpr("Slice", reg(op.Src)),
			op.Dst,
		))
	case *CopyOp:
		block = append(block, opAssign(
			reg(op.Src),
			op.Dst,
		))

	case *TransferOp:
		lhs := make([]ast.Expr, len(op.Dsts))
		for i, dst := range op.Dsts {
			lhs[i] = reg(dst)
		}
		rhs := make([]ast.Expr, len(op.Srcs))
		for i, src := range op.Srcs {
			rhs[i] = reg(src)
		}
		block = append(block, &ast.Assign{
			Targets: lhs,
			Op:      "=",
			Sources: rhs,
		})
	default:
		panic(op)
	}
	return block
}

func generateNode(info *regionInfo, node base.NodeID, ctx *DubToGoContext, block []ast.Stmt) []ast.Stmt {
	g := info.decl.CFG
	op := info.decl.Ops[node]
	switch data := op.(type) {
	case *EntryOp:
		block = gotoNode(info, g.GetExit(node, 0), ctx, block)
	case *FlowExitOp:
		block = append(block, &ast.Return{})
	case *ExitOp:
	case *SwitchOp:
		block = emitSwitch(info, reg(data.Cond), g.GetExit(node, 0), g.GetExit(node, 1), ctx, block)
	case DubOp:
		block = GenerateOp(info.decl, data, ctx, block)
		block = generateFlowSwitch(info, node, ctx, block)
	default:
		panic(data)
	}
	return block
}

type regionInfo struct {
	decl   *LLFunc
	labels map[base.NodeID]int
}

func gotoNode(info *regionInfo, n base.NodeID, ctx *DubToGoContext, block []ast.Stmt) []ast.Stmt {
	label, ok := info.labels[n]
	if ok {
		return append(block, &ast.Goto{Text: blockName(label)})
	} else {
		return generateNode(info, n, ctx, block)
	}
}

func emitSwitch(info *regionInfo, cond ast.Expr, t base.NodeID, f base.NodeID, ctx *DubToGoContext, block []ast.Stmt) []ast.Stmt {
	if t != base.NoNode {
		if f != base.NoNode {
			block = append(block, &ast.If{
				Cond: cond,
				Body: gotoNode(info, t, ctx, nil),
				Else: &ast.BlockStmt{
					Body: gotoNode(info, f, ctx, nil),
				},
			})
			return block
		} else {
			return gotoNode(info, t, ctx, block)
		}
	} else {
		return gotoNode(info, f, ctx, block)
	}
}

func generateFlowSwitch(info *regionInfo, node base.NodeID, ctx *DubToGoContext, block []ast.Stmt) []ast.Stmt {
	g := info.decl.CFG
	numExits := g.NumExits(node)

	if numExits == 2 {
		cond := &ast.BinaryExpr{
			Left:  attr(id("frame"), "Flow"),
			Op:    "==",
			Right: constInt(0),
		}
		t := g.GetExit(node, 0)
		f := g.GetExit(node, 1)
		return emitSwitch(info, cond, t, f, ctx, block)
	} else if numExits == 1 {
		return gotoNode(info, g.GetExit(node, 0), ctx, block)
	} else {
		panic(info.decl.Ops[node])
	}
}

func IsParam(f *LLFunc, r DubRegister) bool {
	for _, p := range f.Params {
		if p == r {
			return true
		}
	}
	return false
}

func GenerateGoFunc(f *LLFunc, ctx *DubToGoContext) ast.Decl {
	g := f.CFG
	order, _ := base.ReversePostorder(g)

	heads := []base.NodeID{}
	labels := map[base.NodeID]int{}
	uid := 0

	nit := base.OrderedIterator(order)
	for nit.Next() {
		n := nit.Value()
		if (n == g.Entry() || g.NumEntries(n) >= 2) && n != g.Exit() {
			heads = append(heads, n)
			labels[n] = uid
			uid = uid + 1
		}
	}
	info := &regionInfo{decl: f, labels: labels}

	stmts := []ast.Stmt{}

	// Declare the variables.
	// It is easier to do this up front than calculate where they need to be defined.
	for i, info := range f.Registers {
		r := DubRegister(i)
		if IsParam(f, r) {
			continue
		}
		stmts = append(stmts, &ast.Var{
			Name: RegisterName(r),
			Type: goTypeName(info.T, ctx),
		})
	}

	// Generate Go code from flow blocks
	for _, node := range heads {
		block := []ast.Stmt{}
		label, _ := info.labels[node]
		// HACK assume label 0 is always the entry node.
		if label != 0 {
			block = append(block, &ast.Label{Text: blockName(label)})
		}
		block = generateNode(info, node, ctx, block)
		// Extend the statement list
		stmts = append(stmts, block...)
	}

	results := []*ast.Param{}
	for i, t := range f.ReturnTypes {
		results = append(results, &ast.Param{
			Name: returnVarName(i),
			Type: goTypeName(t, ctx),
		})
	}

	ref := &ast.TypeRef{Impl: ctx.state}

	params := []*ast.Param{
		&ast.Param{
			Name: "frame",
			Type: &ast.PointerType{Element: ref},
		},
	}

	for _, p := range f.Params {
		params = append(params, &ast.Param{
			Name: RegisterName(p),
			Type: goTypeName(f.Registers[p].T, ctx),
		})
	}

	funcDecl := &ast.FuncDecl{
		Name: f.Name,
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
		Body: stmts,
	}

	return funcDecl
}

func tagName(s *LLStruct) string {
	return fmt.Sprintf("is%s", s.Name)
}

func addTags(base *LLStruct, parent *LLStruct, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if parent != nil {
		decls = addTags(base, parent.Implements, ctx, decls)
		decls = append(decls, &ast.FuncDecl{
			Name: tagName(parent),
			Recv: &ast.Param{
				Name: "node",
				Type: goTypeName(base, ctx),
			},
			Type: &ast.FuncType{},
			Body: []ast.Stmt{},
		})
	}
	return decls
}

func GenerateGoStruct(s *LLStruct, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if s.Abstract {
		if len(s.Fields) != 0 {
			panic(s.Name)
		}
		fields := []*ast.Field{
			&ast.Field{
				Name: tagName(s),
				Type: &ast.FuncType{},
			},
		}

		impl := &ast.InterfaceDecl{
			Name:   s.Name,
			Fields: fields,
		}
		ctx.link.ForwardType(s, impl)
		decls = append(decls, impl)
	} else {
		fields := []*ast.Field{}
		for _, f := range s.Fields {
			fields = append(fields, &ast.Field{
				Name: f.Name,
				Type: goTypeName(f.T, ctx),
			})
		}
		impl := &ast.StructDecl{
			Name:   s.Name,
			Fields: fields,
		}
		ctx.link.ForwardType(s, impl)
		decls = append(decls, impl)
		decls = addTags(s, s.Implements, ctx, decls)
	}
	return decls
}

func ExternParserRuntime() (*ast.Package, *ast.StructDecl) {
	state := &ast.StructDecl{
		Name: "State",
	}
	pkg := &ast.Package{
		Extern: true,
		Path:   []string{"evergreen", "dub", "runtime"},
		Files: []*ast.File{
			&ast.File{
				Decls: []ast.Decl{
					state,
				},
			},
		},
	}
	return pkg, state
}

type BuiltinIndex struct {
	IntType    ast.TypeImpl
	BoolType   ast.TypeImpl
	StringType ast.TypeImpl
	RuneType   ast.TypeImpl
}

func ExternBuiltinRuntime() (*ast.Package, *BuiltinIndex) {
	intType := &ast.ExternalType{Name: "int"}
	boolType := &ast.ExternalType{Name: "bool"}
	stringType := &ast.ExternalType{Name: "string"}
	runeType := &ast.ExternalType{Name: "rune"}

	pkg := &ast.Package{
		Extern: true,
		Path:   []string{},
		Files: []*ast.File{
			&ast.File{
				Decls: []ast.Decl{
					intType,
					boolType,
					stringType,
					runeType,
				},
			},
		},
	}
	index := &BuiltinIndex{
		IntType:    intType,
		BoolType:   boolType,
		StringType: stringType,
		RuneType:   runeType,
	}
	return pkg, index
}

func GenerateGo(module string, structs []*LLStruct, funcs []*LLFunc, index *BuiltinIndex, state *ast.StructDecl, link DubToGoLinker) *ast.File {
	ctx := &DubToGoContext{index: index, state: state, link: link}

	imports := []*ast.Import{}

	decls := []ast.Decl{}
	for _, f := range structs {
		decls = GenerateGoStruct(f, ctx, decls)
	}
	for _, f := range funcs {
		decls = append(decls, GenerateGoFunc(f, ctx))
	}

	file := &ast.File{
		Name:    "generated_parser.go",
		Package: "tree",
		Imports: imports,
		Decls:   decls,
	}
	return file
}

package flow

import (
	"evergreen/base"
	core "evergreen/dub/tree"
	ast "evergreen/go/tree"
	"fmt"
)

const (
	STRUCT = iota
	REF
	SCOPE
)

type DubToGoLinker interface {
	TypeRef(s *core.StructType, subtype int) *ast.TypeRef
	ForwardType(s *core.StructType, subtype int, impl ast.TypeImpl)
	Finish()
}

type linkElement struct {
	impl ast.TypeImpl
	refs []*ast.TypeRef
}

type linkerImpl struct {
	types []map[*core.StructType]*linkElement
}

func (l *linkerImpl) get(s *core.StructType, subtype int) *linkElement {
	e, ok := l.types[subtype][s]
	if !ok {
		e = &linkElement{}
		l.types[subtype][s] = e
	}
	return e
}

func (l *linkerImpl) TypeRef(s *core.StructType, subtype int) *ast.TypeRef {
	r := &ast.TypeRef{Name: subtypeName(s, subtype)}
	e := l.get(s, subtype)
	e.refs = append(e.refs, r)
	return r
}

func (l *linkerImpl) ForwardType(s *core.StructType, subtype int, impl ast.TypeImpl) {
	e := l.get(s, subtype)
	e.impl = impl
}

func (l *linkerImpl) Finish() {
	for subtype, types := range l.types {
		for s, e := range types {
			if e.impl == nil {
				panic(fmt.Sprintf("%s / %d", s.Name, subtype))
			}
			for _, r := range e.refs {
				r.Impl = e.impl
			}
			e.refs = nil
		}
	}
}

func MakeLinker() DubToGoLinker {
	types := []map[*core.StructType]*linkElement{}
	for i := 0; i < 3; i++ {
		types = append(types, map[*core.StructType]*linkElement{})
	}
	return &linkerImpl{
		types: types,
	}
}

func subtypeName(s *core.StructType, subtype int) string {
	name := s.Name.Text
	switch subtype {
	case STRUCT:
		// Nothing
	case REF:
		name += "_Ref"
	case SCOPE:
		name += "_Scope"
	default:
		panic(subtype)
	}
	return name
}

type regionInfo struct {
	decl       *LLFunc
	funcDecl   *ast.FuncDecl
	labels     map[base.NodeID]int
	frameInfo  int
	dubToGo    []int
	returnInfo []int
}

func (info *regionInfo) FrameVar() ast.Expr {
	return info.MakeGetLocal(info.frameInfo)
}

func (info *regionInfo) FrameAttr(name string) ast.Expr {
	return attr(info.FrameVar(), name)
}

func (info *regionInfo) GetLocalInfo(r RegisterInfo_Ref) int {
	return info.dubToGo[int(r)]
}

func (info *regionInfo) MakeParam(idx int) *ast.Param {
	return info.funcDecl.MakeParam(idx)
}

func (info *regionInfo) MakeGetLocal(idx int) ast.Expr {
	return info.funcDecl.MakeGetLocal(idx)
}

func (info *regionInfo) MakeSetLocal(idx int) ast.Target {
	return info.funcDecl.MakeSetLocal(idx)
}

func (info *regionInfo) MakeSetReturn(ret int) ast.Target {
	return info.MakeSetLocal(info.returnInfo[ret])
}

func (info *regionInfo) GetReg(r RegisterInfo_Ref) ast.Expr {
	idx := info.GetLocalInfo(r)
	return info.MakeGetLocal(idx)
}

func (info *regionInfo) SetReg(r RegisterInfo_Ref) ast.Target {
	idx := info.GetLocalInfo(r)
	return info.MakeSetLocal(idx)
}

func (info *regionInfo) Param(r RegisterInfo_Ref) *ast.Param {
	idx := info.GetLocalInfo(r)
	return info.MakeParam(idx)
}

// Begin AST construction wrappers

func discard() ast.Target {
	return &ast.SetDiscard{}
}

func globalRef(name string) ast.Expr {
	return &ast.GetGlobal{Text: name}
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
	index *ast.BuiltinTypeIndex
	state *ast.StructDecl
	graph *ast.StructDecl
	link  DubToGoLinker
}

type blockInfo struct {
	label string
}

func blockName(i int) string {
	return fmt.Sprintf("block%d", i)
}

func builtinStmt(info *regionInfo, name string, args ...ast.Expr) ast.Stmt {
	return &ast.Call{Expr: info.FrameAttr(name), Args: args}
}

func builtinExpr(info *regionInfo, name string, args ...ast.Expr) ast.Expr {
	return &ast.Call{Expr: info.FrameAttr(name), Args: args}
}

func constInt(v int64) ast.Expr {
	// Type converstion hack
	return intLiteral(int(v))
}

func returnVarName(i int) string {
	return fmt.Sprintf("ret%d", i)
}

func opAssign(info *regionInfo, expr ast.Expr, dst RegisterInfo_Ref) ast.Stmt {
	if dst != NoRegisterInfo {
		return &ast.Assign{
			Targets: []ast.Target{info.SetReg(dst)},
			Op:      "=",
			Sources: []ast.Expr{expr},
		}
	} else {
		// TODO fix expr / stmt duality.
		return expr.(ast.Stmt)
	}
}

func opMultiAssign(info *regionInfo, expr ast.Expr, dsts []RegisterInfo_Ref) ast.Stmt {
	if len(dsts) != 0 {
		lhs := make([]ast.Target, len(dsts))
		for i, dst := range dsts {
			if dst != NoRegisterInfo {
				lhs[i] = info.SetReg(dst)
			} else {
				lhs[i] = discard()
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

func goFieldType(t core.DubType, ctx *DubToGoContext) ast.Type {
	switch t := t.(type) {
	case *core.StructType:
		if t.Scoped {
			return ctx.link.TypeRef(t, REF)
		}
	case *core.ListType:
		return &ast.SliceType{Element: goFieldType(t.Type, ctx)}
	}
	return goTypeName(t, ctx)
}

func goTypeName(t core.DubType, ctx *DubToGoContext) ast.Type {
	switch t := t.(type) {
	case *core.BuiltinType:
		switch t.Name {
		case "bool":
			return &ast.TypeRef{Impl: ctx.index.Bool}
		case "int":
			return &ast.TypeRef{Impl: ctx.index.Int}
		case "uint32":
			return &ast.TypeRef{Impl: ctx.index.UInt32}
		case "int64":
			return &ast.TypeRef{Impl: ctx.index.Int64}
		case "rune":
			return &ast.TypeRef{Impl: ctx.index.Rune}
		case "string":
			return &ast.TypeRef{Impl: ctx.index.String}
		case "graph":
			return &ast.PointerType{Element: &ast.TypeRef{Impl: ctx.graph}}
		default:
			panic(t.Name)
		}
	case *core.ListType:
		return &ast.SliceType{Element: goTypeName(t.Type, ctx)}
	case *core.StructType:
		out := ctx.link.TypeRef(t, STRUCT)
		if t.IsParent {
			return out
		} else {
			return &ast.PointerType{Element: out}
		}
	default:
		panic(t)
	}
}

func GenerateOp(info *regionInfo, f *LLFunc, op DubOp, ctx *DubToGoContext, block []ast.Stmt) []ast.Stmt {
	if IsNop(op) {
		return block
	}

	switch op := op.(type) {
	case *BinaryOp:
		// TODO validate Op?
		block = append(block, opAssign(
			info,
			&ast.BinaryExpr{
				Left:  info.GetReg(op.Left),
				Op:    op.Op,
				Right: info.GetReg(op.Right),
			},
			op.Dst,
		))
	case *CallOp:
		args := []ast.Expr{
			info.FrameVar(),
		}
		for _, arg := range op.Args {
			args = append(args, info.GetReg(arg))
		}
		block = append(block, opMultiAssign(
			info,
			&ast.Call{
				Expr: globalRef(op.Name),
				Args: args,
			},
			op.Dsts,
		))
	case *ConstructOp:
		elts := make([]*ast.KeywordExpr, len(op.Args))
		for i, arg := range op.Args {
			elts[i] = &ast.KeywordExpr{
				Name: arg.Key,
				Expr: info.GetReg(arg.Value),
			}
		}
		for _, c := range op.Type.Contains {
			elts = append(elts, &ast.KeywordExpr{
				Name: subtypeName(c, SCOPE),
				Expr: &ast.UnaryExpr{
					Op: "&",
					Expr: &ast.StructLiteral{
						Type: ctx.link.TypeRef(c, SCOPE),
						Args: []*ast.KeywordExpr{},
					},
				},
			})

		}
		block = append(block, opAssign(
			info,
			&ast.UnaryExpr{
				Op: "&",
				Expr: &ast.StructLiteral{
					Type: ctx.link.TypeRef(op.Type, STRUCT),
					Args: elts,
				},
			},
			op.Dst,
		))
	case *ConstructListOp:
		elts := make([]ast.Expr, len(op.Args))
		for i, arg := range op.Args {
			elts[i] = info.GetReg(arg)
		}
		block = append(block, opAssign(
			info,
			&ast.ListLiteral{
				// TODO unhack
				Type: goTypeName(op.Type, ctx).(*ast.SliceType),
				Args: elts,
			},
			op.Dst,
		))
	case *CoerceOp:
		block = append(block, opAssign(
			info,
			&ast.TypeCoerce{
				Type: goTypeName(op.T, ctx),
				Expr: info.GetReg(op.Src),
			},
			op.Dst,
		))
	case *ConstantNilOp:
		block = append(block, opAssign(
			info,
			&ast.NilLiteral{},
			op.Dst,
		))
	case *ConstantBoolOp:
		block = append(block, opAssign(
			info,
			&ast.BoolLiteral{Value: op.Value},
			op.Dst,
		))
	case *ConstantIntOp:
		block = append(block, opAssign(
			info,
			// TODO unhack
			&ast.IntLiteral{Value: int(op.Value)},
			op.Dst,
		))
	case *ConstantRuneOp:
		block = append(block, opAssign(
			info,
			&ast.RuneLiteral{Value: op.Value},
			op.Dst,
		))
	case *ConstantStringOp:
		block = append(block, opAssign(
			info,
			&ast.StringLiteral{Value: op.Value},
			op.Dst,
		))
	case *Peek:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "Peek"),
			op.Dst,
		))
	case *Consume:
		block = append(block,
			builtinStmt(info, "Consume"),
		)
	case *AppendOp:
		block = append(block, opAssign(
			info,
			&ast.Call{
				Expr: globalRef("append"),
				Args: []ast.Expr{
					info.GetReg(op.List),
					info.GetReg(op.Value),
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
				Targets: []ast.Target{info.MakeSetReturn(i)},
				Op:      "=",
				Sources: []ast.Expr{info.GetReg(e)},
			})
		}
	case *Fail:
		block = append(block, builtinStmt(info, "Fail"))
	case *Checkpoint:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "Checkpoint"),
			op.Dst,
		))
	case *Recover:
		block = append(block, builtinStmt(info, "Recover", info.GetReg(op.Src)))
	case *LookaheadBegin:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "LookaheadBegin"),
			op.Dst,
		))
	case *LookaheadEnd:
		if op.Failed {
			block = append(block, builtinStmt(info, "LookaheadFail", info.GetReg(op.Src)))
		} else {
			block = append(block, builtinStmt(info, "LookaheadNormal", info.GetReg(op.Src)))
		}
	case *Slice:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "Slice", info.GetReg(op.Src)),
			op.Dst,
		))
	case *CopyOp:
		block = append(block, opAssign(
			info,
			info.GetReg(op.Src),
			op.Dst,
		))

	case *TransferOp:
		if len(op.Dsts) != len(op.Srcs) {
			panic(op)
		}
		lhs := []ast.Target{}
		rhs := []ast.Expr{}
		for i := 0; i < len(op.Dsts); i++ {
			dst := op.Dsts[i]
			src := op.Srcs[i]
			if dst != src {
				lhs = append(lhs, info.SetReg(dst))
				rhs = append(rhs, info.GetReg(src))
			}
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
		block = emitSwitch(info, info.GetReg(data.Cond), g.GetExit(node, 0), g.GetExit(node, 1), ctx, block)
	case DubOp:
		block = GenerateOp(info, info.decl, data, ctx, block)
		block = generateFlowSwitch(info, node, ctx, block)
	default:
		panic(data)
	}
	return block
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
			body := gotoNode(info, t, ctx, nil)
			if ast.NormalFlowMightExit(body) {
				return append(block, &ast.If{
					Cond: cond,
					Body: body,
					Else: &ast.BlockStmt{
						Body: gotoNode(info, f, ctx, nil),
					},
				})
			} else {
				block = append(block, &ast.If{
					Cond: cond,
					Body: body,
				})
				return gotoNode(info, f, ctx, block)
			}
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
			Left:  info.FrameAttr("Flow"),
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

func ParamIndex(f *LLFunc, r RegisterInfo_Ref) int {
	for i, p := range f.Params {
		if p == r {
			return i
		}
	}
	return -1
}

func IsParam(f *LLFunc, r RegisterInfo_Ref) bool {
	return ParamIndex(f, r) != -1
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

	funcDecl := &ast.FuncDecl{
		Name:            f.Name,
		LocalInfo_Scope: &ast.LocalInfo_Scope{},
	}

	// Make local infos
	frameType := &ast.PointerType{Element: &ast.TypeRef{Impl: ctx.state}}

	frameInfo := funcDecl.CreateLocalInfo("frame", frameType)

	numRegisters := f.RegisterInfo_Scope.Len()
	dubToGo := make([]int, numRegisters)
	for i := 0; i < numRegisters; i++ {
		ref := RegisterInfo_Ref(i)
		info := f.RegisterInfo_Scope.Get(ref)
		dubToGo[i] = funcDecl.CreateLocalInfo(
			RegisterName(ref),
			goTypeName(info.T, ctx),
		)
	}

	info := &regionInfo{
		decl:      f,
		funcDecl:  funcDecl,
		labels:    labels,
		frameInfo: frameInfo,
		dubToGo:   dubToGo,
	}

	// Map function parameters
	params := []*ast.Param{
		info.MakeParam(frameInfo),
	}
	for _, p := range f.Params {
		params = append(params, info.Param(p))
	}

	// Map function results
	resultInfo := make([]int, len(f.ReturnTypes))
	results := []*ast.Param{}
	for i, t := range f.ReturnTypes {
		idx := funcDecl.CreateLocalInfo(returnVarName(i), goTypeName(t, ctx))
		resultInfo[i] = idx
		results = append(results, info.MakeParam(idx))
	}
	info.returnInfo = resultInfo

	// Generate Go code from flow blocks
	stmts := []ast.Stmt{}
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

	funcDecl.Type = &ast.FuncType{
		Params:  params,
		Results: results,
	}
	funcDecl.Body = stmts

	return funcDecl
}

func tagName(s *core.StructType) string {
	return fmt.Sprintf("is%s", s.Name.Text)
}

func addTags(base *core.StructType, parent *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if parent != nil {
		decls = addTags(base, parent.Implements, ctx, decls)
		decl := &ast.FuncDecl{
			Name:            tagName(parent),
			Type:            &ast.FuncType{},
			Body:            []ast.Stmt{},
			LocalInfo_Scope: &ast.LocalInfo_Scope{},
		}
		recv := decl.CreateLocalInfo("node", goTypeName(base, ctx))
		decl.Recv = decl.MakeParam(recv)
		decls = append(decls, decl)
	}
	return decls
}

func GenerateScopeHelpers(s *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	ref := &ast.TypeDef{
		Name: subtypeName(s, REF),
		Type: &ast.TypeRef{Impl: ctx.index.UInt32},
	}
	ctx.link.ForwardType(s, REF, ref)

	noRef := &ast.VarDecl{
		Name: "No" + s.Name.Text,
		Type: &ast.TypeRef{Impl: ref},
		Expr: &ast.UnaryExpr{
			Op: "^",
			Expr: &ast.TypeCoerce{
				Type: &ast.TypeRef{Impl: ref},
				Expr: &ast.IntLiteral{Value: 0},
			},
		},
		Const: true,
	}

	scope := &ast.StructDecl{
		Name: subtypeName(s, SCOPE),
		Fields: []*ast.Field{
			&ast.Field{
				Name: "objects",
				Type: &ast.SliceType{Element: goTypeName(s, ctx)},
			},
		},
	}
	ctx.link.ForwardType(s, SCOPE, scope)

	decls = append(decls, ref, noRef, scope)
	return decls
}

func GenerateGoStruct(s *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if s.IsParent {
		if s.Scoped {
			panic(s.Name)
		}
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
			Name:   s.Name.Text,
			Fields: fields,
		}
		ctx.link.ForwardType(s, STRUCT, impl)
		decls = append(decls, impl)
	} else {
		if s.Scoped {
			decls = GenerateScopeHelpers(s, ctx, decls)
		}

		fields := []*ast.Field{}
		for _, f := range s.Fields {
			fields = append(fields, &ast.Field{
				Name: f.Name.Text,
				Type: goFieldType(f.Type, ctx),
			})
		}
		for _, c := range s.Contains {
			if !c.Scoped {
				panic(c)
			}
			fields = append(fields, &ast.Field{
				Name: subtypeName(c, SCOPE),
				Type: &ast.PointerType{Element: ctx.link.TypeRef(c, SCOPE)},
			})
		}

		impl := &ast.StructDecl{
			Name:   s.Name.Text,
			Fields: fields,
		}
		ctx.link.ForwardType(s, STRUCT, impl)
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

func ExternTestingPackage() (*ast.Package, *ast.StructDecl) {
	t := &ast.StructDecl{
		Name: "T",
	}
	pkg := &ast.Package{
		Extern: true,
		Path:   []string{"testing"},
		Files: []*ast.File{
			&ast.File{
				Decls: []ast.Decl{
					t,
				},
			},
		},
	}
	return pkg, t
}

func ExternGraph() (*ast.Package, *ast.StructDecl) {
	graph := &ast.StructDecl{
		Name: "Graph",
	}
	pkg := &ast.Package{
		Extern: true,
		Path:   []string{"evergreen", "base"},
		Files: []*ast.File{
			&ast.File{
				Decls: []ast.Decl{
					graph,
				},
			},
		},
	}
	return pkg, graph
}

func MakeBuiltinTypes() *ast.BuiltinTypeIndex {
	return &ast.BuiltinTypeIndex{
		Int:    &ast.ExternalType{Name: "int"},
		UInt32: &ast.ExternalType{Name: "uint32"},
		Int64:  &ast.ExternalType{Name: "int64"},
		Bool:   &ast.ExternalType{Name: "bool"},
		String: &ast.ExternalType{Name: "string"},
		Rune:   &ast.ExternalType{Name: "rune"},
	}
}

func GenerateGo(package_name string, structs []*core.StructType, funcs []*LLFunc, index *ast.BuiltinTypeIndex, state *ast.StructDecl, graph *ast.StructDecl, link DubToGoLinker) *ast.File {
	ctx := &DubToGoContext{
		index: index,
		state: state,
		graph: graph,
		link:  link,
	}

	imports := []*ast.Import{}

	decls := []ast.Decl{}
	for _, f := range structs {
		decls = GenerateGoStruct(f, ctx, decls)
	}
	for _, f := range funcs {
		decls = append(decls, GenerateGoFunc(f, ctx))
	}

	file := &ast.File{
		Name:    "generated_dub.go",
		Package: package_name,
		Imports: imports,
		Decls:   decls,
	}
	return file
}

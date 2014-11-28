package golang

import (
	"evergreen/base"
	"evergreen/dub/flow"
	core "evergreen/dub/tree"
	ast "evergreen/go/tree"
	"fmt"
)

type regionInfo struct {
	decl       *flow.LLFunc
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

func (info *regionInfo) GetLocalInfo(r flow.RegisterInfo_Ref) int {
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

func (info *regionInfo) GetReg(r flow.RegisterInfo_Ref) ast.Expr {
	idx := info.GetLocalInfo(r)
	return info.MakeGetLocal(idx)
}

func (info *regionInfo) SetReg(r flow.RegisterInfo_Ref) ast.Target {
	idx := info.GetLocalInfo(r)
	return info.MakeSetLocal(idx)
}

func (info *regionInfo) Param(r flow.RegisterInfo_Ref) *ast.Param {
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
	state *ast.StructType
	graph *ast.StructType
	t     *ast.StructType
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

func opAssign(info *regionInfo, expr ast.Expr, dst flow.RegisterInfo_Ref) ast.Stmt {
	if dst != flow.NoRegisterInfo {
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

func opMultiAssign(info *regionInfo, expr ast.Expr, dsts []flow.RegisterInfo_Ref) ast.Stmt {
	if len(dsts) != 0 {
		lhs := make([]ast.Target, len(dsts))
		for i, dst := range dsts {
			if dst != flow.NoRegisterInfo {
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

func GenerateOp(info *regionInfo, f *flow.LLFunc, op flow.DubOp, ctx *DubToGoContext, block []ast.Stmt) []ast.Stmt {
	if flow.IsNop(op) {
		return block
	}

	switch op := op.(type) {
	case *flow.BinaryOp:
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
	case *flow.CallOp:
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
	case *flow.ConstructOp:
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
	case *flow.ConstructListOp:
		elts := make([]ast.Expr, len(op.Args))
		for i, arg := range op.Args {
			elts[i] = info.GetReg(arg)
		}
		block = append(block, opAssign(
			info,
			&ast.ListLiteral{
				// TODO unhack
				Type: ast.RefForType(goType(op.Type, ctx)).(*ast.SliceRef),
				Args: elts,
			},
			op.Dst,
		))
	case *flow.CoerceOp:
		block = append(block, opAssign(
			info,
			&ast.TypeCoerce{
				Type: ast.RefForType(goType(op.T, ctx)),
				Expr: info.GetReg(op.Src),
			},
			op.Dst,
		))
	case *flow.ConstantNilOp:
		block = append(block, opAssign(
			info,
			&ast.NilLiteral{},
			op.Dst,
		))
	case *flow.ConstantBoolOp:
		block = append(block, opAssign(
			info,
			&ast.BoolLiteral{Value: op.Value},
			op.Dst,
		))
	case *flow.ConstantIntOp:
		block = append(block, opAssign(
			info,
			// TODO unhack
			&ast.IntLiteral{Value: int(op.Value)},
			op.Dst,
		))
	case *flow.ConstantRuneOp:
		block = append(block, opAssign(
			info,
			&ast.RuneLiteral{Value: op.Value},
			op.Dst,
		))
	case *flow.ConstantStringOp:
		block = append(block, opAssign(
			info,
			&ast.StringLiteral{Value: op.Value},
			op.Dst,
		))
	case *flow.Peek:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "Peek"),
			op.Dst,
		))
	case *flow.Consume:
		block = append(block,
			builtinStmt(info, "Consume"),
		)
	case *flow.AppendOp:
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
	case *flow.ReturnOp:
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
	case *flow.Fail:
		block = append(block, builtinStmt(info, "Fail"))
	case *flow.Checkpoint:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "Checkpoint"),
			op.Dst,
		))
	case *flow.Recover:
		block = append(block, builtinStmt(info, "Recover", info.GetReg(op.Src)))
	case *flow.LookaheadBegin:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "LookaheadBegin"),
			op.Dst,
		))
	case *flow.LookaheadEnd:
		if op.Failed {
			block = append(block, builtinStmt(info, "LookaheadFail", info.GetReg(op.Src)))
		} else {
			block = append(block, builtinStmt(info, "LookaheadNormal", info.GetReg(op.Src)))
		}
	case *flow.Slice:
		block = append(block, opAssign(
			info,
			builtinExpr(info, "Slice", info.GetReg(op.Src)),
			op.Dst,
		))
	case *flow.CopyOp:
		block = append(block, opAssign(
			info,
			info.GetReg(op.Src),
			op.Dst,
		))

	case *flow.TransferOp:
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
	case *flow.EntryOp:
		block = gotoNode(info, g.GetExit(node, 0), ctx, block)
	case *flow.FlowExitOp:
		block = append(block, &ast.Return{})
	case *flow.ExitOp:
	case *flow.SwitchOp:
		block = emitSwitch(info, info.GetReg(data.Cond), g.GetExit(node, 0), g.GetExit(node, 1), ctx, block)
	case flow.DubOp:
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

func ParamIndex(f *flow.LLFunc, r flow.RegisterInfo_Ref) int {
	for i, p := range f.Params {
		if p == r {
			return i
		}
	}
	return -1
}

func IsParam(f *flow.LLFunc, r flow.RegisterInfo_Ref) bool {
	return ParamIndex(f, r) != -1
}

func GenerateGoFunc(f *flow.LLFunc, ctx *DubToGoContext) ast.Decl {
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
	frameType := &ast.PointerRef{Element: &ast.NameRef{T: ctx.state}}

	frameInfo := funcDecl.CreateLocalInfo("frame", frameType)

	numRegisters := f.RegisterInfo_Scope.Len()
	dubToGo := make([]int, numRegisters)
	for i := 0; i < numRegisters; i++ {
		ref := flow.RegisterInfo_Ref(i)
		info := f.RegisterInfo_Scope.Get(ref)
		dubToGo[i] = funcDecl.CreateLocalInfo(
			flow.RegisterName(ref),
			ast.RefForType(goType(info.T, ctx)),
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
		idx := funcDecl.CreateLocalInfo(returnVarName(i), ast.RefForType(goType(t, ctx)))
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

	funcDecl.Type = &ast.FuncTypeRef{
		Params:  params,
		Results: results,
	}
	funcDecl.Body = stmts

	return funcDecl
}

func addTags(base *core.StructType, parent *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if parent != nil {
		decls = addTags(base, parent.Implements, ctx, decls)
		decl := &ast.FuncDecl{
			Name:            tagName(parent),
			Type:            &ast.FuncTypeRef{},
			Body:            []ast.Stmt{},
			LocalInfo_Scope: &ast.LocalInfo_Scope{},
		}
		recv := decl.CreateLocalInfo("node", ast.RefForType(goType(base, ctx)))
		decl.Recv = decl.MakeParam(recv)
		decls = append(decls, decl)
	}
	return decls
}

func DeclForType(t ast.GoType, ctx *DubToGoContext) ast.Decl {
	switch t := t.(type) {
	case *ast.TypeDefType:
		return &ast.TypeDefDecl{
			Name: t.Name,
			Type: ast.RefForType(t.Type),
			T:    t,
		}
	case *ast.StructType:
		fields := []*ast.FieldDecl{}
		for _, f := range t.Fields {
			fields = append(fields, &ast.FieldDecl{
				Name: f.Name,
				Type: ast.RefForType(f.Type),
			})
		}

		return &ast.StructDecl{
			Name:   t.Name,
			Fields: fields,
			T:      t,
		}
	case *ast.InterfaceType:
		fields := []*ast.FieldDecl{}
		for _, f := range t.Fields {
			fields = append(fields, &ast.FieldDecl{
				Name: f.Name,
				Type: ast.RefForType(f.Type),
			})
		}
		return &ast.InterfaceDecl{
			Name:   t.Name,
			Fields: fields,
			T:      t,
		}
	default:
		panic(t)
	}
}

func GenerateScopeHelpers(s *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	ref := DeclForType(ctx.link.GetType(s, REF), ctx)

	noRef := &ast.VarDecl{
		Name: "No" + s.Name.Text,
		Type: ctx.link.TypeRef(s, REF),
		Expr: &ast.UnaryExpr{
			Op: "^",
			Expr: &ast.TypeCoerce{
				Type: ctx.link.TypeRef(s, REF),
				Expr: &ast.IntLiteral{Value: 0},
			},
		},
		Const: true,
	}

	scope := DeclForType(ctx.link.GetType(s, SCOPE), ctx)

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
		decls = append(decls, DeclForType(ctx.link.GetType(s, STRUCT), ctx))
	} else {
		if s.Scoped {
			decls = GenerateScopeHelpers(s, ctx, decls)
		}
		decls = append(decls, DeclForType(ctx.link.GetType(s, STRUCT), ctx))
		decls = addTags(s, s.Implements, ctx, decls)
	}
	return decls
}

func externParserRuntime() (*ast.PackageAST, *ast.StructType) {
	stateT := &ast.StructType{
		Name: "State",
	}
	state := &ast.StructDecl{
		Name: "State",
		T:    stateT,
	}
	pkg := &ast.PackageAST{
		Extern: true,
		Path:   []string{"evergreen", "dub", "runtime"},
		Files: []*ast.FileAST{
			&ast.FileAST{
				Decls: []ast.Decl{
					state,
				},
			},
		},
	}
	return pkg, stateT
}

func externTestingPackage() (*ast.PackageAST, *ast.StructType) {
	tT := &ast.StructType{
		Name: "T",
	}
	t := &ast.StructDecl{
		Name: "T",
		T:    tT,
	}
	pkg := &ast.PackageAST{
		Extern: true,
		Path:   []string{"testing"},
		Files: []*ast.FileAST{
			&ast.FileAST{
				Decls: []ast.Decl{
					t,
				},
			},
		},
	}
	return pkg, tT
}

func externGraph() (*ast.PackageAST, *ast.StructType) {
	graphT := &ast.StructType{
		Name: "Graph",
	}
	graph := &ast.StructDecl{
		Name: "Graph",
		T:    graphT,
	}
	pkg := &ast.PackageAST{
		Extern: true,
		Path:   []string{"evergreen", "base"},
		Files: []*ast.FileAST{
			&ast.FileAST{
				Decls: []ast.Decl{
					graph,
				},
			},
		},
	}
	return pkg, graphT
}

func generateGoFile(package_name string, dubPkg *flow.DubPackage, ctx *DubToGoContext) *ast.FileAST {
	imports := []*ast.Import{}

	decls := []ast.Decl{}
	for _, t := range dubPkg.Structs {
		decls = GenerateGoStruct(t, ctx, decls)
	}
	for _, f := range dubPkg.Funcs {
		decls = append(decls, GenerateGoFunc(f, ctx))
	}

	file := &ast.FileAST{
		Name:    "generated_dub.go",
		Package: package_name,
		Imports: imports,
		Decls:   decls,
	}
	return file
}

func GenerateGo(program []*flow.DubPackage, root string, generate_tests bool) *ast.ProgramAST {
	link := makeLinker()

	packages := []*ast.PackageAST{}

	index := makeBuiltinTypes()

	pkg, state := externParserRuntime()
	packages = append(packages, pkg)

	pkg, graph := externGraph()
	packages = append(packages, pkg)

	pkg, t := externTestingPackage()
	packages = append(packages, pkg)

	ctx := &DubToGoContext{
		index: index,
		state: state,
		graph: graph,
		t:     t,
		link:  link,
	}

	createTypeMapping(program, ctx.link)
	createTypes(program, ctx)

	for _, dubPkg := range program {
		path := []string{root}
		path = append(path, dubPkg.Path...)
		leaf := path[len(path)-1]

		files := []*ast.FileAST{}
		files = append(files, generateGoFile(leaf, dubPkg, ctx))

		if generate_tests && len(dubPkg.Tests) != 0 {
			files = append(files, GenerateTests(leaf, dubPkg.Tests, ctx))
		}
		packages = append(packages, &ast.PackageAST{
			Path:  path,
			Files: files,
		})
	}

	return &ast.ProgramAST{
		Builtins: index,
		Packages: packages,
	}
}

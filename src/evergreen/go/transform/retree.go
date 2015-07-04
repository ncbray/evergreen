package transform

import (
	"evergreen/compiler"
	"evergreen/go/core"
	"evergreen/go/flow"
	"evergreen/go/tree"
	"evergreen/graph"
	"fmt"
)

func blockName(i int) string {
	return fmt.Sprintf("block%d", i)
}

func gotoBlock(i int) *tree.Goto {
	return &tree.Goto{Text: blockName(i)}
}

func blockLabel(i int) *tree.Label {
	return &tree.Label{Text: blockName(i)}
}

func findBlockHeads(g *graph.Graph, order []graph.NodeID) ([]graph.NodeID, map[graph.NodeID]int) {
	heads := []graph.NodeID{}
	labels := map[graph.NodeID]int{}
	uid := 0

	nit := graph.OrderedIterator(order)
	for nit.HasNext() {
		n := nit.GetNext()
		if (n == g.Entry() || g.HasMultipleEntries(n)) && n != g.Exit() {
			heads = append(heads, n)
			labels[n] = uid
			uid = uid + 1
		}
	}
	return heads, labels
}

func getLocal(lclMap []*tree.LocalInfo, reg *flow.Register) *tree.GetLocal {
	return &tree.GetLocal{
		Info: lclMap[reg.Index],
	}
}

func getLocalList(lclMap []*tree.LocalInfo, regs []*flow.Register) []tree.Expr {
	lcls := make([]tree.Expr, len(regs))
	for i, reg := range regs {
		lcls[i] = getLocal(lclMap, reg)
	}
	return lcls
}

func setLocal(lclMap []*tree.LocalInfo, reg *flow.Register) tree.Target {
	if reg != nil {
		return &tree.SetLocal{
			Info: lclMap[reg.Index],
		}
	} else {
		return &tree.SetDiscard{}
	}
}

func scalarAssign(expr tree.Expr, lclMap []*tree.LocalInfo, reg *flow.Register) tree.Stmt {
	if reg == nil {
		return expr
	} else {
		return &tree.Assign{
			Sources: []tree.Expr{expr},
			Op:      "=",
			Targets: []tree.Target{setLocal(lclMap, reg)},
		}
	}
}

func multiAssign(expr tree.Expr, lclMap []*tree.LocalInfo, regs []*flow.Register) tree.Stmt {
	if len(regs) == 0 {
		return expr
	} else {
		tgts := make([]tree.Target, len(regs))
		for i, tgt := range regs {
			tgts[i] = setLocal(lclMap, tgt)
		}
		return &tree.Assign{
			Sources: []tree.Expr{expr},
			Op:      "=",
			Targets: tgts,
		}
	}
}

func callableName(c core.Callable) string {
	switch c := c.(type) {
	case *core.Function:
		return c.Name
	case *core.IntrinsicFunction:
		return c.Name
	default:
		panic(c)
	}
}

func generateNode(coreProg *core.CoreProgram, decl *flow.FlowFunc, lclMap []*tree.LocalInfo, labels map[graph.NodeID]int, parent_label int, is_head bool, node graph.NodeID, block []tree.Stmt) ([]tree.Stmt, bool) {
	g := decl.CFG
	for {
		if !is_head {
			label, ok := labels[node]
			if ok {
				can_fallthrough := label == parent_label+1
				can_fallthrough = false // Disabled
				if can_fallthrough {
					return block, true
				} else {

					block = append(block, gotoBlock(label))
					return block, false
				}
			}
		}

		op := decl.Ops[node]
		switch op := op.(type) {
		case *flow.Entry:
			// TODO
		case *flow.Exit:
			block = append(block, &tree.Return{})
			return block, false
		case *flow.ConstantString:
			block = append(block, scalarAssign(&tree.StringLiteral{
				Value: op.Value,
			}, lclMap, op.Dst))
		case *flow.ConstantRune:
			block = append(block, scalarAssign(&tree.RuneLiteral{
				Value: op.Value,
			}, lclMap, op.Dst))
		case *flow.ConstantInt:
			block = append(block, scalarAssign(&tree.IntLiteral{
				Value: int(op.Value),
			}, lclMap, op.Dst))
		case *flow.ConstantFloat32:
			block = append(block, scalarAssign(&tree.Float32Literal{
				Value: op.Value,
			}, lclMap, op.Dst))
		case *flow.ConstantBool:
			block = append(block, scalarAssign(&tree.BoolLiteral{
				Value: op.Value,
			}, lclMap, op.Dst))
		case *flow.ConstantNil:
			block = append(block, scalarAssign(&tree.NilLiteral{}, lclMap, op.Dst))
		case *flow.Call:
			f := op.Target
			block = append(block, multiAssign(&tree.Call{
				Expr: &tree.GetFunction{Func: f},
				Args: getLocalList(lclMap, op.Args),
			}, lclMap, op.Dsts))
		case *flow.MethodCall:
			// TODO simple IR
			block = append(block, multiAssign(&tree.Call{
				Expr: &tree.Selector{
					Expr: getLocal(lclMap, op.Expr),
					Text: op.Name,
				},
				Args: getLocalList(lclMap, op.Args),
			}, lclMap, op.Dsts))
		case *flow.Transfer:
			srcs := []tree.Expr{}
			tgts := []tree.Target{}
			// SSA can cause registers to be transfered to themselves.  Filter this out.
			for i := 0; i < len(op.Srcs); i++ {
				src := op.Srcs[i]
				tgt := op.Dsts[i]
				if src != tgt {
					srcs = append(srcs, getLocal(lclMap, src))
					tgts = append(tgts, setLocal(lclMap, tgt))
				}
			}
			if len(srcs) > 0 {
				block = append(block, &tree.Assign{
					Sources: srcs,
					Op:      "=",
					Targets: tgts,
				})
			}
		case *flow.Coerce:
			block = append(block, scalarAssign(&tree.TypeCoerce{
				Type: tree.RefForType(op.Type),
				Expr: getLocal(lclMap, op.Src),
			}, lclMap, op.Dst))
		case *flow.Attr:
			block = append(block, scalarAssign(&tree.Selector{
				Expr: getLocal(lclMap, op.Expr),
				Text: op.Name,
			}, lclMap, op.Dst))
		case *flow.BinaryOp:
			block = append(block, scalarAssign(&tree.BinaryExpr{
				Left:  getLocal(lclMap, op.Left),
				Op:    op.Op,
				Right: getLocal(lclMap, op.Right),
			}, lclMap, op.Dst))
		case *flow.ConstructStruct:
			args := make([]*tree.KeywordExpr, len(op.Args))
			for i, arg := range op.Args {
				args[i] = &tree.KeywordExpr{
					Name: arg.Name,
					Expr: getLocal(lclMap, arg.Arg),
				}
			}
			var expr tree.Expr = &tree.StructLiteral{
				Type: tree.RefForType(op.Type),
				Args: args,
			}
			if op.AddrTaken {
				expr = &tree.UnaryExpr{
					Op:   "&",
					Expr: expr,
				}
			}
			block = append(block, scalarAssign(expr, lclMap, op.Dst))
		case *flow.ConstructSlice:
			ref := tree.RefForType(op.Type)
			sref, ok := ref.(*tree.SliceRef)
			if !ok {
				panic(op.Type)
			}
			block = append(block, scalarAssign(&tree.ListLiteral{
				Type: sref,
				Args: getLocalList(lclMap, op.Args),
			}, lclMap, op.Dst))
		case *flow.Nop:
			// TODO
		case *flow.Switch:
			var body []tree.Stmt
			bodyFall := false
			var elseBody []tree.Stmt
			elseFall := false

			eit := g.ExitIterator(node)
			for eit.HasNext() {
				e, next := eit.GetNext()
				flowID := decl.Edges[e]
				child, childFall := generateNode(coreProg, decl, lclMap, labels, parent_label, false, next, []tree.Stmt{})
				switch flowID {
				case flow.COND_TRUE:
					body, bodyFall = child, childFall
				case flow.COND_FALSE:
					elseBody, elseFall = child, childFall
				default:
					panic(flowID)
				}
			}
			var else_ *tree.Block = nil
			if len(elseBody) > 0 && bodyFall {
				else_ = &tree.Block{Body: elseBody}
			}
			block = append(block, &tree.If{
				Cond: getLocal(lclMap, op.Cond),
				T:    &tree.Block{Body: body},
				F:    else_,
			})
			if !bodyFall {
				block = append(block, elseBody...)
			}
			return block, bodyFall || elseFall
		case *flow.Return:
			block = append(block, &tree.Return{
				Args: getLocalList(lclMap, op.Args),
			})
			return block, false
		default:
			panic(op)
		}

		_, next := g.GetUniqueExit(node)
		if next == graph.NoNode {
			panic(decl.Ops[node])
		}
		node = next
		is_head = false
	}
}

func RetreeFunc(coreProg *core.CoreProgram, f *core.Function, decl *flow.FlowFunc) *tree.FuncDecl {
	//return NewRetreeFunc(coreProg, f, decl)

	funcDecl := &tree.FuncDecl{
		Name:            f.Name,
		LocalInfo_Scope: &tree.LocalInfo_Scope{},
		Block:           &tree.Block{Body: []tree.Stmt{}},
	}

	// Translate locals
	lclMap := makeLocalMap(decl, funcDecl)

	// Translate receiver
	info := decl.Recv
	if info != nil {
		mapped := lclMap[info.Index]
		funcDecl.Recv = &tree.Param{
			Name: mapped.Name,
			Type: tree.RefForType(info.T),
			Info: mapped,
		}
	}

	funcDecl.Type = makeFuncType(decl, lclMap)

	// Don't reconstruct empty functions.
	_, first := decl.CFG.GetUniqueExit(decl.CFG.Entry())
	if first != decl.CFG.Exit() {
		order, _ := graph.ReversePostorder(decl.CFG)
		heads, labels := findBlockHeads(decl.CFG, order)

		// Generate Go code from flow blocks
		stmts := []tree.Stmt{}
		for _, node := range heads {
			block := []tree.Stmt{}
			label, _ := labels[node]
			// HACK assume label 0 is always the entry node.
			if label != 0 {
				block = append(block, blockLabel(label))
			}
			block, _ = generateNode(coreProg, decl, lclMap, labels, label, true, node, block)
			// Extend the statement list
			stmts = append(stmts, block...)
		}

		funcDecl.Block.Body = stmts
	}
	return funcDecl
}

type TreeBypass struct {
	Tests []*tree.FileAST
}

func pathLeaf(path []string) string {
	return path[len(path)-1]
}

func getPackage(t core.GoType) *core.Package {
	switch t := t.(type) {
	case *core.StructType:
		return t.Package
	case *core.InterfaceType:
		return t.Package
	case *core.TypeDefType:
		return t.Package
	default:
		panic(t)
	}

}

func declForType(t core.GoType) tree.Decl {
	switch t := t.(type) {
	case *core.TypeDefType:
		return &tree.TypeDefDecl{
			Name: t.Name,
			Type: tree.RefForType(t.Type),
			T:    t,
		}
	case *core.StructType:
		fields := []*tree.FieldDecl{}
		for _, f := range t.Fields {
			fields = append(fields, &tree.FieldDecl{
				Name: f.Name,
				Type: tree.RefForType(f.Type),
			})
		}

		return &tree.StructDecl{
			Name:   t.Name,
			Fields: fields,
			T:      t,
		}
	case *core.InterfaceType:
		fields := []*tree.FieldDecl{}
		for _, f := range t.Fields {
			fields = append(fields, &tree.FieldDecl{
				Name: f.Name,
				Type: tree.RefForType(f.Type),
			})
		}
		return &tree.InterfaceDecl{
			Name:   t.Name,
			Fields: fields,
			T:      t,
		}
	default:
		panic(t)
	}
}

func generateGoFile(coreProg *core.CoreProgram, flowProg *flow.FlowProgram, types []core.GoType, funcs []*core.Function, file *tree.FileAST) {
	file.Name = "generated_dub.go"

	for _, t := range types {
		file.Decls = append(file.Decls, declForType(t))

		st, ok := t.(*core.StructType)
		if ok {
			for _, cf := range st.Methods {
				f := flowProg.FlowFunc_Scope.Get(flow.FlowFunc_Ref(cf.Index))
				file.Decls = append(file.Decls, RetreeFunc(coreProg, cf, f))
			}
		}
	}

	for _, cf := range funcs {
		f := flowProg.FlowFunc_Scope.Get(flow.FlowFunc_Ref(cf.Index))
		if f.Recv != nil {
			continue
		}
		file.Decls = append(file.Decls, RetreeFunc(coreProg, cf, f))
	}
}

func FlowToTree(status compiler.PassStatus, program *flow.FlowProgram, coreProg *core.CoreProgram, bypass *TreeBypass) *tree.ProgramAST {
	status.Begin()
	defer status.End()

	// Bucket types for each package.
	packageTypes := make([][]core.GoType, coreProg.Package_Scope.Len())
	for _, t := range program.Types {
		pIndex := getPackage(t).Index
		packageTypes[pIndex] = append(packageTypes[pIndex], t)
	}

	packageDecls := []*tree.PackageAST{}
	piter := coreProg.Package_Scope.Iter()
	for piter.Next() {
		p, pkg := piter.Value()
		if pkg.Extern {
			continue
		}
		leaf := pathLeaf(pkg.Path)
		fileAST := &tree.FileAST{
			Package: leaf,
			Imports: []*tree.Import{},
		}
		fileDecls := []*tree.FileAST{fileAST}

		generateGoFile(coreProg, program, packageTypes[p], pkg.Functions, fileAST)
		if bypass.Tests[p] != nil {
			fileDecls = append(fileDecls, bypass.Tests[p])
		}

		pkgAST := &tree.PackageAST{
			Files: fileDecls,
			P:     pkg,
		}
		packageDecls = append(packageDecls, pkgAST)
	}

	return &tree.ProgramAST{
		Builtins: program.Builtins,
		Packages: packageDecls,
	}
}

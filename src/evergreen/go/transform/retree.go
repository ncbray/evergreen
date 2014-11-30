package transform

import (
	"evergreen/base"
	"evergreen/go/flow"
	"evergreen/go/tree"
	//core "evergreen/go/tree"
	"fmt"
)

func blockName(i int) string {
	return fmt.Sprintf("block%d", i)
}

func GotoBlock(i int) *tree.Goto {
	return &tree.Goto{Text: blockName(i)}
}

func BlockLabel(i int) *tree.Label {
	return &tree.Label{Text: blockName(i)}
}

func FindBlockHeads(g *base.Graph, order []base.NodeID) ([]base.NodeID, map[base.NodeID]int) {
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
	return heads, labels
}

func getLocal(lclMap []tree.LocalInfo_Ref, reg flow.Register_Ref) *tree.GetLocal {
	return &tree.GetLocal{
		Info: lclMap[reg],
	}
}

func getLocalList(lclMap []tree.LocalInfo_Ref, regs []flow.Register_Ref) []tree.Expr {
	lcls := make([]tree.Expr, len(regs))
	for i, reg := range regs {
		lcls[i] = getLocal(lclMap, reg)
	}
	return lcls
}

func setLocal(lclMap []tree.LocalInfo_Ref, reg flow.Register_Ref) tree.Target {
	if reg != flow.NoRegister {
		return &tree.SetLocal{
			Info: lclMap[reg],
		}
	} else {
		return &tree.SetDiscard{}
	}
}

func scalarAssign(expr tree.Expr, lclMap []tree.LocalInfo_Ref, reg flow.Register_Ref) tree.Stmt {
	if reg == flow.NoRegister {
		return expr
	} else {
		return &tree.Assign{
			Sources: []tree.Expr{expr},
			Op:      "=",
			Targets: []tree.Target{setLocal(lclMap, reg)},
		}
	}
}

func multiAssign(expr tree.Expr, lclMap []tree.LocalInfo_Ref, regs []flow.Register_Ref) tree.Stmt {
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

func generateNode(decl *flow.LLFunc, lclMap []tree.LocalInfo_Ref, labels map[base.NodeID]int, parent_label int, is_head bool, node base.NodeID, block []tree.Stmt) ([]tree.Stmt, bool) {
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

					block = append(block, GotoBlock(label))
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
		case *flow.ConstantBool:
			block = append(block, scalarAssign(&tree.BoolLiteral{
				Value: op.Value,
			}, lclMap, op.Dst))
		case *flow.ConstantNil:
			block = append(block, scalarAssign(&tree.NilLiteral{}, lclMap, op.Dst))
		case *flow.Call:
			block = append(block, multiAssign(&tree.Call{
				Expr: &tree.GetGlobal{Text: op.Name}, // HACK
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
			body, bodyFall := generateNode(decl, lclMap, labels, parent_label, false, g.GetExit(node, 0), []tree.Stmt{})
			elseBody, elseFall := generateNode(decl, lclMap, labels, parent_label, false, g.GetExit(node, 1), []tree.Stmt{})
			var elseStmt tree.Stmt = nil
			if len(elseBody) > 0 && bodyFall {
				elseStmt = &tree.BlockStmt{Body: elseBody}
			}
			block = append(block, &tree.If{
				Cond: getLocal(lclMap, op.Cond),
				Body: body,
				Else: elseStmt,
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

		if g.NumExits(node) == 1 {
			node = g.GetExit(node, 0)
			is_head = false
		} else {
			panic(decl.Ops[node])
		}
	}
}

func RetreeFunc(decl *flow.LLFunc) *tree.FuncDecl {
	funcDecl := &tree.FuncDecl{
		Name:            decl.Name,
		LocalInfo_Scope: &tree.LocalInfo_Scope{},
	}

	// Translate locals
	numRegisters := decl.Register_Scope.Len()
	lclMap := make([]tree.LocalInfo_Ref, numRegisters)
	for i := 0; i < numRegisters; i++ {
		ref := flow.Register_Ref(i)
		info := decl.Register_Scope.Get(ref)
		name := info.Name
		if name == "" {
			name = flow.RegisterName(ref - 1) // HACK to match existing code generator
		}
		lclMap[i] = tree.LocalInfo_Ref(funcDecl.LocalInfo_Scope.Register(&tree.LocalInfo{
			Name: name,
			T:    tree.RefForType(info.T),
		}))
	}

	ft := &tree.FuncTypeRef{}

	// Translate parameters
	ft.Params = make([]*tree.Param, len(decl.Params))
	for i, ref := range decl.Params {
		info := decl.Register_Scope.Get(ref)
		index := lclMap[ref]
		mapped := funcDecl.LocalInfo_Scope.Get(index)
		ft.Params[i] = &tree.Param{
			Name: mapped.Name,
			Type: tree.RefForType(info.T),
			Info: index,
		}
	}

	// Translate returns
	ft.Results = make([]*tree.Param, len(decl.Results))
	for i, ref := range decl.Results {
		info := decl.Register_Scope.Get(ref)
		index := lclMap[ref]
		mapped := funcDecl.LocalInfo_Scope.Get(index)
		ft.Results[i] = &tree.Param{
			Name: mapped.Name,
			Type: tree.RefForType(info.T),
			Info: index,
		}
	}

	funcDecl.Type = ft

	order, _ := base.ReversePostorder(decl.CFG)
	heads, labels := FindBlockHeads(decl.CFG, order)

	// Generate Go code from flow blocks
	stmts := []tree.Stmt{}
	for _, node := range heads {
		block := []tree.Stmt{}
		label, _ := labels[node]
		// HACK assume label 0 is always the entry node.
		if label != 0 {
			block = append(block, BlockLabel(label))
		}
		block, _ = generateNode(decl, lclMap, labels, label, true, node, block)
		// Extend the statement list
		stmts = append(stmts, block...)
	}

	funcDecl.Body = stmts

	if false {
		tree.RetreeFunc(funcDecl)

		b, w := base.BufferedCodeWriter()
		tree.GenerateDecl(funcDecl, w)

		fmt.Println(b.String())
	}

	return funcDecl
}

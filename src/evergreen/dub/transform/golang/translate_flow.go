package golang

import (
	srccore "evergreen/dub/core"
	src "evergreen/dub/flow"
	dstcore "evergreen/go/core"
	dst "evergreen/go/flow"
	"evergreen/graph"
)

func simpleFlow(stitcher *graph.FlowStitcher, srcID graph.NodeID, dstID graph.NodeID) {
	stitcher.SetHead(srcID, dstID)
	stitcher.SetEdge(srcID, 0, dstID, 0)
}

func dubFlow(ctx *DubToGoContext, builder *dst.GoFlowBuilder, stitcher *graph.FlowStitcher, frameRef dst.Register_Ref, srcID graph.NodeID, dstID graph.NodeID) {
	stitcher.SetHead(srcID, dstID)

	if stitcher.NumExits(srcID) != 2 {
		panic(srcID)
	}

	normal := stitcher.GetExit(srcID, 0)
	fail := stitcher.GetExit(srcID, 1)

	if normal != graph.NoNode {
		if fail != graph.NoNode {
			flow := builder.MakeRegister("flow", ctx.index.Int)
			reference := builder.MakeRegister("normal", ctx.index.Int)
			cond := builder.MakeRegister("cond", ctx.index.Bool)

			attrID := builder.EmitOp(&dst.Attr{
				Expr: frameRef,
				Name: "Flow",
				Dst:  flow,
			}, 1)

			constID := builder.EmitOp(&dst.ConstantInt{
				Value: 0,
				Dst:   reference,
			}, 1)

			compareID := builder.EmitOp(&dst.BinaryOp{
				Left:  flow,
				Op:    "==",
				Right: reference,
				Dst:   cond,
			}, 1)

			switchID := builder.EmitOp(&dst.Switch{
				Cond: cond,
			}, 2)

			stitcher.Internal(dstID, 0, attrID)
			stitcher.Internal(attrID, 0, constID)
			stitcher.Internal(constID, 0, compareID)
			stitcher.Internal(compareID, 0, switchID)

			stitcher.SetEdge(srcID, 0, switchID, 0)
			stitcher.SetEdge(srcID, 1, switchID, 1)
		} else {
			stitcher.SetEdge(srcID, 0, dstID, 0)
		}
	} else if fail != graph.NoNode {
		stitcher.SetEdge(srcID, 1, dstID, 0)
	} else {
		// Dead end should not happen?
		panic(srcID)
	}
}

func dstReg(regMap []dst.Register_Ref, reg src.RegisterInfo_Ref) dst.Register_Ref {
	if reg == src.NoRegisterInfo {
		return dst.NoRegister
	}
	return regMap[reg]
}

func multiDstReg(regMap []dst.Register_Ref, reg src.RegisterInfo_Ref) []dst.Register_Ref {
	if reg == src.NoRegisterInfo {
		return nil
	}
	return []dst.Register_Ref{regMap[reg]}
}

func regList(regMap []dst.Register_Ref, args []src.RegisterInfo_Ref) []dst.Register_Ref {
	out := make([]dst.Register_Ref, len(args))
	for i, arg := range args {
		out[i] = regMap[arg]
	}
	return out
}

func translateFlow(srcF *src.LLFunc, ctx *DubToGoContext) (*dstcore.Function, *dst.FlowFunc) {
	goCoreFunc := &dstcore.Function{
		Name: srcF.Name,
	}

	goFlowFunc := &dst.FlowFunc{
		Recv:           dst.NoRegister,
		Register_Scope: &dst.Register_Scope{},
	}

	builder := dst.MakeGoFlowBuilder(goFlowFunc)

	frameReg := builder.MakeRegister("frame", &dstcore.PointerType{Element: ctx.state})

	// Remap registers
	num := srcF.RegisterInfo_Scope.Len()
	regMap := make([]dst.Register_Ref, num)
	for i := 0; i < num; i++ {
		r := srcF.RegisterInfo_Scope.Get(src.RegisterInfo_Ref(i))
		regMap[i] = builder.MakeRegister(r.Name, goType(r.T, ctx))
	}

	// Remap parameters
	goFlowFunc.Params = make([]dst.Register_Ref, len(srcF.Params)+1)
	goFlowFunc.Params[0] = frameReg
	for i, p := range srcF.Params {
		goFlowFunc.Params[i+1] = regMap[p]
	}

	// Create result registers
	goFlowFunc.Results = make([]dst.Register_Ref, len(srcF.ReturnTypes))
	for i, rt := range srcF.ReturnTypes {
		goFlowFunc.Results[i] = builder.MakeRegister("ret", goType(rt, ctx))
	}

	stitcher := graph.MakeFlowStitcher(srcF.CFG, goFlowFunc.CFG)

	order, _ := graph.ReversePostorder(srcF.CFG)
	nit := graph.OrderedIterator(order)
	for nit.Next() {
		srcID := nit.Value()
		op := srcF.Ops[srcID]

		switch op := op.(type) {
		case *src.EntryOp:
			// Entry already exists
			dstID := srcID
			stitcher.SetEdge(srcID, 0, dstID, 0)
		case *src.ExitOp:
			// Exit already exists
			dstID := srcID
			stitcher.SetHead(srcID, dstID)
		case *src.SwitchOp:
			dstID := builder.EmitOp(&dst.Switch{
				Cond: regMap[op.Cond],
			}, 2)
			stitcher.SetHead(srcID, dstID)
			stitcher.SetEdge(srcID, 0, dstID, 0)
			stitcher.SetEdge(srcID, 1, dstID, 1)
		case *src.FlowExitOp:
			// TODO is there anything that needs to be done?
			dstID := builder.EmitOp(&dst.Nop{}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.CallOp:
			args := []dst.Register_Ref{frameReg}
			args = append(args, regList(regMap, op.Args)...)
			dstID := builder.EmitOp(&dst.Call{
				// HACK assumes functions are defined in the same order.
				Target: dstcore.Function_Ref(op.Target),
				Args:   args,
				Dsts:   regList(regMap, op.Dsts),
			}, 1)
			dubFlow(ctx, builder, stitcher, frameReg, srcID, dstID)
		case *src.ConstructOp:
			t := ctx.link.GetType(op.Type, STRUCT)
			st, ok := t.(*dstcore.StructType)
			if !ok {
				panic(t)
			}
			args := make([]*dst.NamedArg, len(op.Args))
			for i, arg := range op.Args {
				args[i] = &dst.NamedArg{
					Name: arg.Key,
					Arg:  regMap[arg.Value],
				}
			}

			nodes := []graph.NodeID{}

			for _, c := range op.Type.Contains {
				scopeTA := ctx.link.GetType(c, SCOPE)
				scopeT, ok := scopeTA.(*dstcore.StructType)
				if !ok {
					panic(scopeTA)
				}
				scopeTP := &dstcore.PointerType{Element: scopeT}
				scope := builder.MakeRegister("scope", scopeTP)

				nodes = append(nodes, builder.EmitOp(&dst.ConstructStruct{
					Type:      scopeT,
					AddrTaken: true,
					Dst:       scope,
				}, 1))
				args = append(args, &dst.NamedArg{
					Name: subtypeName(c, SCOPE),
					Arg:  scope,
				})
			}

			nodes = append(nodes, builder.EmitOp(&dst.ConstructStruct{
				Type:      st,
				AddrTaken: true,
				Args:      args,
				Dst:       dstReg(regMap, op.Dst),
			}, 1))

			stitcher.SetHead(srcID, nodes[0])
			for i := 0; i < len(nodes)-1; i++ {
				stitcher.Internal(nodes[i], 0, nodes[i+1])
			}

			stitcher.SetEdge(srcID, 0, nodes[len(nodes)-1], 0)
		case *src.ConstructListOp:
			dstID := builder.EmitOp(&dst.ConstructSlice{
				Type: goSliceType(op.Type, ctx),
				Args: regList(regMap, op.Args),
				Dst:  dstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.TransferOp:
			dstID := builder.EmitOp(&dst.Transfer{
				Srcs: regList(regMap, op.Srcs),
				Dsts: regList(regMap, op.Dsts),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.ConstantRuneOp:
			dstID := builder.EmitOp(&dst.ConstantRune{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.ConstantStringOp:
			dstID := builder.EmitOp(&dst.ConstantString{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.ConstantIntOp:
			dstID := builder.EmitOp(&dst.ConstantInt{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.ConstantBoolOp:
			dstID := builder.EmitOp(&dst.ConstantBool{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.ConstantNilOp:
			dstID := builder.EmitOp(&dst.ConstantNil{
				Dst: dstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.BinaryOp:
			dstID := builder.EmitOp(&dst.BinaryOp{
				Left:  regMap[op.Left],
				Op:    op.Op,
				Right: regMap[op.Right],
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.Checkpoint:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Checkpoint",
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.Fail:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Fail",
			}, 1)
			dubFlow(ctx, builder, stitcher, frameReg, srcID, dstID)
		case *src.Recover:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Recover",
				Args: []dst.Register_Ref{regMap[op.Src]},
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.Peek:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Peek",
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			dubFlow(ctx, builder, stitcher, frameReg, srcID, dstID)
		case *src.Consume:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Consume",
			}, 1)
			dubFlow(ctx, builder, stitcher, frameReg, srcID, dstID)
		case *src.LookaheadBegin:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "LookaheadBegin",
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.LookaheadEnd:
			name := "LookaheadNormal"
			if op.Failed {
				name = "LookaheadFail"
			}
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: name,
				Args: []dst.Register_Ref{regMap[op.Src]},
			}, 1)
			dubFlow(ctx, builder, stitcher, frameReg, srcID, dstID)
		case *src.Slice:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Slice",
				Args: []dst.Register_Ref{regMap[op.Src]},
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			dubFlow(ctx, builder, stitcher, frameReg, srcID, dstID)
		case *src.CoerceOp:
			dstID := builder.EmitOp(&dst.Coerce{
				Src:  regMap[op.Src],
				Type: goType(op.T, ctx),
				Dst:  regMap[op.Dst],
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.AppendOp:
			dstID := builder.EmitOp(&dst.Append{
				Src:  regMap[op.List],
				Args: []dst.Register_Ref{regMap[op.Value]},
				Dst:  regMap[op.Dst],
			}, 1)
			simpleFlow(stitcher, srcID, dstID)
		case *src.ReturnOp:
			transferID := builder.EmitOp(&dst.Transfer{
				Srcs: regList(regMap, op.Exprs),
				Dsts: goFlowFunc.Results, // TODO copy?
			}, 1)

			if true {
				simpleFlow(stitcher, srcID, transferID)
			} else {

				returnID := builder.EmitOp(&dst.Return{}, 1)

				stitcher.SetHead(srcID, transferID)
				stitcher.Internal(transferID, 0, returnID)
				stitcher.SetEdge(srcID, 0, returnID, 0)
			}
		default:
			panic(op)
		}
	}
	return goCoreFunc, goFlowFunc
}

func createTagInternal(base *srccore.StructType, parent *srccore.StructType, goCoreProg *dstcore.CoreProgram, goFlowProg *dst.FlowProgram, p dstcore.Package_Ref, selfType dstcore.GoType) {
	if parent == nil {
		return
	}

	createTagInternal(base, parent.Implements, goCoreProg, goFlowProg, p, selfType)

	goCoreFunc := &dstcore.Function{
		Name:    "is" + parent.Name,
		Package: p,
	}

	goFlowFunc := &dst.FlowFunc{
		Recv: dst.NoRegister,
		CFG:  graph.CreateGraph(),
		Ops: []dst.GoOp{
			&dst.Entry{},
			&dst.Exit{},
		},
		Register_Scope: &dst.Register_Scope{},
	}

	goFlowFunc.Recv = goFlowFunc.Register_Scope.Register(&dst.Register{
		Name: "node",
		T: &dstcore.PointerType{
			Element: selfType,
		},
	})

	// Empty function.
	goFlowFunc.CFG.Connect(0, 0, 1)

	goCoreProg.Function_Scope.Register(goCoreFunc)
	goFlowProg.FlowFunc_Scope.Register(goFlowFunc)

	// TODO attach method to type.
}

// Fake functions for enforcing type relationships.
func createTags(dubCoreProg *srccore.CoreProgram, dubFlowProg *src.DubProgram, goCoreProg *dstcore.CoreProgram, goFlowProg *dst.FlowProgram, packages []dstcore.Package_Ref, ctx *DubToGoContext) {
	for _, s := range dubCoreProg.Structures {
		if s.IsParent || s.Implements == nil {
			continue
		}

		pIndex := dubCoreProg.File_Scope.Get(s.File).Package
		p := packages[pIndex]

		selfType := ctx.link.GetType(s, STRUCT)

		createTagInternal(s, s.Implements, goCoreProg, goFlowProg, p, selfType)
	}
}

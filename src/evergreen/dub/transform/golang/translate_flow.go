package golang

import (
	srccore "evergreen/dub/core"
	src "evergreen/dub/flow"
	dstcore "evergreen/go/core"
	dst "evergreen/go/flow"
	"evergreen/graph"
)

type flowMapper struct {
	ctx      *DubToGoContext
	stitcher *graph.EdgeStitcher
	builder  *dst.GoFlowBuilder
	original *graph.Graph
}

func (mapper *flowMapper) simpleExitFlow(srcID graph.NodeID, dstID graph.NodeID) {
	// Copy edges.
	srcG := mapper.original
	eit := srcG.ExitIterator(srcID)
	for eit.HasNext() {
		e, _ := eit.GetNext()
		flow := srcG.EdgeFlow(e)
		switch flow {
		case src.NORMAL:
			mapper.stitcher.MapEdge(e, mapper.builder.EmitEdge(dstID, 0))
		default:
			panic(flow)
		}
	}
}

func (mapper *flowMapper) simpleFlow(srcID graph.NodeID, dstID graph.NodeID) {
	mapper.stitcher.MapIncomingEdges(srcID, dstID)
	mapper.simpleExitFlow(srcID, dstID)
}

func (mapper *flowMapper) dubFlow(frameRef dst.Register_Ref, srcID graph.NodeID, dstID graph.NodeID) {
	ctx := mapper.ctx
	stitcher := mapper.stitcher
	builder := mapper.builder
	srcG := mapper.original

	stitcher.MapIncomingEdges(srcID, dstID)

	normal := graph.NoEdge
	fail := graph.NoEdge

	eit := srcG.ExitIterator(srcID)
	for eit.HasNext() {
		e, _ := eit.GetNext()
		flow := srcG.EdgeFlow(e)
		switch flow {
		case src.NORMAL:
			normal = e
		case src.FAIL:
			fail = e
		default:
			panic(flow)
		}
	}

	if normal != graph.NoEdge {
		if fail != graph.NoEdge {
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

			builder.EmitConnection(dstID, 0, attrID)
			builder.EmitConnection(attrID, 0, constID)
			builder.EmitConnection(constID, 0, compareID)
			builder.EmitConnection(compareID, 0, switchID)

			stitcher.MapEdge(normal, builder.EmitEdge(switchID, 0))
			stitcher.MapEdge(fail, builder.EmitEdge(switchID, 1))
		} else {
			stitcher.MapEdge(normal, builder.EmitEdge(dstID, 0))
		}
	} else if fail != graph.NoEdge {
		stitcher.MapEdge(fail, builder.EmitEdge(dstID, 0))
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
		Name:    srcF.Name,
		Package: dstcore.NoPackage,
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

	srcG := srcF.CFG
	stitcher := graph.MakeEdgeStitcher(srcG, goFlowFunc.CFG)
	mapper := &flowMapper{
		ctx:      ctx,
		builder:  builder,
		stitcher: stitcher,
		original: srcG,
	}

	order, _ := graph.ReversePostorder(srcG)
	nit := graph.OrderedIterator(order)
	for nit.HasNext() {
		srcID := nit.GetNext()
		op := srcF.Ops[srcID]

		switch op := op.(type) {
		case *src.EntryOp:
			// Entry already exists
			dstID := srcID
			mapper.simpleExitFlow(srcID, dstID)
		case *src.ExitOp:
			// Exit already exists
			dstID := srcID
			stitcher.MapIncomingEdges(srcID, dstID)
		case *src.SwitchOp:
			dstID := builder.EmitOp(&dst.Switch{
				Cond: regMap[op.Cond],
			}, 2)
			stitcher.MapIncomingEdges(srcID, dstID)

			// Copy edges.
			eit := srcG.ExitIterator(srcID)
			for eit.HasNext() {
				e, _ := eit.GetNext()
				flow := srcG.EdgeFlow(e)
				switch flow {
				case src.COND_TRUE:
					stitcher.MapEdge(e, builder.EmitEdge(dstID, 0))
				case src.COND_FALSE:
					stitcher.MapEdge(e, builder.EmitEdge(dstID, 1))
				default:
					panic(flow)
				}
			}
		case *src.FlowExitOp:
			// TODO is there anything that needs to be done?
			dstID := builder.EmitOp(&dst.Nop{}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.CallOp:
			args := []dst.Register_Ref{frameReg}
			args = append(args, regList(regMap, op.Args)...)
			dstID := builder.EmitOp(&dst.Call{
				// HACK assumes functions are defined in the same order.
				Target: dstcore.Function_Ref(op.Target),
				Args:   args,
				Dsts:   regList(regMap, op.Dsts),
			}, 1)
			mapper.dubFlow(frameReg, srcID, dstID)
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

			stitcher.MapIncomingEdges(srcID, nodes[0])
			for i := 0; i < len(nodes)-1; i++ {
				builder.EmitConnection(nodes[i], 0, nodes[i+1])
			}
			mapper.simpleExitFlow(srcID, nodes[len(nodes)-1])
		case *src.ConstructListOp:
			dstID := builder.EmitOp(&dst.ConstructSlice{
				Type: goSliceType(op.Type, ctx),
				Args: regList(regMap, op.Args),
				Dst:  dstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.TransferOp:
			dstID := builder.EmitOp(&dst.Transfer{
				Srcs: regList(regMap, op.Srcs),
				Dsts: regList(regMap, op.Dsts),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantRuneOp:
			dstID := builder.EmitOp(&dst.ConstantRune{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantStringOp:
			dstID := builder.EmitOp(&dst.ConstantString{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantIntOp:
			dstID := builder.EmitOp(&dst.ConstantInt{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantBoolOp:
			dstID := builder.EmitOp(&dst.ConstantBool{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantNilOp:
			dstID := builder.EmitOp(&dst.ConstantNil{
				Dst: dstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.BinaryOp:
			dstID := builder.EmitOp(&dst.BinaryOp{
				Left:  regMap[op.Left],
				Op:    op.Op,
				Right: regMap[op.Right],
				Dst:   dstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.Checkpoint:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Checkpoint",
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.Fail:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Fail",
			}, 1)
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.Recover:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Recover",
				Args: []dst.Register_Ref{regMap[op.Src]},
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.Peek:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Peek",
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.Consume:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Consume",
			}, 1)
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.LookaheadBegin:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "LookaheadBegin",
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			mapper.simpleFlow(srcID, dstID)
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
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.Slice:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Slice",
				Args: []dst.Register_Ref{regMap[op.Src]},
				Dsts: multiDstReg(regMap, op.Dst),
			}, 1)
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.CoerceOp:
			dstID := builder.EmitOp(&dst.Coerce{
				Src:  regMap[op.Src],
				Type: goType(op.T, ctx),
				Dst:  regMap[op.Dst],
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.AppendOp:
			dstID := builder.EmitOp(&dst.Append{
				Src:  regMap[op.List],
				Args: []dst.Register_Ref{regMap[op.Value]},
				Dst:  regMap[op.Dst],
			}, 1)
			mapper.simpleFlow(srcID, dstID)
		case *src.ReturnOp:
			transferID := builder.EmitOp(&dst.Transfer{
				Srcs: regList(regMap, op.Exprs),
				Dsts: goFlowFunc.Results, // TODO copy?
			}, 1)

			if true {
				mapper.simpleFlow(srcID, transferID)
			} else {

				returnID := builder.EmitOp(&dst.Return{}, 1)

				stitcher.MapIncomingEdges(srcID, transferID)
				builder.EmitConnection(transferID, 0, returnID)
				mapper.simpleExitFlow(srcID, returnID)
			}
		default:
			panic(op)
		}
	}
	return goCoreFunc, goFlowFunc
}

func createTagInternal(base *srccore.StructType, parent *srccore.StructType, goCoreProg *dstcore.CoreProgram, goFlowProg *dst.FlowProgram, p dstcore.Package_Ref, selfType *dstcore.StructType) {
	if parent == nil {
		return
	}

	createTagInternal(base, parent.Implements, goCoreProg, goFlowProg, p, selfType)

	goCoreFunc := &dstcore.Function{
		Name:    "is" + parent.Name,
		Package: dstcore.NoPackage,
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
	g := goFlowFunc.CFG
	g.ConnectEdgeExit(g.IndexedExitEdge(g.Entry(), 0), g.Exit())

	f := goCoreProg.Function_Scope.Register(goCoreFunc)
	goFlowProg.FlowFunc_Scope.Register(goFlowFunc)

	// Index.
	dstcore.InsertFunctionIntoPackage(goCoreProg, p, f)
	selfType.Methods = append(selfType.Methods, f)
}

// Fake functions for enforcing type relationships.
func createTags(dubCoreProg *srccore.CoreProgram, dubFlowProg *src.DubProgram, goCoreProg *dstcore.CoreProgram, goFlowProg *dst.FlowProgram, packages []dstcore.Package_Ref, ctx *DubToGoContext) {
	for _, s := range dubCoreProg.Structures {
		if s.IsParent || s.Implements == nil {
			continue
		}

		pIndex := dubCoreProg.File_Scope.Get(s.File).Package
		p := packages[pIndex]

		absSelfType := ctx.link.GetType(s, STRUCT)
		selfType, ok := absSelfType.(*dstcore.StructType)
		if !ok {
			panic(absSelfType)
		}

		createTagInternal(s, s.Implements, goCoreProg, goFlowProg, p, selfType)
	}
}

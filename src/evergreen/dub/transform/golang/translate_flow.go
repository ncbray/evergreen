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
	original *src.LLFunc
}

func (mapper *flowMapper) simpleExitFlow(srcID graph.NodeID, dstID graph.NodeID) {
	// Copy edges.
	srcG := mapper.original.CFG
	eit := srcG.ExitIterator(srcID)
	for eit.HasNext() {
		e, _ := eit.GetNext()
		flow := mapper.original.Edges[e]
		switch flow {
		case src.NORMAL:
			mapper.stitcher.MapEdge(e, mapper.builder.EmitEdge(dstID, dst.NORMAL))
		case src.RETURN:
			mapper.stitcher.MapEdge(e, mapper.builder.EmitEdge(dstID, dst.RETURN))
		default:
			panic(flow)
		}
	}
}

func (mapper *flowMapper) simpleFlow(srcID graph.NodeID, dstID graph.NodeID) {
	mapper.stitcher.MapIncomingEdges(srcID, dstID)
	mapper.simpleExitFlow(srcID, dstID)
}

func (mapper *flowMapper) handleFailEdge(original graph.EdgeID, translated graph.EdgeID, doesExit bool) {
	if doesExit {
		// This fail edge exits the function.
		// Before translation, failiure edges are non-local flow, but after translation they are
		// plain-old normal flow. Cap the edge with a return to exit the function with non-local
		// flow.
		builder := mapper.builder
		returnNode := builder.EmitOp(&dst.Return{})
		returnEdge := builder.EmitEdge(returnNode, dst.RETURN)
		mapper.builder.ConnectEdgeExit(translated, returnNode)
		mapper.stitcher.MapEdge(original, returnEdge)

	} else {
		mapper.stitcher.MapEdge(original, translated)
	}
}

func (mapper *flowMapper) dubFlow(frameRef *dst.Register, srcID graph.NodeID, dstID graph.NodeID) {
	ctx := mapper.ctx
	stitcher := mapper.stitcher
	builder := mapper.builder
	srcG := mapper.original.CFG

	stitcher.MapIncomingEdges(srcID, dstID)

	normal := graph.NoEdge
	fail := graph.NoEdge
	failExits := false

	eit := srcG.ExitIterator(srcID)
	for eit.HasNext() {
		e, dstID := eit.GetNext()
		flow := mapper.original.Edges[e]
		switch flow {
		case src.NORMAL:
			normal = e
		case src.FAIL:
			fail = e
			_, failExits = mapper.original.Ops[dstID].(*src.ExitOp)
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
			})

			constID := builder.EmitOp(&dst.ConstantInt{
				Value: 0,
				Dst:   reference,
			})

			compareID := builder.EmitOp(&dst.BinaryOp{
				Left:  flow,
				Op:    "==",
				Right: reference,
				Dst:   cond,
			})

			switchID := builder.EmitOp(&dst.Switch{
				Cond: cond,
			})

			builder.EmitConnection(dstID, dst.NORMAL, attrID)
			builder.EmitConnection(attrID, dst.NORMAL, constID)
			builder.EmitConnection(constID, dst.NORMAL, compareID)
			builder.EmitConnection(compareID, dst.NORMAL, switchID)

			stitcher.MapEdge(normal, builder.EmitEdge(switchID, dst.COND_TRUE))
			mapper.handleFailEdge(fail, builder.EmitEdge(switchID, dst.COND_FALSE), failExits)
		} else {
			stitcher.MapEdge(normal, builder.EmitEdge(dstID, dst.NORMAL))
		}
	} else if fail != graph.NoEdge {
		mapper.handleFailEdge(fail, builder.EmitEdge(dstID, dst.NORMAL), failExits)
	} else {
		// Dead end should not happen?
		panic(srcID)
	}
}

func dstReg(regMap []*dst.Register, reg *src.RegisterInfo) *dst.Register {
	if reg == nil {
		return nil
	}
	return regMap[reg.Index]
}

func multiDstReg(regMap []*dst.Register, reg *src.RegisterInfo) []*dst.Register {
	if reg == nil {
		return nil
	}
	return []*dst.Register{regMap[reg.Index]}
}

func regList(regMap []*dst.Register, args []*src.RegisterInfo) []*dst.Register {
	out := make([]*dst.Register, len(args))
	for i, arg := range args {
		out[i] = regMap[arg.Index]
	}
	return out
}

func translateFlow(srcF *src.LLFunc, ctx *DubToGoContext) *dst.FlowFunc {
	goFlowFunc := &dst.FlowFunc{
		Recv:           nil,
		Function:       ctx.functionMap[srcF.F.Index],
		Register_Scope: &dst.Register_Scope{},
	}

	builder := dst.MakeGoFlowBuilder(goFlowFunc)

	frameReg := builder.MakeRegister("frame", &dstcore.PointerType{Element: ctx.state})

	// Remap registers
	num := srcF.RegisterInfo_Scope.Len()
	regMap := make([]*dst.Register, num)
	for i := 0; i < num; i++ {
		r := srcF.RegisterInfo_Scope.Get(src.RegisterInfo_Ref(i))
		regMap[i] = builder.MakeRegister(r.Name, goType(r.T, ctx))
	}

	// Remap parameters
	goFlowFunc.Params = make([]*dst.Register, len(srcF.Params)+1)
	goFlowFunc.Params[0] = frameReg
	for i, p := range srcF.Params {
		goFlowFunc.Params[i+1] = regMap[p.Index]
	}

	// Create result registers
	goFlowFunc.Results = make([]*dst.Register, len(srcF.ReturnTypes))
	for i, rt := range srcF.ReturnTypes {
		goFlowFunc.Results[i] = builder.MakeRegister("ret", goType(rt, ctx))
	}

	srcG := srcF.CFG
	dstG := goFlowFunc.CFG
	stitcher := graph.MakeEdgeStitcher(srcG, dstG)
	mapper := &flowMapper{
		ctx:      ctx,
		builder:  builder,
		stitcher: stitcher,
		original: srcF,
	}

	order, _ := graph.ReversePostorder(srcG)
	nit := graph.OrderedIterator(order)
	for nit.HasNext() {
		srcID := nit.GetNext()
		op := srcF.Ops[srcID]

		switch op := op.(type) {
		case *src.EntryOp:
			// Entry node already exists
			dstID := srcID
			mapper.simpleExitFlow(srcID, dstID)
		case *src.ExitOp:
			// Exit node already exists
			dstID := srcID
			stitcher.MapIncomingEdges(srcID, dstID)
		case *src.SwitchOp:
			dstID := builder.EmitOp(&dst.Switch{
				Cond: regMap[op.Cond.Index],
			})
			stitcher.MapIncomingEdges(srcID, dstID)

			// Copy edges.
			eit := srcG.ExitIterator(srcID)
			for eit.HasNext() {
				e, _ := eit.GetNext()
				flow := mapper.original.Edges[e]
				switch flow {
				case src.COND_TRUE:
					stitcher.MapEdge(e, builder.EmitEdge(dstID, dst.COND_TRUE))
				case src.COND_FALSE:
					stitcher.MapEdge(e, builder.EmitEdge(dstID, dst.COND_FALSE))
				default:
					panic(flow)
				}
			}
		case *src.CallOp:
			mappedArgs := regList(regMap, op.Args)
			mappedDsts := regList(regMap, op.Dsts)
			switch c := op.Target.(type) {
			case *srccore.Function:
				args := []*dst.Register{frameReg}
				args = append(args, mappedArgs...)
				dstID := builder.EmitOp(&dst.Call{
					Target: ctx.functionMap[c.Index],
					Args:   args,
					Dsts:   mappedDsts,
				})
				mapper.dubFlow(frameReg, srcID, dstID)
			case *srccore.IntrinsicFunction:
				if c.Parent == ctx.core.Builtins.Append {
					if len(mappedArgs) != 2 {
						panic(op)
					}
					if len(mappedDsts) > 1 {
						panic(op)
					}
					dstID := builder.EmitOp(&dst.Call{
						Target: ctx.index.Append,
						Args:   mappedArgs,
						Dsts:   mappedDsts,
					})
					mapper.dubFlow(frameReg, srcID, dstID)
				} else {

					switch c {
					case ctx.core.Builtins.Position:
						dstID := builder.EmitOp(&dst.MethodCall{
							Expr: frameReg,
							Name: "Checkpoint",
							Args: mappedArgs,
							Dsts: mappedDsts,
						})
						mapper.dubFlow(frameReg, srcID, dstID)

					case ctx.core.Builtins.Slice:
						dstID := builder.EmitOp(&dst.MethodCall{
							Expr: frameReg,
							Name: "Slice",
							Args: mappedArgs,
							Dsts: mappedDsts,
						})
						mapper.dubFlow(frameReg, srcID, dstID)

					default:
						panic(c)
					}
				}
			default:
				panic(op.Target)
			}
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
					Arg:  regMap[arg.Value.Index],
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
				}))
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
			}))

			stitcher.MapIncomingEdges(srcID, nodes[0])
			for i := 0; i < len(nodes)-1; i++ {
				builder.EmitConnection(nodes[i], dst.NORMAL, nodes[i+1])
			}
			mapper.simpleExitFlow(srcID, nodes[len(nodes)-1])
		case *src.ConstructListOp:
			dstID := builder.EmitOp(&dst.ConstructSlice{
				Type: goSliceType(op.Type, ctx),
				Args: regList(regMap, op.Args),
				Dst:  dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.TransferOp:
			dstID := builder.EmitOp(&dst.Transfer{
				Srcs: regList(regMap, op.Srcs),
				Dsts: regList(regMap, op.Dsts),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantRuneOp:
			dstID := builder.EmitOp(&dst.ConstantRune{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantStringOp:
			dstID := builder.EmitOp(&dst.ConstantString{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantIntOp:
			dstID := builder.EmitOp(&dst.ConstantInt{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantFloat32Op:
			dstID := builder.EmitOp(&dst.ConstantFloat32{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantBoolOp:
			dstID := builder.EmitOp(&dst.ConstantBool{
				Value: op.Value,
				Dst:   dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.ConstantNilOp:
			dstID := builder.EmitOp(&dst.ConstantNil{
				Dst: dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.BinaryOp:
			dstID := builder.EmitOp(&dst.BinaryOp{
				Left:  regMap[op.Left.Index],
				Op:    op.Op,
				Right: regMap[op.Right.Index],
				Dst:   dstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.Checkpoint:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Checkpoint",
				Dsts: multiDstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.Fail:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Fail",
			})
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.Recover:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Recover",
				Args: []*dst.Register{regMap[op.Src.Index]},
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.Peek:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Peek",
				Dsts: multiDstReg(regMap, op.Dst),
			})
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.Consume:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "Consume",
			})
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.LookaheadBegin:
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: "LookaheadBegin",
				Dsts: multiDstReg(regMap, op.Dst),
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.LookaheadEnd:
			name := "LookaheadNormal"
			if op.Failed {
				name = "LookaheadFail"
			}
			dstID := builder.EmitOp(&dst.MethodCall{
				Expr: frameReg,
				Name: name,
				Args: []*dst.Register{regMap[op.Src.Index]},
			})
			mapper.dubFlow(frameReg, srcID, dstID)
		case *src.CoerceOp:
			dstID := builder.EmitOp(&dst.Coerce{
				Src:  regMap[op.Src.Index],
				Type: goType(op.T, ctx),
				Dst:  regMap[op.Dst.Index],
			})
			mapper.simpleFlow(srcID, dstID)
		case *src.ReturnOp:
			transferID := builder.EmitOp(&dst.Transfer{
				Srcs: regList(regMap, op.Exprs),
				Dsts: goFlowFunc.Results, // TODO copy?
			})

			stitcher.MapIncomingEdges(srcID, transferID)
			returnID := builder.EmitOp(&dst.Return{})
			builder.EmitConnection(transferID, dst.NORMAL, returnID)
			mapper.simpleExitFlow(srcID, returnID)
		default:
			panic(op)
		}
	}
	return goFlowFunc
}

func createTagInternal(base *srccore.StructType, parent *srccore.StructType, goCoreProg *dstcore.CoreProgram, goFlowProg *dst.FlowProgram, p *dstcore.Package, selfType *dstcore.StructType) {
	if parent == nil {
		return
	}

	createTagInternal(base, parent.Implements, goCoreProg, goFlowProg, p, selfType)

	goCoreFunc := &dstcore.Function{
		Name:    "is" + parent.Name,
		Package: nil,
	}

	goFlowFunc := &dst.FlowFunc{
		Recv: nil,
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
	g.ConnectEdge(g.Entry(), dst.AllocEdge(goFlowFunc, 0), g.Exit())

	f := goCoreProg.Function_Scope.Register(goCoreFunc)
	goFlowProg.FlowFunc_Scope.Register(goFlowFunc)

	// Index.
	dstcore.InsertFunctionIntoPackage(goCoreProg, p, f)
	selfType.Methods = append(selfType.Methods, f)
}

// Fake functions for enforcing type relationships.
func createTags(dubCoreProg *srccore.CoreProgram, dubFlowProg *src.DubProgram, goCoreProg *dstcore.CoreProgram, goFlowProg *dst.FlowProgram, packages []*dstcore.Package, ctx *DubToGoContext) {
	for _, s := range dubCoreProg.Structures {
		if s.IsParent || s.Implements == nil {
			continue
		}

		p := packages[s.File.Package.Index]

		absSelfType := ctx.link.GetType(s, STRUCT)
		selfType, ok := absSelfType.(*dstcore.StructType)
		if !ok {
			panic(absSelfType)
		}

		createTagInternal(s, s.Implements, goCoreProg, goFlowProg, p, selfType)
	}
}

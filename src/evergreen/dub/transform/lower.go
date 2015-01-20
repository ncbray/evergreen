package transform

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"evergreen/dub/flow"
	"evergreen/dub/tree"
	"evergreen/graph"
)

type dubBuilder struct {
	index    *core.BuiltinTypeIndex
	decl     *tree.FuncDecl
	flow     *flow.LLFunc
	localMap []*flow.RegisterInfo
	graph    *graph.Graph
}

func (builder *dubBuilder) EmitOp(op flow.DubOp) graph.NodeID {
	return flow.AllocNode(builder.flow, op)
}

func (builder *dubBuilder) EmitEdge(nid graph.NodeID, flowID int) graph.EdgeID {
	e := flow.AllocEdge(builder.flow, flowID)
	builder.graph.ConnectEdgeEntry(nid, e)
	return e
}

func (builder *dubBuilder) CreateRegister(name string, t core.DubType) *flow.RegisterInfo {
	info := &flow.RegisterInfo{
		Name: name,
		T:    t,
	}
	return builder.flow.RegisterInfo_Scope.Register(info)
}

func (builder *dubBuilder) CreateCheckpointRegister() *flow.RegisterInfo {
	return builder.CreateRegister("checkpoint", builder.index.Int)
}

func (builder *dubBuilder) ZeroRegister(dst *flow.RegisterInfo) flow.DubOp {
	switch t := dst.T.(type) {
	case *core.StructType:
		return &flow.ConstantNilOp{Dst: dst}
	case *core.BuiltinType:
		// TODO switch on object identity
		switch t.Name {
		case "rune":
			return &flow.ConstantRuneOp{Value: 0, Dst: dst}
		case "string":
			return &flow.ConstantStringOp{Value: "", Dst: dst}
		case "int":
			return &flow.ConstantIntOp{Value: 0, Dst: dst}
		case "bool":
			return &flow.ConstantBoolOp{Value: false, Dst: dst}
		default:
			panic(t.Name)
		}
	default:
		panic(dst.T)
	}
}

func makeRuneSwitch(cond *flow.RegisterInfo, op string, value rune, builder *dubBuilder) (graph.NodeID, graph.NodeID) {
	vreg := builder.CreateRegister("other", builder.index.Rune)
	make_value := builder.EmitOp(&flow.ConstantRuneOp{Value: value, Dst: vreg})

	breg := builder.CreateRegister("cond", builder.index.Bool)
	compare := builder.EmitOp(
		&flow.BinaryOp{
			Left:  cond,
			Op:    op,
			Right: vreg,
			Dst:   breg,
		},
	)

	decide := builder.EmitOp(&flow.SwitchOp{Cond: breg})
	builder.graph.ConnectEdgeExit(builder.EmitEdge(make_value, flow.NORMAL), compare)
	builder.graph.ConnectEdgeExit(builder.EmitEdge(compare, flow.NORMAL), decide)
	return make_value, decide
}

// TODO should be COND_TRUE and COND_FALSE, but the flow builder
// assumes normal flow when splitting off an edge.
const CONTINUE_MATCHING = flow.NORMAL
const STOP_MATCHING = flow.FAIL

func lowerRuneMatch(match *tree.RuneRangeMatch, used bool, builder *dubBuilder, fb *graph.FlowBuilder) *flow.RegisterInfo {
	// Read
	var cond *flow.RegisterInfo
	if len(match.Filters) > 0 || used {
		cond = builder.CreateRegister("c", builder.index.Rune)
	}
	body := builder.EmitOp(&flow.Peek{Dst: cond})
	fb.AttachFlow(flow.NORMAL, body)
	filters := fb.SplitOffEdge(builder.EmitEdge(body, flow.NORMAL))

	for _, flt := range match.Filters {
		if flt.Min > flt.Max {
			panic(flt.Min)
		}
		if flt.Min != flt.Max {
			minEntry, minDecide := makeRuneSwitch(cond, ">=", flt.Min, builder)
			maxEntry, maxDecide := makeRuneSwitch(cond, "<=", flt.Max, builder)

			// Check only if we haven't found a match.
			filters.AttachFlow(CONTINUE_MATCHING, minEntry)

			// Match
			builder.graph.ConnectEdgeExit(builder.EmitEdge(minDecide, flow.COND_TRUE), maxEntry)
			filters.RegisterExit(builder.EmitEdge(maxDecide, flow.COND_TRUE), STOP_MATCHING)

			// No match
			filters.RegisterExit(builder.EmitEdge(minDecide, flow.COND_FALSE), CONTINUE_MATCHING)
			filters.RegisterExit(builder.EmitEdge(maxDecide, flow.COND_FALSE), CONTINUE_MATCHING)
		} else {
			entry, decide := makeRuneSwitch(cond, "==", flt.Min, builder)

			// Check only if we haven't found a match.
			filters.AttachFlow(CONTINUE_MATCHING, entry)

			// Match
			filters.RegisterExit(builder.EmitEdge(decide, flow.COND_TRUE), STOP_MATCHING)

			// No match
			filters.RegisterExit(builder.EmitEdge(decide, flow.COND_FALSE), CONTINUE_MATCHING)
		}
	}

	normalCase := STOP_MATCHING
	failCase := CONTINUE_MATCHING
	if match.Invert {
		normalCase, failCase = failCase, normalCase
	}

	// The rune matched, consume it.
	if filters.HasFlow(normalCase) {
		c := builder.EmitOp(&flow.Consume{})
		filters.AttachFlow(normalCase, c)
		fb.RegisterExit(builder.EmitEdge(c, flow.NORMAL), flow.NORMAL)
	}
	// Make the fail official.
	if filters.HasFlow(failCase) {
		f := builder.EmitOp(&flow.Fail{})
		filters.AttachFlow(failCase, f)
		fb.RegisterExit(builder.EmitEdge(f, flow.FAIL), flow.FAIL)
	}

	// The peek can fail.
	fb.RegisterExit(builder.EmitEdge(body, flow.FAIL), flow.FAIL)

	return cond
}

func lowerMatch(match tree.TextMatch, builder *dubBuilder, fb *graph.FlowBuilder) {
	switch match := match.(type) {
	case *tree.RuneRangeMatch:
		lowerRuneMatch(match, false, builder, fb)
	case *tree.StringLiteralMatch:
		// HACK desugar
		for _, c := range []rune(match.Value) {
			lowerRuneMatch(&tree.RuneRangeMatch{Filters: []*tree.RuneFilter{&tree.RuneFilter{Min: c, Max: c}}}, false, builder, fb)
		}
	case *tree.MatchSequence:
		for _, child := range match.Matches {
			lowerMatch(child, builder, fb)
		}
	case *tree.MatchChoice:
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		fb.AttachFlow(flow.NORMAL, head)

		for i, child := range match.Matches {
			block := fb.SplitOffEdge(builder.EmitEdge(head, flow.NORMAL))
			lowerMatch(child, builder, block)

			// Recover if not the last block.
			if i < len(match.Matches)-1 {
				newHead := builder.EmitOp(&flow.Recover{Src: checkpoint})
				block.AttachFlow(flow.FAIL, newHead)
				head = newHead
			}
			fb.AbsorbExits(block)
		}
	case *tree.MatchRepeat:
		// HACK unroll
		for i := 0; i < match.Min; i++ {
			lowerMatch(match.Match, builder, fb)
		}

		// Checkpoint
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		fb.AttachFlow(flow.NORMAL, head)
		child := fb.SplitOffEdge(builder.EmitEdge(head, flow.NORMAL))

		// Handle the body
		lowerMatch(match.Match, builder, child)

		// Normal flow iterates
		child.AttachFlow(flow.NORMAL, head)

		// Stop iterating on failure and recover
		recover := builder.EmitOp(&flow.Recover{Src: checkpoint})
		child.AttachFlow(flow.FAIL, recover)
		child.RegisterExit(builder.EmitEdge(recover, flow.NORMAL), flow.NORMAL)

		fb.AbsorbExits(child)
	case *tree.MatchLookahead:
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.LookaheadBegin{Dst: checkpoint})
		fb.AttachFlow(flow.NORMAL, head)
		child := fb.SplitOffEdge(builder.EmitEdge(head, flow.NORMAL))

		lowerMatch(match.Match, builder, child)

		normal := builder.EmitOp(&flow.LookaheadEnd{Failed: false, Src: checkpoint})
		fail := builder.EmitOp(&flow.LookaheadEnd{Failed: true, Src: checkpoint})

		if match.Invert {
			child.AttachFlow(flow.NORMAL, fail)
			child.AttachFlow(flow.FAIL, normal)
		} else {
			child.AttachFlow(flow.NORMAL, normal)
			child.AttachFlow(flow.FAIL, fail)
		}

		child.RegisterExit(builder.EmitEdge(normal, flow.NORMAL), flow.NORMAL)
		child.RegisterExit(builder.EmitEdge(fail, flow.FAIL), flow.FAIL)

		fb.AbsorbExits(child)
	default:
		panic(match)
	}
}

func getPossibleFlows(c core.Callable, builtins *core.BuiltinTypeIndex) (bool, bool) {
	switch c := c.(type) {
	case *core.IntrinsicFunction:
		if c.Parent == builtins.Append {
			return true, false
		}
		switch c {
		case builtins.Position, builtins.Slice:
			return true, false
		default:
			panic(c)
		}
	}
	return true, true
}

func lowerMultiValueExpr(expr tree.ASTExpr, builder *dubBuilder, used bool, fb *graph.FlowBuilder) []*flow.RegisterInfo {
	switch expr := expr.(type) {

	case *tree.Call:
		args := make([]*flow.RegisterInfo, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, builder, true, fb)
		}
		var dsts []*flow.RegisterInfo
		if used {
			switch t := expr.T.(type) {
			case *core.TupleType:
				dsts = make([]*flow.RegisterInfo, len(t.Types))
				for i, child := range t.Types {
					dsts[i] = builder.CreateRegister("", child)
				}
			default:
				dsts = []*flow.RegisterInfo{builder.CreateRegister("", t)}
			}
		}
		body := builder.EmitOp(&flow.CallOp{Target: expr.Target, Args: args, Dsts: dsts})
		fb.AttachFlow(flow.NORMAL, body)

		canNormal, canFail := getPossibleFlows(expr.Target, builder.index)

		if canNormal {
			fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		}
		if canFail {
			fb.RegisterExit(builder.EmitEdge(body, flow.FAIL), flow.FAIL)
		}
		return dsts
	default:
		return []*flow.RegisterInfo{lowerExpr(expr, builder, used, fb)}
	}
}

func lowerExpr(expr tree.ASTExpr, builder *dubBuilder, used bool, fb *graph.FlowBuilder) *flow.RegisterInfo {
	switch expr := expr.(type) {
	case *tree.If:
		cond := lowerExpr(expr.Expr, builder, true, fb)
		decide := builder.EmitOp(&flow.SwitchOp{Cond: cond})
		fb.AttachFlow(flow.NORMAL, decide)

		block := fb.SplitOffEdge(builder.EmitEdge(decide, flow.COND_TRUE))
		lowerBlock(expr.Block, builder, block)
		fb.AbsorbExits(block)

		block = fb.SplitOffEdge(builder.EmitEdge(decide, flow.COND_FALSE))
		lowerBlock(expr.Else, builder, block)
		fb.AbsorbExits(block)

		return nil

	case *tree.Repeat:
		// HACK unroll
		for i := 0; i < expr.Min; i++ {
			lowerBlock(expr.Block, builder, fb)
		}

		// Checkpoint at head of loop
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		fb.AttachFlow(flow.NORMAL, head)
		block := fb.SplitOffEdge(builder.EmitEdge(head, flow.NORMAL))

		// Handle the body
		lowerBlock(expr.Block, builder, block)

		// Normal flow iterates
		block.AttachFlow(flow.NORMAL, head)

		// Stop iterating on failure and recover
		recover := builder.EmitOp(&flow.Recover{Src: checkpoint})
		block.AttachFlow(flow.FAIL, recover)
		block.RegisterExit(builder.EmitEdge(recover, flow.NORMAL), flow.NORMAL)

		fb.AbsorbExits(block)
		return nil
	case *tree.Choice:
		var checkpoint *flow.RegisterInfo
		if len(expr.Blocks) > 1 {
			checkpoint = builder.CreateCheckpointRegister()
		}
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		fb.AttachFlow(flow.NORMAL, head)

		for i, b := range expr.Blocks {
			block := fb.SplitOffEdge(builder.EmitEdge(head, flow.NORMAL))
			lowerBlock(b, builder, block)

			// Recover if not the last block.
			if i < len(expr.Blocks)-1 {
				newHead := builder.EmitOp(&flow.Recover{Src: checkpoint})
				block.AttachFlow(flow.FAIL, newHead)
				head = newHead
			}
			fb.AbsorbExits(block)
		}
		return nil

	case *tree.Optional:
		// Checkpoint
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		fb.AttachFlow(flow.NORMAL, head)
		block := fb.SplitOffEdge(builder.EmitEdge(head, flow.NORMAL))

		lowerBlock(expr.Block, builder, block)

		if block.HasFlow(flow.FAIL) {
			restore := builder.EmitOp(&flow.Recover{Src: checkpoint})
			block.AttachFlow(flow.FAIL, restore)
			block.RegisterExit(builder.EmitEdge(restore, flow.NORMAL), flow.NORMAL)
		}

		fb.AbsorbExits(block)
		return nil

	case *tree.GetLocal:
		if !used {
			return nil
		}
		dst := builder.CreateRegister(expr.Info.Name, expr.Info.T)
		body := builder.EmitOp(&flow.CopyOp{Src: builder.localMap[expr.Info.Index], Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Assign:
		var srcs []*flow.RegisterInfo
		if expr.Expr != nil {
			srcs = lowerMultiValueExpr(expr.Expr, builder, true, fb)
			if len(expr.Targets) != len(srcs) {
				panic(expr.Targets)
			}
		}
		dsts := make([]*flow.RegisterInfo, len(expr.Targets))
		for i, etgt := range expr.Targets {
			tgt, ok := etgt.(*tree.SetLocal)
			if !ok {
				panic(expr.Targets)
			}
			if tree.IsDiscard(tgt.Info.Name) {
				continue
			}
			dst := builder.localMap[tgt.Info.Index]
			dsts[i] = dst
			var op flow.DubOp
			if srcs != nil {
				src := srcs[i]
				if src == nil {
					panic(expr)
				}
				op = &flow.CopyOp{Src: src, Dst: dst}
				if src.Name == "" {
					src.Name = dst.Name
				}
			} else {
				op = builder.ZeroRegister(dst)
			}
			body := builder.EmitOp(op)
			fb.AttachFlow(flow.NORMAL, body)
			fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		}
		// HACK should actually return a multivalue
		if len(dsts) == 1 {
			return dsts[0]
		}
		return nil

	case *tree.RuneLiteral:
		if !used {
			return nil
		}
		dst := builder.CreateRegister("c_r", builder.index.Rune)
		body := builder.EmitOp(&flow.ConstantRuneOp{Value: expr.Value, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.StringLiteral:
		if !used {
			return nil
		}
		dst := builder.CreateRegister("c_s", builder.index.String)
		body := builder.EmitOp(&flow.ConstantStringOp{Value: expr.Value, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.IntLiteral:
		if !used {
			return nil
		}
		dst := builder.CreateRegister("c_i", builder.index.Int)
		body := builder.EmitOp(&flow.ConstantIntOp{Value: int64(expr.Value), Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Float32Literal:
		if !used {
			return nil
		}
		dst := builder.CreateRegister("c_f", builder.index.Float32)
		body := builder.EmitOp(&flow.ConstantFloat32Op{Value: expr.Value, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.BoolLiteral:
		if !used {
			return nil
		}
		dst := builder.CreateRegister("c_b", builder.index.Bool)
		body := builder.EmitOp(&flow.ConstantBoolOp{Value: expr.Value, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Return:
		exprs := make([]*flow.RegisterInfo, len(expr.Exprs))
		for i, e := range expr.Exprs {
			exprs[i] = lowerExpr(e, builder, true, fb)
		}
		body := builder.EmitOp(&flow.ReturnOp{Exprs: exprs})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.RETURN), flow.RETURN)
		return nil

	case *tree.Fail:
		body := builder.EmitOp(&flow.Fail{})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.FAIL), flow.FAIL)
		return nil

	case *tree.BinaryOp:
		left := lowerExpr(expr.Left, builder, true, fb)
		right := lowerExpr(expr.Right, builder, true, fb)
		var dst *flow.RegisterInfo
		if used {
			dst = builder.CreateRegister("", expr.T)
		}
		body := builder.EmitOp(&flow.BinaryOp{Left: left, Op: expr.Op, Right: right, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Call:
		dsts := lowerMultiValueExpr(expr, builder, true, fb)
		var dst *flow.RegisterInfo
		if used {
			if len(dsts) != 1 {
				panic(expr)
			} else {
				dst = dsts[0]
			}
		}
		return dst

	case *tree.Construct:
		args := make([]*flow.KeyValue, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = &flow.KeyValue{
				Key:   arg.Name.Text,
				Value: lowerExpr(arg.Expr, builder, true, fb),
			}
		}
		t := tree.ResolveType(expr.Type)
		s, ok := t.(*core.StructType)
		if !ok {
			panic(t)
		}
		var dst *flow.RegisterInfo
		if used {
			dst = builder.CreateRegister("", t)
		}
		body := builder.EmitOp(&flow.ConstructOp{Type: s, Args: args, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.ConstructList:
		args := make([]*flow.RegisterInfo, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, builder, true, fb)
		}
		t := tree.ResolveType(expr.Type)
		l, ok := t.(*core.ListType)
		if !ok {
			panic(t)
		}
		var dst *flow.RegisterInfo
		if used {
			dst = builder.CreateRegister("", t)
		}
		body := builder.EmitOp(&flow.ConstructListOp{Type: l, Args: args, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Coerce:
		t := tree.ResolveType(expr.Type)
		src := lowerExpr(expr.Expr, builder, true, fb)
		var dst *flow.RegisterInfo
		if used {
			dst = builder.CreateRegister("", t)
		}
		body := builder.EmitOp(&flow.CoerceOp{Src: src, T: t, Dst: dst})
		fb.AttachFlow(flow.NORMAL, body)
		// TODO can coersion fail?
		fb.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.StringMatch:
		var begin *flow.RegisterInfo

		// Checkpoint
		if used {
			begin = builder.CreateRegister("begin", builder.index.Int)
			// HACK assume checkpoint is just the index
			check := builder.EmitOp(&flow.Checkpoint{Dst: begin})
			fb.AttachFlow(flow.NORMAL, check)
			fb.RegisterExit(builder.EmitEdge(check, flow.NORMAL), flow.NORMAL)
		}

		lowerMatch(expr.Match, builder, fb)

		// Create a slice
		if used {
			end := builder.CreateRegister("end", builder.index.Int)
			// HACK assume checkpoint is just the index
			check := builder.EmitOp(&flow.Checkpoint{Dst: end})
			fb.AttachFlow(flow.NORMAL, check)
			fb.RegisterExit(builder.EmitEdge(check, flow.NORMAL), flow.NORMAL)

			dst := builder.CreateRegister("slice", builder.index.String)
			slice := builder.EmitOp(&flow.CallOp{
				Target: builder.index.Slice,
				Args:   []*flow.RegisterInfo{begin, end},
				Dsts:   []*flow.RegisterInfo{dst},
			})
			fb.AttachFlow(flow.NORMAL, slice)
			fb.RegisterExit(builder.EmitEdge(slice, flow.NORMAL), flow.NORMAL)
			return dst
		} else {
			return nil
		}

	case *tree.RuneMatch:
		return lowerRuneMatch(expr.Match, used, builder, fb)
	default:
		panic(expr)
	}

}

func lowerBlock(block []tree.ASTExpr, builder *dubBuilder, fb *graph.FlowBuilder) {
	for _, expr := range block {
		lowerExpr(expr, builder, false, fb)
	}
}

func lowerAST(program *tree.Program, decl *tree.FuncDecl, funcMap []*flow.LLFunc) *flow.LLFunc {
	f := funcMap[decl.F.Index]

	g := graph.CreateGraph()
	ops := []flow.DubOp{
		&flow.EntryOp{},
		&flow.ExitOp{},
	}

	f.CFG = g
	f.Ops = ops

	entryEdge := flow.AllocEdge(f, 0)
	g.ConnectEdgeEntry(g.Entry(), entryEdge)

	builder := &dubBuilder{
		index: program.Builtins,
		decl:  decl,
		flow:  f,
		graph: g,
	}

	// Allocate register for locals
	numLocals := decl.LocalInfo_Scope.Len()
	builder.localMap = make([]*flow.RegisterInfo, numLocals)
	for i := 0; i < numLocals; i++ {
		lcl := decl.LocalInfo_Scope.Get(tree.LocalInfo_Ref(i))
		builder.localMap[i] = builder.CreateRegister(lcl.Name, lcl.T)
	}

	// Function parameters
	params := make([]*flow.RegisterInfo, len(decl.Params))
	for i, p := range decl.Params {
		params[i] = builder.localMap[p.Info.Index]
	}
	f.Params = params

	// Function returns
	types := make([]core.DubType, len(decl.ReturnTypes))
	for i, node := range decl.ReturnTypes {
		types[i] = tree.ResolveType(node)
	}
	f.ReturnTypes = types

	fb := graph.CreateFlowBuilder(g, entryEdge, flow.NUM_FLOWS)
	lowerBlock(decl.Block, builder, fb)

	// Attach flows to exit.
	if fb.HasFlow(flow.NORMAL) {
		ret := builder.EmitOp(&flow.ReturnOp{Exprs: []*flow.RegisterInfo{}})
		fb.AttachFlow(flow.NORMAL, ret)
		fb.RegisterExit(builder.EmitEdge(ret, flow.RETURN), flow.RETURN)
	}

	if fb.HasFlow(flow.RETURN) {
		fb.AttachFlow(flow.RETURN, g.Exit())
	}

	if fb.HasFlow(flow.FAIL) {
		fb.AttachFlow(flow.FAIL, g.Exit())
	}

	if fb.HasFlow(flow.EXCEPTION) {
		panic("exceptions not supported, yet?")
	}
	// TODO lint other flows.
	return f
}

func lowerPackage(program *tree.Program, pkg *tree.Package, funcMap []*flow.LLFunc) *flow.DubPackage {
	dubPkg := &flow.DubPackage{
		Path:    pkg.Path,
		Structs: []*core.StructType{},
		Funcs:   []*flow.LLFunc{},
		Tests:   []*tree.Test{},
	}

	// Lower to flow IR
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *tree.FuncDecl:
				f := lowerAST(program, decl, funcMap)
				dubPkg.Funcs = append(dubPkg.Funcs, f)
			case *tree.StructDecl:
				dubPkg.Structs = append(dubPkg.Structs, decl.T)
			default:
				panic(decl)
			}
		}
		dubPkg.Tests = append(dubPkg.Tests, file.Tests...)
	}

	return dubPkg
}

func ssiProgram(status compiler.PassStatus, program *flow.DubProgram) {
	status.Begin()
	defer status.End()

	for _, f := range program.LLFuncs {
		flow.SSI(f)
	}
}

func LowerProgram(status compiler.PassStatus, program *tree.Program, coreProg *core.CoreProgram) *flow.DubProgram {
	status.Begin()
	defer status.End()

	n := coreProg.Function_Scope.Len()
	dubFuncs := make([]*flow.LLFunc, n)
	iter := coreProg.Function_Scope.Iter()
	for iter.Next() {
		fref, f := iter.Value()
		df := &flow.LLFunc{
			Name:               f.Name,
			RegisterInfo_Scope: &flow.RegisterInfo_Scope{},
			F:                  f,
		}
		dubFuncs[fref] = df
	}

	dubPackages := []*flow.DubPackage{}
	for _, pkg := range program.Packages {
		dubPackages = append(dubPackages, lowerPackage(program, pkg, dubFuncs))
	}

	dubProg := &flow.DubProgram{
		Core:     coreProg,
		Packages: dubPackages,
		LLFuncs:  dubFuncs,
	}

	ssiProgram(status.Pass("ssi"), dubProg)

	return dubProg
}

package transform

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"evergreen/dub/flow"
	"evergreen/dub/tree"
	"evergreen/graph"
)

const numRegionExits = 4

type dubBuilder struct {
	index    *tree.BuiltinTypeIndex
	decl     *tree.FuncDecl
	flow     *flow.LLFunc
	localMap []flow.RegisterInfo_Ref
	graph    *graph.Graph
	ops      []flow.DubOp
}

func (builder *dubBuilder) EmitOp(op flow.DubOp) graph.NodeID {
	id := builder.graph.CreateNode(2)
	if int(id) != len(builder.ops) {
		panic(op)
	}
	builder.ops = append(builder.ops, op)
	return id
}

func (builder *dubBuilder) EmitEdge(nid graph.NodeID, flow int) graph.EdgeID {
	return builder.graph.IndexedExitEdge(nid, flow)
}

func (builder *dubBuilder) CreateRegister(name string, t core.DubType) flow.RegisterInfo_Ref {
	info := &flow.RegisterInfo{
		Name: name,
		T:    t,
	}
	return builder.flow.RegisterInfo_Scope.Register(info)
}

func (builder *dubBuilder) CreateCheckpointRegister() flow.RegisterInfo_Ref {
	return builder.CreateRegister("checkpoint", builder.index.Int)
}

func (builder *dubBuilder) ZeroRegister(dst flow.RegisterInfo_Ref) flow.DubOp {
	info := builder.flow.RegisterInfo_Scope.Get(dst)
	switch t := info.T.(type) {
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
		panic(info.T)
	}
}

func makeRuneSwitch(cond flow.RegisterInfo_Ref, op string, value rune, builder *dubBuilder) (graph.NodeID, graph.NodeID) {
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

	builder.graph.Connect(make_value, flow.NORMAL, compare)
	builder.graph.Connect(compare, flow.NORMAL, decide)

	return make_value, decide
}

func lowerRuneMatch(match *tree.RuneRangeMatch, used bool, builder *dubBuilder, gr *graph.GraphRegion) flow.RegisterInfo_Ref {
	// Read
	cond := flow.NoRegisterInfo
	if len(match.Filters) > 0 || used {
		cond = builder.CreateRegister("c", builder.index.Rune)
	}
	body := builder.EmitOp(&flow.Peek{Dst: cond})
	gr.AttachFlow(flow.NORMAL, body)
	gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
	gr.RegisterExit(builder.EmitEdge(body, flow.FAIL), flow.FAIL)

	filters := builder.graph.CreateRegion(numRegionExits)

	onMatch := flow.FAIL
	onNoMatch := flow.NORMAL
	if !match.Invert {
		onMatch, onNoMatch = onNoMatch, onMatch
		// Make sure the implicit exit points to "fail"
		filters.Swap(flow.NORMAL, flow.FAIL)
	}

	for _, flt := range match.Filters {
		if flt.Min > flt.Max {
			panic(flt.Min)
		}
		if flt.Min != flt.Max {
			minEntry, minDecide := makeRuneSwitch(cond, ">=", flt.Min, builder)
			maxEntry, maxDecide := makeRuneSwitch(cond, "<=", flt.Max, builder)

			// Check only if we haven't found a match.
			filters.AttachFlow(onNoMatch, minEntry)

			// Match
			builder.graph.Connect(minDecide, 0, maxEntry)
			filters.RegisterExit(builder.EmitEdge(maxDecide, 0), onMatch)

			// No match
			filters.RegisterExit(builder.EmitEdge(minDecide, 1), onNoMatch)
			filters.RegisterExit(builder.EmitEdge(maxDecide, 1), onNoMatch)
		} else {
			entry, decide := makeRuneSwitch(cond, "==", flt.Min, builder)

			// Check only if we haven't found a match.
			filters.AttachFlow(onNoMatch, entry)

			// Match
			filters.RegisterExit(builder.EmitEdge(decide, 0), onMatch)

			// No match
			filters.RegisterExit(builder.EmitEdge(decide, 1), onNoMatch)
		}
	}

	// The rune matched, consume it.
	if filters.HasFlow(flow.NORMAL) {
		c := builder.EmitOp(&flow.Consume{})
		filters.AttachFlow(flow.NORMAL, c)
		filters.RegisterExit(builder.EmitEdge(c, flow.NORMAL), flow.NORMAL)
	}
	// Make the fail official.
	if filters.HasFlow(flow.FAIL) {
		f := builder.EmitOp(&flow.Fail{})
		filters.AttachFlow(flow.FAIL, f)
		filters.RegisterExit(builder.EmitEdge(f, flow.FAIL), flow.FAIL)
	}

	gr.Splice(flow.NORMAL, filters)
	return cond
}

func lowerMatch(match tree.TextMatch, builder *dubBuilder, gr *graph.GraphRegion) {
	switch match := match.(type) {
	case *tree.RuneRangeMatch:
		lowerRuneMatch(match, false, builder, gr)
	case *tree.StringLiteralMatch:
		// HACK desugar
		for _, c := range []rune(match.Value) {
			lowerRuneMatch(&tree.RuneRangeMatch{Filters: []*tree.RuneFilter{&tree.RuneFilter{Min: c, Max: c}}}, false, builder, gr)
		}
	case *tree.MatchSequence:
		for _, child := range match.Matches {
			lowerMatch(child, builder, gr)
		}
	case *tree.MatchChoice:
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		gr.AttachFlow(flow.NORMAL, head)

		for i, child := range match.Matches {
			block := builder.graph.CreateRegion(numRegionExits)
			lowerMatch(child, builder, block)

			// Recover if not the last block.
			if i < len(match.Matches)-1 {
				newHead := builder.EmitOp(&flow.Recover{Src: checkpoint})
				block.AttachFlow(flow.FAIL, newHead)
				gr.SpliceToEdge(builder.EmitEdge(head, flow.NORMAL), block)
				head = newHead
			} else {
				gr.SpliceToEdge(builder.EmitEdge(head, flow.NORMAL), block)
			}
		}
	case *tree.MatchRepeat:
		// HACK unroll
		for i := 0; i < match.Min; i++ {
			lowerMatch(match.Match, builder, gr)
		}

		child := builder.graph.CreateRegion(numRegionExits)

		// Checkpoint
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		child.AttachFlow(flow.NORMAL, head)
		child.RegisterExit(builder.EmitEdge(head, flow.NORMAL), flow.NORMAL)

		// Handle the body
		lowerMatch(match.Match, builder, child)

		// Normal flow iterates
		child.AttachFlow(flow.NORMAL, head)

		// Stop iterating on failure and recover
		{
			body := builder.EmitOp(&flow.Recover{Src: checkpoint})
			child.AttachFlow(flow.FAIL, body)
			child.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		}

		gr.Splice(flow.NORMAL, child)
	case *tree.MatchLookahead:
		child := builder.graph.CreateRegion(numRegionExits)

		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.LookaheadBegin{Dst: checkpoint})
		child.AttachFlow(flow.NORMAL, head)
		child.RegisterExit(builder.EmitEdge(head, flow.NORMAL), flow.NORMAL)

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

		gr.Splice(flow.NORMAL, child)
	default:
		panic(match)
	}
}

func lowerMultiValueExpr(expr tree.ASTExpr, builder *dubBuilder, used bool, gr *graph.GraphRegion) []flow.RegisterInfo_Ref {
	switch expr := expr.(type) {

	case *tree.Call:
		args := make([]flow.RegisterInfo_Ref, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, builder, true, gr)
		}
		var dsts []flow.RegisterInfo_Ref
		if used {
			dsts = make([]flow.RegisterInfo_Ref, len(expr.T))
			for i, t := range expr.T {
				dsts[i] = builder.CreateRegister("", t)
			}
		}
		target, ok := expr.Target.(*tree.FuncDecl)
		if !ok {
			panic(expr.Target)
		}
		body := builder.EmitOp(&flow.CallOp{Target: target.F, Args: args, Dsts: dsts})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		gr.RegisterExit(builder.EmitEdge(body, flow.FAIL), flow.FAIL)

		return dsts
	default:
		return []flow.RegisterInfo_Ref{lowerExpr(expr, builder, used, gr)}
	}
}

func lowerExpr(expr tree.ASTExpr, builder *dubBuilder, used bool, gr *graph.GraphRegion) flow.RegisterInfo_Ref {
	switch expr := expr.(type) {
	case *tree.If:
		cond := lowerExpr(expr.Expr, builder, true, gr)
		decide := builder.EmitOp(&flow.SwitchOp{Cond: cond})
		gr.AttachFlow(flow.NORMAL, decide)

		block := builder.graph.CreateRegion(numRegionExits)
		lowerBlock(expr.Block, builder, block)

		gr.SpliceToEdge(builder.EmitEdge(decide, 0), block)
		gr.RegisterExit(builder.EmitEdge(decide, 1), flow.NORMAL)
		return flow.NoRegisterInfo

	case *tree.Repeat:
		// HACK unroll
		for i := 0; i < expr.Min; i++ {
			lowerBlock(expr.Block, builder, gr)
		}

		block := builder.graph.CreateRegion(numRegionExits)

		// Checkpoint at head of loop
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		block.AttachFlow(flow.NORMAL, head)
		block.RegisterExit(builder.EmitEdge(head, flow.NORMAL), flow.NORMAL)

		// Handle the body
		lowerBlock(expr.Block, builder, block)

		// Normal flow iterates
		block.AttachFlow(flow.NORMAL, head)

		// Stop iterating on failure and recover
		{
			body := builder.EmitOp(&flow.Recover{Src: checkpoint})
			block.AttachFlow(flow.FAIL, body)
			block.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		}

		gr.Splice(flow.NORMAL, block)
		return flow.NoRegisterInfo
	case *tree.Choice:
		checkpoint := flow.NoRegisterInfo
		if len(expr.Blocks) > 1 {
			checkpoint = builder.CreateCheckpointRegister()
		}
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		gr.AttachFlow(flow.NORMAL, head)

		for i, b := range expr.Blocks {
			block := builder.graph.CreateRegion(numRegionExits)
			lowerBlock(b, builder, block)

			// Recover if not the last block.
			if i < len(expr.Blocks)-1 {
				newHead := builder.EmitOp(&flow.Recover{Src: checkpoint})
				block.AttachFlow(flow.FAIL, newHead)
				gr.SpliceToEdge(builder.EmitEdge(head, flow.NORMAL), block)
				head = newHead
			} else {
				gr.SpliceToEdge(builder.EmitEdge(head, flow.NORMAL), block)
			}
		}
		return flow.NoRegisterInfo

	case *tree.Optional:
		// Checkpoint
		checkpoint := builder.CreateCheckpointRegister()
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		gr.AttachFlow(flow.NORMAL, head)
		gr.RegisterExit(builder.EmitEdge(head, flow.NORMAL), flow.NORMAL)

		block := builder.graph.CreateRegion(numRegionExits)

		lowerBlock(expr.Block, builder, block)

		if block.HasFlow(flow.FAIL) {
			restore := builder.EmitOp(&flow.Recover{Src: checkpoint})
			block.AttachFlow(flow.FAIL, restore)
			block.RegisterExit(builder.EmitEdge(restore, flow.NORMAL), flow.NORMAL)
		}

		gr.Splice(flow.NORMAL, block)

		return flow.NoRegisterInfo

	case *tree.NameRef:
		if !used {
			return flow.NoRegisterInfo
		}
		dst := builder.CreateRegister(expr.Name.Text, builder.decl.LocalInfo_Scope.Get(expr.Local).T)
		body := builder.EmitOp(&flow.CopyOp{Src: builder.localMap[expr.Local], Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Assign:
		var srcs []flow.RegisterInfo_Ref
		if expr.Expr != nil {
			srcs = lowerMultiValueExpr(expr.Expr, builder, true, gr)
			if len(expr.Targets) != len(srcs) {
				panic(expr.Targets)
			}
		}
		dsts := make([]flow.RegisterInfo_Ref, len(expr.Targets))
		for i, etgt := range expr.Targets {
			tgt, ok := etgt.(*tree.NameRef)
			if !ok {
				panic(expr.Targets)
			}
			if tree.IsDiscard(tgt.Name.Text) {
				continue
			}
			dst := builder.localMap[tgt.Local]
			dsts[i] = dst
			var op flow.DubOp
			if srcs != nil {
				src := srcs[i]
				if src == flow.NoRegisterInfo {
					panic(expr)
				}
				op = &flow.CopyOp{Src: src, Dst: dst}
			} else {
				op = builder.ZeroRegister(dst)
			}
			body := builder.EmitOp(op)
			gr.AttachFlow(flow.NORMAL, body)
			gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		}
		// HACK should actually return a multivalue
		if len(dsts) == 1 {
			return dsts[0]
		}
		return flow.NoRegisterInfo

	case *tree.RuneLiteral:
		if !used {
			return flow.NoRegisterInfo
		}
		dst := builder.CreateRegister("c_r", builder.index.Rune)
		body := builder.EmitOp(&flow.ConstantRuneOp{Value: expr.Value, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.StringLiteral:
		if !used {
			return flow.NoRegisterInfo
		}
		dst := builder.CreateRegister("c_s", builder.index.String)
		body := builder.EmitOp(&flow.ConstantStringOp{Value: expr.Value, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.IntLiteral:
		if !used {
			return flow.NoRegisterInfo
		}
		dst := builder.CreateRegister("c_i", builder.index.Int)
		body := builder.EmitOp(&flow.ConstantIntOp{Value: int64(expr.Value), Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.BoolLiteral:
		if !used {
			return flow.NoRegisterInfo
		}
		dst := builder.CreateRegister("c_b", builder.index.Bool)
		body := builder.EmitOp(&flow.ConstantBoolOp{Value: expr.Value, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Return:
		exprs := make([]flow.RegisterInfo_Ref, len(expr.Exprs))
		for i, e := range expr.Exprs {
			exprs[i] = lowerExpr(e, builder, true, gr)
		}
		body := builder.EmitOp(&flow.ReturnOp{Exprs: exprs})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.RETURN)
		return flow.NoRegisterInfo

	case *tree.Fail:
		body := builder.EmitOp(&flow.Fail{})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.FAIL), flow.FAIL)

		return flow.NoRegisterInfo

	case *tree.Position:
		if !used {
			return flow.NoRegisterInfo
		}
		pos := builder.CreateRegister("pos", builder.index.Int)
		// HACK assume checkpoint is just the index
		head := builder.EmitOp(&flow.Checkpoint{Dst: pos})
		gr.AttachFlow(flow.NORMAL, head)
		gr.RegisterExit(builder.EmitEdge(head, flow.NORMAL), flow.NORMAL)
		return pos
	case *tree.BinaryOp:
		left := lowerExpr(expr.Left, builder, true, gr)
		right := lowerExpr(expr.Right, builder, true, gr)
		dst := flow.NoRegisterInfo
		if used {
			dst = builder.CreateRegister("", expr.T)
		}
		body := builder.EmitOp(&flow.BinaryOp{Left: left, Op: expr.Op, Right: right, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst
	case *tree.Append:
		l := lowerExpr(expr.List, builder, true, gr)
		v := lowerExpr(expr.Expr, builder, true, gr)
		dst := flow.NoRegisterInfo
		if used {
			dst = builder.CreateRegister("", expr.T)
		}

		body := builder.EmitOp(&flow.AppendOp{List: l, Value: v, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Call:
		dsts := lowerMultiValueExpr(expr, builder, true, gr)
		dst := flow.NoRegisterInfo
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
				Value: lowerExpr(arg.Expr, builder, true, gr),
			}
		}
		t := tree.ResolveType(expr.Type)
		s, ok := t.(*core.StructType)
		if !ok {
			panic(t)
		}
		dst := flow.NoRegisterInfo
		if used {
			dst = builder.CreateRegister("", t)
		}
		body := builder.EmitOp(&flow.ConstructOp{Type: s, Args: args, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.ConstructList:
		args := make([]flow.RegisterInfo_Ref, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, builder, true, gr)
		}
		t := tree.ResolveType(expr.Type)
		l, ok := t.(*core.ListType)
		if !ok {
			panic(t)
		}
		dst := flow.NoRegisterInfo
		if used {
			dst = builder.CreateRegister("", t)
		}
		body := builder.EmitOp(&flow.ConstructListOp{Type: l, Args: args, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Coerce:
		t := tree.ResolveType(expr.Type)
		src := lowerExpr(expr.Expr, builder, true, gr)
		dst := flow.NoRegisterInfo
		if used {
			dst = builder.CreateRegister("", t)
		}
		body := builder.EmitOp(&flow.CoerceOp{Src: src, T: t, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		// TODO can coersion fail?
		gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		return dst

	case *tree.Slice:
		start := builder.CreateRegister("pos", builder.index.Int)
		// HACK assume checkpoint is just the index
		{
			head := builder.EmitOp(&flow.Checkpoint{Dst: start})
			gr.AttachFlow(flow.NORMAL, head)
			gr.RegisterExit(builder.EmitEdge(head, flow.NORMAL), flow.NORMAL)
		}
		lowerBlock(expr.Block, builder, gr)

		// Create a slice
		dst := flow.NoRegisterInfo
		if used {
			dst = builder.CreateRegister("slice", builder.index.String)
		}
		{
			body := builder.EmitOp(&flow.Slice{Src: start, Dst: dst})
			gr.AttachFlow(flow.NORMAL, body)
			gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		}
		return dst

	case *tree.StringMatch:
		dst := flow.NoRegisterInfo
		start := flow.NoRegisterInfo

		// Checkpoint
		if used {
			start = builder.CreateRegister("pos", builder.index.Int)
			// HACK assume checkpoint is just the index
			{
				head := builder.EmitOp(&flow.Checkpoint{Dst: start})
				gr.AttachFlow(flow.NORMAL, head)
				gr.RegisterExit(builder.EmitEdge(head, flow.NORMAL), flow.NORMAL)
			}
		}

		lowerMatch(expr.Match, builder, gr)

		// Create a slice
		if used {
			dst = builder.CreateRegister("slice", builder.index.String)
			body := builder.EmitOp(&flow.Slice{Src: start, Dst: dst})
			gr.AttachFlow(flow.NORMAL, body)
			gr.RegisterExit(builder.EmitEdge(body, flow.NORMAL), flow.NORMAL)
		}
		return dst

	case *tree.RuneMatch:
		return lowerRuneMatch(expr.Match, used, builder, gr)
	default:
		panic(expr)
	}

}

func lowerBlock(block []tree.ASTExpr, builder *dubBuilder, gr *graph.GraphRegion) {
	for _, expr := range block {
		lowerExpr(expr, builder, false, gr)
	}
}

func lowerAST(program *tree.Program, decl *tree.FuncDecl, funcMap []*flow.LLFunc) *flow.LLFunc {
	g := graph.CreateGraph()
	ops := []flow.DubOp{
		&flow.EntryOp{},
		&flow.ExitOp{},
	}

	f := funcMap[decl.F]

	builder := &dubBuilder{
		index: program.Builtins,
		decl:  decl,
		flow:  f,
		graph: g,
		ops:   ops,
	}

	// Allocate register for locals
	numLocals := decl.LocalInfo_Scope.Len()
	builder.localMap = make([]flow.RegisterInfo_Ref, numLocals)
	for i := 0; i < numLocals; i++ {
		lcl := decl.LocalInfo_Scope.Get(tree.LocalInfo_Ref(i))
		builder.localMap[i] = builder.CreateRegister(lcl.Name, lcl.T)
	}

	// Function parameters
	params := make([]flow.RegisterInfo_Ref, len(decl.Params))
	for i, p := range decl.Params {
		params[i] = builder.localMap[p.Name.Local]
	}
	f.Params = params

	// Function returns
	types := make([]core.DubType, len(decl.ReturnTypes))
	for i, node := range decl.ReturnTypes {
		types[i] = tree.ResolveType(node)
	}
	f.ReturnTypes = types

	gr := g.CreateRegion(numRegionExits)
	lowerBlock(decl.Block, builder, gr)
	gr.MergeFlowInto(flow.RETURN, flow.NORMAL)

	// TODO only connect the real exits, assert no virtual exits.
	for i := 0; i < numRegionExits; i++ {
		if gr.HasFlow(i) {
			fe := builder.EmitOp(&flow.FlowExitOp{Flow: i})
			gr.AttachFlow(i, fe)
			gr.RegisterExit(builder.EmitEdge(fe, 0), 0)
		}
	}

	g.ConnectRegion(gr)

	f.CFG = g
	f.Ops = builder.ops
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
			F:                  fref,
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

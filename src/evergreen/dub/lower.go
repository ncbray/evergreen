package dub

import (
	"evergreen/base"
	"evergreen/dub/flow"
	"evergreen/dub/tree"
)

const REGION_EXITS = 4

type GlobalDubBuilder struct {
	Types  map[tree.ASTType]flow.DubType
	String flow.DubType
	Rune   flow.DubType
	Int    flow.DubType
	Bool   flow.DubType
}

func (builder *GlobalDubBuilder) TranslateType(t tree.ASTType) flow.DubType {
	switch t := t.(type) {
	case *tree.StructDecl:
		dt, ok := builder.Types[t]
		if !ok {
			panic(t)
		}
		return dt
	case *tree.BuiltinType:
		dt, ok := builder.Types[t]
		if !ok {
			panic(t)
		}
		return dt
	case *tree.ListType:
		parent := builder.TranslateType(t.Type)
		// TODO memoize
		return &flow.ListType{Type: parent}
	default:
		panic(t)
	}
}

type DubBuilder struct {
	decl      *tree.FuncDecl
	registers []flow.RegisterInfo
	localMap  []flow.DubRegister
	glbl      *GlobalDubBuilder
	graph     *base.Graph
	ops       []flow.DubOp
}

func (builder *DubBuilder) EmitOp(op flow.DubOp) base.NodeID {
	id := builder.graph.CreateNode(2)
	if int(id) != len(builder.ops) {
		panic(op)
	}
	builder.ops = append(builder.ops, op)
	return id
}

func (builder *DubBuilder) CreateRegister(t tree.ASTType) flow.DubRegister {
	return builder.CreateLLRegister(builder.glbl.TranslateType(t))
}

func (builder *DubBuilder) CreateLLRegister(t flow.DubType) flow.DubRegister {
	builder.registers = append(builder.registers, flow.RegisterInfo{T: t})
	return flow.DubRegister(len(builder.registers) - 1)
}

func (builder *DubBuilder) ZeroRegister(dst flow.DubRegister) flow.DubOp {
	info := builder.registers[dst]
	switch t := info.T.(type) {
	case *flow.LLStruct:
		return &flow.ConstantNilOp{Dst: dst}
	case *flow.IntrinsicType:
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

func makeRuneSwitch(cond flow.DubRegister, op string, value rune, builder *DubBuilder) (base.NodeID, base.NodeID) {
	vreg := builder.CreateLLRegister(builder.glbl.Rune)
	make_value := builder.EmitOp(&flow.ConstantRuneOp{Value: value, Dst: vreg})

	breg := builder.CreateLLRegister(builder.glbl.Bool)
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

func lowerRuneMatch(match *tree.RuneRangeMatch, used bool, builder *DubBuilder, gr *base.GraphRegion) flow.DubRegister {
	// Read
	cond := flow.NoRegister
	if len(match.Filters) > 0 || used {
		cond = builder.CreateLLRegister(builder.glbl.Rune)
	}
	body := builder.EmitOp(&flow.Peek{Dst: cond})
	gr.AttachFlow(flow.NORMAL, body)
	gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
	gr.RegisterExit(body, flow.FAIL, flow.FAIL)

	filters := builder.graph.CreateRegion(REGION_EXITS)

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
			filters.RegisterExit(maxDecide, 0, onMatch)

			// No match
			filters.RegisterExit(minDecide, 1, onNoMatch)
			filters.RegisterExit(maxDecide, 1, onNoMatch)
		} else {
			entry, decide := makeRuneSwitch(cond, "==", flt.Min, builder)

			// Check only if we haven't found a match.
			filters.AttachFlow(onNoMatch, entry)

			// Match
			filters.RegisterExit(decide, 0, onMatch)

			// No match
			filters.RegisterExit(decide, 1, onNoMatch)
		}
	}

	// The rune matched, consume it.
	if filters.HasFlow(flow.NORMAL) {
		c := builder.EmitOp(&flow.Consume{})
		filters.AttachFlow(flow.NORMAL, c)
		filters.RegisterExit(c, flow.NORMAL, flow.NORMAL)
	}
	// Make the fail official.
	if filters.HasFlow(flow.FAIL) {
		f := builder.EmitOp(&flow.Fail{})
		filters.AttachFlow(flow.FAIL, f)
		filters.RegisterExit(f, flow.FAIL, flow.FAIL)
	}

	gr.Splice(flow.NORMAL, filters)
	return cond
}

func lowerMatch(match tree.TextMatch, builder *DubBuilder, gr *base.GraphRegion) {
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
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		gr.AttachFlow(flow.NORMAL, head)

		for i, child := range match.Matches {
			block := builder.graph.CreateRegion(REGION_EXITS)
			lowerMatch(child, builder, block)

			// Recover if not the last block.
			if i < len(match.Matches)-1 {
				newHead := builder.EmitOp(&flow.Recover{Src: checkpoint})
				block.AttachFlow(flow.FAIL, newHead)
				gr.SpliceToEdge(head, flow.NORMAL, block)
				head = newHead
			} else {
				gr.SpliceToEdge(head, flow.NORMAL, block)
			}
		}
	case *tree.MatchRepeat:
		// HACK unroll
		for i := 0; i < match.Min; i++ {
			lowerMatch(match.Match, builder, gr)
		}

		child := builder.graph.CreateRegion(REGION_EXITS)

		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		child.AttachFlow(flow.NORMAL, head)
		child.RegisterExit(head, flow.NORMAL, flow.NORMAL)

		// Handle the body
		lowerMatch(match.Match, builder, child)

		// Normal flow iterates
		child.AttachFlow(flow.NORMAL, head)

		// Stop iterating on failure and recover
		{
			body := builder.EmitOp(&flow.Recover{Src: checkpoint})
			child.AttachFlow(flow.FAIL, body)
			child.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		}

		gr.Splice(flow.NORMAL, child)
	case *tree.MatchLookahead:
		child := builder.graph.CreateRegion(REGION_EXITS)

		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := builder.EmitOp(&flow.LookaheadBegin{Dst: checkpoint})
		child.AttachFlow(flow.NORMAL, head)
		child.RegisterExit(head, flow.NORMAL, flow.NORMAL)

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

		child.RegisterExit(normal, flow.NORMAL, flow.NORMAL)
		child.RegisterExit(fail, flow.FAIL, flow.FAIL)

		gr.Splice(flow.NORMAL, child)
	default:
		panic(match)
	}
}

func lowerMultiValueExpr(expr tree.ASTExpr, builder *DubBuilder, used bool, gr *base.GraphRegion) []flow.DubRegister {
	switch expr := expr.(type) {

	case *tree.Call:
		args := make([]flow.DubRegister, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, builder, true, gr)
		}
		var dsts []flow.DubRegister
		if used {
			dsts = make([]flow.DubRegister, len(expr.T))
			for i, t := range expr.T {
				dsts[i] = builder.CreateRegister(t)
			}
		}
		body := builder.EmitOp(&flow.CallOp{Name: expr.Name.Text, Args: args, Dsts: dsts})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		gr.RegisterExit(body, flow.FAIL, flow.FAIL)

		return dsts
	default:
		return []flow.DubRegister{lowerExpr(expr, builder, used, gr)}
	}
}

func lowerExpr(expr tree.ASTExpr, builder *DubBuilder, used bool, gr *base.GraphRegion) flow.DubRegister {
	switch expr := expr.(type) {
	case *tree.If:
		cond := lowerExpr(expr.Expr, builder, true, gr)
		decide := builder.EmitOp(&flow.SwitchOp{Cond: cond})
		gr.AttachFlow(flow.NORMAL, decide)

		block := builder.graph.CreateRegion(REGION_EXITS)
		lowerBlock(expr.Block, builder, block)

		gr.SpliceToEdge(decide, 0, block)
		gr.RegisterExit(decide, 1, flow.NORMAL)
		return flow.NoRegister

	case *tree.Repeat:
		// HACK unroll
		for i := 0; i < expr.Min; i++ {
			lowerBlock(expr.Block, builder, gr)
		}

		block := builder.graph.CreateRegion(REGION_EXITS)

		// Checkpoint at head of loop
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		block.AttachFlow(flow.NORMAL, head)
		block.RegisterExit(head, flow.NORMAL, flow.NORMAL)

		// Handle the body
		lowerBlock(expr.Block, builder, block)

		// Normal flow iterates
		block.AttachFlow(flow.NORMAL, head)

		// Stop iterating on failure and recover
		{
			body := builder.EmitOp(&flow.Recover{Src: checkpoint})
			block.AttachFlow(flow.FAIL, body)
			block.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		}

		gr.Splice(flow.NORMAL, block)
		return flow.NoRegister
	case *tree.Choice:
		checkpoint := flow.NoRegister
		if len(expr.Blocks) > 1 {
			checkpoint = builder.CreateLLRegister(builder.glbl.Int)
		}
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		gr.AttachFlow(flow.NORMAL, head)

		for i, b := range expr.Blocks {
			block := builder.graph.CreateRegion(REGION_EXITS)
			lowerBlock(b, builder, block)

			// Recover if not the last block.
			if i < len(expr.Blocks)-1 {
				newHead := builder.EmitOp(&flow.Recover{Src: checkpoint})
				block.AttachFlow(flow.FAIL, newHead)
				gr.SpliceToEdge(head, flow.NORMAL, block)
				head = newHead
			} else {
				gr.SpliceToEdge(head, flow.NORMAL, block)
			}
		}
		return flow.NoRegister

	case *tree.Optional:
		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := builder.EmitOp(&flow.Checkpoint{Dst: checkpoint})
		gr.AttachFlow(flow.NORMAL, head)
		gr.RegisterExit(head, flow.NORMAL, flow.NORMAL)

		block := builder.graph.CreateRegion(REGION_EXITS)

		lowerBlock(expr.Block, builder, block)

		restore := builder.EmitOp(&flow.Recover{Src: checkpoint})
		block.AttachFlow(flow.FAIL, restore)
		block.RegisterExit(restore, flow.NORMAL, flow.NORMAL)

		gr.Splice(flow.NORMAL, block)

		return flow.NoRegister

	case *tree.NameRef:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateRegister(builder.decl.LocalInfo_Scope.Get(expr.Local).T)
		body := builder.EmitOp(&flow.CopyOp{Src: builder.localMap[expr.Local], Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.Assign:
		var srcs []flow.DubRegister
		if expr.Expr != nil {
			srcs = lowerMultiValueExpr(expr.Expr, builder, true, gr)
			if len(expr.Targets) != len(srcs) {
				panic(expr.Targets)
			}
		}
		for i, etgt := range expr.Targets {
			tgt, ok := etgt.(*tree.NameRef)
			if !ok {
				panic(expr.Targets)
			}
			if tree.IsDiscard(tgt.Name.Text) {
				continue
			}
			dst := builder.localMap[tgt.Local]
			var op flow.DubOp
			if srcs != nil {
				op = &flow.CopyOp{Src: srcs[i], Dst: dst}
			} else {
				op = builder.ZeroRegister(dst)
			}
			body := builder.EmitOp(op)
			gr.AttachFlow(flow.NORMAL, body)
			gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		}
		// HACK should actuall return a multivalue
		return flow.NoRegister

	case *tree.RuneLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Rune)
		body := builder.EmitOp(&flow.ConstantRuneOp{Value: expr.Value, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.StringLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.String)
		body := builder.EmitOp(&flow.ConstantStringOp{Value: expr.Value, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.IntLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Int)
		body := builder.EmitOp(&flow.ConstantIntOp{Value: int64(expr.Value), Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.BoolLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Bool)
		body := builder.EmitOp(&flow.ConstantBoolOp{Value: expr.Value, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.Return:
		exprs := make([]flow.DubRegister, len(expr.Exprs))
		for i, e := range expr.Exprs {
			exprs[i] = lowerExpr(e, builder, true, gr)
		}
		body := builder.EmitOp(&flow.ReturnOp{Exprs: exprs})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.RETURN)
		return flow.NoRegister

	case *tree.Fail:
		body := builder.EmitOp(&flow.Fail{})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.FAIL, flow.FAIL)

		return flow.NoRegister

	case *tree.Position:
		if !used {
			return flow.NoRegister
		}
		pos := builder.CreateLLRegister(builder.glbl.Int)
		// HACK assume checkpoint is just the index
		head := builder.EmitOp(&flow.Checkpoint{Dst: pos})
		gr.AttachFlow(flow.NORMAL, head)
		gr.RegisterExit(head, flow.NORMAL, flow.NORMAL)
		return pos
	case *tree.BinaryOp:
		left := lowerExpr(expr.Left, builder, true, gr)
		right := lowerExpr(expr.Right, builder, true, gr)
		dst := flow.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := builder.EmitOp(&flow.BinaryOp{Left: left, Op: expr.Op, Right: right, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst
	case *tree.Append:
		l := lowerExpr(expr.List, builder, true, gr)
		v := lowerExpr(expr.Expr, builder, true, gr)
		dst := flow.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}

		body := builder.EmitOp(&flow.AppendOp{List: l, Value: v, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.Call:
		dsts := lowerMultiValueExpr(expr, builder, true, gr)
		dst := flow.NoRegister
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
		t := builder.glbl.TranslateType(tree.ResolveType(expr.Type))
		s, ok := t.(*flow.LLStruct)
		if !ok {
			panic(t)
		}
		dst := flow.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := builder.EmitOp(&flow.ConstructOp{Type: s, Args: args, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.ConstructList:
		args := make([]flow.DubRegister, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, builder, true, gr)
		}
		t := builder.glbl.TranslateType(tree.ResolveType(expr.Type))
		l, ok := t.(*flow.ListType)
		if !ok {
			panic(t)
		}
		dst := flow.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := builder.EmitOp(&flow.ConstructListOp{Type: l, Args: args, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.Coerce:
		t := builder.glbl.TranslateType(tree.ResolveType(expr.Type))
		src := lowerExpr(expr.Expr, builder, true, gr)
		dst := flow.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := builder.EmitOp(&flow.CoerceOp{Src: src, T: t, Dst: dst})
		gr.AttachFlow(flow.NORMAL, body)
		// TODO can coersion fail?
		gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		return dst

	case *tree.Slice:
		start := builder.CreateLLRegister(builder.glbl.Int)
		// HACK assume checkpoint is just the index
		{
			head := builder.EmitOp(&flow.Checkpoint{Dst: start})
			gr.AttachFlow(flow.NORMAL, head)
			gr.RegisterExit(head, flow.NORMAL, flow.NORMAL)
		}
		lowerBlock(expr.Block, builder, gr)

		// Create a slice
		dst := flow.NoRegister
		if used {
			dst = builder.CreateLLRegister(builder.glbl.String)
		}
		{
			body := builder.EmitOp(&flow.Slice{Src: start, Dst: dst})
			gr.AttachFlow(flow.NORMAL, body)
			gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		}
		return dst

	case *tree.StringMatch:
		dst := flow.NoRegister
		start := flow.NoRegister

		// Checkpoint
		if used {
			start = builder.CreateLLRegister(builder.glbl.Int)
			// HACK assume checkpoint is just the index
			{
				head := builder.EmitOp(&flow.Checkpoint{Dst: start})
				gr.AttachFlow(flow.NORMAL, head)
				gr.RegisterExit(head, flow.NORMAL, flow.NORMAL)
			}
		}

		lowerMatch(expr.Match, builder, gr)

		// Create a slice
		if used {
			dst = builder.CreateLLRegister(builder.glbl.String)
			body := builder.EmitOp(&flow.Slice{Src: start, Dst: dst})
			gr.AttachFlow(flow.NORMAL, body)
			gr.RegisterExit(body, flow.NORMAL, flow.NORMAL)
		}
		return dst

	case *tree.RuneMatch:
		return lowerRuneMatch(expr.Match, used, builder, gr)
	default:
		panic(expr)
	}

}

func lowerBlock(block []tree.ASTExpr, builder *DubBuilder, gr *base.GraphRegion) {
	for _, expr := range block {
		lowerExpr(expr, builder, false, gr)
	}
}

func LowerAST(decl *tree.FuncDecl, glbl *GlobalDubBuilder) *flow.LLFunc {
	g := base.CreateGraph()
	ops := []flow.DubOp{
		&flow.EntryOp{},
		&flow.ExitOp{},
	}
	builder := &DubBuilder{decl: decl, glbl: glbl, graph: g, ops: ops}

	f := &flow.LLFunc{Name: decl.Name.Text}

	// Allocate register for locals
	numLocals := decl.LocalInfo_Scope.Len()
	builder.localMap = make([]flow.DubRegister, numLocals)
	for i := 0; i < numLocals; i++ {
		builder.localMap[i] = builder.CreateRegister(decl.LocalInfo_Scope.Get(tree.LocalInfo_Ref(i)).T)
	}

	// Function parameters
	params := make([]flow.DubRegister, len(decl.Params))
	for i, p := range decl.Params {
		params[i] = builder.localMap[p.Name.Local]
	}
	f.Params = params

	// Function returns
	types := make([]flow.DubType, len(decl.ReturnTypes))
	for i, node := range decl.ReturnTypes {
		types[i] = builder.glbl.TranslateType(tree.ResolveType(node))
	}
	f.ReturnTypes = types

	gr := g.CreateRegion(REGION_EXITS)
	lowerBlock(decl.Block, builder, gr)
	gr.MergeFlowInto(flow.RETURN, flow.NORMAL)

	// TODO only connect the real exits, assert no virtual exits.
	for i := 0; i < REGION_EXITS; i++ {
		if gr.HasFlow(i) {
			fe := builder.EmitOp(&flow.FlowExitOp{Flow: i})
			gr.AttachFlow(i, fe)
			gr.RegisterExit(fe, 0, 0)
		}
	}

	g.ConnectRegion(gr)

	f.CFG = g
	f.Ops = builder.ops
	f.Registers = builder.registers
	return f
}

func LowerStruct(decl *tree.StructDecl, s *flow.LLStruct, gbuilder *GlobalDubBuilder) *flow.LLStruct {
	fields := []*flow.LLField{}
	var implements *flow.LLStruct

	contains := []*flow.LLStruct{}
	for _, t := range decl.Contains {
		tt := gbuilder.TranslateType(tree.ResolveType(t))
		contained, ok := tt.(*flow.LLStruct)
		if !ok {
			panic(tt)
		}
		contains = append(contains, contained)

	}
	if decl.Implements != nil {
		t := gbuilder.TranslateType(tree.ResolveType(decl.Implements))
		var ok bool
		implements, ok = t.(*flow.LLStruct)
		if !ok {
			panic(t)
		}
	}
	for _, field := range decl.Fields {
		fields = append(fields, &flow.LLField{
			Name: field.Name.Text,
			T:    gbuilder.TranslateType(tree.ResolveType(field.Type)),
		})
	}
	*s = flow.LLStruct{
		Name:       decl.Name.Text,
		Implements: implements,
		Fields:     fields,
		Scoped:     decl.Scoped,
		Contains:   contains,
	}
	return s
}

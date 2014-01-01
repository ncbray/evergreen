package dub

import (
	"evergreen/base"
	"evergreen/dub/flow"
	"evergreen/dub/tree"
)

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
	switch info.T.(type) {
	case *flow.LLStruct:
		return &flow.ConstantNilOp{Dst: dst}
	case *flow.RuneType:
		return &flow.ConstantRuneOp{Value: 0, Dst: dst}
	case *flow.StringType:
		return &flow.ConstantStringOp{Value: "", Dst: dst}
	case *flow.IntType:
		return &flow.ConstantIntOp{Value: 0, Dst: dst}
	case *flow.BoolType:
		return &flow.ConstantBoolOp{Value: false, Dst: dst}
	default:
		panic(info.T)
	}
}

func makeRuneSwitch(cond flow.DubRegister, op string, value rune, builder *DubBuilder) (*base.Node, *base.Node) {
	vreg := builder.CreateLLRegister(builder.glbl.Rune)
	make_value := flow.CreateBlock([]flow.DubOp{
		&flow.ConstantRuneOp{Value: value, Dst: vreg},
	})

	breg := builder.CreateLLRegister(builder.glbl.Bool)
	compare := flow.CreateBlock([]flow.DubOp{
		&flow.BinaryOp{
			Left:  cond,
			Op:    op,
			Right: vreg,
			Dst:   breg,
		},
	})

	decide := flow.CreateSwitch(breg)

	make_value.SetExit(flow.NORMAL, compare)
	compare.SetExit(flow.NORMAL, decide)

	return make_value, decide
}

func lowerRuneMatch(match *tree.RuneRangeMatch, r *base.Region, builder *DubBuilder) flow.DubRegister {
	// Read
	cond := builder.CreateLLRegister(builder.glbl.Rune)
	body := flow.CreateBlock([]flow.DubOp{
		&flow.Peek{Dst: cond},
	})
	r.Connect(flow.NORMAL, body)
	r.AttachDefaultExits(body)

	filters := flow.CreateRegion()

	onMatch := flow.NORMAL
	onNoMatch := flow.FAIL

	if match.Invert {
		onMatch, onNoMatch = onNoMatch, onMatch
	} else {
		// Fail by default.
		filters.GetExit(flow.NORMAL).TransferEntries(filters.GetExit(flow.FAIL))
	}

	for _, flt := range match.Filters {
		if flt.Min > flt.Max {
			panic(flt.Min)
		}
		if flt.Min != flt.Max {
			minEntry, minDecide := makeRuneSwitch(cond, ">=", flt.Min, builder)
			maxEntry, maxDecide := makeRuneSwitch(cond, "<=", flt.Max, builder)

			// Check only if we haven't found a match.
			filters.Connect(onNoMatch, minEntry)

			// Match
			minDecide.SetExit(0, maxEntry)
			maxDecide.SetExit(0, filters.GetExit(onMatch))

			// No match
			minDecide.SetExit(1, filters.GetExit(onNoMatch))
			maxDecide.SetExit(1, filters.GetExit(onNoMatch))
		} else {
			entry, decide := makeRuneSwitch(cond, "==", flt.Min, builder)

			// Check only if we haven't found a match.
			filters.Connect(onNoMatch, entry)

			// Match
			decide.SetExit(0, filters.GetExit(onMatch))

			// No match
			decide.SetExit(1, filters.GetExit(onNoMatch))
		}
	}

	{
		// The rune matched, consume it.
		body = flow.CreateBlock([]flow.DubOp{
			&flow.Consume{},
		})
		filters.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, filters.GetExit(flow.NORMAL))
	}
	{
		// Make the fail official.
		body = flow.CreateBlock([]flow.DubOp{
			&flow.Fail{},
		})
		filters.Connect(flow.FAIL, body)
		body.SetExit(flow.FAIL, filters.GetExit(flow.FAIL))
	}
	r.Splice(0, filters)

	return cond
}

func lowerMatch(match tree.TextMatch, r *base.Region, builder *DubBuilder) {
	switch match := match.(type) {
	case *tree.RuneRangeMatch:
		lowerRuneMatch(match, r, builder)
	case *tree.StringLiteralMatch:
		// HACK desugar
		for _, c := range []rune(match.Value) {
			lowerRuneMatch(&tree.RuneRangeMatch{Filters: []*tree.RuneFilter{&tree.RuneFilter{Min: c, Max: c}}}, r, builder)
		}
	case *tree.MatchSequence:
		for _, child := range match.Matches {
			lowerMatch(child, r, builder)
		}
	case *tree.MatchChoice:
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := flow.CreateBlock([]flow.DubOp{
			&flow.Checkpoint{Dst: checkpoint},
		})
		r.Connect(flow.NORMAL, head)

		for i, child := range match.Matches {
			block := flow.CreateRegion()
			lowerMatch(child, block, builder)

			// Connect the head to the entry of this block
			entryNode := block.Head()
			entryEdge := block.Entry.GetExit(0)
			entryNode.ReplaceEntry(entryEdge, []*base.Edge{head.GetExit(0)})

			// Recover if not the last block.
			if i < len(match.Matches)-1 {
				head = flow.CreateBlock([]flow.DubOp{
					&flow.Recover{Src: checkpoint},
				})
				block.Connect(flow.FAIL, head)
			} else {
				head = nil
			}

			// Absorb the exits that have not been directed to head.
			r.AbsorbExits(block)
		}
	case *tree.MatchRepeat:
		// HACK unroll
		for i := 0; i < match.Min; i++ {
			lowerMatch(match.Match, r, builder)
		}

		child := flow.CreateRegion()

		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := flow.CreateBlock([]flow.DubOp{
			&flow.Checkpoint{Dst: checkpoint},
		})
		child.Connect(flow.NORMAL, head)
		head.SetExit(flow.NORMAL, child.GetExit(flow.NORMAL))

		// Handle the body
		lowerMatch(match.Match, child, builder)

		// Normal flow iterates
		child.GetExit(flow.NORMAL).TransferEntries(head)

		// Stop iterating on failure and recover
		{
			body := flow.CreateBlock([]flow.DubOp{
				&flow.Recover{Src: checkpoint},
			})

			child.Connect(flow.FAIL, body)
			body.SetExit(flow.NORMAL, child.GetExit(flow.NORMAL))
		}

		r.Splice(flow.NORMAL, child)
	default:
		panic(match)
	}
}

func lowerExpr(expr tree.ASTExpr, r *base.Region, builder *DubBuilder, used bool) flow.DubRegister {
	switch expr := expr.(type) {
	case *tree.If:
		// TODO Min
		//l := flow.CreateRegion()

		cond := lowerExpr(expr.Expr, r, builder, true)
		block := lowerBlock(expr.Block, builder)

		// TODO conditional
		decide := flow.CreateSwitch(cond)

		r.Connect(flow.NORMAL, decide)
		decide.SetExit(0, r.GetExit(flow.NORMAL))
		r.Splice(flow.NORMAL, block)
		decide.SetExit(1, r.GetExit(flow.NORMAL))

		return flow.NoRegister

	case *tree.Repeat:
		// HACK unroll
		for i := 0; i < expr.Min; i++ {
			block := lowerBlock(expr.Block, builder)
			r.Splice(flow.NORMAL, block)
		}

		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := flow.CreateBlock([]flow.DubOp{
			&flow.Checkpoint{Dst: checkpoint},
		})

		// Handle the body
		block := lowerBlock(expr.Block, builder)

		// Splice in the checkpoint as the first operation.
		// Note: block may be empty, but this code is carefully designed to work in that case.
		// Sure it's an infinite loop, but it would stange if that loop vanished.
		oldHead := block.Head()
		oldHead.TransferEntries(head)
		head.SetExit(flow.NORMAL, oldHead)

		// Normal flow iterates
		block.GetExit(flow.NORMAL).TransferEntries(head)

		// Stop iterating on failure and recover
		block.GetExit(flow.FAIL).TransferEntries(block.GetExit(flow.NORMAL))
		{
			body := flow.CreateBlock([]flow.DubOp{
				&flow.Recover{Src: checkpoint},
			})

			block.Connect(flow.NORMAL, body)
			body.SetExit(flow.NORMAL, block.GetExit(flow.NORMAL))
		}

		r.Splice(flow.NORMAL, block)

		return flow.NoRegister
	case *tree.Choice:
		checkpoint := flow.NoRegister
		if len(expr.Blocks) > 1 {
			checkpoint = builder.CreateLLRegister(builder.glbl.Int)
		}
		head := flow.CreateBlock([]flow.DubOp{
			&flow.Checkpoint{Dst: checkpoint},
		})
		r.Connect(flow.NORMAL, head)

		for i, block := range expr.Blocks {
			block := lowerBlock(block, builder)

			// Connect the head to the entry of this block
			entryNode := block.Head()
			entryEdge := block.Entry.GetExit(0)
			entryNode.ReplaceEntry(entryEdge, []*base.Edge{head.GetExit(0)})

			// Recover if not the last block.
			if i < len(expr.Blocks)-1 {
				head = flow.CreateBlock([]flow.DubOp{
					&flow.Recover{Src: checkpoint},
				})
				block.Connect(flow.FAIL, head)
			} else {
				head = nil
			}

			// Absorb the exits that have not been directed to head.
			r.AbsorbExits(block)
		}
		return flow.NoRegister

	case *tree.Optional:
		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := flow.CreateBlock([]flow.DubOp{
			&flow.Checkpoint{Dst: checkpoint},
		})
		r.Connect(flow.NORMAL, head)
		head.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))

		block := lowerBlock(expr.Block, builder)

		restore := flow.CreateBlock([]flow.DubOp{
			&flow.Recover{Src: checkpoint},
		})

		block.Connect(flow.FAIL, restore)
		restore.SetExit(flow.NORMAL, block.GetExit(flow.NORMAL))

		r.Splice(flow.NORMAL, block)

		return flow.NoRegister

	case *tree.GetName:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateRegister(builder.decl.Locals[expr.Info].T)
		body := flow.CreateBlock([]flow.DubOp{
			&flow.CopyOp{Src: builder.localMap[expr.Info], Dst: dst},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.Assign:
		dst := builder.localMap[expr.Info]
		var op flow.DubOp
		if expr.Expr != nil {
			src := lowerExpr(expr.Expr, r, builder, true)
			op = &flow.CopyOp{Src: src, Dst: dst}
		} else {
			op = builder.ZeroRegister(dst)
		}
		body := flow.CreateBlock([]flow.DubOp{op})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.RuneLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Rune)
		body := flow.CreateBlock([]flow.DubOp{
			&flow.ConstantRuneOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.StringLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.String)
		body := flow.CreateBlock([]flow.DubOp{
			&flow.ConstantStringOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.IntLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Int)
		body := flow.CreateBlock([]flow.DubOp{
			&flow.ConstantIntOp{Value: int64(expr.Value), Dst: dst},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.BoolLiteral:
		if !used {
			return flow.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Bool)
		body := flow.CreateBlock([]flow.DubOp{
			&flow.ConstantBoolOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.Return:
		exprs := make([]flow.DubRegister, len(expr.Exprs))
		for i, e := range expr.Exprs {
			exprs[i] = lowerExpr(e, r, builder, true)
		}
		body := flow.CreateBlock([]flow.DubOp{
			&flow.ReturnOp{Exprs: exprs},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.RETURN))
		return flow.NoRegister

	case *tree.Fail:
		body := flow.CreateBlock([]flow.DubOp{
			&flow.Fail{},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.FAIL, r.GetExit(flow.FAIL))

		return flow.NoRegister

	case *tree.BinaryOp:
		left := lowerExpr(expr.Left, r, builder, true)
		right := lowerExpr(expr.Right, r, builder, true)
		dst := flow.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := flow.CreateBlock([]flow.DubOp{
			&flow.BinaryOp{
				Left:  left,
				Op:    expr.Op,
				Right: right,
				Dst:   dst,
			},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst
	case *tree.Append:
		l := lowerExpr(expr.List, r, builder, true)
		v := lowerExpr(expr.Expr, r, builder, true)
		dst := flow.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}

		body := flow.CreateBlock([]flow.DubOp{
			&flow.AppendOp{
				List:  l,
				Value: v,
				Dst:   dst,
			},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.Call:
		dst := flow.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := flow.CreateBlock([]flow.DubOp{
			&flow.CallOp{
				Name: expr.Name,
				Dst:  dst,
			},
		})
		r.Connect(flow.NORMAL, body)
		r.AttachDefaultExits(body)
		return dst
	case *tree.Construct:
		args := make([]*flow.KeyValue, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = &flow.KeyValue{
				Key:   arg.Name,
				Value: lowerExpr(arg.Expr, r, builder, true),
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
		body := flow.CreateBlock([]flow.DubOp{
			&flow.ConstructOp{
				Type: s,
				Args: args,
				Dst:  dst,
			},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.ConstructList:
		args := make([]flow.DubRegister, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, r, builder, true)
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
		body := flow.CreateBlock([]flow.DubOp{
			&flow.ConstructListOp{
				Type: l,
				Args: args,
				Dst:  dst,
			},
		})
		r.Connect(flow.NORMAL, body)
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.Coerce:
		t := builder.glbl.TranslateType(tree.ResolveType(expr.Type))
		src := lowerExpr(expr.Expr, r, builder, true)
		dst := flow.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := flow.CreateBlock([]flow.DubOp{
			&flow.CoerceOp{
				Src: src,
				T:   t,
				Dst: dst,
			},
		})
		r.Connect(flow.NORMAL, body)
		// TODO can coersion fail?
		body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		return dst

	case *tree.Slice:
		start := builder.CreateLLRegister(builder.glbl.Int)
		// HACK assume checkpoint is just the index
		{
			head := flow.CreateBlock([]flow.DubOp{
				&flow.Checkpoint{Dst: start},
			})
			r.Connect(flow.NORMAL, head)
			head.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		}
		block := lowerBlock(expr.Block, builder)
		r.Splice(flow.NORMAL, block)

		// Create a slice
		dst := flow.NoRegister
		if used {
			dst = builder.CreateLLRegister(builder.glbl.String)
		}
		{
			body := flow.CreateBlock([]flow.DubOp{
				&flow.Slice{Src: start, Dst: dst},
			})

			r.Connect(flow.NORMAL, body)
			body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
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
				head := flow.CreateBlock([]flow.DubOp{
					&flow.Checkpoint{Dst: start},
				})
				r.Connect(flow.NORMAL, head)
				head.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
			}
		}

		lowerMatch(expr.Match, r, builder)

		// Create a slice
		if used {
			dst = builder.CreateLLRegister(builder.glbl.String)
			body := flow.CreateBlock([]flow.DubOp{
				&flow.Slice{Src: start, Dst: dst},
			})

			r.Connect(flow.NORMAL, body)
			body.SetExit(flow.NORMAL, r.GetExit(flow.NORMAL))
		}
		return dst

	case *tree.RuneMatch:
		return lowerRuneMatch(expr.Match, r, builder)
	default:
		panic(expr)
	}

}

func lowerBlock(block []tree.ASTExpr, builder *DubBuilder) *base.Region {
	r := flow.CreateRegion()
	for _, expr := range block {
		lowerExpr(expr, r, builder, false)
	}
	return r
}

func LowerAST(decl *tree.FuncDecl, glbl *GlobalDubBuilder) *flow.LLFunc {
	builder := &DubBuilder{decl: decl, glbl: glbl}

	f := &flow.LLFunc{Name: decl.Name}
	types := make([]flow.DubType, len(decl.ReturnTypes))
	for i, node := range decl.ReturnTypes {
		types[i] = builder.glbl.TranslateType(tree.ResolveType(node))
	}
	f.ReturnTypes = types
	// Allocate register for locals
	builder.localMap = make([]flow.DubRegister, len(decl.Locals))
	for i, info := range decl.Locals {
		builder.localMap[i] = builder.CreateRegister(info.T)
	}
	f.Region = lowerBlock(decl.Block, builder)
	f.Region.GetExit(flow.RETURN).TransferEntries(f.Region.GetExit(flow.NORMAL))
	f.Registers = builder.registers
	return f
}

func LowerStruct(decl *tree.StructDecl, s *flow.LLStruct, gbuilder *GlobalDubBuilder) *flow.LLStruct {
	fields := []*flow.LLField{}
	var implements *flow.LLStruct
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
			Name: field.Name,
			T:    gbuilder.TranslateType(tree.ResolveType(field.Type)),
		})
	}
	*s = flow.LLStruct{
		Name:       decl.Name,
		Implements: implements,
		Fields:     fields,
	}
	return s
}

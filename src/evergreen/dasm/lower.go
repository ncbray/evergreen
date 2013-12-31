package dasm

import (
	"evergreen/base"
	"evergreen/dub"
	"evergreen/dubx"
)

type GlobalDubBuilder struct {
	Types  map[ASTType]dub.DubType
	String dub.DubType
	Rune   dub.DubType
	Int    dub.DubType
	Bool   dub.DubType
}

func (builder *GlobalDubBuilder) TranslateType(t ASTType) dub.DubType {
	switch t := t.(type) {
	case *StructDecl, *BuiltinType:
		dt, ok := builder.Types[t]
		if !ok {
			panic(t)
		}
		return dt
	case *ListType:
		parent := builder.TranslateType(t.Type)
		// TODO memoize
		return &dub.ListType{Type: parent}
	default:
		panic(t)
	}
}

type DubBuilder struct {
	decl      *FuncDecl
	registers []dub.RegisterInfo
	localMap  []dub.DubRegister
	glbl      *GlobalDubBuilder
}

func (builder *DubBuilder) CreateRegister(t ASTType) dub.DubRegister {
	return builder.CreateLLRegister(builder.glbl.TranslateType(t))
}

func (builder *DubBuilder) CreateLLRegister(t dub.DubType) dub.DubRegister {
	builder.registers = append(builder.registers, dub.RegisterInfo{T: t})
	return dub.DubRegister(len(builder.registers) - 1)
}

func (builder *DubBuilder) ZeroRegister(dst dub.DubRegister) dub.DubOp {
	info := builder.registers[dst]
	switch info.T.(type) {
	case *dub.LLStruct:
		return &dub.ConstantNilOp{Dst: dst}
	case *dub.RuneType:
		return &dub.ConstantRuneOp{Value: 0, Dst: dst}
	case *dub.StringType:
		return &dub.ConstantStringOp{Value: "", Dst: dst}
	case *dub.IntType:
		return &dub.ConstantIntOp{Value: 0, Dst: dst}
	case *dub.BoolType:
		return &dub.ConstantBoolOp{Value: false, Dst: dst}
	default:
		panic(info.T)
	}
}

func makeRuneSwitch(cond dub.DubRegister, op string, value rune, builder *DubBuilder) (*base.Node, *base.Node) {
	vreg := builder.CreateLLRegister(builder.glbl.Rune)
	make_value := dub.CreateBlock([]dub.DubOp{
		&dub.ConstantRuneOp{Value: value, Dst: vreg},
	})

	breg := builder.CreateLLRegister(builder.glbl.Bool)
	compare := dub.CreateBlock([]dub.DubOp{
		&dub.BinaryOp{
			Left:  cond,
			Op:    op,
			Right: vreg,
			Dst:   breg,
		},
	})

	decide := dub.CreateSwitch(breg)

	make_value.SetExit(dub.NORMAL, compare)
	compare.SetExit(dub.NORMAL, decide)

	return make_value, decide
}

func lowerRuneMatch(match *dubx.RuneMatch, r *base.Region, builder *DubBuilder) dub.DubRegister {
	// Read
	cond := builder.CreateLLRegister(builder.glbl.Rune)
	body := dub.CreateBlock([]dub.DubOp{
		&dub.Peek{Dst: cond},
	})
	r.Connect(dub.NORMAL, body)
	r.AttachDefaultExits(body)

	filters := dub.CreateRegion()

	onMatch := dub.NORMAL
	onNoMatch := dub.FAIL

	if match.Invert {
		onMatch, onNoMatch = onNoMatch, onMatch
	} else {
		// Fail by default.
		filters.GetExit(dub.NORMAL).TransferEntries(filters.GetExit(dub.FAIL))
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
		body = dub.CreateBlock([]dub.DubOp{
			&dub.Consume{},
		})
		filters.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, filters.GetExit(dub.NORMAL))
	}
	{
		// Make the fail official.
		body = dub.CreateBlock([]dub.DubOp{
			&dub.Fail{},
		})
		filters.Connect(dub.FAIL, body)
		body.SetExit(dub.FAIL, filters.GetExit(dub.FAIL))
	}
	r.Splice(0, filters)

	return cond
}

func lowerMatch(match dubx.TextMatch, r *base.Region, builder *DubBuilder) {
	switch match := match.(type) {
	case *dubx.RuneMatch:
		lowerRuneMatch(match, r, builder)
	case *dubx.MatchSequence:
		for _, child := range match.Matches {
			lowerMatch(child, r, builder)
		}
	case *dubx.MatchChoice:
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})
		r.Connect(dub.NORMAL, head)

		for i, child := range match.Matches {
			block := dub.CreateRegion()
			lowerMatch(child, block, builder)

			// Connect the head to the entry of this block
			entryNode := block.Head()
			entryEdge := block.Entry.GetExit(0)
			entryNode.ReplaceEntry(entryEdge, []*base.Edge{head.GetExit(0)})

			// Recover if not the last block.
			if i < len(match.Matches)-1 {
				head = dub.CreateBlock([]dub.DubOp{
					&dub.Recover{Src: checkpoint},
				})
				block.Connect(dub.FAIL, head)
			} else {
				head = nil
			}

			// Absorb the exits that have not been directed to head.
			r.AbsorbExits(block)
		}
	case *dubx.MatchRepeat:
		// HACK unroll
		for i := 0; i < match.Min; i++ {
			lowerMatch(match.Match, r, builder)
		}

		child := dub.CreateRegion()

		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})
		child.Connect(dub.NORMAL, head)
		head.SetExit(dub.NORMAL, child.GetExit(dub.NORMAL))

		// Handle the body
		lowerMatch(match.Match, child, builder)

		// Normal flow iterates
		child.GetExit(dub.NORMAL).TransferEntries(head)

		// Stop iterating on failure and recover
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Recover{Src: checkpoint},
			})

			child.Connect(dub.FAIL, body)
			body.SetExit(dub.NORMAL, child.GetExit(dub.NORMAL))
		}

		r.Splice(dub.NORMAL, child)
	default:
		panic(match)
	}
}

func lowerExpr(expr ASTExpr, r *base.Region, builder *DubBuilder, used bool) dub.DubRegister {
	switch expr := expr.(type) {
	case *If:
		// TODO Min
		//l := dub.CreateRegion()

		cond := lowerExpr(expr.Expr, r, builder, true)
		block := lowerBlock(expr.Block, builder)

		// TODO conditional
		decide := dub.CreateSwitch(cond)

		r.Connect(dub.NORMAL, decide)
		decide.SetExit(0, r.GetExit(dub.NORMAL))
		r.Splice(dub.NORMAL, block)
		decide.SetExit(1, r.GetExit(dub.NORMAL))

		return dub.NoRegister

	case *Repeat:
		// HACK unroll
		for i := 0; i < expr.Min; i++ {
			block := lowerBlock(expr.Block, builder)
			r.Splice(dub.NORMAL, block)
		}

		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})

		// Handle the body
		block := lowerBlock(expr.Block, builder)

		// Splice in the checkpoint as the first operation.
		// Note: block may be empty, but this code is carefully designed to work in that case.
		// Sure it's an infinite loop, but it would stange if that loop vanished.
		oldHead := block.Head()
		oldHead.TransferEntries(head)
		head.SetExit(dub.NORMAL, oldHead)

		// Normal flow iterates
		block.GetExit(dub.NORMAL).TransferEntries(head)

		// Stop iterating on failure and recover
		block.GetExit(dub.FAIL).TransferEntries(block.GetExit(dub.NORMAL))
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Recover{Src: checkpoint},
			})

			block.Connect(dub.NORMAL, body)
			body.SetExit(dub.NORMAL, block.GetExit(dub.NORMAL))
		}

		r.Splice(dub.NORMAL, block)

		return dub.NoRegister
	case *Choice:
		checkpoint := dub.NoRegister
		if len(expr.Blocks) > 1 {
			checkpoint = builder.CreateLLRegister(builder.glbl.Int)
		}
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})
		r.Connect(dub.NORMAL, head)

		for i, block := range expr.Blocks {
			block := lowerBlock(block, builder)

			// Connect the head to the entry of this block
			entryNode := block.Head()
			entryEdge := block.Entry.GetExit(0)
			entryNode.ReplaceEntry(entryEdge, []*base.Edge{head.GetExit(0)})

			// Recover if not the last block.
			if i < len(expr.Blocks)-1 {
				head = dub.CreateBlock([]dub.DubOp{
					&dub.Recover{Src: checkpoint},
				})
				block.Connect(dub.FAIL, head)
			} else {
				head = nil
			}

			// Absorb the exits that have not been directed to head.
			r.AbsorbExits(block)
		}
		return dub.NoRegister

	case *Optional:
		// Checkpoint
		checkpoint := builder.CreateLLRegister(builder.glbl.Int)
		head := dub.CreateBlock([]dub.DubOp{
			&dub.Checkpoint{Dst: checkpoint},
		})
		r.Connect(dub.NORMAL, head)
		head.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))

		block := lowerBlock(expr.Block, builder)

		restore := dub.CreateBlock([]dub.DubOp{
			&dub.Recover{Src: checkpoint},
		})

		block.Connect(dub.FAIL, restore)
		restore.SetExit(dub.NORMAL, block.GetExit(dub.NORMAL))

		r.Splice(dub.NORMAL, block)

		return dub.NoRegister

	case *GetName:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateRegister(builder.decl.Locals[expr.Info].T)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CopyOp{Src: builder.localMap[expr.Info], Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Assign:
		dst := builder.localMap[expr.Info]
		var op dub.DubOp
		if expr.Expr != nil {
			src := lowerExpr(expr.Expr, r, builder, true)
			op = &dub.CopyOp{Src: src, Dst: dst}
		} else {
			op = builder.ZeroRegister(dst)
		}
		body := dub.CreateBlock([]dub.DubOp{op})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *dubx.RuneLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Rune)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantRuneOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *dubx.StringLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.String)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantStringOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *dubx.IntLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Int)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantIntOp{Value: int64(expr.Value), Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *dubx.BoolLiteral:
		if !used {
			return dub.NoRegister
		}
		dst := builder.CreateLLRegister(builder.glbl.Bool)
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstantBoolOp{Value: expr.Value, Dst: dst},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Return:
		exprs := make([]dub.DubRegister, len(expr.Exprs))
		for i, e := range expr.Exprs {
			exprs[i] = lowerExpr(e, r, builder, true)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ReturnOp{Exprs: exprs},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.RETURN))
		return dub.NoRegister

	case *Fail:
		body := dub.CreateBlock([]dub.DubOp{
			&dub.Fail{},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.FAIL, r.GetExit(dub.FAIL))

		return dub.NoRegister

	case *BinaryOp:
		left := lowerExpr(expr.Left, r, builder, true)
		right := lowerExpr(expr.Right, r, builder, true)
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.BinaryOp{
				Left:  left,
				Op:    expr.Op,
				Right: right,
				Dst:   dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst
	case *Append:
		l := lowerExpr(expr.List, r, builder, true)
		v := lowerExpr(expr.Value, r, builder, true)
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}

		body := dub.CreateBlock([]dub.DubOp{
			&dub.AppendOp{
				List:  l,
				Value: v,
				Dst:   dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Call:
		dst := dub.NoRegister
		if used {
			dst = builder.CreateRegister(expr.T)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CallOp{
				Name: expr.Name,
				Dst:  dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		r.AttachDefaultExits(body)
		return dst
	case *Construct:
		args := make([]*dub.KeyValue, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = &dub.KeyValue{
				Key:   arg.Key,
				Value: lowerExpr(arg.Value, r, builder, true),
			}
		}
		t := builder.glbl.TranslateType(ResolveType(expr.Type))
		s, ok := t.(*dub.LLStruct)
		if !ok {
			panic(t)
		}
		dst := dub.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstructOp{
				Type: s,
				Args: args,
				Dst:  dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *ConstructList:
		args := make([]dub.DubRegister, len(expr.Args))
		for i, arg := range expr.Args {
			args[i] = lowerExpr(arg, r, builder, true)
		}
		t := builder.glbl.TranslateType(ResolveType(expr.Type))
		l, ok := t.(*dub.ListType)
		if !ok {
			panic(t)
		}
		dst := dub.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.ConstructListOp{
				Type: l,
				Args: args,
				Dst:  dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Coerce:
		t := builder.glbl.TranslateType(ResolveType(expr.Type))
		src := lowerExpr(expr.Expr, r, builder, true)
		dst := dub.NoRegister
		if used {
			dst = builder.CreateLLRegister(t)
		}
		body := dub.CreateBlock([]dub.DubOp{
			&dub.CoerceOp{
				Src: src,
				T:   t,
				Dst: dst,
			},
		})
		r.Connect(dub.NORMAL, body)
		// TODO can coersion fail?
		body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		return dst

	case *Slice:
		start := builder.CreateLLRegister(builder.glbl.Int)
		// HACK assume checkpoint is just the index
		{
			head := dub.CreateBlock([]dub.DubOp{
				&dub.Checkpoint{Dst: start},
			})
			r.Connect(dub.NORMAL, head)
			head.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		}
		block := lowerBlock(expr.Block, builder)
		r.Splice(dub.NORMAL, block)

		// Create a slice
		dst := dub.NoRegister
		if used {
			dst = builder.CreateLLRegister(builder.glbl.String)
		}
		{
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Slice{Src: start, Dst: dst},
			})

			r.Connect(dub.NORMAL, body)
			body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		}
		return dst

	case *StringMatch:
		dst := dub.NoRegister
		start := dub.NoRegister

		// Checkpoint
		if used {
			start = builder.CreateLLRegister(builder.glbl.Int)
			// HACK assume checkpoint is just the index
			{
				head := dub.CreateBlock([]dub.DubOp{
					&dub.Checkpoint{Dst: start},
				})
				r.Connect(dub.NORMAL, head)
				head.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
			}
		}

		lowerMatch(expr.Expr, r, builder)

		// Create a slice
		if used {
			dst = builder.CreateLLRegister(builder.glbl.String)
			body := dub.CreateBlock([]dub.DubOp{
				&dub.Slice{Src: start, Dst: dst},
			})

			r.Connect(dub.NORMAL, body)
			body.SetExit(dub.NORMAL, r.GetExit(dub.NORMAL))
		}
		return dst

	case *RuneMatch:
		return lowerRuneMatch(expr.Expr, r, builder)
	default:
		panic(expr)
	}

}

func lowerBlock(block []ASTExpr, builder *DubBuilder) *base.Region {
	r := dub.CreateRegion()
	for _, expr := range block {
		lowerExpr(expr, r, builder, false)
	}
	return r
}

func LowerAST(decl *FuncDecl, glbl *GlobalDubBuilder) *dub.LLFunc {
	builder := &DubBuilder{decl: decl, glbl: glbl}

	f := &dub.LLFunc{Name: decl.Name}
	types := make([]dub.DubType, len(decl.ReturnTypes))
	for i, node := range decl.ReturnTypes {
		types[i] = builder.glbl.TranslateType(ResolveType(node))
	}
	f.ReturnTypes = types
	// Allocate register for locals
	builder.localMap = make([]dub.DubRegister, len(decl.Locals))
	for i, info := range decl.Locals {
		builder.localMap[i] = builder.CreateRegister(info.T)
	}
	f.Region = lowerBlock(decl.Block, builder)
	f.Region.GetExit(dub.RETURN).TransferEntries(f.Region.GetExit(dub.NORMAL))
	f.Registers = builder.registers
	return f
}

func LowerStruct(decl *StructDecl, s *dub.LLStruct, gbuilder *GlobalDubBuilder) *dub.LLStruct {
	fields := []*dub.LLField{}
	var implements *dub.LLStruct
	if decl.Implements != nil {
		t := gbuilder.TranslateType(ResolveType(decl.Implements))
		var ok bool
		implements, ok = t.(*dub.LLStruct)
		if !ok {
			panic(t)
		}
	}
	for _, field := range decl.Fields {
		fields = append(fields, &dub.LLField{
			Name: field.Name,
			T:    gbuilder.TranslateType(ResolveType(field.Type)),
		})
	}
	*s = dub.LLStruct{
		Name:       decl.Name,
		Implements: implements,
		Fields:     fields,
	}
	return s
}

package flow

import (
	"evergreen/base"
)

func AddDef(reg DubRegister, node int, defuse *base.DefUseCollector) {
	if reg != NoRegister {
		defuse.AddDef(node, int(reg))
	}
}

func AddUse(reg DubRegister, node int, defuse *base.DefUseCollector) {
	defuse.AddUse(node, int(reg))
}

func collectDefUse(node *base.Node, defuse *base.DefUseCollector) {
	switch op := node.Data.(type) {
	case *DubEntry, *DubExit:
	case *Consume, *Fail:
	case *Checkpoint:
		AddDef(op.Dst, node.Name, defuse)
	case *Peek:
		AddDef(op.Dst, node.Name, defuse)
	case *LookaheadBegin:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantRuneOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantStringOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantIntOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantBoolOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantNilOp:
		AddDef(op.Dst, node.Name, defuse)
	case *CallOp:
		AddDef(op.Dst, node.Name, defuse)
	case *Slice:
		AddUse(op.Src, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *BinaryOp:
		AddUse(op.Left, node.Name, defuse)
		AddUse(op.Right, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *AppendOp:
		AddUse(op.List, node.Name, defuse)
		AddUse(op.Value, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *CopyOp:
		AddUse(op.Src, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *CoerceOp:
		AddUse(op.Src, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *Recover:
		AddUse(op.Src, node.Name, defuse)
	case *LookaheadEnd:
		AddUse(op.Src, node.Name, defuse)
	case *DubSwitch:
		AddUse(op.Cond, node.Name, defuse)
	case *ReturnOp:
		for _, arg := range op.Exprs {
			AddUse(arg, node.Name, defuse)
		}
	case *ConstructOp:
		for _, arg := range op.Args {
			AddUse(arg.Value, node.Name, defuse)
		}
		AddDef(op.Dst, node.Name, defuse)
	case *ConstructListOp:
		for _, arg := range op.Args {
			AddUse(arg, node.Name, defuse)
		}
		AddDef(op.Dst, node.Name, defuse)
	default:
		panic(op)
	}
}

type Transfer struct {
	Src DubRegister
	Dst DubRegister
}

type NameMap struct {
	names []map[int]int
	idoms []int
}

func (nm *NameMap) GetName(n int, v int) int {
	name, ok := nm.names[n][v]
	if !ok {
		if n == nm.idoms[n] {
			panic(v)
		}
		name = nm.GetName(nm.idoms[n], v)
		nm.names[n][v] = name
	}
	return name
}

func (nm *NameMap) SetName(n int, v int, name int) {
	nm.names[n][v] = name
}

func CreateNameMap(numNodes int, idoms []int) *NameMap {
	nm := &NameMap{
		names: make([]map[int]int, numNodes),
		idoms: idoms,
	}
	for i := 0; i < numNodes; i++ {
		nm.names[i] = map[int]int{}
	}
	return nm
}

type RegisterReallocator struct {
	decl *LLFunc
	info []RegisterInfo
	nm   *NameMap
	live base.LivenessOracle
}

func (r *RegisterReallocator) Allocate(n int, reg DubRegister) DubRegister {
	v := int(reg)
	name := len(r.info)
	r.info = append(r.info, r.decl.Registers[v])
	r.nm.SetName(n, v, name)
	return DubRegister(name)
}

func (r *RegisterReallocator) MakeOutput(n int, reg DubRegister) DubRegister {
	if reg != NoRegister && r.live.LiveAtExit(n, int(reg)) {
		return r.Allocate(n, reg)
	}
	return NoRegister
}

func (r *RegisterReallocator) Transfer(n int, reg DubRegister) DubRegister {
	v := int(reg)
	name, ok := r.nm.names[n][v]
	if ok {
		return DubRegister(name)
	} else {
		return r.Allocate(n, reg)
	}
}

func (r *RegisterReallocator) Get(n int, reg DubRegister) DubRegister {
	return DubRegister(r.nm.GetName(n, int(reg)))
}

func (r *RegisterReallocator) Set(n int, reg DubRegister, name DubRegister) {
	r.nm.SetName(n, int(reg), int(name))
}

func renameOp(n int, data interface{}, ra *RegisterReallocator) {
	switch op := data.(type) {
	case *DubEntry, *DubExit:
	case *Consume, *Fail:
	case *Checkpoint:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *Peek:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *LookaheadBegin:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *ConstantRuneOp:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *ConstantStringOp:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *ConstantIntOp:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *ConstantBoolOp:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *ConstantNilOp:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *CallOp:
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *Slice:
		op.Src = ra.Get(n, op.Src)
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *BinaryOp:
		op.Left = ra.Get(n, op.Left)
		op.Right = ra.Get(n, op.Right)
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *AppendOp:
		op.List = ra.Get(n, op.List)
		op.Value = ra.Get(n, op.Value)
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *CopyOp:
		// Copy propagation
		op.Src = ra.Get(n, op.Src)
		ra.Set(n, op.Dst, op.Src)
		op.Dst = op.Src
	case *CoerceOp:
		op.Src = ra.Get(n, op.Src)
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *Recover:
		op.Src = ra.Get(n, op.Src)
	case *LookaheadEnd:
		op.Src = ra.Get(n, op.Src)
	case *DubSwitch:
		op.Cond = ra.Get(n, op.Cond)
	case *ReturnOp:
		for i, arg := range op.Exprs {
			op.Exprs[i] = ra.Get(n, arg)
		}
	case *ConstructOp:
		for i, arg := range op.Args {
			op.Args[i].Value = ra.Get(n, arg.Value)
		}
		op.Dst = ra.MakeOutput(n, op.Dst)
	case *ConstructListOp:
		for i, arg := range op.Args {
			op.Args[i] = ra.Get(n, arg)
		}
		op.Dst = ra.MakeOutput(n, op.Dst)
	default:
		panic(data)
	}
}

func rename(decl *LLFunc, order []*base.Node, builder *base.SSIBuilder) {
	nm := CreateNameMap(len(order), builder.Idoms)
	ra := &RegisterReallocator{decl: decl, nm: nm, live: builder.Live}
	for _, node := range order {
		n := node.Name
		// Rename destinations of incoming transfers.
		if len(nm.names[n]) != 0 {
			panic(nm.names[n])
		}
		for i := 0; i < node.NumEntries(); i++ {
			e := node.GetEntry(i)
			data := e.Data
			transfers, ok := data.([]Transfer)
			if !ok {
				continue
			}
			for t := 0; t < len(transfers); t++ {
				transfers[t].Dst = ra.Transfer(n, transfers[t].Dst)
			}
		}

		// Rename op
		// NOTE: may overwrite incoming transfers.
		renameOp(n, node.Data, ra)

		// Rename outgoing transfers.
		for i := 0; i < node.NumExits(); i++ {
			e := node.GetExit(i)
			data := e.Data
			transfers, ok := data.([]Transfer)
			if !ok {
				continue
			}
			for t := 0; t < len(transfers); t++ {
				transfers[t].Src = ra.Get(n, transfers[t].Src)
			}
		}

	}
	//fmt.Println(decl.Name, len(decl.Registers), len(ra.info))
	decl.Registers = ra.info
}

func SSI(decl *LLFunc) {
	order := base.ReversePostorder(decl.Region)
	defuse := base.FindDefUse(order, len(decl.Registers), collectDefUse)
	live := base.FindLiveVars(order, defuse)

	builder := base.CreateSSIBuilder(decl.Region, order, live)
	for i := 0; i < len(decl.Registers); i++ {
		base.SSI(builder, i, defuse.VarDefAt[i])
	}

	//fmt.Println(decl.Name)
	//for i := 0; i < len(order); i++ {
	//	fmt.Println(i, live.LiveSet(i), builder.PhiFuncs[i])
	//}
	//fmt.Println()

	// Place the transfer functions on edges.
	for _, node := range order {
		for i := 0; i < node.NumExits(); i++ {
			dst := node.GetNext(i)
			if dst == nil {
				continue
			}
			phiFuncs := builder.PhiFuncs[dst.Name]
			if len(phiFuncs) == 0 {
				continue
			}
			transfers := []Transfer{}
			for _, v := range phiFuncs {
				transfers = append(transfers, Transfer{Src: DubRegister(v), Dst: DubRegister(v)})
			}
			node.GetExit(i).Data = transfers
		}
	}

	rename(decl, order, builder)
}

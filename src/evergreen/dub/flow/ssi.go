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

type NameMap struct {
	names     []map[int]int
	transfers []map[int]int
	idoms     []int
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
		names:     make([]map[int]int, numNodes),
		transfers: make([]map[int]int, numNodes),
		idoms:     idoms,
	}
	for i := 0; i < numNodes; i++ {
		nm.names[i] = map[int]int{}
		nm.transfers[i] = map[int]int{}
	}
	return nm
}

type RegisterReallocator struct {
	decl *LLFunc
	info []RegisterInfo
	nm   *NameMap
}

func (r *RegisterReallocator) Allocate(v int) int {
	name := len(r.info)
	r.info = append(r.info, r.decl.Registers[v])
	return name
}

func (r *RegisterReallocator) MakeOutput(n int, reg DubRegister) DubRegister {
	if reg != NoRegister {
		v := int(reg)
		name := r.Allocate(v)
		r.nm.SetName(n, v, name)
		return DubRegister(name)
	}
	return NoRegister
}

func (r *RegisterReallocator) Transfer(dst int, reg DubRegister) DubRegister {
	v := int(reg)
	name, ok := r.nm.transfers[dst][v]
	if !ok {
		name = r.Allocate(v)
		r.nm.transfers[dst][v] = name
		r.nm.SetName(dst, v, name)
	}
	return DubRegister(name)
}

func (r *RegisterReallocator) Get(n int, reg DubRegister) DubRegister {
	return DubRegister(r.nm.GetName(n, int(reg)))
}

func (r *RegisterReallocator) Set(n int, reg DubRegister, name DubRegister) {
	r.nm.SetName(n, int(reg), int(name))
}

func renameOp(node *base.Node, data interface{}, ra *RegisterReallocator) {
	n := node.Name
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
	case *TransferOp:
		for i, src := range op.Srcs {
			op.Srcs[i] = ra.Get(n, src)
		}
		// The destinations need to be consistent based on the target
		tgt := node.GetNext(0).Name
		for i, dst := range op.Dsts {
			op.Dsts[i] = ra.Transfer(tgt, dst)
		}
	default:
		panic(data)
	}
}

func rename(decl *LLFunc) {
	order := base.ReversePostorder(decl.Region)
	idoms := base.FindIdoms(order)

	nm := CreateNameMap(len(order), idoms)
	ra := &RegisterReallocator{decl: decl, nm: nm}
	for _, node := range order {
		renameOp(node, node.Data, ra)
		_, is_copy := node.Data.(*CopyOp)
		if is_copy {
			node.Remove()
		}
	}
	//fmt.Println(decl.Name, len(decl.Registers), len(ra.info))
	decl.Registers = ra.info
}

func killUnusedOutputs(n int, data interface{}, live base.LivenessOracle) {
	switch op := data.(type) {
	case *DubEntry, *DubExit:
	case *Consume, *Fail:
	case *Checkpoint:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *Peek:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *LookaheadBegin:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *ConstantRuneOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *ConstantStringOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *ConstantIntOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *ConstantBoolOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *ConstantNilOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *CallOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *Slice:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *BinaryOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *AppendOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *CopyOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *CoerceOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *Recover:
	case *LookaheadEnd:
	case *DubSwitch:
	case *ReturnOp:
	case *ConstructOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	case *ConstructListOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegister
		}
	default:
		panic(data)
	}
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
			op := &TransferOp{
				Srcs: make([]DubRegister, len(phiFuncs)),
				Dsts: make([]DubRegister, len(phiFuncs)),
			}
			for j, v := range phiFuncs {
				op.Srcs[j] = DubRegister(v)
				op.Dsts[j] = DubRegister(v)
			}
			n := CreateNode(op)
			n.InsertAt(0, node.GetExit(i))
		}
		// Do this while the order and liveness info are still good.
		op, ok := node.Data.(DubOp)
		if ok {
			killUnusedOutputs(node.Name, op, live)
			if IsNop(op) {
				node.Remove()
			}
		}
	}

	rename(decl)
}

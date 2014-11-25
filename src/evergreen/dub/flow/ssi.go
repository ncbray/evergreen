package flow

import (
	"evergreen/base"
)

func AddDef(reg RegisterInfo_Ref, node base.NodeID, defuse *base.DefUseCollector) {
	if reg != NoRegisterInfo {
		defuse.AddDef(node, int(reg))
	}
}

func AddUse(reg RegisterInfo_Ref, node base.NodeID, defuse *base.DefUseCollector) {
	if reg == NoRegisterInfo {
		panic("Tried to use non-existant register.")
	}
	defuse.AddUse(node, int(reg))
}

func collectDefUse(decl *LLFunc, node base.NodeID, op DubOp, defuse *base.DefUseCollector) {
	switch op := op.(type) {
	case *EntryOp:
		for _, p := range decl.Params {
			AddDef(p, node, defuse)
		}
	case *FlowExitOp, *ExitOp:
	case *Consume, *Fail:
	case *Checkpoint:
		AddDef(op.Dst, node, defuse)
	case *Peek:
		AddDef(op.Dst, node, defuse)
	case *LookaheadBegin:
		AddDef(op.Dst, node, defuse)
	case *ConstantRuneOp:
		AddDef(op.Dst, node, defuse)
	case *ConstantStringOp:
		AddDef(op.Dst, node, defuse)
	case *ConstantIntOp:
		AddDef(op.Dst, node, defuse)
	case *ConstantBoolOp:
		AddDef(op.Dst, node, defuse)
	case *ConstantNilOp:
		AddDef(op.Dst, node, defuse)
	case *CallOp:
		for _, arg := range op.Args {
			AddUse(arg, node, defuse)
		}
		for _, dst := range op.Dsts {
			AddDef(dst, node, defuse)
		}
	case *Slice:
		AddUse(op.Src, node, defuse)
		AddDef(op.Dst, node, defuse)
	case *BinaryOp:
		AddUse(op.Left, node, defuse)
		AddUse(op.Right, node, defuse)
		AddDef(op.Dst, node, defuse)
	case *AppendOp:
		AddUse(op.List, node, defuse)
		AddUse(op.Value, node, defuse)
		AddDef(op.Dst, node, defuse)
	case *CopyOp:
		AddUse(op.Src, node, defuse)
		AddDef(op.Dst, node, defuse)
	case *CoerceOp:
		AddUse(op.Src, node, defuse)
		AddDef(op.Dst, node, defuse)
	case *Recover:
		AddUse(op.Src, node, defuse)
	case *LookaheadEnd:
		AddUse(op.Src, node, defuse)
	case *SwitchOp:
		AddUse(op.Cond, node, defuse)
	case *ReturnOp:
		for _, arg := range op.Exprs {
			AddUse(arg, node, defuse)
		}
	case *ConstructOp:
		for _, arg := range op.Args {
			AddUse(arg.Value, node, defuse)
		}
		AddDef(op.Dst, node, defuse)
	case *ConstructListOp:
		for _, arg := range op.Args {
			AddUse(arg, node, defuse)
		}
		AddDef(op.Dst, node, defuse)
	default:
		panic(op)
	}
}

type NameMap struct {
	names     []map[int]int
	transfers []map[int]int
	idoms     []base.NodeID
}

func (nm *NameMap) GetName(n base.NodeID, v int) int {
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

func (nm *NameMap) SetName(n base.NodeID, v int, name int) {
	nm.names[n][v] = name
}

func CreateNameMap(numNodes int, idoms []base.NodeID) *NameMap {
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
	info []*RegisterInfo
	nm   *NameMap
}

func (r *RegisterReallocator) Allocate(v int) int {
	name := len(r.info)
	ref := RegisterInfo_Ref(v)
	r.info = append(r.info, r.decl.RegisterInfo_Scope.Get(ref))
	return name
}

func (r *RegisterReallocator) MakeOutput(n base.NodeID, reg RegisterInfo_Ref) RegisterInfo_Ref {
	if reg != NoRegisterInfo {
		v := int(reg)
		name := r.Allocate(v)
		r.nm.SetName(n, v, name)
		return RegisterInfo_Ref(name)
	}
	return NoRegisterInfo
}

func (r *RegisterReallocator) Transfer(dst base.NodeID, reg RegisterInfo_Ref) RegisterInfo_Ref {
	v := int(reg)
	name, ok := r.nm.transfers[dst][v]
	if !ok {
		name = r.Allocate(v)
		r.nm.transfers[dst][v] = name
		r.nm.SetName(dst, v, name)
	}
	return RegisterInfo_Ref(name)
}

func (r *RegisterReallocator) Get(n base.NodeID, reg RegisterInfo_Ref) RegisterInfo_Ref {
	return RegisterInfo_Ref(r.nm.GetName(n, int(reg)))
}

func (r *RegisterReallocator) Set(n base.NodeID, reg RegisterInfo_Ref, name RegisterInfo_Ref) {
	r.nm.SetName(n, int(reg), int(name))
}

func renameOp(n base.NodeID, data DubOp, ra *RegisterReallocator) {
	switch op := data.(type) {
	case *EntryOp, *FlowExitOp, *ExitOp:
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
		for i, arg := range op.Args {
			op.Args[i] = ra.Get(n, arg)
		}
		for i, dst := range op.Dsts {
			op.Dsts[i] = ra.MakeOutput(n, dst)
		}
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
	case *SwitchOp:
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
		tgt := ra.decl.CFG.GetExit(n, 0)
		for i, dst := range op.Dsts {
			op.Dsts[i] = ra.Transfer(tgt, dst)
		}
	default:
		panic(data)
	}
}

func rename(decl *LLFunc) {
	g := decl.CFG
	order, index := base.ReversePostorder(g)
	idoms := base.FindDominators(g, order, index)

	nm := CreateNameMap(g.NumNodes(), idoms)
	ra := &RegisterReallocator{decl: decl, nm: nm}

	// Define the function parameters.
	for i, p := range decl.Params {
		decl.Params[i] = ra.MakeOutput(g.Entry(), p)
	}

	nit := base.OrderedIterator(order)
	for nit.Next() {
		n := nit.Value()
		op := decl.Ops[n]
		renameOp(n, op, ra)
		_, is_copy := op.(*CopyOp)
		if is_copy {
			decl.CFG.Remove(n)
		}
	}
	decl.RegisterInfo_Scope.Replace(ra.info)
}

func killUnusedOutputs(n base.NodeID, op DubOp, live base.LivenessOracle) {
	switch op := op.(type) {
	case *EntryOp, *FlowExitOp, *ExitOp:
	case *Consume, *Fail:
	case *Checkpoint:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *Peek:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *LookaheadBegin:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *ConstantRuneOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *ConstantStringOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *ConstantIntOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *ConstantBoolOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *ConstantNilOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *CallOp:
		anyLive := false
		for i, dst := range op.Dsts {
			if !live.LiveAtExit(n, int(dst)) {
				op.Dsts[i] = NoRegisterInfo
			} else {
				anyLive = true
			}
		}
		if !anyLive {
			op.Dsts = []RegisterInfo_Ref{}
		}
	case *Slice:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *BinaryOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *AppendOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *CopyOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *CoerceOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *Recover:
	case *LookaheadEnd:
	case *SwitchOp:
	case *ReturnOp:
	case *ConstructOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	case *ConstructListOp:
		if !live.LiveAtExit(n, int(op.Dst)) {
			op.Dst = NoRegisterInfo
		}
	default:
		panic(op)
	}
}

func CreateTransfer(decl *LLFunc, size int) (base.NodeID, *TransferOp) {
	n := decl.CFG.CreateNode(1)
	if int(n) != len(decl.Ops) {
		panic("desync")
	}

	op := &TransferOp{
		Srcs: make([]RegisterInfo_Ref, size),
		Dsts: make([]RegisterInfo_Ref, size),
	}
	decl.Ops = append(decl.Ops, op)

	return n, op
}

func SSI(decl *LLFunc) {
	g := decl.CFG

	defuse := base.CreateDefUse(len(decl.Ops), decl.RegisterInfo_Scope.Len())
	nit := base.NodeIterator(g)
	for nit.Next() {
		n := nit.Value()
		collectDefUse(decl, n, decl.Ops[n], defuse)
	}

	live := base.FindLiveVars(g, defuse)

	builder := base.CreateSSIBuilder(g, live)
	for i := 0; i < decl.RegisterInfo_Scope.Len(); i++ {
		base.SSI(builder, i, defuse.VarDefAt[i])
	}

	//fmt.Println(decl.Name)
	//for i := 0; i < len(order); i++ {
	//	fmt.Println(i, live.LiveSet(i), builder.PhiFuncs[i])
	//}
	//fmt.Println()

	// Place the transfer functions on edges.

	nit = base.NodeIterator(g)
	for nit.Next() {
		n := nit.Value()
		eit := base.ExitIterator(g, n)
		for eit.Next() {
			dst := eit.Value()
			phiFuncs := builder.PhiFuncs[dst]
			if len(phiFuncs) == 0 {
				continue
			}
			t, op := CreateTransfer(decl, len(phiFuncs))
			for j, v := range phiFuncs {
				op.Srcs[j] = RegisterInfo_Ref(v)
				op.Dsts[j] = RegisterInfo_Ref(v)
			}
			g.InsertAt(t, 0, eit.Edge())
		}
		// Do this while the order and liveness info are still good.
		op := decl.Ops[n]
		killUnusedOutputs(n, op, live)
		if IsNop(op) {
			g.Remove(n)
		}
	}

	rename(decl)
}

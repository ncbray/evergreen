package flow

import (
	"evergreen/graph"
	"evergreen/ssi"
)

func addDef(reg *RegisterInfo, node graph.NodeID, defuse *ssi.DefUseCollector) {
	if reg != nil {
		defuse.AddDef(node, int(reg.Index))
	}
}

func addUse(reg *RegisterInfo, node graph.NodeID, defuse *ssi.DefUseCollector) {
	if reg == nil {
		panic("Tried to use non-existant register.")
	}
	defuse.AddUse(node, int(reg.Index))
}

func collectDefUse(decl *LLFunc, node graph.NodeID, op DubOp, defuse *ssi.DefUseCollector) {
	switch op := op.(type) {
	case *EntryOp:
		for _, p := range decl.Params {
			addDef(p, node, defuse)
		}
	case *ExitOp:
	case *Consume, *Fail:
	case *Checkpoint:
		addDef(op.Dst, node, defuse)
	case *Peek:
		addDef(op.Dst, node, defuse)
	case *LookaheadBegin:
		addDef(op.Dst, node, defuse)
	case *ConstantRuneOp:
		addDef(op.Dst, node, defuse)
	case *ConstantStringOp:
		addDef(op.Dst, node, defuse)
	case *ConstantIntOp:
		addDef(op.Dst, node, defuse)
	case *ConstantFloat32Op:
		addDef(op.Dst, node, defuse)
	case *ConstantBoolOp:
		addDef(op.Dst, node, defuse)
	case *ConstantNilOp:
		addDef(op.Dst, node, defuse)
	case *CallOp:
		for _, arg := range op.Args {
			addUse(arg, node, defuse)
		}
		for _, dst := range op.Dsts {
			addDef(dst, node, defuse)
		}
	case *Slice:
		addUse(op.Src, node, defuse)
		addDef(op.Dst, node, defuse)
	case *BinaryOp:
		addUse(op.Left, node, defuse)
		addUse(op.Right, node, defuse)
		addDef(op.Dst, node, defuse)
	case *CopyOp:
		addUse(op.Src, node, defuse)
		addDef(op.Dst, node, defuse)
	case *CoerceOp:
		addUse(op.Src, node, defuse)
		addDef(op.Dst, node, defuse)
	case *Recover:
		addUse(op.Src, node, defuse)
	case *LookaheadEnd:
		addUse(op.Src, node, defuse)
	case *SwitchOp:
		addUse(op.Cond, node, defuse)
	case *ReturnOp:
		for _, arg := range op.Exprs {
			addUse(arg, node, defuse)
		}
	case *ConstructOp:
		for _, arg := range op.Args {
			addUse(arg.Value, node, defuse)
		}
		addDef(op.Dst, node, defuse)
	case *ConstructListOp:
		for _, arg := range op.Args {
			addUse(arg, node, defuse)
		}
		addDef(op.Dst, node, defuse)
	default:
		panic(op)
	}
}

type NameMap struct {
	names     []map[RegisterInfo_Ref]*RegisterInfo
	transfers []map[RegisterInfo_Ref]*RegisterInfo
	idoms     []graph.NodeID
}

func (nm *NameMap) GetName(n graph.NodeID, reg *RegisterInfo) *RegisterInfo {
	newReg, ok := nm.names[n][reg.Index]
	if !ok {
		if n == nm.idoms[n] {
			panic(reg)
		}
		newReg = nm.GetName(nm.idoms[n], reg)
		nm.names[n][reg.Index] = newReg
	}
	return newReg
}

func (nm *NameMap) SetName(n graph.NodeID, reg *RegisterInfo, newReg *RegisterInfo) {
	nm.names[n][reg.Index] = newReg
}

func CreateNameMap(numNodes int, idoms []graph.NodeID) *NameMap {
	nm := &NameMap{
		names:     make([]map[RegisterInfo_Ref]*RegisterInfo, numNodes),
		transfers: make([]map[RegisterInfo_Ref]*RegisterInfo, numNodes),
		idoms:     idoms,
	}
	for i := 0; i < numNodes; i++ {
		nm.names[i] = map[RegisterInfo_Ref]*RegisterInfo{}
		nm.transfers[i] = map[RegisterInfo_Ref]*RegisterInfo{}
	}
	return nm
}

type RegisterReallocator struct {
	decl *LLFunc
	info []*RegisterInfo
	nm   *NameMap
}

func (r *RegisterReallocator) Allocate(reg *RegisterInfo) *RegisterInfo {
	newReg := &RegisterInfo{
		Name:  reg.Name,
		T:     reg.T,
		Index: RegisterInfo_Ref(len(r.info)),
	}
	r.info = append(r.info, newReg)
	return newReg
}

func (r *RegisterReallocator) MakeOutput(n graph.NodeID, reg *RegisterInfo) *RegisterInfo {
	if reg == nil {
		return nil
	}
	newReg := r.Allocate(reg)
	r.nm.SetName(n, reg, newReg)
	return newReg
}

func (r *RegisterReallocator) Transfer(dst graph.NodeID, reg *RegisterInfo) *RegisterInfo {
	newReg, ok := r.nm.transfers[dst][reg.Index]
	if ok {
		return newReg
	}
	newReg = r.Allocate(reg)
	r.nm.transfers[dst][reg.Index] = newReg
	r.nm.SetName(dst, reg, newReg)
	return newReg
}

func (r *RegisterReallocator) Get(n graph.NodeID, reg *RegisterInfo) *RegisterInfo {
	return r.nm.GetName(n, reg)
}

func (r *RegisterReallocator) Set(n graph.NodeID, reg *RegisterInfo, name *RegisterInfo) {
	r.nm.SetName(n, reg, name)
}

func renameOp(n graph.NodeID, data DubOp, ra *RegisterReallocator) {
	switch op := data.(type) {
	case *EntryOp, *ExitOp:
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
	case *ConstantFloat32Op:
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
		_, tgt := ra.decl.CFG.GetUniqueExit(n)
		for i, dst := range op.Dsts {
			op.Dsts[i] = ra.Transfer(tgt, dst)
		}
	default:
		panic(data)
	}
}

func rename(decl *LLFunc) {
	g := decl.CFG
	order, index := graph.ReversePostorder(g)
	idoms := graph.FindDominators(g, order, index)

	nm := CreateNameMap(g.NumNodes(), idoms)
	ra := &RegisterReallocator{decl: decl, nm: nm}

	// Define the function parameters.
	for i, p := range decl.Params {
		decl.Params[i] = ra.MakeOutput(g.Entry(), p)
	}

	nit := graph.OrderedIterator(order)
	for nit.HasNext() {
		n := nit.GetNext()
		op := decl.Ops[n]
		renameOp(n, op, ra)
		_, is_copy := op.(*CopyOp)
		if is_copy {
			decl.CFG.KillNode(n)
		}
	}
	decl.RegisterInfo_Scope.Replace(ra.info)
}

func deadAtExit(live ssi.LivenessOracle, n graph.NodeID, reg *RegisterInfo) bool {
	if reg == nil {
		return true
	}
	return !live.LiveAtExit(n, int(reg.Index))
}

func killUnusedOutputs(n graph.NodeID, op DubOp, live ssi.LivenessOracle) {
	switch op := op.(type) {
	case *EntryOp, *ExitOp:
	case *Consume, *Fail:
	case *Checkpoint:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *Peek:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *LookaheadBegin:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *ConstantRuneOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *ConstantStringOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *ConstantIntOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *ConstantFloat32Op:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *ConstantBoolOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *ConstantNilOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *CallOp:
		anyLive := false
		for i, dst := range op.Dsts {
			if !live.LiveAtExit(n, int(dst.Index)) {
				op.Dsts[i] = nil
			} else {
				anyLive = true
			}
		}
		if !anyLive {
			op.Dsts = []*RegisterInfo{}
		}
	case *Slice:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *BinaryOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *CopyOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *CoerceOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *Recover:
	case *LookaheadEnd:
	case *SwitchOp:
	case *ReturnOp:
	case *ConstructOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	case *ConstructListOp:
		if deadAtExit(live, n, op.Dst) {
			op.Dst = nil
		}
	default:
		panic(op)
	}
}

func createTransfer(decl *LLFunc, size int) (graph.NodeID, graph.EdgeID, *TransferOp) {
	op := &TransferOp{
		Srcs: make([]*RegisterInfo, size),
		Dsts: make([]*RegisterInfo, size),
	}
	n := AllocNode(decl, op)
	e := AllocEdge(decl, 0)
	decl.CFG.ConnectEdgeEntry(n, e)
	return n, e, op
}

func place(decl *LLFunc, builder *ssi.SSIBuilder, live *ssi.LiveVars) {
	g := decl.CFG
	// Place the transfer functions on edges.
	nit := g.NodeIterator()
	for nit.HasNext() {
		n := nit.GetNext()
		eit := g.ExitIterator(n)
		for eit.HasNext() {
			edge, dst := eit.GetNext()
			phiFuncs := builder.PhiFuncs[dst]
			if len(phiFuncs) == 0 {
				continue
			}
			_, te, op := createTransfer(decl, len(phiFuncs))
			for j, v := range phiFuncs {
				info := decl.RegisterInfo_Scope.Get(RegisterInfo_Ref(v))
				op.Srcs[j] = info
				op.Dsts[j] = info
			}
			g.InsertInEdge(te, edge)
		}
		// Do this while the order and liveness info are still good.
		op := decl.Ops[n]
		killUnusedOutputs(n, op, live)
		if IsNop(op) {
			g.KillNode(n)
		}
	}
}

func makeDefUse(decl *LLFunc) *ssi.DefUseCollector {
	defuse := ssi.CreateDefUse(len(decl.Ops), decl.RegisterInfo_Scope.Len())
	nit := decl.CFG.NodeIterator()
	for nit.HasNext() {
		n := nit.GetNext()
		collectDefUse(decl, n, decl.Ops[n], defuse)
	}
	return defuse
}

func SSI(decl *LLFunc) {
	g := decl.CFG

	defuse := makeDefUse(decl)

	live := ssi.FindLiveVars(g, defuse)

	builder := ssi.CreateSSIBuilder(g, live)
	for i := 0; i < decl.RegisterInfo_Scope.Len(); i++ {
		ssi.SSI(builder, i, defuse.VarDefAt[i])
	}

	place(decl, builder, live)
	rename(decl)
}

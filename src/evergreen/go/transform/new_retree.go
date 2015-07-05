package transform

import (
	"evergreen/go/core"
	"evergreen/go/flow"
	"evergreen/go/tree"
	"evergreen/graph"
	"fmt"
)

type retreeGo struct {
	Decl   *flow.FlowFunc
	LclMap []*tree.LocalInfo
}

func switchExits(decl *flow.FlowFunc, nid graph.NodeID) (graph.NodeID, graph.NodeID) {
	t := graph.NoNode
	f := graph.NoNode

	eit := decl.CFG.ExitIterator(nid)
	for eit.HasNext() {
		e, next := eit.GetNext()
		flowID := decl.Edges[e]
		switch flowID {
		case flow.COND_TRUE:
			if t != graph.NoNode {
				panic(e)
			}
			t = next
		case flow.COND_FALSE:
			if f != graph.NoNode {
				panic(e)
			}
			f = next
		default:
			panic(flowID)
		}
	}
	if t == graph.NoNode || f == graph.NoNode {
		panic(nid)
	}
	return t, f
}

func linearExit(decl *flow.FlowFunc, nid graph.NodeID) graph.NodeID {
	exit := graph.NoNode

	eit := decl.CFG.ExitIterator(nid)
	for eit.HasNext() {
		e, next := eit.GetNext()
		flowID := decl.Edges[e]
		switch flowID {
		case flow.NORMAL, flow.RETURN:
			if exit != graph.NoNode {
				panic(e)
			}
			exit = next
		default:
			panic(flowID)
		}
	}
	return exit
}

func opToStmts(decl *flow.FlowFunc, lcl_map []*tree.LocalInfo, nid graph.NodeID, block []tree.Stmt) ([]tree.Stmt, bool) {
	op := decl.Ops[nid]
	terminal := false
	switch op := op.(type) {
	case *flow.Entry:
		// TODO
	case *flow.Exit:
		// TODO
		block = append(block, &tree.Return{})
		terminal = true
	case *flow.ConstantString:
		block = append(block, scalarAssign(&tree.StringLiteral{
			Value: op.Value,
		}, lcl_map, op.Dst))
	case *flow.ConstantRune:
		block = append(block, scalarAssign(&tree.RuneLiteral{
			Value: op.Value,
		}, lcl_map, op.Dst))
	case *flow.ConstantInt:
		block = append(block, scalarAssign(&tree.IntLiteral{
			Value: int(op.Value),
		}, lcl_map, op.Dst))
	case *flow.ConstantFloat32:
		block = append(block, scalarAssign(&tree.Float32Literal{
			Value: op.Value,
		}, lcl_map, op.Dst))
	case *flow.ConstantBool:
		block = append(block, scalarAssign(&tree.BoolLiteral{
			Value: op.Value,
		}, lcl_map, op.Dst))
	case *flow.ConstantNil:
		block = append(block, scalarAssign(&tree.NilLiteral{}, lcl_map, op.Dst))
	case *flow.Call:
		f := op.Target
		block = append(block, multiAssign(&tree.Call{
			Expr: &tree.GetFunction{Func: f},
			Args: getLocalList(lcl_map, op.Args),
		}, lcl_map, op.Dsts))
	case *flow.MethodCall:
		// TODO simple IR
		block = append(block, multiAssign(&tree.Call{
			Expr: &tree.Selector{
				Expr: getLocal(lcl_map, op.Expr),
				Text: op.Name,
			},
			Args: getLocalList(lcl_map, op.Args),
		}, lcl_map, op.Dsts))
	case *flow.Transfer:
		srcs := []tree.Expr{}
		tgts := []tree.Target{}
		// SSA can cause registers to be transfered to themselves.  Filter this out.
		for i := 0; i < len(op.Srcs); i++ {
			src := op.Srcs[i]
			tgt := op.Dsts[i]
			if src != tgt {
				srcs = append(srcs, getLocal(lcl_map, src))
				tgts = append(tgts, setLocal(lcl_map, tgt))
			}
		}
		if len(srcs) > 0 {
			block = append(block, &tree.Assign{
				Sources: srcs,
				Op:      "=",
				Targets: tgts,
			})
		}
	case *flow.Coerce:
		block = append(block, scalarAssign(&tree.TypeCoerce{
			Type: tree.RefForType(op.Type),
			Expr: getLocal(lcl_map, op.Src),
		}, lcl_map, op.Dst))
	case *flow.Attr:
		block = append(block, scalarAssign(&tree.Selector{
			Expr: getLocal(lcl_map, op.Expr),
			Text: op.Name,
		}, lcl_map, op.Dst))
	case *flow.BinaryOp:
		block = append(block, scalarAssign(&tree.BinaryExpr{
			Left:  getLocal(lcl_map, op.Left),
			Op:    op.Op,
			Right: getLocal(lcl_map, op.Right),
		}, lcl_map, op.Dst))
	case *flow.ConstructStruct:
		args := make([]*tree.KeywordExpr, len(op.Args))
		for i, arg := range op.Args {
			args[i] = &tree.KeywordExpr{
				Name: arg.Name,
				Expr: getLocal(lcl_map, arg.Arg),
			}
		}
		var expr tree.Expr = &tree.StructLiteral{
			Type: tree.RefForType(op.Type),
			Args: args,
		}
		if op.AddrTaken {
			expr = &tree.UnaryExpr{
				Op:   "&",
				Expr: expr,
			}
		}
		block = append(block, scalarAssign(expr, lcl_map, op.Dst))
	case *flow.ConstructSlice:
		ref := tree.RefForType(op.Type)
		sref, ok := ref.(*tree.SliceRef)
		if !ok {
			panic(op.Type)
		}
		block = append(block, scalarAssign(&tree.ListLiteral{
			Type: sref,
			Args: getLocalList(lcl_map, op.Args),
		}, lcl_map, op.Dst))
	case *flow.Nop:
		// TODO
	case *flow.Return:
		block = append(block, &tree.Return{
			Args: getLocalList(lcl_map, op.Args),
		})
		terminal = true
	case *flow.Switch:
		tid, fid := switchExits(decl, nid)

		t := gotoBlock(int(tid))
		f := gotoBlock(int(fid))

		block = append(block, &tree.If{
			Cond: getLocal(lcl_map, op.Cond),
			T: &tree.Block{
				Body: []tree.Stmt{t},
			},
			F: &tree.Block{
				Body: []tree.Stmt{f},
			},
		})
		terminal = true // TODO correct ???
	default:
		panic(op)
	}
	return block, terminal
}

func retreeBlock(retree *retreeGo, nodes []graph.NodeID, block []tree.Stmt) []tree.Stmt {
	// HACK assume returns have been inserted before all exits.
	if nodes[0] == retree.Decl.CFG.Exit() {
		return block
	}

	// Entry block is not labeled.
	if nodes[0] != retree.Decl.CFG.Entry() {
		block = append(block, blockLabel(int(nodes[0])))
	}

	// Handle the nodes.
	terminal := false
	for _, nid := range nodes {
		block, terminal = opToStmts(retree.Decl, retree.LclMap, nid, block)
		if terminal {
			return block
		}
	}

	eit := retree.Decl.CFG.ExitIterator(nodes[len(nodes)-1])
	for eit.HasNext() {
		_, next := eit.GetNext()
		block = append(block, gotoBlock(int(next)))
		return block
	}
	return block
}

func retreeCluster(retree *retreeGo, cluster graph.Cluster, stmts []tree.Stmt) []tree.Stmt {
	switch cluster := cluster.(type) {
	case *graph.ClusterLeaf:
		stmts = retreeBlock(retree, cluster.Nodes, stmts)
	case *graph.ClusterLinear:
		for _, c := range cluster.Clusters {
			stmts = retreeCluster(retree, c, stmts)
		}
	case *graph.ClusterSwitch:
		for _, c := range cluster.Children {
			stmts = retreeCluster(retree, c, stmts)
		}
	case *graph.ClusterLoop:
		stmts = append(stmts, &tree.For{
			Block: &tree.Block{Body: retreeCluster(retree, cluster.Body, []tree.Stmt{})},
		})
	default:
		panic(cluster)
	}
	return stmts
}

func makeLocalMap(decl *flow.FlowFunc, funcDecl *tree.FuncDecl) []*tree.LocalInfo {
	numRegisters := decl.Register_Scope.Len()
	lclMap := make([]*tree.LocalInfo, numRegisters)
	for i := 0; i < numRegisters; i++ {
		info := decl.Register_Scope.Get(flow.Register_Ref(i))
		lclMap[i] = funcDecl.LocalInfo_Scope.Register(&tree.LocalInfo{
			Name: info.Name,
			T:    tree.RefForType(info.T),
		})
	}
	return lclMap
}

func makeFuncType(decl *flow.FlowFunc, lclMap []*tree.LocalInfo) *tree.FuncTypeRef {
	ft := &tree.FuncTypeRef{}

	// Translate parameters
	ft.Params = make([]*tree.Param, len(decl.Params))
	for i, info := range decl.Params {
		mapped := lclMap[info.Index]
		ft.Params[i] = &tree.Param{
			Name: mapped.Name,
			Type: tree.RefForType(info.T),
			Info: mapped,
		}
	}

	// Translate returns
	ft.Results = make([]*tree.Param, len(decl.Results))
	for i, info := range decl.Results {
		mapped := lclMap[info.Index]
		ft.Results[i] = &tree.Param{
			Name: mapped.Name,
			Type: tree.RefForType(info.T),
			Info: mapped,
		}
	}
	return ft
}

func NewRetreeFunc(coreProg *core.CoreProgram, f *core.Function, decl *flow.FlowFunc) *tree.FuncDecl {
	funcDecl := &tree.FuncDecl{
		Name:            f.Name,
		LocalInfo_Scope: &tree.LocalInfo_Scope{},
		Block:           &tree.Block{Body: []tree.Stmt{}},
	}

	lclMap := makeLocalMap(decl, funcDecl)

	// Translate receiver
	info := decl.Recv
	if info != nil {
		mapped := lclMap[info.Index]
		funcDecl.Recv = &tree.Param{
			Name: mapped.Name,
			Type: tree.RefForType(info.T),
			Info: mapped,
		}
	}

	funcDecl.Type = makeFuncType(decl, lclMap)

	// Don't reconstruct empty functions.
	_, first := decl.CFG.GetUniqueExit(decl.CFG.Entry())
	if first == decl.CFG.Exit() {
		return funcDecl
	}

	//tree.DumpDecl(RetreeFunc3(f, decl))
	//panic(f.Name)

	cluster := graph.MakeCluster(decl.CFG)
	retree := &retreeGo{
		Decl:   decl,
		LclMap: lclMap,
	}
	funcDecl.Block.Body = retreeCluster(retree, cluster, funcDecl.Block.Body)

	return funcDecl
}

// Start new

type GoASTBuilderExit interface {
}

type LinearExit struct {
	dst graph.NodeID
}

type BranchExit struct {
	cond tree.Expr
	t    graph.NodeID
	f    graph.NodeID
}

type DanglingExit struct {
	block *tree.Block
	dst   graph.NodeID
}

type MultiExit struct {
	exits []DanglingExit
}

func (self *MultiExit) Merge(other *MultiExit) GoASTBuilderExit {
	exits := make([]DanglingExit, len(self.exits)+len(other.exits))
	for i := 0; i < len(self.exits); i++ {
		exits[i] = self.exits[i]
	}
	for i := 0; i < len(other.exits); i++ {
		exits[len(self.exits)+i] = other.exits[i]
	}
	// TODO infer linear exit.
	return &MultiExit{exits: exits}
}

func singleExit(block *tree.Block, dst graph.NodeID) *MultiExit {
	return &MultiExit{
		exits: []DanglingExit{
			DanglingExit{
				block: block,
				dst:   dst,
			},
		},
	}
}

type GoASTBuilder struct {
	head     graph.NodeID
	funcDecl *tree.FuncDecl
	block    []tree.Stmt
	exit     GoASTBuilderExit
}

func (builder *GoASTBuilder) finalize() (*tree.Block, *MultiExit) {
	block := &tree.Block{Body: builder.block}
	switch exit := builder.exit.(type) {
	case *LinearExit:
		return block, &MultiExit{
			exits: []DanglingExit{
				DanglingExit{
					block: block,
					dst:   exit.dst,
				},
			},
		}
	case *BranchExit:
		tBlock := &tree.Block{Body: []tree.Stmt{}}
		fBlock := &tree.Block{Body: []tree.Stmt{}}
		block.Body = append(block.Body, &tree.If{
			Cond: exit.cond,
			T:    tBlock,
			F:    fBlock,
		})
		return block, &MultiExit{
			exits: []DanglingExit{
				DanglingExit{
					block: tBlock,
					dst:   exit.t,
				},
				DanglingExit{
					block: fBlock,
					dst:   exit.f,
				},
			},
		}
	case *MultiExit:
		return block, exit
	default:
		panic(exit)
	}
}

func (builder *GoASTBuilder) AppendChildren(children []ASTBuilder) {
	switch len(children) {
	case 1:
		c0 := children[0].(*GoASTBuilder)

		switch exit := builder.exit.(type) {
		case *LinearExit:
			if exit.dst != c0.head {
				panic(exit)
			}
			builder.block = append(builder.block, c0.block...)
			builder.exit = c0.exit

		case *MultiExit:
			for _, dangling := range exit.exits {
				if dangling.dst != c0.head {
					panic(c0.head)
				}
			}

			builder.block = append(builder.block, c0.block...)
			builder.exit = c0.exit
		default:
			panic(exit)
		}
	case 2:
		c0 := children[0].(*GoASTBuilder)
		c1 := children[1].(*GoASTBuilder)

		var cond tree.Expr

		switch exit := builder.exit.(type) {
		case *BranchExit:
			// HACK assume order is preserved
			if c0.head != exit.t || c1.head != exit.f {
				panic(exit)
			}
			cond = exit.cond

		case *MultiExit:
			lclInfo := builder.funcDecl.LocalInfo_Scope.Register(&tree.LocalInfo{
				Name: "next",
				T:    &tree.NameRef{Name: "int"},
			})

			for _, dangling := range exit.exits {
				label := -1
				if dangling.dst == c0.head {
					label = 1
				} else if dangling.dst == c1.head {
					label = 0
				} else {
					panic(c0.head)
				}
				dangling.block.Body = append(dangling.block.Body, &tree.Assign{
					Targets: []tree.Target{
						&tree.SetLocal{Info: lclInfo},
					},
					Op: "=",
					Sources: []tree.Expr{
						&tree.IntLiteral{Value: label},
					},
				})
			}

			cond = &tree.GetLocal{Info: lclInfo}
		default:
			panic(exit)
		}

		tBlock, tExit := c0.finalize()
		fBlock, fExit := c1.finalize()

		builder.block = append(builder.block, &tree.If{
			Cond: cond,
			T:    tBlock,
			F:    fBlock,
		})

		builder.exit = tExit.Merge(fExit)

	default:
		panic(len(children))
	}
}

type GoASTBuilderFactory struct {
	decl     *flow.FlowFunc
	funcDecl *tree.FuncDecl
	lclMap   []*tree.LocalInfo
}

func (factory *GoASTBuilderFactory) CreateInitial(n graph.NodeID) ASTBuilder {
	builder := &GoASTBuilder{head: n, funcDecl: factory.funcDecl}
	switch op := factory.decl.Ops[n].(type) {
	case *flow.Switch:
		cond := getLocal(factory.lclMap, op.Cond)
		t, f := switchExits(factory.decl, n)
		builder.exit = &BranchExit{cond: cond, t: t, f: f}
	case *flow.Exit:
		builder.exit = &MultiExit{}
	default:
		builder.block, _ = opToStmts(factory.decl, factory.lclMap, n, []tree.Stmt{})
		dst := linearExit(factory.decl, n)
		if dst == graph.NoNode || dst == factory.decl.CFG.Exit() {
			builder.exit = &MultiExit{}
		} else {
			builder.exit = &LinearExit{dst: dst}
		}
	}
	return builder
}

func (factory *GoASTBuilderFactory) Placeholder(n graph.NodeID) ASTBuilder {
	return &GoASTBuilder{head: n, funcDecl: factory.funcDecl, exit: &LinearExit{dst: n}}
}

type ASTBuilder interface {
	AppendChildren(children []ASTBuilder)
}

type ASTBuilderFactory interface {
	CreateInitial(n graph.NodeID) ASTBuilder
	Placeholder(n graph.NodeID) ASTBuilder
}

type nodeSequenceInfo struct {
	builder           ASTBuilder
	incomingEdgeCount int
	edgeExists        bool
	consume           bool
}

type retreeSequencer struct {
	cfg      *graph.Graph
	nodeInfo []graph.NodeInfo
	edgeType []graph.EdgeType

	factory ASTBuilderFactory

	nodeSequence []nodeSequenceInfo

	current graph.NodeID
	loop    graph.NodeID

	edgeCount  int
	shouldStop bool
}

func (seq *retreeSequencer) Init(n graph.NodeID, loop graph.NodeID) {
	fmt.Println("init", n, loop)
	seq.current = n
	seq.loop = loop
	seq.edgeCount = 0
	seq.shouldStop = false

	seq.nodeSequence[n].builder = seq.factory.CreateInitial(n)

	seq.prepChild(n, n)
}

func (seq *retreeSequencer) prepChild(n graph.NodeID, child graph.NodeID) {
	init := n == child

	xit := seq.cfg.ExitIterator(child)
	for xit.HasNext() {
		e, dst := xit.GetNext()

		// Clean up redundant edges
		if seq.nodeSequence[dst].edgeExists {
			fmt.Println("exists", n, child, dst)
			if !init {
				seq.nodeSequence[dst].incomingEdgeCount -= 1
			}
			seq.cfg.KillEdge(e)
			continue
		}

		fmt.Println("?", n, child, dst, graph.Dominates(seq.nodeInfo, n, dst))
		seq.edgeCount += 1
		seq.nodeSequence[dst].edgeExists = true

		if init {
			seq.nodeSequence[dst].incomingEdgeCount += 1
		}

		// If this node has a dominator cross edge in this loop, stop processing.
		if seq.nodeInfo[dst].LoopHead == seq.loop && !graph.Dominates(seq.nodeInfo, n, dst) {
			seq.shouldStop = true
		}
	}
}

func (seq *retreeSequencer) Done(n graph.NodeID) {
	cfg := seq.cfg
	xit := cfg.ExitIterator(n)
	for xit.HasNext() {
		_, dst := xit.GetNext()
		fmt.Println("Unmarking", dst)
		// Clean up
		if !seq.nodeSequence[dst].edgeExists {
			panic(n)
		}
		seq.nodeSequence[dst].edgeExists = false
		seq.edgeCount -= 1
	}
	if seq.edgeCount != 0 {
		panic(n)
	}
}

func (seq *retreeSequencer) Contract() bool {
	if seq.shouldStop || seq.edgeCount == 0 {
		return false
	}

	cfg := seq.cfg
	n := seq.current

	fmt.Println(">>>>", n, seq.edgeCount)

	children := []ASTBuilder{}

	xit := cfg.ExitIterator(n)
	for xit.HasNext() {
		_, dst := xit.GetNext()
		seq.nodeSequence[dst].consume = false

		if seq.nodeInfo[dst].LoopHead != seq.loop {
			fmt.Println(n, dst, "out of loop")
			panic(n)
		} else if seq.nodeSequence[dst].incomingEdgeCount == 1 {
			fmt.Println(n, dst, "consume")
			seq.nodeSequence[dst].consume = true
			children = append(children, seq.nodeSequence[dst].builder)
			seq.nodeSequence[dst].builder = nil
		} else {
			fmt.Println(n, dst, "defer")
			children = append(children, seq.factory.Placeholder(dst))
		}
	}

	current := seq.nodeSequence[n].builder

	current.AppendChildren(children)

	// Contract the graph
	xit = cfg.ExitIterator(n)
	for xit.HasNext() {
		e, dst := xit.GetNext()
		if seq.nodeSequence[dst].consume {
			seq.prepChild(n, dst)
			cfg.ReplaceEdgeWithExits(e, dst)
			seq.edgeCount -= 1
		}
	}

	fmt.Println("<<<<", n)

	return true
}

func RetreeFunc3(f *core.Function, decl *flow.FlowFunc) *tree.FuncDecl {
	funcDecl := &tree.FuncDecl{
		Name:            f.Name,
		LocalInfo_Scope: &tree.LocalInfo_Scope{},
		Block:           &tree.Block{Body: []tree.Stmt{}},
	}

	lclMap := makeLocalMap(decl, funcDecl)
	funcDecl.Type = makeFuncType(decl, lclMap)

	cfg := decl.CFG
	nodes, edges, postorder := graph.AnalyzeStructure(cfg)

	// Eliminate edges to the exit node, flow does not need to remerge to exit the function.
	// If the exit node is not disconnected, restructuring will try to re-merge after return, etc.
	cfg = cfg.Copy()
	eit := cfg.EntryIterator(cfg.Exit())
	for eit.HasNext() {
		_, e := eit.GetNext()
		cfg.KillEdge(e)
	}

	seq := &retreeSequencer{
		cfg:          cfg,
		nodeInfo:     nodes,
		edgeType:     edges,
		nodeSequence: make([]nodeSequenceInfo, len(nodes)),
		factory:      &GoASTBuilderFactory{decl: decl, funcDecl: funcDecl, lclMap: lclMap},
	}

	for _, n := range postorder {
		if seq.nodeInfo[n].IsHead {
			panic(n)
		}

		seq.Init(n, seq.nodeInfo[n].LoopHead)
		for seq.Contract() {
		}
		seq.Done(n)
		fmt.Println()
	}

	funcDecl.Block.Body = seq.nodeSequence[decl.CFG.Entry()].builder.(*GoASTBuilder).block

	return funcDecl
}

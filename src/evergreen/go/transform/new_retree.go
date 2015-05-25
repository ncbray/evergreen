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
	for _, nid := range nodes {
		op := retree.Decl.Ops[nid]
		switch op := op.(type) {
		case *flow.Entry:
			// TODO
		case *flow.Exit:
			// TODO
			block = append(block, &tree.Return{})
			return block
		case *flow.ConstantString:
			block = append(block, scalarAssign(&tree.StringLiteral{
				Value: op.Value,
			}, retree.LclMap, op.Dst))
		case *flow.ConstantRune:
			block = append(block, scalarAssign(&tree.RuneLiteral{
				Value: op.Value,
			}, retree.LclMap, op.Dst))
		case *flow.ConstantInt:
			block = append(block, scalarAssign(&tree.IntLiteral{
				Value: int(op.Value),
			}, retree.LclMap, op.Dst))
		case *flow.ConstantFloat32:
			block = append(block, scalarAssign(&tree.Float32Literal{
				Value: op.Value,
			}, retree.LclMap, op.Dst))
		case *flow.ConstantBool:
			block = append(block, scalarAssign(&tree.BoolLiteral{
				Value: op.Value,
			}, retree.LclMap, op.Dst))
		case *flow.ConstantNil:
			block = append(block, scalarAssign(&tree.NilLiteral{}, retree.LclMap, op.Dst))
		case *flow.Call:
			f := op.Target
			block = append(block, multiAssign(&tree.Call{
				Expr: &tree.GetFunction{Func: f},
				Args: getLocalList(retree.LclMap, op.Args),
			}, retree.LclMap, op.Dsts))
		case *flow.MethodCall:
			// TODO simple IR
			block = append(block, multiAssign(&tree.Call{
				Expr: &tree.Selector{
					Expr: getLocal(retree.LclMap, op.Expr),
					Text: op.Name,
				},
				Args: getLocalList(retree.LclMap, op.Args),
			}, retree.LclMap, op.Dsts))
		case *flow.Transfer:
			srcs := []tree.Expr{}
			tgts := []tree.Target{}
			// SSA can cause registers to be transfered to themselves.  Filter this out.
			for i := 0; i < len(op.Srcs); i++ {
				src := op.Srcs[i]
				tgt := op.Dsts[i]
				if src != tgt {
					srcs = append(srcs, getLocal(retree.LclMap, src))
					tgts = append(tgts, setLocal(retree.LclMap, tgt))
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
				Expr: getLocal(retree.LclMap, op.Src),
			}, retree.LclMap, op.Dst))
		case *flow.Attr:
			block = append(block, scalarAssign(&tree.Selector{
				Expr: getLocal(retree.LclMap, op.Expr),
				Text: op.Name,
			}, retree.LclMap, op.Dst))
		case *flow.BinaryOp:
			block = append(block, scalarAssign(&tree.BinaryExpr{
				Left:  getLocal(retree.LclMap, op.Left),
				Op:    op.Op,
				Right: getLocal(retree.LclMap, op.Right),
			}, retree.LclMap, op.Dst))
		case *flow.ConstructStruct:
			args := make([]*tree.KeywordExpr, len(op.Args))
			for i, arg := range op.Args {
				args[i] = &tree.KeywordExpr{
					Name: arg.Name,
					Expr: getLocal(retree.LclMap, arg.Arg),
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
			block = append(block, scalarAssign(expr, retree.LclMap, op.Dst))
		case *flow.ConstructSlice:
			ref := tree.RefForType(op.Type)
			sref, ok := ref.(*tree.SliceRef)
			if !ok {
				panic(op.Type)
			}
			block = append(block, scalarAssign(&tree.ListLiteral{
				Type: sref,
				Args: getLocalList(retree.LclMap, op.Args),
			}, retree.LclMap, op.Dst))
		case *flow.Nop:
			// TODO
		case *flow.Return:
			block = append(block, &tree.Return{
				Args: getLocalList(retree.LclMap, op.Args),
			})
			return block
		case *flow.Switch:
			var t tree.Stmt
			var f tree.Stmt

			eit := retree.Decl.CFG.ExitIterator(nid)
			for eit.HasNext() {
				e, next := eit.GetNext()
				flowID := retree.Decl.Edges[e]
				switch flowID {
				case flow.COND_TRUE:
					t = gotoBlock(int(next))
				case flow.COND_FALSE:
					f = gotoBlock(int(next))
				default:
					panic(flowID)
				}
			}
			if t == nil || f == nil {
				panic(op)
			}
			block = append(block, &tree.If{
				Cond: getLocal(retree.LclMap, op.Cond),
				Body: []tree.Stmt{t},
				Else: &tree.BlockStmt{
					Body: []tree.Stmt{f},
				},
			})
			return block
		default:
			panic(op)
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
			Body: retreeCluster(retree, cluster.Body, []tree.Stmt{}),
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
		Body:            []tree.Stmt{},
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

	//RetreeFunc3(decl)
	//panic(f.Name)

	cluster := graph.MakeCluster(decl.CFG)
	retree := &retreeGo{
		Decl:   decl,
		LclMap: lclMap,
	}
	funcDecl.Body = retreeCluster(retree, cluster, funcDecl.Body)

	return funcDecl
}

// Start new

type GoASTBuilder struct {
	head graph.NodeID
}

type GoASTBuilderFactory struct {
	cfg *graph.Graph
}

func (factory *GoASTBuilderFactory) CreateInitial(n graph.NodeID) ASTBuilder {
	return &GoASTBuilder{head: n}
}

type ASTBuilder interface {
}

type ASTBuilderFactory interface {
	CreateInitial(n graph.NodeID) ASTBuilder
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
		} else {
			fmt.Println(n, dst, "defer")
			panic(n)
		}
	}

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

func RetreeFunc3(decl *flow.FlowFunc) {
	nodes, edges, postorder := graph.AnalyzeStructure(decl.CFG)
	seq := &retreeSequencer{
		cfg:          decl.CFG.Copy(),
		nodeInfo:     nodes,
		edgeType:     edges,
		nodeSequence: make([]nodeSequenceInfo, len(nodes)),
		factory:      &GoASTBuilderFactory{},
	}

	for _, n := range postorder {
		if seq.nodeInfo[n].IsHead {
			panic(n)
		}

		seq.Init(n, seq.nodeInfo[n].LoopHead)
		for seq.Contract() {
		}
		fmt.Println()
	}
}

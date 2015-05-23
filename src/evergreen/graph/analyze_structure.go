package graph

type NodeInfo struct {
	Pre           int
	Post          int
	IDom          NodeID
	LoopHead      NodeID
	IsHead        bool
	IsIrreducible bool
	Live          bool
}

type EdgeType int

const (
	DEAD EdgeType = iota
	FORWARD
	BACKWARD
	CROSS
	REENTRY
)

type loopFinder struct {
	graph     *Graph
	node      []NodeInfo
	edge      []EdgeType
	postorder []NodeID
	depth     []int
	current   int
}

func (lf *loopFinder) beginTraversingNode(n NodeID, prev NodeID) {
	lf.current += 1

	lf.node[n] = NodeInfo{
		Pre:      lf.current,
		IDom:     prev, // A reasonable inital guess.
		LoopHead: NoNode,
		Live:     true,
	}

	lf.depth[n] = lf.current
}

func (lf *loopFinder) endTraversingNode(n NodeID) {
	lf.current += 1
	lf.node[n].Post = lf.current

	lf.depth[n] = 0

	lf.postorder = append(lf.postorder, n)
}

func (lf *loopFinder) isUnprocessed(n NodeID) bool {
	return !lf.node[n].Live
}

func (lf *loopFinder) isBeingProcessed(dst NodeID) bool {
	return lf.depth[dst] > 0
}

func (lf *loopFinder) markLoopHeader(child NodeID, head NodeID) {
	if child == head {
		// Trivial loop.
		return
	}
	childHead := lf.node[child].LoopHead
	for childHead != NoNode {
		if childHead == head {
			return
		}
		if lf.depth[childHead] < lf.depth[head] {
			// Found a closer head for this child, adopt it.
			lf.node[child].LoopHead = head
			child, head = head, childHead
		} else {
			child = childHead
		}
		childHead = lf.node[child].LoopHead
	}
	lf.node[child].LoopHead = head
}

func (lf *loopFinder) process(n NodeID, prev NodeID) {
	lf.beginTraversingNode(n, prev)
	xit := lf.graph.ExitIterator(n)
	for xit.HasNext() {
		e, next := xit.GetNext()
		if lf.isUnprocessed(next) {
			lf.edge[e] = FORWARD
			lf.process(next, n)

			// Propagage loop headers upwards.
			head := lf.node[next].LoopHead
			if head != NoNode {
				lf.markLoopHeader(n, head)
			}
		} else if lf.isBeingProcessed(next) {
			lf.edge[e] = BACKWARD
			lf.node[next].IsHead = true
			lf.markLoopHeader(n, next)
		} else {
			lf.edge[e] = CROSS
			if lf.node[next].LoopHead != NoNode {
				// Propagate loop header from cross edge.
				otherHead := lf.node[next].LoopHead
				if lf.isBeingProcessed(otherHead) {
					lf.markLoopHeader(n, otherHead)
				} else {
					lf.edge[e] = REENTRY
					lf.node[otherHead].IsIrreducible = true
					// Find and mark the common loop head.
					otherHead = lf.node[otherHead].LoopHead
					for otherHead != NoNode {
						if lf.isBeingProcessed(otherHead) {
							lf.markLoopHeader(n, otherHead)
							break
						}
						otherHead = lf.node[otherHead].LoopHead
					}
				}
			}
		}
	}
	lf.endTraversingNode(n)
}

func Intersect(nodes []NodeInfo, n0 NodeID, n1 NodeID) NodeID {
	i0 := nodes[n0].Post
	i1 := nodes[n1].Post
	for i0 != i1 {
		for i0 < i1 {
			n0 = nodes[n0].IDom
			i0 = nodes[n0].Post
		}
		for i0 > i1 {
			n1 = nodes[n1].IDom
			i1 = nodes[n1].Post
		}
	}
	return n0
}

func (lf *loopFinder) intersect(n0 NodeID, n1 NodeID) NodeID {
	return Intersect(lf.node, n0, n1)
}

func (lf *loopFinder) cleanDeadEdges() {
	numEdges := lf.graph.NumEdges()
	for i := 0; i < numEdges; i++ {
		if lf.edge[i] == DEAD {
			lf.graph.KillEdge(EdgeID(i))
		}
	}
}

func (lf *loopFinder) findIdoms() {
	g := lf.graph
	changed := true
	for changed {
		changed = false
		for i := len(lf.postorder) - 1; i >= 0; i-- {
			n := lf.postorder[i]
			original := lf.node[n].IDom
			idom := original
			eit := g.EntryIterator(n)
			for eit.HasNext() {
				src, _ := eit.GetNext()
				idom = lf.intersect(idom, src)
			}
			if idom != original {
				lf.node[n].IDom = idom
				changed = true
			}
		}
	}
}

func AnalyzeStructure(g *Graph) ([]NodeInfo, []EdgeType, []NodeID) {
	numNodes := g.NumNodes()
	lf := &loopFinder{
		graph: g,
		node:  make([]NodeInfo, numNodes),
		edge:  make([]EdgeType, g.NumEdges()),
		depth: make([]int, numNodes),
	}

	e := g.Entry()
	lf.process(e, e)
	lf.cleanDeadEdges()
	lf.findIdoms()
	return lf.node, lf.edge, lf.postorder
}

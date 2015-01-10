package graph

import (
	"fmt"
	"strings"
)

type Cluster interface {
	isCluster()
	Dump(margin string)
	DumpShort() string
}

type ClusterLeaf struct {
	Head  NodeID
	Nodes []NodeID
}

func (cluster *ClusterLeaf) isCluster() {
}

func (cluster *ClusterLeaf) Dump(margin string) {
	fmt.Printf("%sleaf %d\n", margin, cluster.Head)
	childMargin := margin + ".   "
	for _, n := range cluster.Nodes {
		fmt.Printf("%s%d\n", childMargin, n)
	}
}

func (cluster *ClusterLeaf) DumpShort() string {
	text := make([]string, len(cluster.Nodes))
	for i, n := range cluster.Nodes {
		text[i] = fmt.Sprintf("%d", n)
	}
	return fmt.Sprintf("(%s)", strings.Join(text, " "))
}

type ClusterLinear struct {
	Head     NodeID
	Clusters []Cluster
}

func (cluster *ClusterLinear) isCluster() {
}

func (cluster *ClusterLinear) Dump(margin string) {
	fmt.Printf("%slinear %d\n", margin, cluster.Head)
	childMargin := margin + ".   "
	for _, c := range cluster.Clusters {
		c.Dump(childMargin)
	}
}

func (cluster *ClusterLinear) DumpShort() string {
	text := make([]string, len(cluster.Clusters))
	for i, c := range cluster.Clusters {
		text[i] = c.DumpShort()
	}
	return fmt.Sprintf("[%s]", strings.Join(text, " "))
}

type ClusterSwitch struct {
	Head     NodeID
	Children []Cluster
}

func (cluster *ClusterSwitch) isCluster() {
}

func (cluster *ClusterSwitch) Dump(margin string) {
	fmt.Printf("%sswitch %d\n", margin, cluster.Head)
	childMargin := margin + ".   "
	for _, c := range cluster.Children {
		c.Dump(childMargin)
	}
}

func (cluster *ClusterSwitch) DumpShort() string {
	text := make([]string, len(cluster.Children))
	for i, c := range cluster.Children {
		text[i] = c.DumpShort()
	}
	return fmt.Sprintf("<%s>", strings.Join(text, " "))
}

type ClusterLoop struct {
	Head NodeID
	Body Cluster
}

func (cluster *ClusterLoop) isCluster() {
}

func (cluster *ClusterLoop) Dump(margin string) {
	fmt.Printf("%sloop %d\n", margin, cluster.Head)
	childMargin := margin + ".   "
	cluster.Body.Dump(childMargin)
}

func (cluster *ClusterLoop) DumpShort() string {
	return fmt.Sprintf("{%s}", cluster.Body.DumpShort())
}

type Candidate struct {
	Node        NodeID
	Cluster     Cluster
	CrossEdgeIn bool
}

func isUniqueSource(g *Graph, src NodeID, dst NodeID, idoms []NodeID) bool {
	it := g.EntryIterator(dst)
	for it.HasNext() {
		n, e := it.GetNext()
		if idoms[n] == NoNode {
			// This edge is actually dead.
			g.KillEdge(e)
			continue
		}
		if n != src {
			return false
		}
	}
	return true
}

type nodeInfo struct {
	cluster            Cluster
	incomingCrossEdges int
	crossEdgeSource    bool
	mergePoint         bool
	splitPoint         bool
}

type edgeInfo struct {
	crossEdge bool
	backEdge  bool
}

type clusterBuilder struct {
	graph    *Graph
	idoms    []NodeID
	nodeInfo []nodeInfo
	edgeInfo []edgeInfo

	currentHead    NodeID
	currentCluster Cluster
	ready          []Candidate

	isLoop                bool
	numBackedgesRemaining int
}

func (cb *clusterBuilder) popReady() []Candidate {
	ready := cb.ready
	cb.ready = []Candidate{}
	return ready
}

func (cb *clusterBuilder) consider(child NodeID, cross bool) bool {
	// If this (potentially transitive) child has no cross edges, it can be
	// infered that it is dominated by the current head.
	if cb.nodeInfo[child].incomingCrossEdges == 0 {
		unique := isUniqueSource(cb.graph, cb.currentHead, child, cb.idoms)
		if !unique {
			panic(child)
		}
		cluster := cb.nodeInfo[child].cluster
		if cluster == nil {
			panic(child)
		}
		cb.nodeInfo[child].cluster = nil
		cb.ready = append(cb.ready, Candidate{
			Node:        child,
			Cluster:     cluster,
			CrossEdgeIn: cross,
		})
		return true
	}
	return false
}

func (cb *clusterBuilder) contract(target NodeID) {
	head := cb.currentHead
	xit := cb.graph.ExitIterator(target)
	for xit.HasNext() {
		e, child := xit.GetNext()
		if head == child {
			cb.numBackedgesRemaining -= 1
			// Will become a trivial backedge, kill rather than moving.
			cb.graph.KillEdge(e)
			continue
		} else if head == cb.idoms[child] {
			// This move will eliminate a cross edge.
			cb.nodeInfo[child].incomingCrossEdges -= 1

			// HACK
			cb.graph.MoveEdgeEntry(head, e)

			if cb.consider(child, true) {
				cb.graph.KillEdge(e)
				continue
			}
		} else if target == cb.idoms[child] {
			cb.idoms[child] = head

			// HACK
			cb.graph.MoveEdgeEntry(head, e)

			if cb.consider(child, true) {
				cb.graph.KillEdge(e)
				continue
			}
		}
		cb.graph.MoveEdgeEntry(head, e)
	}
	// All of the entries should be dead after contraction.
	eit := cb.graph.EntryIterator(target)
	for eit.HasNext() {
		src, e := eit.GetNext()
		if src != cb.currentHead {
			panic(src)
		}
		cb.graph.KillEdge(e)
	}
}

func (cb *clusterBuilder) InitNode(head NodeID) {
	cb.currentHead = head
	cb.currentCluster = &ClusterLeaf{Head: head, Nodes: []NodeID{head}}

	cb.isLoop = false
	cb.numBackedgesRemaining = 0

	entryCount := 0
	it := cb.graph.EntryIterator(head)
	for it.HasNext() {
		src, e := it.GetNext()
		if cb.idoms[src] == NoNode {
			// Unreachable.
			cb.graph.KillEdge(e)
			continue
		}
		entryCount += 1
		if src == head {
			// Kill trivial backedges so that we don't need to
			// consider them again.
			cb.graph.KillEdge(e)
			cb.isLoop = true
		} else if cb.idoms[src] == head {
			// Non-trivial backedge.
			cb.edgeInfo[e].backEdge = true
			cb.numBackedgesRemaining += 1
			cb.isLoop = true
		} else if cb.idoms[head] != src {
			// Cross edge.
			cb.edgeInfo[e].crossEdge = true
			cb.nodeInfo[src].crossEdgeSource = true
			cb.nodeInfo[head].incomingCrossEdges += 1
		}
	}

	cb.nodeInfo[head].mergePoint = entryCount > 1

	// TODO not quite right.
	cb.nodeInfo[head].splitPoint = cb.graph.HasMultipleExits(head)
}

func (cb *clusterBuilder) BeginNode() {
	head := cb.currentHead
	// Enqueue ready nodes for processing.
	cb.ready = []Candidate{}
	xit := cb.graph.ExitIterator(head)
	for xit.HasNext() {
		_, dst := xit.GetNext()
		if head == cb.idoms[dst] {
			cb.consider(dst, false)
		}
	}
}

func (cb *clusterBuilder) EndNode() {
	cb.nodeInfo[cb.currentHead].cluster = cb.currentCluster
}

func mergeLoop(cb *clusterBuilder) {
	src := cb.currentHead
	xit := cb.graph.ExitIterator(src)
	for xit.HasNext() {
		e, dst := xit.GetNext()
		if src == dst {
			cb.graph.KillEdge(e)
		}
	}
	cb.isLoop = false
	cb.currentCluster = &ClusterLoop{Head: src, Body: cb.currentCluster}
}

func (cb *clusterBuilder) makeLinear(head NodeID, src Cluster, dst Cluster) Cluster {
	switch src := src.(type) {
	case *ClusterLeaf:
		switch dst := dst.(type) {
		case *ClusterLeaf:
			if !cb.nodeInfo[src.Nodes[len(src.Nodes)-1]].splitPoint && !cb.nodeInfo[dst.Nodes[0]].mergePoint {
				src.Nodes = append(src.Nodes, dst.Nodes...)
				return src
			}
		case *ClusterLinear:
			dst.Head = head

			first := dst.Clusters[0]
			switch first := first.(type) {
			case *ClusterLeaf:
				if !cb.nodeInfo[first.Head].mergePoint {
					first.Head = head
					first.Nodes = append(src.Nodes, first.Nodes...)
					return dst
				}

			}
			// Just put it at the front.
			dst.Clusters = append([]Cluster{src}, dst.Clusters...)
			return dst
		}
	case *ClusterLinear:
		switch dst := dst.(type) {
		case *ClusterLinear:
			src.Clusters = append(src.Clusters, dst.Clusters...)
		default:
			src.Clusters = append(src.Clusters, dst)
		}
		return src
	}
	return &ClusterLinear{
		Head: head,
		Clusters: []Cluster{
			src,
			dst,
		},
	}
}

func mergeLinear(cb *clusterBuilder) {
	ready := cb.popReady()
	candidate := ready[0]
	cb.contract(candidate.Node)

	if candidate.CrossEdgeIn && false {
		// If there's a cross edge pointing to this node, merging it into a linear block would cause problems.
		cb.currentCluster = cb.makeLinear(cb.currentHead, cb.currentCluster, &ClusterSwitch{
			Head:     cb.currentHead,
			Children: []Cluster{candidate.Cluster},
		})
	} else {
		cb.currentCluster = cb.makeLinear(cb.currentHead, cb.currentCluster, candidate.Cluster)
	}
}

func mergeSwitch(cb *clusterBuilder) {
	ready := cb.popReady()
	children := make([]Cluster, len(ready))
	for i := 0; i < len(ready); i++ {
		cb.contract(ready[i].Node)
		children[i] = ready[i].Cluster
	}

	cb.currentCluster = cb.makeLinear(cb.currentHead, cb.currentCluster, &ClusterSwitch{
		Head:     cb.currentHead,
		Children: children,
	})
}

func merge(cb *clusterBuilder) bool {
	head := cb.currentHead
	if cb.isLoop && cb.numBackedgesRemaining == 0 {
		mergeLoop(cb)
	}
	if len(cb.ready) == 0 {
		xit := cb.graph.ExitIterator(head)
		// Check for irreducible graphs.
		for xit.HasNext() {
			_, dst := xit.GetNext()
			if head == cb.idoms[dst] {
				// TODO support irreducible graphs.
				panic("Irreducible graph")
			}
		}
		// No more children to merge.
		// If we didn't collapse the loop, something is wrong.
		if cb.numBackedgesRemaining != 0 {
			panic(head)
		}
		return false
	} else if len(cb.ready) == 1 {
		mergeLinear(cb)
	} else {
		mergeSwitch(cb)
	}
	return true
}

// TODO incrementalize into one pass.
func gatherInfo(cb *clusterBuilder, order []NodeID) {
	g := cb.graph
	idoms := cb.idoms

	for j := len(order) - 1; j >= 0; j-- {
		n := order[j]
		if idoms[n] == NoNode {
			// Dead node.
			continue
		}

		numEntries := 0
		eit := g.EntryIterator(n)
		for eit.HasNext() {
			src, e := eit.GetNext()
			if idoms[src] == NoNode {
				g.KillEdge(e)
				continue
			}

			numEntries += 1
		}

		numExits := 0
		xit := g.ExitIterator(n)
		for xit.HasNext() {
			e, dst := xit.GetNext()
			if idoms[dst] == NoNode {
				g.KillEdge(e)
				continue
			}

			numExits += 1
		}

	}
}

func makeCluster(g *Graph) Cluster {
	order, index := ReversePostorder(g)
	cb := &clusterBuilder{
		graph:    g.Copy(),
		idoms:    FindDominators(g, order, index),
		nodeInfo: make([]nodeInfo, g.NumNodes()),
		edgeInfo: make([]edgeInfo, g.NumEdges()),
	}

	gatherInfo(cb, order)

	for j := len(order) - 1; j >= 0; j-- {
		head := order[j]
		cb.InitNode(head)
		// Defer processing nodes with cross edges.
		if cb.nodeInfo[head].crossEdgeSource {
			cb.nodeInfo[head].cluster = &ClusterLeaf{Head: head, Nodes: []NodeID{head}}
			continue
		}
		cb.BeginNode()
		for {
			if !merge(cb) {
				break
			}
		}
		cb.EndNode()
	}
	result := cb.nodeInfo[g.Entry()].cluster
	//result.Dump("")
	//fmt.Println()
	return result
}

type lfNodeInfo struct {
	pre           int
	post          int
	idom          NodeID
	loopHead      NodeID
	isHead        bool
	isIrreducible bool
	live          bool
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
	node      []lfNodeInfo
	edge      []EdgeType
	postorder []NodeID
	depth     []int
	current   int
}

func (lf *loopFinder) beginTraversingNode(n NodeID, prev NodeID) {
	lf.current += 1

	lf.node[n] = lfNodeInfo{
		pre:      lf.current,
		idom:     prev, // A reasonable inital guess.
		loopHead: NoNode,
		live:     true,
	}

	lf.depth[n] = lf.current
}

func (lf *loopFinder) endTraversingNode(n NodeID) {
	lf.current += 1
	lf.node[n].post = lf.current

	lf.depth[n] = 0

	lf.postorder = append(lf.postorder, n)
}

func (lf *loopFinder) isUnprocessed(n NodeID) bool {
	return !lf.node[n].live
}

func (lf *loopFinder) isBeingProcessed(dst NodeID) bool {
	return lf.depth[dst] > 0
}

func (lf *loopFinder) markLoopHeader(child NodeID, head NodeID) {
	if child == head {
		// Trivial loop.
		return
	}
	childHead := lf.node[child].loopHead
	for childHead != NoNode {
		if childHead == head {
			return
		}
		if lf.depth[childHead] < lf.depth[head] {
			// Found a closer head for this child, adopt it.
			lf.node[child].loopHead = head
			child, head = head, childHead
		} else {
			child = childHead
		}
		childHead = lf.node[child].loopHead
	}
	lf.node[child].loopHead = head
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
			head := lf.node[next].loopHead
			if head != NoNode {
				lf.markLoopHeader(n, head)
			}
		} else if lf.isBeingProcessed(next) {
			lf.edge[e] = BACKWARD
			lf.node[next].isHead = true
			lf.markLoopHeader(n, next)
		} else {
			lf.edge[e] = CROSS
			if lf.node[next].loopHead != NoNode {
				// Propagate loop header from cross edge.
				otherHead := lf.node[next].loopHead
				if lf.isBeingProcessed(otherHead) {
					lf.markLoopHeader(n, otherHead)
				} else {
					lf.edge[e] = REENTRY
					lf.node[otherHead].isIrreducible = true
					// Find and mark the common loop head.
					otherHead = lf.node[otherHead].loopHead
					for otherHead != NoNode {
						if lf.isBeingProcessed(otherHead) {
							lf.markLoopHeader(n, otherHead)
							break
						}
						otherHead = lf.node[otherHead].loopHead
					}
				}
			}
		}
	}
	lf.endTraversingNode(n)
}

func (lf *loopFinder) intersect(n0 NodeID, n1 NodeID) NodeID {
	i0 := lf.node[n0].post
	i1 := lf.node[n1].post
	for i0 != i1 {
		for i0 < i1 {
			n0 = lf.node[n0].idom
			i0 = lf.node[n0].post
		}
		for i0 > i1 {
			n1 = lf.node[n1].idom
			i1 = lf.node[n1].post
		}
	}
	return n0
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
		fmt.Println("idom iteration")
		changed = false
		for i := len(lf.postorder) - 1; i >= 0; i-- {
			n := lf.postorder[i]
			original := lf.node[n].idom
			idom := original
			eit := g.EntryIterator(n)
			for eit.HasNext() {
				src, _ := eit.GetNext()
				idom = lf.intersect(idom, src)
			}
			if idom != original {
				lf.node[n].idom = idom
				changed = true
			}
		}
	}
}

func analyzeStructure(g *Graph) ([]lfNodeInfo, []EdgeType, []NodeID) {
	numNodes := g.NumNodes()
	lf := &loopFinder{
		graph: g.Copy(),
		node:  make([]lfNodeInfo, numNodes),
		edge:  make([]EdgeType, g.NumEdges()),
		depth: make([]int, numNodes),
	}

	e := g.Entry()
	lf.process(e, e)
	lf.cleanDeadEdges()
	lf.findIdoms()
	return lf.node, lf.edge, lf.postorder
}

func uniqueEntry(g *Graph, src NodeID, dst NodeID) bool {
	eit := g.EntryIterator(dst)
	for eit.HasNext() {
		prev, _ := eit.GetNext()
		if prev != src {
			return false
		}
	}
	return true
}

func contract(g *Graph, src NodeID, dst NodeID, nodes []lfNodeInfo) {
	// TODO preserve edge order.

	eit := g.EntryIterator(dst)
	for eit.HasNext() {
		prev, e := eit.GetNext()
		if prev != src {
			panic(prev)
		}
		g.KillEdge(e)
	}

	xit := g.ExitIterator(dst)
	for xit.HasNext() {
		e, _ := xit.GetNext()
		g.MoveEdgeEntry(src, e)
		// TODO update idoms?
	}
}

func contractLoop(g *Graph, n NodeID) {
	xit := g.ExitIterator(n)
	for xit.HasNext() {
		e, dst := xit.GetNext()
		if n == dst {
			g.KillEdge(e)
		}
	}
}

func appendCluster(src Cluster, dst Cluster, singleEntry bool) Cluster {
	switch src := src.(type) {
	case *ClusterLeaf:
		switch dst := dst.(type) {
		case *ClusterLeaf:
			if singleEntry {
				src.Nodes = append(src.Nodes, dst.Nodes...)
				return src
			}
		case *ClusterLinear:
			if singleEntry {
				other, ok := dst.Clusters[0].(*ClusterLeaf)
				if ok {
					other.Nodes = append(src.Nodes, other.Nodes...)
					return dst
				}
			}
			dst.Clusters = append([]Cluster{src}, dst.Clusters...)
			return dst
		}
	case *ClusterLinear:
		switch dst := dst.(type) {
		case *ClusterLinear:
			src.Clusters = append(src.Clusters, dst.Clusters...)
			return src
		}
		src.Clusters = append(src.Clusters, dst)
		return src
	default:
		switch dst := dst.(type) {
		case *ClusterLinear:
			dst.Clusters = append([]Cluster{src}, dst.Clusters...)
			return dst
		}
	}
	return &ClusterLinear{Clusters: []Cluster{src, dst}}
}

func makeCluster2(g *Graph, nodes []lfNodeInfo, edges []EdgeType, postorder []NodeID) {
	g = g.Copy()

	clusters := make([]Cluster, g.NumNodes())

	for _, n := range postorder {
		var cluster Cluster = &ClusterLeaf{
			Nodes: []NodeID{n},
		}
		currentHead := nodes[n].loopHead
		if nodes[n].isHead {
			currentHead = n
		}
		for {
			fmt.Println(cluster.DumpShort())

			ready := []NodeID{}
			readyClusters := []Cluster{}
			xit := g.ExitIterator(n)
			for xit.HasNext() {
				e, dst := xit.GetNext()
				if nodes[dst].loopHead != currentHead {
					continue
				}
				if edges[e] == BACKWARD {
					continue
				}
				if clusters[dst] != nil && uniqueEntry(g, n, dst) {
					ready = append(ready, dst)
					readyClusters = append(readyClusters, clusters[dst])
					clusters[dst] = nil
				}
			}
			fmt.Println(n, ready)
			if len(ready) > 0 {
				if len(ready) > 1 {
					cluster = appendCluster(cluster, &ClusterSwitch{Children: readyClusters}, false)
				} else {
					singleEntry := !g.HasMultipleEntries(ready[0])
					cluster = appendCluster(cluster, readyClusters[0], singleEntry)
				}
				for _, dst := range ready {
					contract(g, n, dst, nodes)
				}
			} else {
				if currentHead == n {
					fmt.Println("loop done")
					contractLoop(g, n)
					cluster = &ClusterLoop{Body: cluster}
					currentHead = nodes[n].loopHead
					continue
				} else {
					break
				}
			}
		}
		clusters[n] = cluster
	}
	fmt.Println(clusters[0].DumpShort())
}

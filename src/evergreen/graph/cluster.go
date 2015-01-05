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

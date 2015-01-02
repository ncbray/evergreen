package graph

import (
	"fmt"
)

type Cluster interface {
	isCluster()
	Dump(margin string)
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

type ClusterSwitch struct {
	Head     NodeID
	Cond     Cluster
	Children []Cluster
}

func (cluster *ClusterSwitch) isCluster() {
}

func (cluster *ClusterSwitch) Dump(margin string) {
	fmt.Printf("%sswitch %d\n", margin, cluster.Head)
	childMargin := margin + ".   "
	cluster.Cond.Dump(childMargin)
	for _, c := range cluster.Children {
		c.Dump(childMargin)
	}
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

type ClusterComplex struct {
	Head     NodeID
	Clusters []Cluster
}

func (cluster *ClusterComplex) isCluster() {
}

func (cluster *ClusterComplex) Dump(margin string) {
	fmt.Printf("%scomplex %d\n", margin, cluster.Head)
	childMargin := margin + ".   "
	for _, c := range cluster.Clusters {
		c.Dump(childMargin)
	}
}

func isLoopHead(g *Graph, n NodeID, index []int) bool {
	it := g.EntryIterator(n)
	for it.HasNext() {
		src, _ := it.GetNext()
		if index[src] >= index[n] {
			return true
		}
	}
	return false
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

type clusterBuilder struct {
	graph          *Graph
	cluster        []Cluster
	idoms          []NodeID
	currentHead    NodeID
	currentCluster Cluster

	ready   []Candidate
	pending []Candidate

	hasLoopEdge  bool
	hasCrossEdge bool
}

func (cb *clusterBuilder) considerNext(src NodeID, dst NodeID) {
	if src == dst {
		cb.hasLoopEdge = true
	} else if src != cb.idoms[dst] {
		cb.hasCrossEdge = true
	} else if cb.cluster[dst] != nil {
		candidate := Candidate{
			Node:    dst,
			Cluster: cb.cluster[dst],
		}
		cb.cluster[dst] = nil
		cb.enqueue(candidate)
	}
}

func (cb *clusterBuilder) popReady() []Candidate {
	ready := cb.ready
	cb.ready = []Candidate{}
	return ready
}

func (cb *clusterBuilder) contract(dst NodeID) {
	xit := cb.graph.ExitIterator(dst)
	for xit.HasNext() {
		e, _ := xit.GetNext()
		cb.graph.MoveEdgeEntry(cb.currentHead, e)
	}
	eit := cb.graph.EntryIterator(dst)
	for eit.HasNext() {
		src, e := eit.GetNext()
		if src != cb.currentHead {
			panic(src)
		}
		cb.graph.KillEdge(e)
	}

}

func (cb *clusterBuilder) enqueue(candidate Candidate) {
	if isUniqueSource(cb.graph, cb.currentHead, candidate.Node, cb.idoms) {
		cb.ready = append(cb.ready, candidate)
	} else {
		// If the node is not immediately ready, there must be a cross edge pointing to it.
		candidate.CrossEdgeIn = true
		cb.pending = append(cb.pending, candidate)
	}
}

func (cb *clusterBuilder) PromotePending() {
	pending := cb.pending
	cb.pending = []Candidate{}
	for _, p := range pending {
		cb.enqueue(p)
	}
}

func (cb *clusterBuilder) ScanExits(src NodeID) {
	it := cb.graph.ExitIterator(src)
	for it.HasNext() {
		_, dst := it.GetNext()
		cb.considerNext(src, dst)
	}
}

func (cb *clusterBuilder) BeginNode(head NodeID) {
	cb.currentHead = head
	cb.currentCluster = &ClusterLeaf{Head: head, Nodes: []NodeID{head}}
	cb.ready = []Candidate{}
	cb.pending = []Candidate{}
	cb.hasLoopEdge = false
	cb.hasCrossEdge = false
	cb.ScanExits(head)
}

func (cb *clusterBuilder) EndNode() {
	cb.cluster[cb.currentHead] = cb.currentCluster
}

func mergeIrreducible(cb *clusterBuilder) {
	// TODO implement
	panic("irreducible graph merging not implemented.")
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
	cb.hasLoopEdge = false
	cb.currentCluster = &ClusterLoop{Head: src, Body: cb.currentCluster}
}

func makeLinear(head NodeID, src Cluster, dst Cluster) Cluster {
	switch src := src.(type) {
	case *ClusterLeaf:
		switch dst := dst.(type) {
		case *ClusterLeaf:
			src.Nodes = append(src.Nodes, dst.Nodes...)
			return src
		case *ClusterLinear:
			dst.Head = head
			dst.Clusters[0] = makeLinear(head, src, dst.Clusters[0])
			return dst
		case *ClusterSwitch:
			dst.Head = head
			dst.Cond = makeLinear(head, src, dst.Cond)
			return dst
		}
	case *ClusterLinear:
		switch dst := dst.(type) {
		case *ClusterLinear:
			src.Clusters = append(src.Clusters, dst.Clusters...)
		case *ClusterSwitch:
			dst.Head = head
			dst.Cond = makeLinear(head, src, dst.Cond)
			return dst
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

	if candidate.CrossEdgeIn {
		// If there's a cross edge pointing to this node, merging it into a linear block would cause problems.
		cb.currentCluster = &ClusterSwitch{
			Head:     cb.currentHead,
			Cond:     cb.currentCluster,
			Children: []Cluster{candidate.Cluster},
		}
	} else {
		cb.currentCluster = makeLinear(cb.currentHead, cb.currentCluster, candidate.Cluster)
	}

	//cb.currentCluster.Dump("")
	cb.PromotePending()
	cb.ScanExits(cb.currentHead)
}

func mergeSwitch(cb *clusterBuilder) {
	ready := cb.popReady()
	children := make([]Cluster, len(ready))
	for i := 0; i < len(ready); i++ {
		cb.contract(ready[i].Node)
		children[i] = ready[i].Cluster
	}

	cb.currentCluster = &ClusterSwitch{
		Head:     cb.currentHead,
		Cond:     cb.currentCluster,
		Children: children,
	}

	//cb.currentCluster.Dump("")
	cb.PromotePending()
	cb.ScanExits(cb.currentHead)
}

func merge(cb *clusterBuilder) bool {
	if len(cb.ready) == 0 {
		if len(cb.pending) != 0 {
			mergeIrreducible(cb)
		} else if cb.hasLoopEdge {
			// TODO merge loop ASAP
			mergeLoop(cb)
		} else {
			// No more children to merge.
			return false
		}
	} else if len(cb.ready) == 1 {
		mergeLinear(cb)
	} else {
		mergeSwitch(cb)
	}
	return true
}

func makeCluster(g *Graph, styler DotStyler) Cluster {
	order, index := ReversePostorder(g)
	cb := &clusterBuilder{
		graph:   g.Copy(),
		cluster: make([]Cluster, g.NumNodes()),
		idoms:   FindDominators(g, order, index),
	}

	for j := len(order) - 1; j >= 0; j-- {
		n := order[j]
		cb.BeginNode(n)
		for {
			if !merge(cb) {
				break
			}
		}
		cb.EndNode()
	}
	//cluster[0].Dump("")
	//fmt.Println()
	return cb.cluster[g.Entry()]
}

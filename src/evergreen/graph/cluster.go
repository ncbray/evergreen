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
	Nodes []NodeID
}

func (cluster *ClusterLeaf) isCluster() {
}

func (cluster *ClusterLeaf) Dump(margin string) {
	fmt.Printf("%sleaf %d\n", margin, len(cluster.Nodes))
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
	Clusters []Cluster
}

func (cluster *ClusterLinear) isCluster() {
}

func (cluster *ClusterLinear) Dump(margin string) {
	fmt.Printf("%slinear %d\n", margin, len(cluster.Clusters))
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
	Children []Cluster
}

func (cluster *ClusterSwitch) isCluster() {
}

func (cluster *ClusterSwitch) Dump(margin string) {
	fmt.Printf("%sswitch %d\n", margin, len(cluster.Children))
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
	Body Cluster
}

func (cluster *ClusterLoop) isCluster() {
}

func (cluster *ClusterLoop) Dump(margin string) {
	fmt.Printf("%sloop\n", margin)
	childMargin := margin + ".   "
	cluster.Body.Dump(childMargin)
}

func (cluster *ClusterLoop) DumpShort() string {
	return fmt.Sprintf("{%s}", cluster.Body.DumpShort())
}

func makeCluster(g *Graph) Cluster {
	info, edges, postorder := analyzeStructure(g)
	return makeCluster2(g, info, edges, postorder)
}

type nodeInfo struct {
	pre           int
	post          int
	idom          NodeID
	loopHead      NodeID
	isHead        bool
	isIrreducible bool
	live          bool
}

type edgeType int

const (
	DEAD edgeType = iota
	FORWARD
	BACKWARD
	CROSS
	REENTRY
)

type loopFinder struct {
	graph     *Graph
	node      []nodeInfo
	edge      []edgeType
	postorder []NodeID
	depth     []int
	current   int
}

func (lf *loopFinder) beginTraversingNode(n NodeID, prev NodeID) {
	lf.current += 1

	lf.node[n] = nodeInfo{
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

func analyzeStructure(g *Graph) ([]nodeInfo, []edgeType, []NodeID) {
	numNodes := g.NumNodes()
	lf := &loopFinder{
		graph: g.Copy(),
		node:  make([]nodeInfo, numNodes),
		edge:  make([]edgeType, g.NumEdges()),
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

func contract(g *Graph, src NodeID, dst NodeID, nodes []nodeInfo) {
	transferedExits := false
	eit := g.EntryIterator(dst)
	for eit.HasNext() {
		prev, e := eit.GetNext()
		if prev != src {
			panic(prev)
		}
		if !transferedExits {
			g.ReplaceEdgeWithExits(e, dst)
			transferedExits = true
		}
		g.KillEdge(e)
	}
	if !transferedExits {
		panic(dst)
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

func appendCluster(src Cluster, dst Cluster, shouldFuse bool) Cluster {
	if shouldFuse {
		switch src := src.(type) {
		case *ClusterLeaf:
			switch dst := dst.(type) {
			case *ClusterLeaf:
				src.Nodes = append(src.Nodes, dst.Nodes...)
				return src
			case *ClusterLinear:
				other, ok := dst.Clusters[0].(*ClusterLeaf)
				if ok {
					other.Nodes = append(src.Nodes, other.Nodes...)
					return dst
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
	} else {
		switch src := src.(type) {
		case *ClusterLinear:
			src.Clusters = append(src.Clusters, dst)
			return src
		}
	}
	return &ClusterLinear{Clusters: []Cluster{src, dst}}
}

func isClusterHead(g *Graph, src NodeID, nodes []nodeInfo, edges []edgeType) bool {
	xit := g.ExitIterator(src)
	for xit.HasNext() {
		e, dst := xit.GetNext()
		if edges[e] != BACKWARD && nodes[dst].idom != src {
			return false
		}
	}
	return true
}

func clusterRegion(g *Graph, n NodeID, currentHead NodeID, cluster Cluster, clusters []Cluster, nodes []nodeInfo, edges []edgeType) Cluster {
	for {
		ready := []NodeID{}
		readyClusters := []Cluster{}
		//readyEdges := []EdgeID{}
		//pendingEdges := []EdgeID{}
		numExits := 0
		xit := g.ExitIterator(n)
		for xit.HasNext() {
			e, dst := xit.GetNext()
			numExits += 1
			if nodes[dst].loopHead != currentHead {
				continue
			}
			if edges[e] == BACKWARD {
				continue
			}
			if uniqueEntry(g, n, dst) {
				if clusters[dst] != nil {
					ready = append(ready, dst)
					readyClusters = append(readyClusters, clusters[dst])
					clusters[dst] = nil
				}
			}
		}
		if len(ready) > 0 {
			if len(ready) > 1 {
				cluster = appendCluster(cluster, &ClusterSwitch{Children: readyClusters}, false)
			} else {
				cluster = appendCluster(cluster, readyClusters[0], numExits == 1)
			}
			for _, dst := range ready {
				contract(g, n, dst, nodes)
			}
		} else {
			return cluster
		}
	}
}

func makeCluster2(g *Graph, nodes []nodeInfo, edges []edgeType, postorder []NodeID) Cluster {
	g = g.Copy()

	clusters := make([]Cluster, g.NumNodes())

	for _, n := range postorder {
		var cluster Cluster = &ClusterLeaf{
			Nodes: []NodeID{n},
		}
		if nodes[n].isHead {
			cluster = clusterRegion(g, n, n, cluster, clusters, nodes, edges)
			contractLoop(g, n)
			cluster = &ClusterLoop{Body: cluster}
		}
		if isClusterHead(g, n, nodes, edges) {
			cluster = clusterRegion(g, n, nodes[n].loopHead, cluster, clusters, nodes, edges)
		}
		clusters[n] = cluster
	}
	return clusters[0]
}

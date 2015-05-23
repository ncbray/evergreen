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

func MakeCluster(g *Graph) Cluster {
	g = g.Copy()
	info, edges, postorder := AnalyzeStructure(g)
	return contractClusters(g, info, edges, postorder)
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

func contract(g *Graph, src NodeID, dst NodeID, nodes []NodeInfo) {
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

func isClusterHead(g *Graph, src NodeID, nodes []NodeInfo, edges []EdgeType) bool {
	xit := g.ExitIterator(src)
	for xit.HasNext() {
		e, dst := xit.GetNext()
		if edges[e] != BACKWARD && nodes[dst].IDom != src {
			return false
		}
	}
	return true
}

func clusterRegion(g *Graph, n NodeID, currentHead NodeID, cluster Cluster, clusters []Cluster, nodes []NodeInfo, edges []EdgeType) Cluster {
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
			if nodes[dst].LoopHead != currentHead {
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

func contractClusters(g *Graph, nodes []NodeInfo, edges []EdgeType, postorder []NodeID) Cluster {
	clusters := make([]Cluster, g.NumNodes())

	for _, n := range postorder {
		var cluster Cluster = &ClusterLeaf{
			Nodes: []NodeID{n},
		}
		if nodes[n].IsHead {
			cluster = clusterRegion(g, n, n, cluster, clusters, nodes, edges)
			contractLoop(g, n)
			cluster = &ClusterLoop{Body: cluster}
		}
		if isClusterHead(g, n, nodes, edges) {
			cluster = clusterRegion(g, n, nodes[n].LoopHead, cluster, clusters, nodes, edges)
		}
		clusters[n] = cluster
	}
	for i := 1; i < len(clusters); i++ {
		if clusters[i] != nil {
			panic(i)
		}
	}
	return clusters[0]
}

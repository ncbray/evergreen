package graph

import (
	"fmt"
)

type Destination struct {
	Edge EdgeID
	Node NodeID
}

type Cluster interface {
	isCluster()
	Dump(margin string)
}

type ClusterLinear struct {
	Head     NodeID
	Nodes    []NodeID
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

func makeCluster(g *Graph, styler DotStyler) Cluster {
	order, index := ReversePostorder(g)
	idoms := FindDominators(g, order, index)

	dominates := make([][]NodeID, g.NumNodes())
	cluster := make([]Cluster, g.NumNodes())

	for j := len(order) - 1; j >= 0; j-- {
		n := order[j]
		it := g.ExitIterator(n)
		localFlow := []Destination{}
		nonLocalFlow := []Destination{}
		dominatesLocalEdges := true
		for it.HasNext() {
			e, dst := it.GetNext()
			d := Destination{Edge: e, Node: dst}
			if styler.IsLocalFlow(e) {
				localFlow = append(localFlow, d)
				if idoms[dst] != n {
					dominatesLocalEdges = false
				}
			} else {
				nonLocalFlow = append(nonLocalFlow, d)
			}
		}

		root := &ClusterLinear{
			Head:     n,
			Nodes:    []NodeID{n},
			Clusters: []Cluster{},
		}

		if len(localFlow) == 1 && dominatesLocalEdges {
			// Attempt to fuse
			dst := localFlow[0].Node
			next := cluster[dst]
			switch next := next.(type) {
			case *ClusterLinear:
				if g.HasMultipleEntries(next.Head) {
					break
				}
				root.Nodes = append(root.Nodes, next.Nodes...)
				root.Clusters = next.Clusters
				cluster[dst] = nil
			default:
				root.Clusters = []Cluster{next}
				cluster[dst] = nil
			}
		}

		var current Cluster = root

		if isLoopHead(g, n, index) {
			current = &ClusterLoop{
				Head: n,
				Body: current,
			}
		}

		// See if there are any dominated nodes we haven't been able to fuse.
		region := []Cluster{current}
		d := dominates[n]
		for i := len(d) - 1; i >= 0; i-- {
			child := d[i]
			if cluster[child] != nil {
				region = append(region, cluster[child])
				cluster[child] = nil
			}
		}
		if len(region) > 1 {
			cluster[n] = &ClusterComplex{Head: n, Clusters: region}
		} else {
			cluster[n] = current
		}

		// Incrementally build the dominator list.
		idom := idoms[n]
		if idom != n {
			dominates[idom] = append(dominates[idom], n)
		}
	}

	//cluster[0].Dump("")
	//fmt.Println()

	return cluster[0]
}

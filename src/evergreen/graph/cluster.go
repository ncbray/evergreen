package graph

import (
	"fmt"
)

type Cluster interface {
	isCluster()
	Dump(margin string)
	Head() NodeID
}

type ClusterLeaf struct {
	Nodes []NodeID
}

func (cluster *ClusterLeaf) isCluster() {
}

func (cluster *ClusterLeaf) Head() NodeID {
	return cluster.Nodes[0]
}

func (cluster *ClusterLeaf) Dump(margin string) {
	fmt.Printf("%sleaf %v\n", margin, cluster.Nodes)
}

type ClusterLinear struct {
	Clusters []Cluster
}

func (cluster *ClusterLinear) isCluster() {
}

func (cluster *ClusterLinear) Head() NodeID {
	return cluster.Clusters[0].Head()
}

func (cluster *ClusterLinear) Dump(margin string) {
	fmt.Printf("%slinear\n", margin)
	childMargin := margin + ".   "
	for _, c := range cluster.Clusters {
		c.Dump(childMargin)
	}
}

type ClusterComplex struct {
	Entry    NodeID
	Clusters []Cluster
}

func (cluster *ClusterComplex) isCluster() {
}

func (cluster *ClusterComplex) Head() NodeID {
	return cluster.Entry
}

func (cluster *ClusterComplex) Dump(margin string) {
	fmt.Printf("%scomplex %d\n", margin, cluster.Entry)
	childMargin := margin + ".   "
	for _, c := range cluster.Clusters {
		c.Dump(childMargin)
	}
}

func leafSeed(n NodeID) Cluster {
	return &ClusterLeaf{Nodes: []NodeID{n}}
}

// TODO more efficient list construction
func prependLinear(n NodeID, next Cluster) Cluster {
	switch next := next.(type) {
	case *ClusterLeaf:
		next.Nodes = append([]NodeID{n}, next.Nodes...)
		return next
	case *ClusterLinear:
		head := next.Clusters[0]
		switch head := head.(type) {
		case *ClusterLeaf:
			head.Nodes = append([]NodeID{n}, head.Nodes...)
		default:
			next.Clusters = append([]Cluster{leafSeed(n)}, next.Clusters...)
		}
		return next
	case *ClusterComplex:
		return &ClusterLinear{Clusters: []Cluster{leafSeed(n), next}}
	default:
		panic(next)
	}

}

func makeCluster(g *Graph) Cluster {
	order, index := ReversePostorder(g)
	idoms := FindDominators(g, order, index)

	dominates := make([][]NodeID, g.NumNodes())
	cluster := make([]Cluster, g.NumNodes())

	for j := len(order) - 1; j >= 0; j-- {
		n := order[j]

		numDominated := len(dominates[n])

		if numDominated == 0 {
			cluster[n] = leafSeed(n)
		} else if numDominated == 1 {
			// TODO more than one exit.
			// TODO could be loop?
			child := dominates[n][0]
			prev := cluster[child]
			cluster[n] = prependLinear(n, prev)
			cluster[child] = nil
		} else {
			children := make([]Cluster, numDominated)
			for i, child := range dominates[n] {
				children[numDominated-i-1] = cluster[child]
				cluster[child] = nil
			}
			cluster[n] = &ClusterComplex{Entry: n, Clusters: children}
		}

		idom := idoms[n]
		if idom != n {
			dominates[idom] = append(dominates[idom], n)
		}
	}

	//cluster[0].Dump("")
	//fmt.Println()

	return cluster[0]
}

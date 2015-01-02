package graph

import (
	"bytes"
	"fmt"
)

func nodeDotID(node NodeID) string {
	return fmt.Sprintf("n%d", node)
}

type DotStyler interface {
	NodeStyle(node NodeID) string
	EdgeStyle(src NodeID, edge EdgeID, dst NodeID) string
	IsLocalFlow(edge EdgeID) bool
}

func drawNode(buf *bytes.Buffer, node NodeID, styler DotStyler) {
	buf.WriteString("  ")
	buf.WriteString(nodeDotID(node))
	buf.WriteString("[")
	buf.WriteString(styler.NodeStyle(node))
	buf.WriteString("];\n")
}

func drawUnclusteredNodes(buf *bytes.Buffer, order []NodeID, styler DotStyler) {
	nit := OrderedIterator(order)
	for nit.HasNext() {
		drawNode(buf, nit.GetNext(), styler)
	}
}

func drawCluster(buf *bytes.Buffer, cluster Cluster, styler DotStyler) {
	switch cluster := cluster.(type) {
	case *ClusterLinear:
		buf.WriteString(fmt.Sprintf("subgraph cluster_%d {\n", cluster.Head))
		buf.WriteString("  labeljust=l;\n")
		buf.WriteString("  label=linear;\n")
		buf.WriteString("  color=lightgrey;\n")

		for _, n := range cluster.Nodes {
			drawNode(buf, n, styler)
		}
		for _, c := range cluster.Clusters {
			drawCluster(buf, c, styler)
		}

		buf.WriteString("}\n")
	case *ClusterLoop:
		buf.WriteString(fmt.Sprintf("subgraph cluster_%d {\n", cluster.Head))
		buf.WriteString("  labeljust=l;\n")
		buf.WriteString("  label=loop;\n")
		buf.WriteString("  color=lightgrey;\n")
		drawCluster(buf, cluster.Body, styler)
		buf.WriteString("}\n")
	case *ClusterComplex:
		buf.WriteString(fmt.Sprintf("subgraph cluster_%d {\n", cluster.Head))
		buf.WriteString("  labeljust=l;\n")
		buf.WriteString("  label=complex;\n")
		buf.WriteString("  color=lightgrey;\n")

		for _, c := range cluster.Clusters {
			drawCluster(buf, c, styler)
		}

		buf.WriteString("}\n")
	default:
		panic(cluster)
	}
}

func drawClusteredNodes(buf *bytes.Buffer, g *Graph, styler DotStyler) {
	cluster := makeCluster(g, styler)
	drawCluster(buf, cluster, styler)
}

func GraphToDot(g *Graph, styler DotStyler) string {
	order, index := ReversePostorder(g)

	var idoms []NodeID
	visualize_idoms := false
	if visualize_idoms {
		idoms = FindDominators(g, order, index)
	}

	buf := &bytes.Buffer{}
	buf.WriteString("digraph G {\n")
	buf.WriteString("  nslimit = 3;\n") // Make big graphs render faster.

	//drawUnclusteredNodes(buf, order, styler)
	drawClusteredNodes(buf, g, styler)

	// Draw edges.
	nit := OrderedIterator(order)
	for nit.HasNext() {
		node := nit.GetNext()
		eit := g.ExitIterator(node)
		for eit.HasNext() {
			edge, dst := eit.GetNext()
			buf.WriteString("  ")
			buf.WriteString(nodeDotID(node))
			buf.WriteString(" -> ")
			buf.WriteString(nodeDotID(dst))
			buf.WriteString("[")
			buf.WriteString(styler.EdgeStyle(node, edge, dst))
			if index[node] >= index[dst] {
				buf.WriteString(",weight=0")
			}
			buf.WriteString("];\n")
		}
	}
	if visualize_idoms {
		nit := OrderedIterator(order)
		for nit.HasNext() {
			src := nit.GetNext()
			dst := idoms[src]
			if src != dst {
				buf.WriteString("  ")
				buf.WriteString(nodeDotID(src))
				buf.WriteString(" -> ")
				buf.WriteString(nodeDotID(dst))
				buf.WriteString("[")
				buf.WriteString("style=dotted")
				buf.WriteString("];\n")
			}
		}
	}
	buf.WriteString("}\n")
	return buf.String()
}

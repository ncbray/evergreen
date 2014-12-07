package graph

import (
	"bytes"
	"fmt"
)

func NodeDotID(node NodeID) string {
	return fmt.Sprintf("n%d", node)
}

type DotStyler interface {
	NodeStyle(node NodeID) string
	EdgeStyle(node NodeID, flow int) string
}

func GraphToDot(g *Graph, styler DotStyler) string {
	order, index := ReversePostorder(g)

	var idoms []NodeID
	visualize_idoms := false
	if visualize_idoms {
		idoms = FindDominators(g, order, index)
	}

	var buf bytes.Buffer
	buf.WriteString("digraph G {\n")
	buf.WriteString("  nslimit = 3;\n") // Make big graphs render faster.
	nit := OrderedIterator(order)
	for nit.Next() {
		node := nit.Value()
		buf.WriteString("  ")
		buf.WriteString(NodeDotID(node))
		buf.WriteString("[")
		buf.WriteString(styler.NodeStyle(node))
		buf.WriteString("];\n")

		eit := ExitIterator(g, node)
		for eit.Next() {
			dst := eit.Value()
			buf.WriteString("  ")
			buf.WriteString(NodeDotID(node))
			buf.WriteString(" -> ")
			buf.WriteString(NodeDotID(dst))
			buf.WriteString("[")
			buf.WriteString(styler.EdgeStyle(node, eit.Label()))
			buf.WriteString("];\n")
		}
	}
	if visualize_idoms {
		nit := OrderedIterator(order)
		for nit.Next() {
			src := nit.Value()
			dst := idoms[src]
			if src != dst {
				buf.WriteString("  ")
				buf.WriteString(NodeDotID(src))
				buf.WriteString(" -> ")
				buf.WriteString(NodeDotID(dst))
				buf.WriteString("[")
				buf.WriteString("style=dotted")
				buf.WriteString("];\n")
			}
		}
	}
	buf.WriteString("}\n")
	return buf.String()
}

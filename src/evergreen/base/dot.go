package base

import (
	"bytes"
	"fmt"
)

type activeNode struct {
	node *Node
	flow int
}

type DFSListener interface {
	PreNode(n *Node)
	PostNode(n *Node)
}

type Postorder struct {
	nodes []*Node
}

func (v *Postorder) PreNode(n *Node) {
}

func (v *Postorder) PostNode(n *Node) {
	v.nodes = append(v.nodes, n)
}

func DFS(r *Region, visitor DFSListener) {
	visited := make(map[*Node]bool)
	stack := make([]activeNode, 0)
	current := activeNode{nil, 0}
	visited[nil] = true
	push := func(node *Node) {
		visitor.PreNode(node)
		stack = append(stack, current)
		current = activeNode{node, 0}
		visited[node] = true
	}
	pop := func() {
		visitor.PostNode(current.node)
		current, stack = stack[len(stack)-1], stack[:len(stack)-1]
	}
	push(r.GetEntry())
	for current.node != nil {
		num := current.node.NumExits()
		if current.flow < num {
			// Reverse iteration gives expected order for postorder
			e := current.node.GetExit(num - current.flow - 1)
			current.flow += 1
			if !visited[e.dst] {
				push(e.dst)
			}
		} else {
			pop()
		}
	}
}

func ReversePostorder(r *Region) []*Node {
	visitor := &Postorder{}
	DFS(r, visitor)
	n := len(visitor.nodes)
	for i := 0; i < n/2; i++ {
		visitor.nodes[i], visitor.nodes[n-1-i] = visitor.nodes[n-1-i], visitor.nodes[i]
	}
	for i := 0; i < n; i++ {
		visitor.nodes[i].Name = i
	}
	return visitor.nodes
}

func intersect(idoms []int, finger1 int, finger2 int) int {
	for finger1 != finger2 {
		for finger1 > finger2 {
			finger1 = idoms[finger1]
		}
		for finger1 < finger2 {
			finger2 = idoms[finger2]
		}
	}
	return finger1
}

// Assumes reverse postorder.
func FindIdoms(ordered []*Node) []int {
	idoms := make([]int, len(ordered))
	earliest := make([]int, len(ordered))

	idoms[0] = 0

	n := len(ordered)

	for i := 1; i < n; i++ {
		idoms[i] = NoNode

		// Find the earliest use of this node.
		e := n
		node := ordered[i]
		for j := 0; j < node.NumExits(); j++ {
			next := node.GetNext(j)
			if next != nil && next.Name < e {
				e = next.Name
			}
		}
		earliest[i] = e
	}
	start := 1
	for start < n {
		i := start
		start = len(ordered)
		for ; i < n; i++ {
			// Note: assumes there are no dead entries.
			entries := ordered[i].peekEntries()
			new_idom := NoNode
			for j := 0; j < len(entries); j++ {
				other := entries[j].src.Name
				// Is it available, yet?
				if idoms[other] == NoNode {
					continue
				}
				// Is it the first we've found?
				if new_idom == NoNode {
					new_idom = other
				} else {
					new_idom = intersect(idoms, other, new_idom)
				}
			}
			if idoms[i] != new_idom {
				idoms[i] = new_idom
				if earliest[i] < start && earliest[i] < i {
					start = earliest[i]
				}
			}
		}
	}
	return idoms
}

func NodeID(node *Node) string {
	return fmt.Sprintf("n%d", node.Name)
}

type DotStyler interface {
	NodeStyle(data interface{}) string
	EdgeStyle(data interface{}, flow int) string
}

func RegionToDot(region *Region, styler DotStyler) string {
	nodes := ReversePostorder(region)

	var idoms []int
	visualize_idoms := false
	if visualize_idoms {
		idoms = FindIdoms(nodes)
	}

	var buf bytes.Buffer
	buf.WriteString("digraph G {\n")
	for _, node := range nodes {
		buf.WriteString("  ")
		buf.WriteString(NodeID(node))
		buf.WriteString("[")
		buf.WriteString(styler.NodeStyle(node.Data))
		buf.WriteString("];\n")

		for i := 0; i < node.NumExits(); i++ {
			dst := node.GetNext(i)
			if dst != nil {
				buf.WriteString("  ")
				buf.WriteString(NodeID(node))
				buf.WriteString(" -> ")
				buf.WriteString(NodeID(dst))
				buf.WriteString("[")
				buf.WriteString(styler.EdgeStyle(node.Data, i))
				buf.WriteString("];\n")
			}
		}
	}
	if visualize_idoms {
		for i, idom := range idoms {
			if i != idom {
				buf.WriteString("  ")
				buf.WriteString(NodeID(nodes[i]))
				buf.WriteString(" -> ")
				buf.WriteString(NodeID(nodes[idom]))
				buf.WriteString("[")
				buf.WriteString("style=dotted")
				buf.WriteString("];\n")
			}
		}
	}
	buf.WriteString("}\n")
	return buf.String()
}

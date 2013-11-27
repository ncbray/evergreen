package dub

import (
	"bytes"
	"fmt"
	"reflect"
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
	push(r.entry)
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
	return visitor.nodes
}

func NodeID(node *Node) string {
	// TODO make deterministic
	return string(fmt.Sprint(reflect.ValueOf(node).Elem().UnsafeAddr()))
}

func RegionToDot(region *Region) string {
	nodes := ReversePostorder(region)

	var buf bytes.Buffer
	buf.WriteString("digraph G {\n")
	for _, node := range nodes {
		buf.WriteString("  ")
		buf.WriteString(NodeID(node))
		buf.WriteString("[")
		buf.WriteString(node.Data.DotNodeStyle())
		buf.WriteString("];\n")

		for i := 0; i < node.NumExits(); i++ {
			dst := node.GetNext(i)
			if dst != nil {
				buf.WriteString("  ")
				buf.WriteString(NodeID(node))
				buf.WriteString(" -> ")
				buf.WriteString(NodeID(dst))
				buf.WriteString("[")
				buf.WriteString(node.Data.DotEdgeStyle(i))
				buf.WriteString("];\n")
			}
		}
	}
	buf.WriteString("}\n")
	return buf.String()
}

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

type DefUseCollector struct {
	Uses [][]int
	Defs [][]int
}

func MakeDefUse(numVars int) *DefUseCollector {
	return &DefUseCollector{Uses: make([][]int, numVars), Defs: make([][]int, numVars)}
}

func (c *DefUseCollector) AddUse(v int, n int) {
	c.Uses[v] = append(c.Uses[v], n)
}

func (c *DefUseCollector) AddDef(v int, n int) {
	c.Defs[v] = append(c.Defs[v], n)
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

func FindFrontiers(ordered []*Node, idoms []int) [][]int {
	n := len(ordered)
	frontiers := make([][]int, n)
	for i := 0; i < n; i++ {
		// Assumes no dead entries.
		entries := ordered[i].peekEntries()
		if len(entries) >= 2 {
			target := idoms[i]
			for _, edge := range entries {
				runner := edge.src.Name
				for runner != target {
					frontiers[runner] = append(frontiers[runner], i)
					runner = idoms[runner]
				}
			}
		}
	}
	return frontiers
}

type SSIBuilder struct {
	nodes    []*Node
	idoms    []int
	df       [][]int
	phiFuncs [][]int
}

func CreateSSIBuilder(r *Region, nodes []*Node) *SSIBuilder {
	idoms := FindIdoms(nodes)
	df := FindFrontiers(nodes, idoms)
	phiFuncs := make([][]int, len(nodes))
	return &SSIBuilder{
		nodes:    nodes,
		idoms:    idoms,
		df:       df,
		phiFuncs: phiFuncs,
	}
}

type SSIState struct {
	builder *SSIBuilder
	uid     int

	phiPlaced   map[int]bool
	defEnqueued map[int]bool
	defQueue    []int

	sigmaPlaced map[int]bool
	useEnqueued map[int]bool
	useQueue    []int
}

func CreateSSIState(builder *SSIBuilder, uid int) *SSIState {
	return &SSIState{
		builder:     builder,
		uid:         uid,
		defEnqueued: map[int]bool{},
		phiPlaced:   map[int]bool{},
		useEnqueued: map[int]bool{},
		sigmaPlaced: map[int]bool{},
	}
}

func (state *SSIState) DiscoveredDef(node int) {
	enqueued, _ := state.defEnqueued[node]
	if !enqueued {
		state.defEnqueued[node] = true
		state.defQueue = append(state.defQueue, node)
	}
}

func (state *SSIState) GetNextDef() int {
	current := state.defQueue[len(state.defQueue)-1]
	state.defQueue = state.defQueue[:len(state.defQueue)-1]
	return current
}

func (state *SSIState) PlacePhi(node int) {
	placed, _ := state.phiPlaced[node]
	if !placed {
		state.builder.phiFuncs[node] = append(state.builder.phiFuncs[node], state.uid)
		state.phiPlaced[node] = true
		state.DiscoveredDef(node)
		for _, e := range state.builder.nodes[node].peekEntries() {
			state.DiscoveredUse(e.src.Name)
		}
	}

}

func (state *SSIState) DiscoveredUse(node int) {
	enqueued, _ := state.useEnqueued[node]
	if !enqueued {
		state.useEnqueued[node] = true
		state.useQueue = append(state.useQueue, node)
	}
}

func SSI(builder *SSIBuilder, uid int, defs []int) {
	state := CreateSSIState(builder, uid)
	for _, def := range defs {
		state.DiscoveredDef(def)
	}
	// TODO pump use queue, place sigmas.
	for len(state.defQueue) > 0 {
		current := state.GetNextDef()
		for _, f := range builder.df[current] {
			state.PlacePhi(f)
		}
	}
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

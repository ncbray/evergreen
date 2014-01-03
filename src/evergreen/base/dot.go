package base

import (
	"bytes"
	"fmt"
	"sort"
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
	VarUseAt [][]int
	VarDefAt [][]int
	NodeUses [][]int
	NodeDefs [][]int
}

func (c *DefUseCollector) AddUse(n int, v int) {
	c.VarUseAt[v] = append(c.VarUseAt[v], n)
	c.NodeUses[n] = append(c.NodeUses[n], v)
}

func (c *DefUseCollector) AddDef(n int, v int) {
	c.VarDefAt[v] = append(c.VarDefAt[v], n)
	c.NodeDefs[n] = append(c.NodeDefs[n], v)
}

func (c *DefUseCollector) IsDefined(n int, v int) bool {
	for _, d := range c.NodeDefs[n] {
		if d == v {
			return true
		}
	}
	return false
}

func CreateDefUse(numNodes int, numVars int) *DefUseCollector {
	return &DefUseCollector{
		VarUseAt: make([][]int, numVars),
		VarDefAt: make([][]int, numVars),
		NodeUses: make([][]int, numNodes),
		NodeDefs: make([][]int, numNodes),
	}
}

func FindDefUse(nodes []*Node, numVars int, visit func(*Node, *DefUseCollector)) *DefUseCollector {
	defuse := CreateDefUse(len(nodes), numVars)
	for _, node := range nodes {
		visit(node, defuse)
	}
	return defuse
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

type LivenessOracle interface {
	LiveAtEntry(n int, v int) bool
	LiveAtExit(n int, v int) bool
}

type SimpleLivenessOracle struct {
}

func (l *SimpleLivenessOracle) LiveAtEntry(n int, v int) bool {
	return true
}

func (l *SimpleLivenessOracle) LiveAtExit(n int, v int) bool {
	return true
}

type LiveVars struct {
	liveIn  []map[int]bool
	liveOut []map[int]bool
}

func canonicalSet(set map[int]bool) []int {
	out := make([]int, len(set))
	i := 0
	for k, _ := range set {
		out[i] = k
		i++
	}
	sort.Ints(out)
	return out
}

func (l *LiveVars) LiveSet(n int) []int {
	return canonicalSet(l.liveIn[n])
}

func (l *LiveVars) LiveAtEntry(n int, v int) bool {
	live, _ := l.liveIn[n][v]
	return live
}

func (l *LiveVars) LiveAtExit(n int, v int) bool {
	live, _ := l.liveOut[n][v]
	return live
}

func FindLiveVars(order []*Node, defuse *DefUseCollector) *LiveVars {
	n := len(order)
	liveIn := make([]map[int]bool, n)
	liveOut := make([]map[int]bool, n)
	// Initialize with the uses for each node.
	for i := 0; i < n; i++ {
		liveIn[i] = map[int]bool{}
		liveOut[i] = map[int]bool{}
		for _, v := range defuse.NodeUses[i] {
			liveIn[i][v] = true
		}
	}
	// Iterate until a stable state is reached.
	changed := true
	for changed {
		changed = false
		// Propagate the uses backwards.
		for i := n - 1; i >= 0; i-- {
			for s := 0; s < order[i].NumExits(); s++ {
				next := order[i].GetNext(s)
				// Not all exits are connected.
				if next == nil {
					continue
				}
				// Merge sets from predecessors.
				for v, _ := range liveIn[next.Name] {
					_, exists := liveOut[i][v]
					if !exists {
						liveOut[i][v] = true
						changed = true
						// Filter out local defines
						if !defuse.IsDefined(i, v) {
							liveIn[i][v] = true
						}
					}
				}
			}
		}
	}
	return &LiveVars{liveIn: liveIn, liveOut: liveOut}
}

type SSIBuilder struct {
	nodes    []*Node
	idoms    []int
	df       [][]int
	PhiFuncs [][]int
	live     LivenessOracle
}

func CreateSSIBuilder(r *Region, nodes []*Node, live LivenessOracle) *SSIBuilder {
	idoms := FindIdoms(nodes)
	df := FindFrontiers(nodes, idoms)
	phiFuncs := make([][]int, len(nodes))
	return &SSIBuilder{
		nodes:    nodes,
		idoms:    idoms,
		df:       df,
		PhiFuncs: phiFuncs,
		live:     live,
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
	if !state.builder.live.LiveAtEntry(node, state.uid) {
		return
	}
	placed, _ := state.phiPlaced[node]
	if !placed {
		state.builder.PhiFuncs[node] = append(state.builder.PhiFuncs[node], state.uid)
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
				buf.WriteString("[")
				buf.WriteString("style=dotted")
				buf.WriteString("];\n")
			}
		}
	}
	buf.WriteString("}\n")
	return buf.String()
}

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
	push(r.Entry)
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

type rpoSearch struct {
	graph *Graph
	order []NodeID
	index []int
}

func (s *rpoSearch) search(n NodeID) {
	if s.index[n] != 0 {
		return
	}
	// Prevent processing it again, the correct index will be computed later.
	s.index[n] = 1
	numExits := s.graph.NumExits(n)
	for i := numExits - 1; i >= 0; i-- {
		dst := s.graph.GetExit(n, i)
		if dst != NoNode {
			s.search(dst)
		}
	}
	s.order = append(s.order, n)
	// Zero is reserved, so use the actual index + 1.
	s.index[n] = len(s.order)
}

func ReversePostorderX(g *Graph) ([]NodeID, []int) {
	numNodes := len(g.nodes)
	s := &rpoSearch{
		graph: g,
		order: make([]NodeID, 1, numNodes),
		index: make([]int, numNodes),
	}
	// Implicit edge from entry to exit ensures exit will always be the last node.
	s.order[0] = g.Exit()
	s.index[g.Exit()] = 1
	s.search(g.Entry())

	o := len(s.order)
	for i := 0; i < o/2; i++ {
		s.order[i], s.order[o-1-i] = s.order[o-1-i], s.order[i]
	}
	for i := 0; i < numNodes; i++ {
		if s.index[i] != 0 {
			s.index[i] = o - s.index[i]
		} else {
			s.index[i] = -1
		}
	}
	return s.order, s.index
}

func ToNodeListHACK(g *Graph, order []NodeID) []*Node {
	nodes := make([]*Node, len(order))
	for i, n := range order {
		nodes[i] = g.ResolveNodeHACK(n)
		nodes[i].Name = i
	}
	return nodes
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

func intersectDom(idoms []NodeID, index []int, n0 NodeID, n1 NodeID) NodeID {
	i0 := index[n0]
	i1 := index[n1]
	for i0 != i1 {
		for i0 > i1 {
			n0 = idoms[n0]
			i0 = index[n0]
		}
		for i0 < i1 {
			n1 = idoms[n1]
			i1 = index[n1]
		}
	}
	return n0
}

func isBackedge(index []int, src NodeID, dst NodeID) bool {
	return index[src] > index[dst]
}

func FindDominators(g *Graph, order []NodeID, index []int) []NodeID {
	numNodes := len(g.nodes)
	idoms := make([]NodeID, numNodes)
	changed := false
	nit := OrderedIterator(order)
	for nit.Next() {
		n := nit.Value()
		// If there are no forward entries into the node, assume an impossible edge.
		new_idom := g.Entry()
		first := true
		eit := EntryIterator(g, n)
		for eit.Next() {
			e := eit.Value()
			if !isBackedge(index, e, n) {
				if first {
					new_idom = e
					first = false
				} else {
					new_idom = intersectDom(idoms, index, new_idom, e)
				}
			} else {
				changed = true
			}
		}
		idoms[n] = new_idom
	}
	for changed {
		changed = false
		nit := OrderedIterator(order)
		for nit.Next() {
			n := nit.Value()
			numEntries := g.NumEntries(n)
			// 0 and 1 entry nodes should be stable after the first pass.
			if numEntries >= 2 {
				newIdom := idoms[n]
				eit := EntryIterator(g, n)
				for eit.Next() {
					e := eit.Value()
					newIdom = intersectDom(idoms, index, newIdom, e)
				}
				if idoms[n] != newIdom {
					idoms[n] = newIdom
					changed = true
				}
			}
		}
	}
	return idoms
}

// Assumes no dead entries.
func FindDominanceFrontiers(g *Graph, idoms []NodeID) [][]NodeID {
	n := len(g.nodes)
	frontiers := make([][]NodeID, n)
	nit := NodeIterator(g)
	for nit.Next() {
		n := nit.Value()
		fmt.Println(n)
		numEntries := g.NumEntries(n)
		if numEntries >= 2 {
			target := idoms[n]
			eit := EntryIterator(g, n)
			for eit.Next() {
				runner := eit.Value()
				for runner != target {
					frontiers[runner] = append(frontiers[runner], n)
					runner = idoms[runner]
				}
			}
		}
	}
	return frontiers
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

	n := len(ordered)

	for i := 0; i < n; i++ {
		node := ordered[i]
		if node.NumEntries() == 0 {
			idoms[i] = 0
		} else {
			idoms[i] = NoNodeIndex
		}

		// Find the earliest use of this node.
		e := n
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
			new_idom := idoms[i]
			for j := 0; j < len(entries); j++ {
				other := entries[j].src.Name
				// Is it available, yet?
				if idoms[other] == NoNodeIndex {
					continue
				}
				// Is it the first we've found?
				if new_idom == NoNodeIndex {
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

func OldFindFrontiers(ordered []*Node, idoms []int) [][]int {
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

func CanonicalSet(set map[int]bool) []int {
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
	return CanonicalSet(l.liveIn[n])
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
	Idoms    []int
	df       [][]int
	PhiFuncs [][]int
	Live     LivenessOracle
}

func CreateSSIBuilder(r *Region, nodes []*Node, live LivenessOracle) *SSIBuilder {
	idoms := FindIdoms(nodes)
	df := OldFindFrontiers(nodes, idoms)
	phiFuncs := make([][]int, len(nodes))
	return &SSIBuilder{
		nodes:    nodes,
		Idoms:    idoms,
		df:       df,
		PhiFuncs: phiFuncs,
		Live:     live,
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
	if !state.builder.Live.LiveAtEntry(node, state.uid) {
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

func NodeDotID(node *Node) string {
	return fmt.Sprintf("n%d", node.Name)
}

type DotStyler interface {
	NodeStyle(node interface{}) string
	EdgeStyle(node interface{}, flow int) string
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
		buf.WriteString(NodeDotID(node))
		buf.WriteString("[")
		buf.WriteString(styler.NodeStyle(node.Data))
		buf.WriteString("];\n")

		for i := 0; i < node.NumExits(); i++ {
			dst := node.GetNext(i)
			if dst != nil {
				buf.WriteString("  ")
				buf.WriteString(NodeDotID(node))
				buf.WriteString(" -> ")
				buf.WriteString(NodeDotID(dst))
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
				buf.WriteString(NodeDotID(nodes[i]))
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

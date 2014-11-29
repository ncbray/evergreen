package base

import (
	"bytes"
	"fmt"
	"sort"
)

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

func ReverseOrder(order []NodeID) {
	o := len(order)
	for i := 0; i < o/2; i++ {
		order[i], order[o-1-i] = order[o-1-i], order[i]
	}
}

func ReversePostorder(g *Graph) ([]NodeID, []int) {
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

	ReverseOrder(s.order)
	o := len(s.order)
	for i := 0; i < numNodes; i++ {
		if s.index[i] != 0 {
			s.index[i] = o - s.index[i]
		} else {
			s.index[i] = -1
		}
	}
	return s.order, s.index
}

type DefUseCollector struct {
	VarUseAt [][]NodeID
	VarDefAt [][]NodeID
	NodeUses [][]int
	NodeDefs [][]int
}

func (c *DefUseCollector) AddUse(n NodeID, v int) {
	c.VarUseAt[v] = append(c.VarUseAt[v], n)
	c.NodeUses[n] = append(c.NodeUses[n], v)
}

func (c *DefUseCollector) AddDef(n NodeID, v int) {
	c.VarDefAt[v] = append(c.VarDefAt[v], n)
	c.NodeDefs[n] = append(c.NodeDefs[n], v)
}

func (c *DefUseCollector) IsDefined(n NodeID, v int) bool {
	for _, d := range c.NodeDefs[n] {
		if d == v {
			return true
		}
	}
	return false
}

func CreateDefUse(numNodes int, numVars int) *DefUseCollector {
	return &DefUseCollector{
		VarUseAt: make([][]NodeID, numVars),
		VarDefAt: make([][]NodeID, numVars),
		NodeUses: make([][]int, numNodes),
		NodeDefs: make([][]int, numNodes),
	}
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

type LivenessOracle interface {
	LiveAtEntry(n NodeID, v int) bool
	LiveAtExit(n NodeID, v int) bool
}

type SimpleLivenessOracle struct {
}

func (l *SimpleLivenessOracle) LiveAtEntry(n NodeID, v int) bool {
	return true
}

func (l *SimpleLivenessOracle) LiveAtExit(n NodeID, v int) bool {
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

func (l *LiveVars) LiveSet(n NodeID) []int {
	return CanonicalSet(l.liveIn[n])
}

func (l *LiveVars) LiveAtEntry(n NodeID, v int) bool {
	live, _ := l.liveIn[n][v]
	return live
}

func (l *LiveVars) LiveAtExit(n NodeID, v int) bool {
	live, _ := l.liveOut[n][v]
	return live
}

func FindLiveVars(g *Graph, defuse *DefUseCollector) *LiveVars {
	// TODO actual backwards order?
	order, _ := ReversePostorder(g)
	ReverseOrder(order)

	n := len(order)
	liveIn := make([]map[int]bool, n)
	liveOut := make([]map[int]bool, n)
	// Initialize with the uses for each node.
	nit := NodeIterator(g)
	for nit.Next() {
		n := nit.Value()
		liveIn[n] = map[int]bool{}
		liveOut[n] = map[int]bool{}
		for _, v := range defuse.NodeUses[n] {
			liveIn[n][v] = true
		}
	}
	// Iterate until a stable state is reached.
	changed := true
	for changed {
		changed = false
		// Propagate the uses backwards.
		nit := OrderedIterator(order)
		for nit.Next() {
			n := nit.Value()
			eit := ExitIterator(g, n)
			for eit.Next() {
				dst := eit.Value()
				// Merge sets from predecessors.
				for v, _ := range liveIn[dst] {
					_, exists := liveOut[n][v]
					if !exists {
						liveOut[n][v] = true
						changed = true
						// Filter out local defines
						if !defuse.IsDefined(n, v) {
							liveIn[n][v] = true
						}
					}
				}
			}
		}
	}
	return &LiveVars{liveIn: liveIn, liveOut: liveOut}
}

type SSIBuilder struct {
	graph    *Graph
	order    []NodeID
	Idoms    []NodeID
	df       [][]NodeID
	PhiFuncs [][]int
	Live     LivenessOracle
}

func CreateSSIBuilder(g *Graph, live LivenessOracle) *SSIBuilder {
	order, index := ReversePostorder(g)
	idoms := FindDominators(g, order, index)
	df := FindDominanceFrontiers(g, idoms)
	phiFuncs := make([][]int, len(g.nodes))
	return &SSIBuilder{
		graph:    g,
		order:    order,
		Idoms:    idoms,
		df:       df,
		PhiFuncs: phiFuncs,
		Live:     live,
	}
}

type SSIState struct {
	builder *SSIBuilder
	uid     int

	phiPlaced   map[NodeID]bool
	defEnqueued map[NodeID]bool
	defQueue    []NodeID

	sigmaPlaced map[NodeID]bool
	useEnqueued map[NodeID]bool
	useQueue    []NodeID
}

func CreateSSIState(builder *SSIBuilder, uid int) *SSIState {
	return &SSIState{
		builder:     builder,
		uid:         uid,
		defEnqueued: map[NodeID]bool{},
		phiPlaced:   map[NodeID]bool{},
		useEnqueued: map[NodeID]bool{},
		sigmaPlaced: map[NodeID]bool{},
	}
}

func (state *SSIState) DiscoveredDef(node NodeID) {
	enqueued, _ := state.defEnqueued[node]
	if !enqueued {
		state.defEnqueued[node] = true
		state.defQueue = append(state.defQueue, node)
	}
}

func (state *SSIState) GetNextDef() NodeID {
	current := state.defQueue[len(state.defQueue)-1]
	state.defQueue = state.defQueue[:len(state.defQueue)-1]
	return current
}

func (state *SSIState) PlacePhi(node NodeID) {
	if !state.builder.Live.LiveAtEntry(node, state.uid) {
		return
	}
	placed, _ := state.phiPlaced[node]
	if !placed {
		state.builder.PhiFuncs[node] = append(state.builder.PhiFuncs[node], state.uid)
		state.phiPlaced[node] = true
		state.DiscoveredDef(node)
		eit := EntryIterator(state.builder.graph, node)
		for eit.Next() {
			state.DiscoveredUse(eit.Value())
		}
	}

}

func (state *SSIState) DiscoveredUse(node NodeID) {
	enqueued, _ := state.useEnqueued[node]
	if !enqueued {
		state.useEnqueued[node] = true
		state.useQueue = append(state.useQueue, node)
	}
}

func SSI(builder *SSIBuilder, uid int, defs []NodeID) {
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

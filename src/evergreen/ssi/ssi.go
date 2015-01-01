// Utilities for performing a SSI transform.
package ssi

import (
	"evergreen/graph"
	"sort"
)

type DefUseCollector struct {
	VarUseAt [][]graph.NodeID
	VarDefAt [][]graph.NodeID
	NodeUses [][]int
	NodeDefs [][]int
}

func (c *DefUseCollector) AddUse(n graph.NodeID, v int) {
	c.VarUseAt[v] = append(c.VarUseAt[v], n)
	c.NodeUses[n] = append(c.NodeUses[n], v)
}

func (c *DefUseCollector) AddDef(n graph.NodeID, v int) {
	c.VarDefAt[v] = append(c.VarDefAt[v], n)
	c.NodeDefs[n] = append(c.NodeDefs[n], v)
}

func (c *DefUseCollector) IsDefined(n graph.NodeID, v int) bool {
	for _, d := range c.NodeDefs[n] {
		if d == v {
			return true
		}
	}
	return false
}

func CreateDefUse(numNodes int, numVars int) *DefUseCollector {
	return &DefUseCollector{
		VarUseAt: make([][]graph.NodeID, numVars),
		VarDefAt: make([][]graph.NodeID, numVars),
		NodeUses: make([][]int, numNodes),
		NodeDefs: make([][]int, numNodes),
	}
}

type LivenessOracle interface {
	LiveAtEntry(n graph.NodeID, v int) bool
	LiveAtExit(n graph.NodeID, v int) bool
}

type SimpleLivenessOracle struct {
}

func (l *SimpleLivenessOracle) LiveAtEntry(n graph.NodeID, v int) bool {
	return true
}

func (l *SimpleLivenessOracle) LiveAtExit(n graph.NodeID, v int) bool {
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

func (l *LiveVars) LiveSet(n graph.NodeID) []int {
	return canonicalSet(l.liveIn[n])
}

func (l *LiveVars) LiveAtEntry(n graph.NodeID, v int) bool {
	live, _ := l.liveIn[n][v]
	return live
}

func (l *LiveVars) LiveAtExit(n graph.NodeID, v int) bool {
	live, _ := l.liveOut[n][v]
	return live
}

func FindLiveVars(g *graph.Graph, defuse *DefUseCollector) *LiveVars {
	// TODO actual backwards order?
	order, _ := graph.ReversePostorder(g)
	graph.ReverseOrder(order)

	n := len(order)
	liveIn := make([]map[int]bool, n)
	liveOut := make([]map[int]bool, n)
	// Initialize with the uses for each node.
	nit := g.NodeIterator()
	for nit.HasNext() {
		n := nit.GetNext()
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
		nit := graph.OrderedIterator(order)
		for nit.HasNext() {
			n := nit.GetNext()
			eit := g.ExitIterator(n)
			for eit.HasNext() {
				_, dst := eit.GetNext()
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
	graph    *graph.Graph
	order    []graph.NodeID
	Idoms    []graph.NodeID
	df       [][]graph.NodeID
	PhiFuncs [][]int
	Live     LivenessOracle
}

func CreateSSIBuilder(g *graph.Graph, live LivenessOracle) *SSIBuilder {
	order, index := graph.ReversePostorder(g)
	idoms := graph.FindDominators(g, order, index)
	df := graph.FindDominanceFrontiers(g, idoms)
	phiFuncs := make([][]int, g.NumNodes())
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

	phiPlaced   map[graph.NodeID]bool
	defEnqueued map[graph.NodeID]bool
	defQueue    []graph.NodeID

	sigmaPlaced map[graph.NodeID]bool
	useEnqueued map[graph.NodeID]bool
	useQueue    []graph.NodeID
}

func CreateSSIState(builder *SSIBuilder, uid int) *SSIState {
	return &SSIState{
		builder:     builder,
		uid:         uid,
		defEnqueued: map[graph.NodeID]bool{},
		phiPlaced:   map[graph.NodeID]bool{},
		useEnqueued: map[graph.NodeID]bool{},
		sigmaPlaced: map[graph.NodeID]bool{},
	}
}

func (state *SSIState) DiscoveredDef(node graph.NodeID) {
	enqueued, _ := state.defEnqueued[node]
	if !enqueued {
		state.defEnqueued[node] = true
		state.defQueue = append(state.defQueue, node)
	}
}

func (state *SSIState) GetNextDef() graph.NodeID {
	current := state.defQueue[len(state.defQueue)-1]
	state.defQueue = state.defQueue[:len(state.defQueue)-1]
	return current
}

func (state *SSIState) PlacePhi(node graph.NodeID) {
	if !state.builder.Live.LiveAtEntry(node, state.uid) {
		return
	}
	placed, _ := state.phiPlaced[node]
	if !placed {
		state.builder.PhiFuncs[node] = append(state.builder.PhiFuncs[node], state.uid)
		state.phiPlaced[node] = true
		state.DiscoveredDef(node)
		eit := state.builder.graph.EntryIterator(node)
		for eit.HasNext() {
			e, _ := eit.GetNext()
			state.DiscoveredUse(e)
		}
	}

}

func (state *SSIState) DiscoveredUse(node graph.NodeID) {
	enqueued, _ := state.useEnqueued[node]
	if !enqueued {
		state.useEnqueued[node] = true
		state.useQueue = append(state.useQueue, node)
	}
}

func SSI(builder *SSIBuilder, uid int, defs []graph.NodeID) {
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

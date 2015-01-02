// Package graph implements a directed graph and common operations on directed graphs.
package graph

type NodeID uint32

const NoNode NodeID = ^NodeID(0)

type EdgeID uint32

const NoEdge EdgeID = ^EdgeID(0)

type node struct {
	entries entryEdges
	exits   exitEdges
}

type edge struct {
	src      NodeID
	nextExit EdgeID
	prevExit EdgeID

	dst       NodeID
	nextEntry EdgeID
	prevEntry EdgeID
}

type entryEdges struct {
	head EdgeID
	tail EdgeID
}

func emptyEntry() entryEdges {
	return entryEdges{head: NoEdge, tail: NoEdge}
}

func singleEntry(g *Graph, eid EdgeID) entryEdges {
	entry := emptyEntry()
	entry.Append(g, eid)
	return entry
}

// An incomplete sanity check, this edge could still be in a list that
// is one element long.
func orphanedEntry(g *Graph, eid EdgeID) bool {
	return g.edges[eid].nextEntry == NoEdge && g.edges[eid].prevEntry == NoEdge
}

func (l *entryEdges) Append(g *Graph, other EdgeID) {
	if !orphanedEntry(g, other) {
		panic(other)
	}
	if l.tail != NoEdge {
		g.edges[other].prevEntry = l.tail
		g.edges[l.tail].nextEntry = other
		l.tail = other
	} else {
		l.head = other
		l.tail = other
	}
}

func (l *entryEdges) Extend(g *Graph, other entryEdges) {
	if other.HasEdges() {
		if l.HasEdges() {
			g.edges[other.head].prevEntry = l.tail
			g.edges[l.tail].nextEntry = other.head
			l.tail = other.tail
		} else {
			l.head = other.head
			l.tail = other.tail
		}
	}
}

func (l *entryEdges) SetDst(g *Graph, n NodeID) {
	current := l.head
	if current != NoEdge && g.edges[current].dst != n {
		for current != NoEdge {
			g.edges[current].dst = n
			current = g.edges[current].nextEntry
		}
	}
}

func (l *entryEdges) Transfer() entryEdges {
	temp := *l
	*l = emptyEntry()
	return temp
}

func (l *entryEdges) Remove(g *Graph, target EdgeID) {
	oldNext := g.edges[target].nextEntry
	oldPrev := g.edges[target].prevEntry

	// Disconnect the entry.
	g.edges[target].prevEntry = NoEdge
	g.edges[target].nextEntry = NoEdge
	g.edges[target].dst = NoNode

	if oldPrev == NoEdge {
		l.head = oldNext
	} else {
		g.edges[oldPrev].nextEntry = oldNext
	}
	if oldNext == NoEdge {
		l.tail = oldPrev
	} else {
		g.edges[oldNext].prevEntry = oldPrev
	}
}

func (l *entryEdges) ReplaceEdge(g *Graph, target EdgeID, replacement EdgeID) {
	if !orphanedEntry(g, replacement) {
		panic(replacement)
	}
	l.ReplaceEdgeWithMultiple(g, target, singleEntry(g, replacement))
}

func (l *entryEdges) ReplaceEdgeWithMultiple(g *Graph, target EdgeID, replacements entryEdges) {
	oldNext := g.edges[target].nextEntry
	oldPrev := g.edges[target].prevEntry
	dst := g.edges[target].dst

	// Disconnect the entry.
	g.edges[target].prevEntry = NoEdge
	g.edges[target].nextEntry = NoEdge
	g.edges[target].dst = NoNode

	var newNext EdgeID
	var newPrev EdgeID
	if !replacements.HasEdges() {
		// Empty replacement.
		newNext = oldNext
		newPrev = oldPrev
	} else {
		// Real replacement.
		newNext = replacements.head
		newPrev = replacements.tail
		g.edges[replacements.head].prevEntry = oldPrev
		g.edges[replacements.tail].nextEntry = oldNext

		// Point the replacements to the new destination.
		replacements.SetDst(g, dst)
	}
	if oldPrev == NoEdge {
		l.head = newNext
	} else {
		g.edges[oldPrev].nextEntry = newNext
	}
	if oldNext == NoEdge {
		l.tail = newPrev
	} else {
		g.edges[oldNext].prevEntry = newPrev
	}
}

func (l *entryEdges) HasEdges() bool {
	return l.head != NoEdge
}

func (l *entryEdges) HasMultipleEdges() bool {
	return l.head != l.tail
}

type entryIterator struct {
	graph   *Graph
	current EdgeID
}

func (it *entryIterator) HasNext() bool {
	return it.current != NoEdge
}

func (it *entryIterator) GetNext() (NodeID, EdgeID) {
	edge := it.current
	it.current = it.graph.edges[edge].nextEntry
	return it.graph.edges[edge].src, edge
}

type exitEdges struct {
	head EdgeID
	tail EdgeID
}

func emptyExit() exitEdges {
	return exitEdges{head: NoEdge, tail: NoEdge}
}

// An incomplete sanity check, this edge could still be in a list that
// is one element long.
func orphanedExit(g *Graph, eid EdgeID) bool {
	return g.edges[eid].nextExit == NoEdge && g.edges[eid].prevExit == NoEdge
}

func (l *exitEdges) Append(g *Graph, other EdgeID) {
	if !orphanedEntry(g, other) {
		panic(other)
	}
	if l.tail != NoEdge {
		g.edges[other].prevExit = l.tail
		g.edges[l.tail].nextExit = other
		l.tail = other
	} else {
		l.head = other
		l.tail = other
	}
}

func (l *exitEdges) Extend(g *Graph, other exitEdges) {
	if other.HasEdges() {
		if l.HasEdges() {
			g.edges[other.head].prevExit = l.tail
			g.edges[l.tail].nextExit = other.head
			l.tail = other.tail
		} else {
			l.head = other.head
			l.tail = other.tail
		}
	}
}

func (l *exitEdges) Remove(g *Graph, target EdgeID) {
	oldNext := g.edges[target].nextExit
	oldPrev := g.edges[target].prevExit

	// Disconnect the exit.
	g.edges[target].prevExit = NoEdge
	g.edges[target].nextExit = NoEdge
	g.edges[target].src = NoNode

	if oldPrev == NoEdge {
		l.head = oldNext
	} else {
		g.edges[oldPrev].nextExit = oldNext
	}
	if oldNext == NoEdge {
		l.tail = oldPrev
	} else {
		g.edges[oldNext].prevExit = oldPrev
	}
}

func (l *exitEdges) HasEdges() bool {
	return l.head != NoEdge
}

func (l *exitEdges) HasMultipleEdges() bool {
	return l.head != l.tail
}

type exitIterator struct {
	graph   *Graph
	current EdgeID
}

func (it *exitIterator) HasNext() bool {
	return it.current != NoEdge
}

func (it *exitIterator) GetNext() (EdgeID, NodeID) {
	edge := it.current
	it.current = it.graph.edges[edge].nextExit
	return edge, it.graph.edges[edge].dst
}

type reverseExitIterator struct {
	graph   *Graph
	current EdgeID
}

func (it *reverseExitIterator) HasNext() bool {
	return it.current != NoEdge
}

func (it *reverseExitIterator) GetNext() (EdgeID, NodeID) {
	edge := it.current
	it.current = it.graph.edges[edge].prevExit
	return edge, it.graph.edges[edge].dst
}

type Graph struct {
	nodes []node
	edges []edge
}

func (g *Graph) Entry() NodeID {
	return 0
}

func (g *Graph) Exit() NodeID {
	return 1
}

func (g *Graph) CreateEdge() EdgeID {
	id := EdgeID(len(g.edges))
	g.edges = append(g.edges, edge{
		src:       NoNode,
		nextExit:  NoEdge,
		prevExit:  NoEdge,
		dst:       NoNode,
		nextEntry: NoEdge,
		prevEntry: NoEdge,
	})
	return id
}

func (g *Graph) CreateNode() NodeID {
	id := NodeID(len(g.nodes))
	g.nodes = append(g.nodes, node{
		entries: emptyEntry(),
		exits:   emptyExit(),
	})
	return id
}

// Number of nodes in existance, some may be dead.
func (g *Graph) NumNodes() int {
	return len(g.nodes)
}

// Number of edges in existance, some may be dead.
func (g *Graph) NumEdges() int {
	return len(g.edges)
}

func (g *Graph) EdgeEntry(e EdgeID) NodeID {
	return g.edges[e].src
}

func (g *Graph) EdgeExit(e EdgeID) NodeID {
	return g.edges[e].dst
}

func (g *Graph) setEdgeExit(e EdgeID, n NodeID) {
	if g.edges[e].dst != NoNode {
		panic(e)
	}
	g.edges[e].dst = n
	g.nodes[n].entries.Append(g, e)
}

func (g *Graph) setEdgeEntry(n NodeID, e EdgeID) {
	if g.edges[e].src != NoNode {
		panic(e)
	}
	g.edges[e].src = n
	g.nodes[n].exits.Append(g, e)
}

func (g *Graph) ConnectEdge(src NodeID, e EdgeID, dst NodeID) {
	g.setEdgeEntry(src, e)
	g.setEdgeExit(e, dst)
}

func (g *Graph) ConnectEdgeEntry(src NodeID, e EdgeID) {
	g.setEdgeEntry(src, e)
}

func (g *Graph) ConnectEdgeExit(e EdgeID, dst NodeID) {
	g.setEdgeExit(e, dst)
}

func (g *Graph) KillEdge(eid EdgeID) {
	src := g.edges[eid].src
	if src != NoNode {
		g.nodes[src].exits.Remove(g, eid)
	}
	dst := g.edges[eid].dst
	if dst != NoNode {
		g.nodes[dst].entries.Remove(g, eid)
	}
}

func getUniqueExit(exits exitEdges) EdgeID {
	if exits.head == exits.tail {
		return exits.head
	}
	return NoEdge
}

func (g *Graph) KillNode(n NodeID) {
	singleEdge := getUniqueExit(g.nodes[n].exits)
	if singleEdge == NoEdge {
		panic(n)
	}
	dst := g.edges[singleEdge].dst
	g.nodes[dst].entries.ReplaceEdgeWithMultiple(g, singleEdge, g.nodes[n].entries.Transfer())
	g.nodes[n].exits.Remove(g, singleEdge)
}

// Insert a dangling node (with a single out edge) in the middle of an existing edge.
func (g *Graph) InsertInEdge(replacement EdgeID, existing EdgeID) {
	insertedNode := g.edges[replacement].src
	dstNode := g.edges[existing].dst
	if dstNode != NoNode {
		// Replace the exising edge with the new edge.
		g.nodes[dstNode].entries.ReplaceEdge(g, existing, replacement)
	}
	// Attach the existing edge to the dangling node.
	g.edges[existing].dst = insertedNode
	g.nodes[insertedNode].entries.Append(g, existing)
}

func (g *Graph) HasMultipleEntries(dst NodeID) bool {
	return g.nodes[dst].entries.HasMultipleEdges()
}

func (g *Graph) GetUniqueExit(src NodeID) (EdgeID, NodeID) {
	e := getUniqueExit(g.nodes[src].exits)
	if e != NoEdge {
		dst := g.EdgeExit(e)
		if dst != NoNode {
			return e, dst
		}
	}
	return NoEdge, NoNode
}

func (g *Graph) extendEntries(n NodeID, entries entryEdges) {
	g.nodes[n].entries.Extend(g, entries)
}

func CreateGraph() *Graph {
	g := &Graph{}
	// Entry
	g.CreateNode()
	// Exit
	g.CreateNode()
	return g
}

type FlowBuilder struct {
	graph *Graph
	flows []entryEdges
}

func CreateFlowBuilder(g *Graph, e EdgeID, numExits int) *FlowBuilder {
	return createFlowBuilder(g, numExits, singleEntry(g, e))
}

func createFlowBuilder(g *Graph, numExits int, edges entryEdges) *FlowBuilder {
	flows := make([]entryEdges, numExits)
	for i := 0; i < numExits; i++ {
		flows[i] = emptyEntry()
	}
	fb := &FlowBuilder{
		graph: g,
		flows: flows,
	}
	fb.flows[0] = edges
	return fb
}

func (fb *FlowBuilder) HasFlow(flow int) bool {
	return fb.flows[flow].HasEdges()
}

func (fb *FlowBuilder) AttachFlow(flow int, dst NodeID) {
	if !fb.HasFlow(flow) {
		panic("Tried to attach non-existant flow")
	}
	g := fb.graph
	fb.flows[flow].SetDst(g, dst)
	g.extendEntries(dst, fb.flows[flow].Transfer())
}

func (fb *FlowBuilder) RegisterExit(eid EdgeID, flow int) {
	g := fb.graph
	if g.edges[eid].src == NoNode {
		panic(eid)
	}
	if g.edges[eid].dst != NoNode {
		panic(eid)
	}
	fb.flows[flow].Append(g, eid)
}

func (fb *FlowBuilder) SplitOffFlow(flow int) *FlowBuilder {
	return createFlowBuilder(fb.graph, len(fb.flows), fb.flows[flow].Transfer())
}

func (fb *FlowBuilder) SplitOffEdge(eid EdgeID) *FlowBuilder {
	return createFlowBuilder(fb.graph, len(fb.flows), singleEntry(fb.graph, eid))
}

func (fb *FlowBuilder) AbsorbExits(other *FlowBuilder) {
	for i := 0; i < len(fb.flows); i++ {
		fb.flows[i].Extend(fb.graph, other.flows[i].Transfer())
	}
}

type nodeIterator struct {
	current int
	count   int
}

func (it *nodeIterator) HasNext() bool {
	return it.current < it.count
}

func (it *nodeIterator) GetNext() NodeID {
	next := NodeID(it.current)
	it.current += 1
	return next
}

// The node length is intentionally snapshotted incase the iteration creates new nodes.
func (g *Graph) NodeIterator() nodeIterator {
	return nodeIterator{count: len(g.nodes), current: 0}
}

type orderedNodeIterator struct {
	current int
	order   []NodeID
}

func (it *orderedNodeIterator) HasNext() bool {
	return it.current < len(it.order)
}

func (it *orderedNodeIterator) GetNext() NodeID {
	value := it.order[it.current]
	it.current += 1
	return value
}

func OrderedIterator(order []NodeID) orderedNodeIterator {
	return orderedNodeIterator{order: order, current: 0}
}

func (g *Graph) EntryIterator(n NodeID) entryIterator {
	return entryIterator{graph: g, current: g.nodes[n].entries.head}
}

func (g *Graph) ExitIterator(n NodeID) exitIterator {
	iter := exitIterator{graph: g, current: g.nodes[n].exits.head}
	return iter
}

func (g *Graph) ReverseExitIterator(n NodeID) reverseExitIterator {
	iter := reverseExitIterator{graph: g, current: g.nodes[n].exits.tail}
	return iter
}

// Package graph implements a directed graph and common operations on directed graphs.
package graph

type edge struct {
	src      *node
	nextExit *edge
	prevExit *edge

	dst       *node
	nextEntry *edge
	prevEntry *edge

	id   EdgeID
	flow int
}

func (e *edge) setExit(other *node) {
	if e.dst != nil {
		panic(e)
	}
	e.dst = other
	other.entries.Append(e)
}

func (e *edge) setEntry(other *node) {
	if e.src != nil {
		panic(e)
	}
	e.src = other
	other.exits.Append(e)
}

func (e *edge) detach() {
	if e.src != nil {
		e.src.exits.Remove(e)
	}
	if e.dst != nil {
		e.dst.entries.Remove(e)
	}
}

type entryEdges struct {
	head *edge
	tail *edge
}

func emptyEntry() entryEdges {
	return entryEdges{}
}

func singleEntry(e *edge) entryEdges {
	entry := emptyEntry()
	entry.Append(e)
	return entry
}

// An incomplete sanity check, this edge could still be in a list that
// is one element long.
func orphanedEntry(e *edge) bool {
	return e.nextEntry == nil && e.prevEntry == nil
}

func (l *entryEdges) Append(other *edge) {
	if !orphanedEntry(other) {
		panic(other)
	}
	if l.tail != nil {
		other.prevEntry = l.tail
		l.tail.nextEntry = other
		l.tail = other
	} else {
		l.head = other
		l.tail = other
	}
}

func (l *entryEdges) Extend(other entryEdges) {
	if other.HasEdges() {
		if l.HasEdges() {
			other.head.prevEntry = l.tail
			l.tail.nextEntry = other.head
			l.tail = other.tail
		} else {
			l.head = other.head
			l.tail = other.tail
		}
	}
}

func (l *entryEdges) SetDst(n *node) {
	current := l.head
	if current != nil && current.dst != n {
		for current != nil {
			current.dst = n
			current = current.nextEntry
		}
	}
}

func (l *entryEdges) Transfer() entryEdges {
	temp := *l
	*l = emptyEntry()
	return temp
}

func (l *entryEdges) Remove(target *edge) {
	oldNext := target.nextEntry
	oldPrev := target.prevEntry

	// Disconnect the entry.
	target.prevEntry = nil
	target.nextEntry = nil
	target.dst = nil

	if oldPrev == nil {
		l.head = oldNext
	} else {
		oldPrev.nextEntry = oldNext
	}
	if oldNext == nil {
		l.tail = oldPrev
	} else {
		oldNext.prevEntry = oldPrev
	}
}

func (l *entryEdges) ReplaceEdge(target *edge, replacement *edge) {
	if !orphanedEntry(replacement) {
		panic(replacement)
	}
	l.ReplaceEdgeWithMultiple(target, singleEntry(replacement))
}

func (l *entryEdges) ReplaceEdgeWithMultiple(target *edge, replacements entryEdges) {
	oldNext := target.nextEntry
	oldPrev := target.prevEntry
	dst := target.dst

	// Disconnect the entry.
	target.prevEntry = nil
	target.nextEntry = nil
	target.dst = nil

	var newNext *edge
	var newPrev *edge
	if !replacements.HasEdges() {
		// Empty replacement.
		newNext = oldNext
		newPrev = oldPrev
	} else {
		// Real replacement.
		newNext = replacements.head
		newPrev = replacements.tail
		replacements.head.prevEntry = oldPrev
		replacements.tail.nextEntry = oldNext

		// Point the replacements to the new destination.
		replacements.SetDst(dst)
	}
	if oldPrev == nil {
		l.head = newNext
	} else {
		oldPrev.nextEntry = newNext
	}
	if oldNext == nil {
		l.tail = newPrev
	} else {
		oldNext.prevEntry = newPrev
	}
}

func (l *entryEdges) HasEdges() bool {
	return l.head != nil
}

func (l *entryEdges) HasMultipleEdges() bool {
	return l.head != l.tail
}

type entryIterator struct {
	graph   *Graph
	current *edge
}

func (it *entryIterator) HasNext() bool {
	return it.current != nil
}

func (it *entryIterator) GetNext() (NodeID, EdgeID) {
	edge := it.current
	it.current = it.current.nextEntry
	n := NoNode
	if edge.src != nil {
		n = edge.src.id
	}
	return n, edge.id
}

type exitEdges struct {
	head *edge
	tail *edge
}

func emptyExit() exitEdges {
	return exitEdges{}
}

// An incomplete sanity check, this edge could still be in a list that
// is one element long.
func orphanedExit(e *edge) bool {
	return e.nextExit == nil && e.prevExit == nil
}

func (l *exitEdges) Append(other *edge) {
	if !orphanedEntry(other) {
		panic(other)
	}
	if l.tail != nil {
		other.prevExit = l.tail
		l.tail.nextExit = other
		l.tail = other
	} else {
		l.head = other
		l.tail = other
	}
}

func (l *exitEdges) Extend(other exitEdges) {
	if other.HasEdges() {
		if l.HasEdges() {
			other.head.prevExit = l.tail
			l.tail.nextExit = other.head
			l.tail = other.tail
		} else {
			l.head = other.head
			l.tail = other.tail
		}
	}
}

func (l *exitEdges) Remove(target *edge) {
	oldNext := target.nextExit
	oldPrev := target.prevExit

	// Disconnect the exit.
	target.prevExit = nil
	target.nextExit = nil
	target.src = nil

	if oldPrev == nil {
		l.head = oldNext
	} else {
		oldPrev.nextExit = oldNext
	}
	if oldNext == nil {
		l.tail = oldPrev
	} else {
		oldNext.prevExit = oldPrev
	}
}

func (l *exitEdges) HasEdges() bool {
	return l.head != nil
}

func (l *exitEdges) HasMultipleEdges() bool {
	return l.head != l.tail
}

type exitIterator struct {
	graph   *Graph
	current *edge
}

func (it *exitIterator) HasNext() bool {
	return it.current != nil
}

func (it *exitIterator) GetNext() (EdgeID, NodeID) {
	edge := it.current
	it.current = edge.nextExit
	n := NoNode
	if edge.dst != nil {
		n = edge.dst.id
	}
	return edge.id, n
}

type reverseExitIterator struct {
	graph   *Graph
	current *edge
}

func (it *reverseExitIterator) HasNext() bool {
	return it.current != nil
}

func (it *reverseExitIterator) GetNext() (EdgeID, NodeID) {
	edge := it.current
	it.current = edge.prevExit
	n := NoNode
	if edge.dst != nil {
		n = edge.dst.id
	}
	return edge.id, n
}

type node struct {
	entries entryEdges
	exits   exitEdges
	id      NodeID
}

func (n *node) replaceSingleEntry(target *edge, replacement *edge) {
	n.entries.ReplaceEdge(target, replacement)
}

func (n *node) replaceEntry(target *edge, replacements entryEdges) {
	n.entries.ReplaceEdgeWithMultiple(target, replacements)
}

type NodeID int

const NoNode NodeID = ^NodeID(0)

type EdgeID int

const NoEdge EdgeID = ^EdgeID(0)

type Graph struct {
	nodes []*node
	edges []*edge
}

func (g *Graph) Entry() NodeID {
	return 0
}

func (g *Graph) Exit() NodeID {
	return 1
}

func (g *Graph) CreateEdge(flow int) EdgeID {
	id := EdgeID(len(g.edges))
	e := &edge{id: id, flow: flow}
	g.edges = append(g.edges, e)
	return id
}

func (g *Graph) CreateNode() NodeID {
	id := NodeID(len(g.nodes))
	n := &node{
		id:      id,
		entries: emptyEntry(),
		exits:   emptyExit(),
	}
	g.nodes = append(g.nodes, n)
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

func (g *Graph) ConnectEdge(src NodeID, e EdgeID, dst NodeID) {
	tgt := g.edges[e]
	tgt.setEntry(g.nodes[src])
	tgt.setExit(g.nodes[dst])
}

func (g *Graph) ConnectEdgeEntry(src NodeID, e EdgeID) {
	g.edges[e].setEntry(g.nodes[src])
}

func (g *Graph) ConnectEdgeExit(e EdgeID, dst NodeID) {
	g.edges[e].setExit(g.nodes[dst])
}

func (g *Graph) KillEdge(eid EdgeID) {
	g.edges[eid].detach()
}

func getUniqueExit(exits exitEdges) *edge {
	if exits.head == exits.tail {
		// Note: could be nil
		return exits.head
	}
	return nil
}

func (g *Graph) KillNode(nid NodeID) {
	n := g.nodes[nid]
	singleEdge := getUniqueExit(n.exits)
	if singleEdge == nil {
		panic(n.id)
	}
	singleEdge.dst.replaceEntry(singleEdge, n.entries.Transfer())
	n.exits.Remove(singleEdge)
}

// Insert a dangling node (with a single out edge) in the middle of an existing edge.
func (g *Graph) InsertInEdge(replacementID EdgeID, existingID EdgeID) {
	replacement := g.edges[replacementID]
	existing := g.edges[existingID]
	insertedNode := replacement.src
	dstNode := existing.dst
	if dstNode != nil {
		// Replace the exising edge with the new edge.
		dstNode.replaceSingleEntry(existing, replacement)
	}
	// Attach the existing edge to the dangling node.
	existing.dst = insertedNode
	insertedNode.entries.Append(existing)
}

func (g *Graph) HasMultipleEntries(dst NodeID) bool {
	return g.nodes[dst].entries.HasMultipleEdges()
}

func (g *Graph) GetUniqueExit(src NodeID) (EdgeID, NodeID) {
	e := getUniqueExit(g.nodes[src].exits)
	if e != nil && e.dst != nil {
		return e.id, e.dst.id
	}
	return NoEdge, NoNode
}

func (g *Graph) EdgeFlow(e EdgeID) int {
	return g.edges[e].flow
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
	return createFlowBuilder(g, numExits, singleEntry(g.edges[e]))
}

func createFlowBuilder(g *Graph, numExits int, edges entryEdges) *FlowBuilder {
	fb := &FlowBuilder{
		graph: g,
		flows: make([]entryEdges, numExits),
	}
	fb.flows[0] = edges
	return fb
}

func (fb *FlowBuilder) HasFlow(flow int) bool {
	return fb.flows[flow].HasEdges()
}

func (fb *FlowBuilder) AttachFlow(flow int, dst NodeID) {
	dstNode := fb.graph.nodes[dst]
	if !fb.HasFlow(flow) {
		panic("Tried to attach non-existant flow")
	}
	fb.flows[flow].SetDst(dstNode)
	dstNode.entries.Extend(fb.flows[flow].Transfer())
}

func (fb *FlowBuilder) RegisterExit(eid EdgeID, flow int) {
	e := fb.graph.edges[eid]
	if e.src == nil {
		panic(e)
	}
	if e.dst != nil {
		panic(e)
	}
	fb.flows[flow].Append(e)
}

func (fb *FlowBuilder) SplitOffFlow(flow int) *FlowBuilder {
	return createFlowBuilder(fb.graph, len(fb.flows), fb.flows[flow].Transfer())
}

func (fb *FlowBuilder) SplitOffEdge(eid EdgeID) *FlowBuilder {
	return createFlowBuilder(fb.graph, len(fb.flows), singleEntry(fb.graph.edges[eid]))
}

func (fb *FlowBuilder) AbsorbExits(other *FlowBuilder) {
	for i := 0; i < len(fb.flows); i++ {
		fb.flows[i].Extend(other.flows[i].Transfer())
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

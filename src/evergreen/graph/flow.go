// Package graph implements a directed graph and common operations on directed graphs.
package graph

type edge struct {
	src      *node
	nextExit *edge
	prevExit *edge

	dst       *node
	nextEntry *edge
	prevEntry *edge

	id EdgeID
}

type entryEdges struct {
	head *edge
	tail *edge
}

func emptyEntry() entryEdges {
	return entryEdges{}
}

func singleEntry(g *Graph, eid EdgeID) entryEdges {
	entry := emptyEntry()
	entry.Append(g, eid)
	return entry
}

// An incomplete sanity check, this edge could still be in a list that
// is one element long.
func orphanedEntry(g *Graph, eid EdgeID) bool {
	e := g.edges[eid]
	return e.nextEntry == nil && e.prevEntry == nil
}

func (l *entryEdges) Append(g *Graph, other EdgeID) {
	if !orphanedEntry(g, other) {
		panic(other)
	}
	e := g.edges[other]
	if l.tail != nil {
		e.prevEntry = l.tail
		l.tail.nextEntry = e
		l.tail = e
	} else {
		l.head = e
		l.tail = e
	}
}

func (l *entryEdges) Extend(g *Graph, other entryEdges) {
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

func (l *entryEdges) SetDst(g *Graph, nid NodeID) {
	var n *node
	if nid != NoNode {
		n = g.nodes[nid]
	}
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

func (l *entryEdges) Remove(g *Graph, target EdgeID) {
	e := g.edges[target]

	oldNext := e.nextEntry
	oldPrev := e.prevEntry

	// Disconnect the entry.
	e.prevEntry = nil
	e.nextEntry = nil
	e.dst = nil

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

func (l *entryEdges) ReplaceEdge(g *Graph, target EdgeID, replacement EdgeID) {
	if !orphanedEntry(g, replacement) {
		panic(replacement)
	}
	l.ReplaceEdgeWithMultiple(g, target, singleEntry(g, replacement))
}

func (l *entryEdges) ReplaceEdgeWithMultiple(g *Graph, target EdgeID, replacements entryEdges) {
	e := g.edges[target]

	oldNext := e.nextEntry
	oldPrev := e.prevEntry
	dst := e.dst

	// Disconnect the entry.
	e.prevEntry = nil
	e.nextEntry = nil
	e.dst = nil

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
		replacements.SetDst(g, nodeToID(dst))
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
	return nodeToID(edge.src), edge.id
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
func orphanedExit(g *Graph, eid EdgeID) bool {
	e := g.edges[eid]
	return e.nextExit == nil && e.prevExit == nil
}

func (l *exitEdges) Append(g *Graph, other EdgeID) {
	if !orphanedEntry(g, other) {
		panic(other)
	}
	e := g.edges[other]
	if l.tail != nil {
		e.prevExit = l.tail
		l.tail.nextExit = e
		l.tail = e
	} else {
		l.head = e
		l.tail = e
	}
}

func (l *exitEdges) Extend(g *Graph, other exitEdges) {
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

func (l *exitEdges) Remove(g *Graph, target EdgeID) {
	e := g.edges[target]

	oldNext := e.nextExit
	oldPrev := e.prevExit

	// Disconnect the exit.
	e.prevExit = nil
	e.nextExit = nil
	e.src = nil

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
	return edge.id, nodeToID(edge.dst)
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
	return edge.id, nodeToID(edge.dst)
}

type node struct {
	entries entryEdges
	exits   exitEdges
	id      NodeID
}

func nodeToID(n *node) NodeID {
	if n != nil {
		return n.id
	} else {
		return NoNode
	}
}

func edgeToID(e *edge) EdgeID {
	if e != nil {
		return e.id
	} else {
		return NoEdge
	}
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

func (g *Graph) CreateEdge() EdgeID {
	id := EdgeID(len(g.edges))
	e := &edge{id: id}
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

func (g *Graph) EdgeEntry(e EdgeID) NodeID {
	return nodeToID(g.edges[e].src)
}

func (g *Graph) EdgeExit(e EdgeID) NodeID {
	return nodeToID(g.edges[e].dst)
}

func (g *Graph) setEdgeExit(eid EdgeID, nid NodeID) {
	e := g.edges[eid]
	if e.dst != nil {
		panic(e)
	}
	dst := g.nodes[nid]
	e.dst = dst
	dst.entries.Append(g, eid)
}

func (g *Graph) setEdgeEntry(nid NodeID, eid EdgeID) {
	e := g.edges[eid]
	if e.src != nil {
		panic(e)
	}
	src := g.nodes[nid]
	e.src = src
	src.exits.Append(g, eid)
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
	e := g.edges[eid]
	if e.src != nil {
		e.src.exits.Remove(g, eid)
	}
	if e.dst != nil {
		e.dst.entries.Remove(g, eid)
	}
}

func getUniqueExit(exits exitEdges) EdgeID {
	if exits.head == exits.tail {
		return edgeToID(exits.head)
	}
	return NoEdge
}

func (g *Graph) KillNode(nid NodeID) {
	n := g.nodes[nid]
	singleEdge := getUniqueExit(n.exits)
	if singleEdge == NoEdge {
		panic(n.id)
	}
	dst := g.edges[singleEdge].dst
	dst.entries.ReplaceEdgeWithMultiple(g, singleEdge, n.entries.Transfer())
	n.exits.Remove(g, singleEdge)
}

// Insert a dangling node (with a single out edge) in the middle of an existing edge.
func (g *Graph) InsertInEdge(replacementID EdgeID, existingID EdgeID) {
	replacement := g.edges[replacementID]
	existing := g.edges[existingID]
	insertedNode := replacement.src
	dstNode := existing.dst
	if dstNode != nil {
		// Replace the exising edge with the new edge.
		dstNode.entries.ReplaceEdge(g, existingID, replacementID)
	}
	// Attach the existing edge to the dangling node.
	existing.dst = insertedNode
	insertedNode.entries.Append(g, existingID)
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
	if !fb.HasFlow(flow) {
		panic("Tried to attach non-existant flow")
	}
	g := fb.graph
	fb.flows[flow].SetDst(g, dst)
	g.extendEntries(dst, fb.flows[flow].Transfer())
}

func (fb *FlowBuilder) RegisterExit(eid EdgeID, flow int) {
	g := fb.graph
	e := g.edges[eid]
	if e.src == nil {
		panic(e)
	}
	if e.dst != nil {
		panic(e)
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

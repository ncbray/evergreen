// Package graph implements a directed graph and common operations on directed graphs.
package graph

type edge struct {
	src   *node
	dst   *node
	index int
	id    EdgeID
}

func (e *edge) attach(other *node) {
	if e.dst != nil {
		panic(e)
	}
	e.dst = other
	other.addEntry(e)
}

func (e *edge) detach() {
	if e.dst != nil {
		e.dst.replaceEntry(e, nil)
	}
}

type entryEdges []*edge

func (l *entryEdges) Append(other *edge) {
	*l = append(*l, other)
}

func (l *entryEdges) Pop() entryEdges {
	temp := *l
	*l = entryEdges{}
	return temp
}

func (l *entryEdges) find(target *edge) int {
	for i, e := range *l {
		if e == target {
			return i
		}
	}
	panic(target)
}

func (l *entryEdges) ReplaceEdge(target *edge, replacement *edge) {
	i := l.find(target)
	(*l)[i] = replacement
	replacement.dst, target.dst = target.dst, nil
}

func (l *entryEdges) ReplaceEdgeWithMultiple(target *edge, replacements entryEdges) {
	i := l.find(target)
	j := i + len(replacements)

	old := l.Pop()
	// Caution: splicing arrays is tricky because "append" may mutate the first argument.
	// Avoid appending.
	result := make(entryEdges, len(old)+len(replacements)-1)
	copy(result[:i], old[:i])
	copy(result[i:j], replacements)
	copy(result[j:], old[i+1:])

	// Redirect the replacement edges to the new destination.
	dst := target.dst
	target.dst = nil
	for _, r := range replacements {
		r.dst = dst
	}

	// Commit.
	*l = result
}

func (l *entryEdges) HasMultipleEdges() bool {
	return len(*l) >= 2
}

type entryIterator struct {
	graph   *Graph
	node    *node
	current int
}

func (it *entryIterator) HasNext() bool {
	return it.current < len(it.node.entries)
}

func (it *entryIterator) GetNext() (NodeID, EdgeID) {
	edge := it.node.entries[it.current]
	it.current += 1
	return edge.src.id, edge.id
}

type exitEdges []*edge

func (l *exitEdges) Pop() exitEdges {
	temp := *l
	*l = exitEdges{}
	return temp
}

type exitIterator struct {
	graph   *Graph
	node    *node
	current int
}

func (it *exitIterator) skipDeadEdges() {
	for it.HasNext() && it.node.exits[it.current].dst == nil {
		it.current += 1
	}
}

func (it *exitIterator) HasNext() bool {
	return it.current < len(it.node.exits)
}

func (it *exitIterator) GetNext() (EdgeID, NodeID) {
	edge := it.node.exits[it.current]
	it.current += 1
	it.skipDeadEdges()
	return edge.id, edge.dst.id
}

type reverseExitIterator struct {
	graph   *Graph
	node    *node
	current int
}

func (it *reverseExitIterator) skipDeadEdges() {
	for it.HasNext() && it.node.exits[it.current].dst == nil {
		it.current -= 1
	}
}

func (it *reverseExitIterator) HasNext() bool {
	return it.current >= 0
}

func (it *reverseExitIterator) GetNext() (EdgeID, NodeID) {
	edge := it.node.exits[it.current]
	it.current -= 1
	it.skipDeadEdges()
	return edge.id, edge.dst.id
}

type node struct {
	entries entryEdges
	exits   exitEdges
	id      NodeID
}

func (n *node) getExit(flow int) *edge {
	return n.exits[flow]
}

func (n *node) addEntry(e *edge) {
	n.entries.Append(e)
}

func (n *node) popEntries() entryEdges {
	return n.entries.Pop()
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

func (g *Graph) createEdge(flow int) *edge {
	id := EdgeID(len(g.edges))
	e := &edge{index: flow, id: id}
	g.edges = append(g.edges, e)
	return e
}

func (g *Graph) CreateNode(exits int) NodeID {
	id := NodeID(len(g.nodes))
	n := &node{
		id:    id,
		exits: make(exitEdges, exits),
	}
	for i := 0; i < exits; i++ {
		e := g.createEdge(i)
		e.src = n
		n.exits[i] = e
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

func (g *Graph) ConnectEdgeExit(src EdgeID, dst NodeID) {
	g.edges[src].attach(g.nodes[dst])
}

func (g *Graph) KillEdge(eid EdgeID) {
	g.edges[eid].detach()
}

func (g *Graph) KillNode(nid NodeID) {
	n := g.nodes[nid]
	for i := 0; i < len(n.exits); i++ {
		e := n.getExit(i)
		// Find the active exit.
		if e.dst == nil {
			continue
		}
		// Make sure there are no other active exits.
		for j := i + 1; j < len(n.exits); j++ {
			if n.exits[j].dst != nil {
				panic(n.id)
			}
		}
		e.dst.replaceEntry(e, n.popEntries())
		break
	}
}

// Insert a dangling node (with a single out edge) in the middle of an existing edge.
func (g *Graph) InsertInEdge(dangling EdgeID, existing EdgeID) {
	replacement := g.edges[dangling]
	target := g.edges[existing]
	n := replacement.src
	if target.dst != nil {
		// Replace the exising edge with the new edge.
		target.dst.replaceSingleEntry(target, replacement)
	}
	// Attach the existing edge to the dangling node.
	target.dst = n
	n.addEntry(target)
}

func (g *Graph) HasMultipleEntries(dst NodeID) bool {
	return g.nodes[dst].entries.HasMultipleEdges()
}

func (g *Graph) GetUniqueExit(src NodeID) (EdgeID, NodeID) {
	if len(g.nodes[src].exits) == 1 {
		e := g.nodes[src].exits[0]
		if e.dst != nil {
			return e.id, e.dst.id
		}
	}
	return NoEdge, NoNode
}

func (g *Graph) EdgeFlow(e EdgeID) int {
	return g.edges[e].index
}

// TODO remove
func (g *Graph) IndexedExitEdge(nid NodeID, flow int) EdgeID {
	return g.nodes[nid].exits[flow].id
}

func CreateGraph() *Graph {
	g := &Graph{}
	// Entry
	g.CreateNode(1)
	// Exit
	g.CreateNode(0)
	return g
}

type FlowBuilder struct {
	graph *Graph
	exits []exitEdges
}

func CreateFlowBuilder(g *Graph, numExits int) *FlowBuilder {
	return createFlowBuilder(g, numExits, exitEdges{g.nodes[g.Entry()].exits[0]})
}

func createFlowBuilder(g *Graph, numExits int, edges []*edge) *FlowBuilder {
	fb := &FlowBuilder{
		graph: g,
		exits: make([]exitEdges, numExits),
	}
	fb.exits[0] = edges
	return fb
}

func (fb *FlowBuilder) HasFlow(flow int) bool {
	return len(fb.exits[flow]) > 0
}

func (fb *FlowBuilder) AttachFlow(flow int, dst NodeID) {
	dstNode := fb.graph.nodes[dst]
	if !fb.HasFlow(flow) {
		panic("Tried to attach non-existant flow")
	}
	// TODO extend entries directly.
	for _, e := range fb.exits[flow] {
		e.attach(dstNode)
	}
	fb.exits[flow] = nil
}

func (fb *FlowBuilder) RegisterExit(eid EdgeID, flow int) {
	e := fb.graph.edges[eid]
	if e.src == nil {
		panic(e)
	}
	if e.dst != nil {
		panic(e)
	}
	fb.exits[flow] = append(fb.exits[flow], e)
}

func (fb *FlowBuilder) popExits(flow int) exitEdges {
	return fb.exits[flow].Pop()
}

func (fb *FlowBuilder) SplitOffFlow(flow int) *FlowBuilder {
	return createFlowBuilder(fb.graph, len(fb.exits), fb.popExits(flow))
}

func (fb *FlowBuilder) SplitOffEdge(eid EdgeID) *FlowBuilder {
	g := fb.graph
	return createFlowBuilder(g, len(fb.exits), exitEdges{g.edges[eid]})
}

func (fb *FlowBuilder) AbsorbExits(other *FlowBuilder) {
	for i := 0; i < len(fb.exits); i++ {
		otherExits := other.exits[i]
		other.exits[i] = nil
		fb.exits[i] = append(fb.exits[i], otherExits...)
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
	return entryIterator{graph: g, node: g.nodes[n], current: 0}
}

func (g *Graph) ExitIterator(n NodeID) exitIterator {
	iter := exitIterator{graph: g, node: g.nodes[n], current: 0}
	iter.skipDeadEdges()
	return iter
}

func (g *Graph) ReverseExitIterator(n NodeID) reverseExitIterator {
	node := g.nodes[n]
	iter := reverseExitIterator{graph: g, node: node, current: len(node.exits) - 1}
	iter.skipDeadEdges()
	return iter
}

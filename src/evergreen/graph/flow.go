// Package graph implements a directed graph and common operations on directed graphs.
package graph

type entryList []*edge

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

type node struct {
	entries entryList
	exits   []*edge
	id      NodeID
}

func (n *node) GetExit(flow int) *edge {
	return n.exits[flow]
}

func (n *node) SetExit(flow int, other *node) {
	if flow >= len(n.exits) {
		panic(flow)
	}
	e := n.GetExit(flow)
	e.attach(other)
}

func (n *node) RemoveExit(flow int) {
	n.GetExit(flow).detach()
}

func (n *node) NumExits() int {
	return len(n.exits)
}

func (n *node) addEntry(e *edge) {
	n.entries = append(n.entries, e)
}

func (n *node) addEntries(e entryList) {
	n.entries = append(n.entries, e...)
}

func (n *node) popEntries() entryList {
	temp := n.entries
	n.entries = nil
	return temp
}

func (n *node) peekEntries() entryList {
	return n.entries
}

func (n *node) TransferEntries(other *node) {
	entries := n.popEntries()
	for _, e := range entries {
		e.dst = other
	}
	other.addEntries(entries)
}

func (n *node) InsertAt(flow int, target *edge) {
	if target.dst != nil {
		target.dst.replaceSingleEntry(target, n.GetExit(flow))
	}
	target.dst = n
	n.addEntry(target)
}

func (n *node) replaceSingleEntry(target *edge, replacement *edge) {
	n.replaceEntry(target, []*edge{replacement})
}

func (n *node) replaceEntry(target *edge, replacements entryList) {
	old := n.popEntries()
	for i, e := range old {
		if e == target {
			// Caution: splicing arrays is tricky because "append" may mutate the first argument.
			// Avoid appending.
			entries := make([]*edge, len(old)+len(replacements)-1)
			copy(entries[:i], old[:i])
			copy(entries[i:i+len(replacements)], replacements)
			copy(entries[i+len(replacements):], old[i+1:])
			n.entries = entries
			dst := target.dst
			target.dst = nil
			for _, r := range replacements {
				r.dst = dst
			}
			return
		}
	}
	panic(target)
}

type NodeID int

const NoNode NodeID = ^NodeID(0)

type EdgeID int

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
		exits: make([]*edge, exits),
	}
	for i := 0; i < exits; i++ {
		e := g.createEdge(i)
		e.src = n
		n.exits[i] = e
	}
	g.nodes = append(g.nodes, n)
	return id
}

func (g *Graph) NumNodes() int {
	return len(g.nodes)
}

func (g *Graph) Connect(src NodeID, edge int, dst NodeID) {
	g.nodes[src].SetExit(edge, g.nodes[dst])
}

func (g *Graph) RemoveEdge(eid EdgeID) {
	g.edges[eid].detach()
}

func (g *Graph) RemoveNode(nid NodeID) {
	n := g.nodes[nid]
	for i := 0; i < len(n.exits); i++ {
		e := n.GetExit(i)
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

func (g *Graph) InsertAt(nid NodeID, flow int, eid EdgeID) {
	g.nodes[nid].InsertAt(flow, g.edges[eid])
}

func (g *Graph) NumEntries(dst NodeID) int {
	return len(g.nodes[dst].entries)
}

func (g *Graph) NumExits(src NodeID) int {
	return len(g.nodes[src].exits)
}

func (g *Graph) GetExit(src NodeID, edge int) NodeID {
	next := g.nodes[src].exits[edge].dst
	if next != nil {
		return next.id
	} else {
		return NoNode
	}
}

func (g *Graph) CreateRegion(exits int) *GraphRegion {
	gr := &GraphRegion{
		graph: g,
		exits: make([][]*edge, exits),
	}
	gr.exits[0] = []*edge{&gr.entryEdge}
	return gr
}

func (g *Graph) ConnectRegion(gr *GraphRegion) {
	regionHead := gr.entryEdge.dst
	if regionHead == nil {
		g.Connect(g.Entry(), 0, g.Exit())
	} else {
		entry := g.nodes[g.Entry()]
		regionHead.replaceSingleEntry(&gr.entryEdge, entry.GetExit(0))
		for i := 0; i < len(gr.exits); i++ {
			if len(gr.exits[i]) > 0 {
				gr.AttachFlow(i, g.Exit())
			}
		}
	}
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

type GraphRegion struct {
	graph     *Graph
	exits     [][]*edge
	entryEdge edge
}

func (gr *GraphRegion) HasFlow(flow int) bool {
	return len(gr.exits[flow]) > 0
}

func (gr *GraphRegion) AttachFlow(flow int, dst NodeID) {
	dstNode := gr.graph.nodes[dst]
	if !gr.HasFlow(flow) {
		panic("Tried to attach non-existant flow")
	}
	// TODO extend entries directly.
	for _, e := range gr.exits[flow] {
		e.attach(dstNode)
	}
	gr.exits[flow] = nil
}

func (gr *GraphRegion) RegisterExit(eid EdgeID, flow int) {
	e := gr.graph.edges[eid]
	if e.src == nil {
		panic(e)
	}
	if e.dst != nil {
		panic(e)
	}
	gr.exits[flow] = append(gr.exits[flow], e)
}

func (gr *GraphRegion) Swap(flow0 int, flow1 int) {
	gr.exits[flow0], gr.exits[flow1] = gr.exits[flow1], gr.exits[flow0]
}

func (gr *GraphRegion) popExits(flow int) []*edge {
	exits := gr.exits[flow]
	gr.exits[flow] = nil
	return exits
}

func (gr *GraphRegion) MergeFlowInto(srcFlow int, dstFlow int) {
	if srcFlow != dstFlow {
		gr.exits[dstFlow] = append(gr.exits[dstFlow], gr.popExits(srcFlow)...)
	}
}

func (gr *GraphRegion) findEntryEdge() int {
	if gr.entryEdge.dst != nil {
		panic(gr.entryEdge.dst)
	}
	for i, exits := range gr.exits {
		if len(exits) == 1 {
			return i
		}
	}
	panic(gr.exits)
}

func (gr *GraphRegion) Splice(flow int, other *GraphRegion) {
	if !gr.HasFlow(flow) {
		panic("Sloppy: tried to splice to nothing.")
	}
	otherHead := other.entryEdge.dst
	if otherHead != nil {
		edges := gr.popExits(flow)
		otherHead.replaceEntry(&other.entryEdge, edges)
		gr.absorbExits(other)
	} else {
		gr.MergeFlowInto(flow, other.findEntryEdge())
	}
}

func (gr *GraphRegion) SpliceToEdge(eid EdgeID, other *GraphRegion) {
	e := gr.graph.edges[eid]
	otherHead := other.entryEdge.dst
	if otherHead != nil {
		otherHead.replaceSingleEntry(&other.entryEdge, e)
		gr.absorbExits(other)
	} else {
		gr.RegisterExit(eid, other.findEntryEdge())
	}
}

func (gr *GraphRegion) absorbExits(other *GraphRegion) {
	for i := 0; i < len(gr.exits); i++ {
		otherExits := other.exits[i]
		other.exits[i] = nil
		gr.exits[i] = append(gr.exits[i], otherExits...)
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
func NodeIterator(g *Graph) nodeIterator {
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

func EntryIterator(g *Graph, n NodeID) entryIterator {
	return entryIterator{graph: g, node: g.nodes[n], current: 0}
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

func ExitIterator(g *Graph, n NodeID) exitIterator {
	iter := exitIterator{graph: g, node: g.nodes[n], current: 0}
	iter.skipDeadEdges()
	return iter
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

func ReverseExitIterator(g *Graph, n NodeID) reverseExitIterator {
	node := g.nodes[n]
	iter := reverseExitIterator{graph: g, node: node, current: len(node.exits) - 1}
	iter.skipDeadEdges()
	return iter
}

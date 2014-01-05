package base

type NodeData interface {
}

type EntryList []*Edge

type Edge struct {
	src   *Node
	dst   *Node
	index int
}

func (e *Edge) attach(other *Node) {
	if e.dst != nil {
		panic(e)
	}
	e.dst = other
	other.addEntry(e)
}

const NoNodeIndex = ^int(0)

type Node struct {
	entries EntryList
	exits   []Edge
	Id      NodeID
	Name    int
	Data    NodeData
}

func CreateNode(data NodeData, numExits int) *Node {
	n := &Node{Data: data, exits: make([]Edge, numExits)}
	for i := 0; i < numExits; i++ {
		n.exits[i] = Edge{src: n, index: i}
	}
	return n
}

func (n *Node) GetNext(flow int) *Node {
	return n.exits[flow].dst
}

func (n *Node) GetExit(flow int) *Edge {
	return &n.exits[flow]
}

func (n *Node) SetExit(flow int, other *Node) {
	if flow >= len(n.exits) {
		panic(flow)
	}
	e := n.GetExit(flow)
	e.attach(other)
}

func (n *Node) NumExits() int {
	return len(n.exits)
}

func (n *Node) SetDefaultExits(exits []*Node) {
	for i, e := range n.exits {
		if e.dst == nil {
			n.SetExit(i, exits[i])
		}
	}
}

func (n *Node) addEntry(e *Edge) {
	n.entries = append(n.entries, e)
}

func (n *Node) addEntries(e EntryList) {
	n.entries = append(n.entries, e...)
}

func (n *Node) popEntries() EntryList {
	temp := n.entries
	n.entries = nil
	return temp
}

func (n *Node) peekEntries() EntryList {
	return n.entries
}

func (n *Node) NumEntries() int {
	return len(n.entries)
}

func (n *Node) HasEntries() bool {
	return len(n.entries) > 0
}

func (n *Node) GetEntry(i int) *Edge {
	return n.entries[i]
}

func (n *Node) TransferEntries(other *Node) {
	entries := n.popEntries()
	for _, e := range entries {
		e.dst = other
	}
	other.addEntries(entries)
}

func (n *Node) InsertAt(flow int, target *Edge) {
	if target.dst != nil {
		target.dst.replaceSingleEntry(target, n.GetExit(flow))
	}
	target.dst = n
	n.addEntry(target)
}

func (n *Node) Remove() {
	for i := 0; i < len(n.exits); i++ {
		e := n.GetExit(i)
		// Find the active exit.
		if e.dst == nil {
			continue
		}
		// Make sure there are no other active exits.
		for j := i + 1; j < len(n.exits); j++ {
			if n.exits[j].dst != nil {
				panic(n.Data)
			}
		}
		e.dst.replaceEntry(e, n.popEntries())
		break
	}
}

func (n *Node) replaceSingleEntry(target *Edge, replacement *Edge) {
	n.replaceEntry(target, []*Edge{replacement})
}

func (n *Node) replaceEntry(target *Edge, replacements EntryList) {
	old := n.popEntries()
	for i, e := range old {
		if e == target {
			// Caution: splicing arrays is tricky because "append" may mutate the first argument.
			// Avoid appending.
			entries := make([]*Edge, len(old)+len(replacements)-1)
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

type Region struct {
	Entry *Node
	Exits []*Node
}

type NodeID int

const NoNode NodeID = ^NodeID(0)

type Graph struct {
	nodes []*Node
}

func (g *Graph) Entry() NodeID {
	return 0
}

func (g *Graph) Exit() NodeID {
	return 1
}

func (g *Graph) CreateNode(data interface{}, exits int) NodeID {
	id := NodeID(len(g.nodes))
	n := &Node{
		Id:    id,
		Data:  data,
		exits: make([]Edge, exits),
	}
	for i := 0; i < exits; i++ {
		n.exits[i] = Edge{src: n, index: i}
	}
	g.nodes = append(g.nodes, n)
	return id
}

func (g *Graph) ResolveNodeHACK(node NodeID) *Node {
	return g.nodes[node]
}

func (g *Graph) Connect(src NodeID, edge int, dst NodeID) {
	g.nodes[src].SetExit(edge, g.nodes[dst])
}

func (g *Graph) NumEntries(dst NodeID) int {
	return len(g.nodes[dst].entries)
}

func (g *Graph) GetEntry(dst NodeID, edge int) NodeID {
	return g.nodes[dst].entries[edge].src.Id
}

func (g *Graph) NumExits(src NodeID) int {
	return len(g.nodes[src].exits)
}

func (g *Graph) GetExit(src NodeID, edge int) NodeID {
	next := g.nodes[src].exits[edge].dst
	if next != nil {
		return next.Id
	} else {
		return NoNode
	}
}

func (g *Graph) CreateRegion(exits int) *GraphRegion {
	gr := &GraphRegion{
		graph: g,
		exits: make([][]*Edge, exits),
	}
	gr.exits[0] = []*Edge{&gr.entryEdge}
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

func CreateGraph(entry interface{}, exit interface{}) *Graph {
	g := &Graph{}
	g.CreateNode(entry, 1)
	g.CreateNode(exit, 0)
	return g
}

type GraphRegion struct {
	graph     *Graph
	exits     [][]*Edge
	entryEdge Edge
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

func (gr *GraphRegion) RegisterExit(src NodeID, edge int, flow int) {
	srcNode := gr.graph.nodes[src]
	e := srcNode.GetExit(edge)
	if e.dst != nil {
		panic(e)
	}
	gr.exits[flow] = append(gr.exits[flow], e)
}

func (gr *GraphRegion) Swap(flow0 int, flow1 int) {
	gr.exits[flow0], gr.exits[flow1] = gr.exits[flow1], gr.exits[flow0]
}

func (gr *GraphRegion) popExits(flow int) []*Edge {
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

func (gr *GraphRegion) SpliceToEdge(src NodeID, flow int, other *GraphRegion) {
	srcNode := gr.graph.nodes[src]
	otherHead := other.entryEdge.dst
	if otherHead != nil {
		otherHead.replaceSingleEntry(&other.entryEdge, srcNode.GetExit(flow))
		gr.absorbExits(other)
	} else {
		gr.RegisterExit(src, flow, other.findEntryEdge())
	}
}

func (gr *GraphRegion) absorbExits(other *GraphRegion) {
	for i := 0; i < len(gr.exits); i++ {
		otherExits := other.exits[i]
		other.exits[i] = nil
		gr.exits[i] = append(gr.exits[i], otherExits...)
	}
}

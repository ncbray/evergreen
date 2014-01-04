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

const NoNode = ^int(0)

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
		target.dst.ReplaceEntry(target, []*Edge{n.GetExit(flow)})
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
		e.dst.ReplaceEntry(e, n.popEntries())
		break
	}
}

func (n *Node) ReplaceEntry(target *Edge, replacements EntryList) {
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

func (r *Region) Head() *Node {
	return r.Entry.GetExit(0).dst
}

func (r *Region) Connect(flow int, n *Node) {
	r.Exits[flow].TransferEntries(n)
}

func (r *Region) AttachDefaultExits(n *Node) {
	n.SetDefaultExits(r.Exits)
}

func (r *Region) Splice(flow int, other *Region) {
	otherEntry := other.Entry.GetExit(0)
	otherHead := otherEntry.dst
	otherHead.ReplaceEntry(otherEntry, r.Exits[flow].popEntries())
	r.AbsorbExits(other)
}

func (r *Region) AbsorbExits(other *Region) {
	for i, exit := range r.Exits {
		other.Exits[i].TransferEntries(exit)
	}
}

func (r *Region) GetEntry() *Node {
	return r.Entry
}

func (r *Region) GetExit(flow int) *Node {
	return r.Exits[flow]
}

type NodeID int

type Graph struct {
	nodes []*Node
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

func (g *Graph) Connect(src NodeID, edge int, dst NodeID) {
	g.nodes[src].SetExit(edge, g.nodes[dst])
}

func (g *Graph) CreateRegion(exits int) *GraphRegion {
	gr := &GraphRegion{
		implicitExit: 0,
		graph:        g,
		exits:        make([][]*Edge, exits),
	}
	return gr
}

func CreateGraph() *Graph {
	return &Graph{}
}

type GraphRegion struct {
	graph        *Graph
	entry        *Node
	exits        [][]*Edge
	implicitExit int
}

func (gr *GraphRegion) HasFlow(flow int) bool {
	if gr.entry == nil {
		return flow == gr.implicitExit
	} else {
		return len(gr.exits[flow]) > 0
	}
}

func (gr *GraphRegion) AttachFlow(flow int, dst NodeID) {
	dstNode := gr.graph.nodes[dst]
	gr.AttachFlowHACK(flow, dstNode)
}

func (gr *GraphRegion) AttachFlowHACK(flow int, dstNode *Node) {
	if gr.entry == nil {
		if gr.implicitExit == flow {
			gr.entry = dstNode
		} else {
			panic("bad first node?")
		}
	} else {
		// TODO extend entries directly.
		for _, e := range gr.exits[flow] {
			e.attach(dstNode)
		}
		gr.exits[flow] = nil
	}
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
	if gr.implicitExit == flow0 {
		gr.implicitExit = flow1
	} else if gr.implicitExit == flow1 {
		gr.implicitExit = flow0
	}
}

func (gr *GraphRegion) MergeFlowInto(srcFlow int, dstFlow int) {
	gr.exits[dstFlow] = append(gr.exits[dstFlow], gr.exits[srcFlow]...)
	gr.exits[srcFlow] = nil
}

func (gr *GraphRegion) Splice(flow int, other *GraphRegion) {
	if !gr.HasFlow(flow) {
		panic("Sloppy: tried to splice to nothing.")
	}
	if other.entry == nil {
		panic("TODO: empty splice")
	} else {
		gr.AttachFlowHACK(flow, other.entry)
		gr.absorbExits(other)
	}
}

func (gr *GraphRegion) SpliceToEdge(src NodeID, flow int, other *GraphRegion) {
	if other.entry != nil {
		gr.graph.Connect(src, flow, other.entry.Id)
		gr.absorbExits(other)
	} else {
		gr.RegisterExit(src, flow, other.implicitExit)
	}
}

func (gr *GraphRegion) absorbExits(other *GraphRegion) {
	for i := 0; i < len(gr.exits); i++ {
		otherExits := other.exits[i]
		other.exits[i] = nil
		gr.exits[i] = append(gr.exits[i], otherExits...)
	}
}

func (gr *GraphRegion) HeadHACK() *Node {
	return gr.entry
}

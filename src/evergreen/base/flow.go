package base

type NodeData interface {
	NumExits() int
	DotNodeStyle() string
	DotEdgeStyle(flow int) string
}

type EntryList []*Edge

type Edge struct {
	src   *Node
	dst   *Node
	index int
}

type Node struct {
	entries EntryList
	exits   []Edge
	Data    NodeData
}

func CreateNode(data NodeData) *Node {
	numExits := data.NumExits()
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
	if e.dst != nil {
		panic(e)
	}
	e.dst = other
	other.addEntry(e)
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

func (n *Node) TransferEntries(other *Node) {
	entries := n.popEntries()
	for _, e := range entries {
		e.dst = other
	}
	other.addEntries(entries)
}

func (n *Node) ReplaceEntry(target *Edge, replacements EntryList) {
	old := n.popEntries()
	for i, e := range old {
		if e == target {
			n.entries = append(append(old[:i], replacements...), old[i+1:]...)
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
	entry *Node
	exits []*Node
}

func (r *Region) Head() *Node {
	return r.entry.GetExit(0).dst
}

func (r *Region) Connect(flow int, n *Node) {
	r.exits[flow].TransferEntries(n)
}

func (r *Region) AttachDefaultExits(n *Node) {
	n.SetDefaultExits(r.exits)
}

func (r *Region) Splice(flow int, other *Region) {
	otherEntry := other.entry.GetExit(0)
	otherHead := otherEntry.dst
	otherHead.ReplaceEntry(otherEntry, r.exits[flow].popEntries())
	for i, exit := range r.exits {
		other.exits[i].TransferEntries(exit)
	}
}

func (r *Region) GetExit(flow int) *Node {
	return r.exits[flow]
}

func CreateRegion(entry *Node, exits []*Node) *Region {
	r := &Region{
		entry: entry,
		exits: exits,
	}
	r.entry.SetExit(0, r.exits[0])
	return r
}

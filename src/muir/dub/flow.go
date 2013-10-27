package dub

type NodeData interface {
	NumExits() int
}

type EntryList []*Edge

type Edge struct {
	src   *NodeImpl
	dst   *NodeImpl
	index int
}

type NodeImpl struct {
	entries EntryList
	exits   []Edge
	data    NodeData
}

func CreateNode(data NodeData) *NodeImpl {
	numExits := data.NumExits()
	n := &NodeImpl{data: data, exits: make([]Edge, numExits)}
	for i := 0; i < numExits; i++ {
		n.exits[i] = Edge{src: n, index: i}
	}
	return n
}

func (n *NodeImpl) GetExit(flow int) *Edge {
	return &n.exits[flow]
}

func (n *NodeImpl) SetExit(flow int, other *NodeImpl) {
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

func (n *NodeImpl) NumExits() int {
	return len(n.exits)
}

func (n *NodeImpl) SetDefaultExits(exits []*NodeImpl) {
	for i, e := range n.exits {
		if e.dst == nil {
			n.SetExit(i, exits[i])
		}
	}
}

func (n *NodeImpl) addEntry(e *Edge) {
	n.entries = append(n.entries, e)
}

func (n *NodeImpl) addEntries(e EntryList) {
	n.entries = append(n.entries, e...)
}

func (n *NodeImpl) popEntries() EntryList {
	temp := n.entries
	n.entries = nil
	return temp
}

func (n *NodeImpl) peekEntries() EntryList {
	return n.entries
}

func (n *NodeImpl) TransferEntries(other *NodeImpl) {
	entries := n.popEntries()
	for _, e := range entries {
		e.dst = other
	}
	other.addEntries(entries)
}

func (n *NodeImpl) ReplaceEntry(target *Edge, replacements EntryList) {
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
	entry *NodeImpl
	exits []*NodeImpl
}

func (r *Region) Head() *NodeImpl {
	return r.entry.GetExit(0).dst
}

func (r *Region) Connect(flow int, n *NodeImpl) {
	r.exits[flow].TransferEntries(n)
}

func (r *Region) AttachDefaultExits(n *NodeImpl) {
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

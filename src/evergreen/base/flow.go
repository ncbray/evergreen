package base

type NodeData interface {
}

type EntryList []*Edge

type Edge struct {
	src   *Node
	dst   *Node
	index int
	Data  interface{}
}

const NoNode = ^int(0)

type Node struct {
	entries EntryList
	exits   []Edge
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

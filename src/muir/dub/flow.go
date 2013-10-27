package dub

type Edge struct {
	src *Node
  dst *Node
  index int
}

func (e *Edge) Connect(other *Node) {
	if e.dst != nil {
		panic(e)
	}
	e.dst = other
  other.entries = append(other.entries, e)
}

type Node struct {
  entries []*Edge
  exits []*Edge
}

func (n *Node) Connect(flow int, other *Node) {
	if flow >= len(n.exits) {
		panic(flow)
	}
	n.exits[flow].Connect(other)
}

func (n *Node) popEntries() []*Edge {
	temp := n.entries
	n.entries = nil
	return temp;
}

func (n *Node) StealEntries(other *Node) {
  if n == other {
		panic(n)
	}
	entries := other.popEntries()
	for _, e := range entries {
		e.dst = n
	}
	n.entries = append(n.entries, entries...)
}

func (n *Node) ReplaceEntry(target *Edge, replacements []*Edge) {
	old := n.popEntries()
	for i, e := range old {
		if (e == target) {
			n.entries = append(append(old[:i], replacements...), old[i+1:]...)
			target.dst = nil
			for _, r := range replacements {
				r.dst = n
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
	return r.entry.exits[0].dst
}

func (r *Region) Connect(flow int, n *Node) {
	n.StealEntries(r.exits[flow])
}

func (r *Region) AttachDefaultExits(n *Node) {
	for i, e := range n.exits {
		if (e.dst == nil) {
			e.Connect(r.exits[i])
		}
	}
}

func (r *Region) Splice(flow int, other *Region) {
	otherEntry := other.entry.exits[0]
	otherHead := otherEntry.dst
	otherHead.ReplaceEntry(otherEntry, r.exits[flow].popEntries())
	for i, exit := range r.exits {
		exit.StealEntries(other.exits[i])
	}
}

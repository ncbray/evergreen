package base

type rpoSearch struct {
	graph *Graph
	order []NodeID
	index []int
}

func (s *rpoSearch) search(n NodeID) {
	if s.index[n] != 0 {
		return
	}
	// Prevent processing it again, the correct index will be computed later.
	s.index[n] = 1
	numExits := s.graph.NumExits(n)
	for i := numExits - 1; i >= 0; i-- {
		dst := s.graph.GetExit(n, i)
		if dst != NoNode {
			s.search(dst)
		}
	}
	s.order = append(s.order, n)
	// Zero is reserved, so use the actual index + 1.
	s.index[n] = len(s.order)
}

func ReverseOrder(order []NodeID) {
	o := len(order)
	for i := 0; i < o/2; i++ {
		order[i], order[o-1-i] = order[o-1-i], order[i]
	}
}

func ReversePostorder(g *Graph) ([]NodeID, []int) {
	numNodes := len(g.nodes)
	s := &rpoSearch{
		graph: g,
		order: make([]NodeID, 1, numNodes),
		index: make([]int, numNodes),
	}
	// Implicit edge from entry to exit ensures exit will always be the last node.
	s.order[0] = g.Exit()
	s.index[g.Exit()] = 1
	s.search(g.Entry())

	ReverseOrder(s.order)
	o := len(s.order)
	for i := 0; i < numNodes; i++ {
		if s.index[i] != 0 {
			s.index[i] = o - s.index[i]
		} else {
			s.index[i] = -1
		}
	}
	return s.order, s.index
}

func intersectDom(idoms []NodeID, index []int, n0 NodeID, n1 NodeID) NodeID {
	i0 := index[n0]
	i1 := index[n1]
	for i0 != i1 {
		for i0 > i1 {
			n0 = idoms[n0]
			i0 = index[n0]
		}
		for i0 < i1 {
			n1 = idoms[n1]
			i1 = index[n1]
		}
	}
	return n0
}

func isBackedge(index []int, src NodeID, dst NodeID) bool {
	return index[src] > index[dst]
}

func FindDominators(g *Graph, order []NodeID, index []int) []NodeID {
	numNodes := len(g.nodes)
	idoms := make([]NodeID, numNodes)
	changed := false
	nit := OrderedIterator(order)
	for nit.Next() {
		n := nit.Value()
		// If there are no forward entries into the node, assume an impossible edge.
		new_idom := g.Entry()
		first := true
		eit := EntryIterator(g, n)
		for eit.Next() {
			e := eit.Value()
			if !isBackedge(index, e, n) {
				if first {
					new_idom = e
					first = false
				} else {
					new_idom = intersectDom(idoms, index, new_idom, e)
				}
			} else {
				changed = true
			}
		}
		idoms[n] = new_idom
	}
	for changed {
		changed = false
		nit := OrderedIterator(order)
		for nit.Next() {
			n := nit.Value()
			numEntries := g.NumEntries(n)
			// 0 and 1 entry nodes should be stable after the first pass.
			if numEntries >= 2 {
				newIdom := idoms[n]
				eit := EntryIterator(g, n)
				for eit.Next() {
					e := eit.Value()
					newIdom = intersectDom(idoms, index, newIdom, e)
				}
				if idoms[n] != newIdom {
					idoms[n] = newIdom
					changed = true
				}
			}
		}
	}
	return idoms
}

// Assumes no dead entries.
func FindDominanceFrontiers(g *Graph, idoms []NodeID) [][]NodeID {
	n := len(g.nodes)
	frontiers := make([][]NodeID, n)
	nit := NodeIterator(g)
	for nit.Next() {
		n := nit.Value()
		numEntries := g.NumEntries(n)
		if numEntries >= 2 {
			target := idoms[n]
			eit := EntryIterator(g, n)
			for eit.Next() {
				runner := eit.Value()
				for runner != target {
					frontiers[runner] = append(frontiers[runner], n)
					runner = idoms[runner]
				}
			}
		}
	}
	return frontiers
}

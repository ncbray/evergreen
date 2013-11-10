package dub

import (
	"testing"
)

func checkEdge(e *Edge, src *Node, dst *Node, t *testing.T) {
	if e.src != src {
		t.Errorf("Got src of %v, expected %v", e.src, src)
	}
	if e.dst != dst {
		t.Errorf("Got dst of %v, expected %v", e.dst, dst)
	}
	if src.GetExit(e.index) != e {
		t.Errorf("Inconsistent indexing for %v, found %v", e, src.GetExit(e.index))
	}
}

func checkTopology(node *Node, entries []*Node, exits []*Node, t *testing.T) {
	if node == nil {
		t.Error("Node should not be nil")
		return
	}
	oentries := node.peekEntries()
	if len(entries) != len(oentries) {
		t.Errorf("Expected %d entries, got %d", len(entries), len(oentries))
	} else {
		for i, entry := range entries {
			checkEdge(oentries[i], entry, node, t)
		}
	}
	if len(exits) != node.NumExits() {
		t.Errorf("Expected %d exits, got %d", len(entries), node.NumExits())
	} else {
		for i, exit := range exits {
			checkEdge(node.GetExit(i), node, exit, t)
		}
	}
}

type TestEntry struct {
}

func (n *TestEntry) NumExits() int {
	return 1
}

type TestExit struct {
	flow int
}

func (n *TestExit) NumExits() int {
	return 0
}

type TestNode struct {
	name string
}

func (n *TestNode) NumExits() int {
	return 2
}

func CreateTestEntry() *Node {
	return CreateNode(&TestEntry{})
}

func CreateTestNode(name string, numExits int) *Node {
	return CreateNode(&TestNode{name: name})
}

func CreateTestExit(flow int) *Node {
	return CreateNode(&TestExit{flow: flow})
}

func CreateTestRegion() *Region {
	r := CreateRegion(
		CreateTestEntry(),
		[]*Node{
			CreateTestExit(0),
			CreateTestExit(1),
		},
	)
	return r
}

func TestSimpleFlow(t *testing.T) {
	r := CreateTestRegion()
	n := CreateTestNode("n", 2)
	r.Connect(0, n)
	r.AttachDefaultExits(n)

	checkTopology(r.entry, []*Node{}, []*Node{n}, t)
	checkTopology(n, []*Node{r.entry}, []*Node{r.exits[0], r.exits[1]}, t)
	checkTopology(r.exits[0], []*Node{n}, []*Node{}, t)
	checkTopology(r.exits[1], []*Node{n}, []*Node{}, t)
}

func TestRepeatFlow(t *testing.T) {
	l := CreateTestRegion()
	n := CreateTestNode("n", 2)
	l.Connect(0, n)
	l.AttachDefaultExits(n)
	head := l.Head()
	// Normal flow iterates
	l.exits[0].TransferEntries(head)
	// Stop iterating on failure
	l.exits[1].TransferEntries(l.exits[0])

	checkTopology(l.entry, []*Node{}, []*Node{n}, t)
	checkTopology(n, []*Node{l.entry, n}, []*Node{n, l.exits[0]}, t)
	checkTopology(l.exits[0], []*Node{n}, []*Node{}, t)
	checkTopology(l.exits[1], []*Node{}, []*Node{}, t)

	r := CreateTestRegion()

	r.Splice(0, l)

	checkTopology(l.entry, []*Node{}, []*Node{nil}, t)
	checkTopology(l.exits[0], []*Node{}, []*Node{}, t)
	checkTopology(l.exits[1], []*Node{}, []*Node{}, t)

	checkTopology(r.entry, []*Node{}, []*Node{n}, t)
	checkTopology(n, []*Node{r.entry, n}, []*Node{n, r.exits[0]}, t)
	checkTopology(r.exits[0], []*Node{n}, []*Node{}, t)
	checkTopology(r.exits[1], []*Node{}, []*Node{}, t)
}

func TestWhileFlow(t *testing.T) {
	l := CreateTestRegion()
	cond := CreateTestNode("cond", 2)
	decide := CreateTestNode("decide", 2)
	body := CreateTestNode("body", 2)

	l.Connect(0, cond)
	l.AttachDefaultExits(cond)

	l.Connect(0, decide)
	decide.SetExit(0, body)

	l.AttachDefaultExits(body)
	l.Connect(0, cond)
	decide.SetExit(1, l.exits[0])

	r := CreateTestRegion()
	r.Splice(0, l)

	checkTopology(r.entry, []*Node{}, []*Node{cond}, t)
	checkTopology(cond, []*Node{r.entry, body}, []*Node{decide, r.exits[1]}, t)
	checkTopology(decide, []*Node{cond}, []*Node{body, r.exits[0]}, t)
	checkTopology(body, []*Node{decide}, []*Node{cond, r.exits[1]}, t)
	checkTopology(r.exits[0], []*Node{decide}, []*Node{}, t)
	checkTopology(r.exits[1], []*Node{cond, body}, []*Node{}, t)
}

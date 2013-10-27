package dub

import (
	"testing"
)

func checkEdge(e *Edge, src *Node, dst *Node, t *testing.T) {
	if (e.src != src) {
		t.Errorf("Got src of %v, expected %v", e.src, src)
	}
	if (e.dst != dst) {
		t.Errorf("Got dst of %v, expected %v", e.dst, dst)
	}
	if (src.exits[e.index] != e) {
		t.Errorf("Inconsistent indexing for %v, found %v", e, src.exits[e.index])
	}
}

func checkTopology(node *Node, entries []*Node, exits []*Node, t *testing.T) {
	if (node == nil) {
		t.Error("Node should not be nil")
		return;
	}
	if len(entries) != len(node.entries) {
		t.Errorf("Expected %d entries, got %d", len(entries), len(node.entries))
	} else {
		for i, entry := range entries {
			checkEdge(node.entries[i], entry, node, t)
		}
	}
	if len(exits) != len(node.exits) {
		t.Errorf("Expected %d exits, got %d", len(entries), len(node.entries))
	} else {
		for i, exit := range exits {
			checkEdge(node.exits[i], node, exit, t)
		}
	}
}

func CreateTestNode(name string, numExits int) *Node {
  n := &Node{exits: make([]*Edge, numExits)}
	for i := 0; i < numExits; i++ {
		n.exits[i] = &Edge{src: n, index: i}
	}
  return n
}

func CreateTestRegion() *Region {
  r := &Region{
		entry: CreateTestNode("entry", 1),
		exits: []*Node{
			CreateTestNode("exit", 0),
			CreateTestNode("exit", 0),
		},
	}
	r.entry.Connect(0, r.exits[0])
	return r;
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
	head.StealEntries(l.exits[0])
	// Stop iterating on failure
	l.exits[0].StealEntries(l.exits[1])

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
	decide.Connect(0, body)

	l.AttachDefaultExits(body)
	l.Connect(0, cond)
	decide.Connect(1, l.exits[0])

	r := CreateTestRegion()
	r.Splice(0, l)

  checkTopology(r.entry, []*Node{}, []*Node{cond}, t)
  checkTopology(cond, []*Node{r.entry, body}, []*Node{decide, r.exits[1]}, t)
  checkTopology(decide, []*Node{cond}, []*Node{body, r.exits[0]}, t)
  checkTopology(body, []*Node{decide}, []*Node{cond, r.exits[1]}, t)
  checkTopology(r.exits[0], []*Node{decide}, []*Node{}, t)
  checkTopology(r.exits[1], []*Node{cond, body}, []*Node{}, t)
}

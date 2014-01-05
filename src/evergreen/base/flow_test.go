package base

import (
	"fmt"
	"testing"
)

func checkEdge(e *Edge, src *Node, dst *Node, t *testing.T) {
	if e.src != src {
		t.Errorf("Got src of %v, expected %v", e.src, src)
	}
	if e.dst != dst {
		t.Errorf("Got dst of %v, expected %v", e.dst, dst)
	}
	if e.src.GetExit(e.index) != e {
		t.Errorf("Inconsistent indexing for %v, found %v", e, src.GetExit(e.index))
	}
}

func checkTopology(g *Graph, id NodeID, entries []NodeID, exits []NodeID, t *testing.T) {
	node := g.nodes[id]
	if node == nil {
		t.Error("Node should not be nil")
		return
	}
	oentries := node.peekEntries()
	if len(entries) != len(oentries) {
		t.Errorf("Expected %d entries, got %d", len(entries), len(oentries))
	} else {
		for i, entry := range entries {
			checkEdge(oentries[i], g.nodes[entry], node, t)
		}
	}
	if len(exits) != node.NumExits() {
		t.Errorf("Expected %d exits, got %d", len(entries), node.NumExits())
	} else {
		for i, exit := range exits {
			checkEdge(node.GetExit(i), node, g.nodes[exit], t)
		}
	}
}

func TestSimpleFlow(t *testing.T) {
	g := CreateGraph(nil, nil)
	n := g.CreateNode(nil, 1)
	g.Connect(g.Entry(), 0, n)
	g.Connect(n, 0, g.Exit())

	checkTopology(g, g.Entry(), []NodeID{}, []NodeID{n}, t)
	checkTopology(g, n, []NodeID{g.Entry()}, []NodeID{g.Exit()}, t)
	checkTopology(g, g.Exit(), []NodeID{n}, []NodeID{}, t)
}

func TestSliceEmptySplice(t *testing.T) {
	g := CreateGraph(nil, nil)
	gr0 := g.CreateRegion(2)
	gr1 := g.CreateRegion(2)
	checkInt("sanity", len(gr0.exits[0]), 1, t)
	gr1.Swap(0, 1)
	gr0.Splice(0, gr1)
	checkInt("swaped", len(gr0.exits[1]), 1, t)
}

func TestSliceEdgeEmptySplice(t *testing.T) {
	g := CreateGraph(nil, nil)
	gr0 := g.CreateRegion(2)
	gr1 := g.CreateRegion(2)
	n := g.CreateNode(nil, 1)
	gr0.AttachFlow(0, n)

	checkInt("sanity", len(gr0.exits[0]), 0, t)
	gr1.Swap(0, 1)
	gr0.SpliceToEdge(n, 0, gr1)
	checkInt("swaped", len(gr0.exits[1]), 1, t)
}

func TestRepeatFlow(t *testing.T) {
	g := CreateGraph(nil, nil)
	gr := g.CreateRegion(2)

	n := g.CreateNode("n", 2)
	gr.AttachFlow(0, n)
	gr.RegisterExit(n, 0, 0)
	gr.RegisterExit(n, 1, 1)

	// Normal flow iterates
	gr.AttachFlow(0, n)

	// Stop iterating on failure
	gr.Swap(0, 1)

	g.ConnectRegion(gr)

	e := g.Entry()
	x := g.Exit()

	checkTopology(g, e, []NodeID{}, []NodeID{n}, t)
	checkTopology(g, n, []NodeID{e, n}, []NodeID{n, x}, t)
	checkTopology(g, x, []NodeID{n}, []NodeID{}, t)
}

func TestWhileFlow(t *testing.T) {
	g := CreateGraph(nil, nil)
	gr := g.CreateRegion(2)
	c := g.CreateNode("c", 1)
	d := g.CreateNode("d", 2)
	b := g.CreateNode("b", 1)

	gr.AttachFlow(0, c)
	gr.RegisterExit(c, 0, 0)

	gr.AttachFlow(0, d)
	gr.RegisterExit(d, 0, 0)
	gr.RegisterExit(d, 1, 1)

	gr.AttachFlow(0, b)
	gr.RegisterExit(b, 0, 0)

	gr.AttachFlow(0, c)

	gr.Swap(0, 1)

	g.ConnectRegion(gr)

	e := g.Entry()
	x := g.Exit()

	checkTopology(g, e, []NodeID{}, []NodeID{c}, t)
	checkTopology(g, c, []NodeID{e, b}, []NodeID{d}, t)
	checkTopology(g, d, []NodeID{c}, []NodeID{b, x}, t)
	checkTopology(g, b, []NodeID{d}, []NodeID{c}, t)
	checkTopology(g, x, []NodeID{d}, []NodeID{}, t)
}

func checkInt(name string, actual int, expected int, t *testing.T) {
	if actual != expected {
		t.Fatalf("%s: %d != %d", name, actual, expected)
	}
}

func checkOrder(actualOrder []*Node, expectedOrder []NodeID, t *testing.T) {
	checkInt("len", len(actualOrder), len(expectedOrder), t)
	for i, expected := range expectedOrder {
		if actualOrder[i].Id != expected {
			t.Fatalf("%d: %#v != %#v", i, actualOrder[i].Id, expected)
		}
		checkInt(fmt.Sprint(i), actualOrder[i].Name, i, t)
	}
}

func checkIntList(actualList []int, expectedList []int, t *testing.T) {
	checkInt("len", len(actualList), len(expectedList), t)
	for i, expected := range expectedList {
		checkInt(fmt.Sprint(i), actualList[i], expected, t)
	}
}

func checkIntListList(actualList [][]int, expectedList [][]int, t *testing.T) {
	checkInt("len", len(actualList), len(expectedList), t)
	for i, expected := range expectedList {
		checkIntList(actualList[i], expected, t)
	}
}

func CreateTestRegion(g *Graph) *Region {
	return &Region{
		Entry: g.ResolveNodeHACK(g.Entry()),
		Exits: []*Node{
			g.ResolveNodeHACK(g.Exit()),
		},
	}
}

func TestSanity(t *testing.T) {
	g := CreateGraph(nil, nil)

	n1 := g.CreateNode("1", 1)
	n2 := g.CreateNode("2", 1)
	n3 := g.CreateNode("3", 1)

	g.Connect(g.Entry(), 0, n1)
	g.Connect(n1, 0, n2)
	g.Connect(n2, 0, n3)
	g.Connect(n3, 0, g.Exit())

	r := CreateTestRegion(g)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []NodeID{g.Entry(), n1, n2, n3, g.Exit()}, t)

	idoms := FindIdoms(ordered)
	checkIntList(idoms, []int{0, 0, 1, 2, 3}, t)
}

func TestLoop(t *testing.T) {
	g := CreateGraph(nil, nil)
	n1 := g.CreateNode("1", 1)
	n2 := g.CreateNode("2", 1)
	n3 := g.CreateNode("3", 1)

	g.Connect(g.Entry(), 0, n1)
	g.Connect(n1, 0, n2)
	g.Connect(n2, 0, n3)
	g.Connect(n3, 0, n1)

	r := CreateTestRegion(g)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []NodeID{g.Entry(), n1, n2, n3}, t)

	idoms := FindIdoms(ordered)
	checkIntList(idoms, []int{0, 0, 1, 2}, t)
}

func TestIrreducible(t *testing.T) {
	g := CreateGraph(nil, nil)
	n1 := g.CreateNode("1", 1)
	n2 := g.CreateNode("2", 2)
	n3 := g.CreateNode("3", 1)
	n4 := g.CreateNode("4", 2)
	n5 := g.CreateNode("5", 1)
	n6 := g.CreateNode("6", 2)

	g.Connect(g.Entry(), 0, n6)

	g.Connect(n6, 0, n5)
	g.Connect(n6, 1, n4)

	g.Connect(n5, 0, n1)

	g.Connect(n4, 0, n2)
	g.Connect(n4, 1, n3)

	g.Connect(n3, 0, n2)

	g.Connect(n2, 0, n1)
	g.Connect(n2, 1, n3)

	g.Connect(n1, 0, n2)

	r := CreateTestRegion(g)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []NodeID{g.Entry(), n6, n5, n4, n3, n2, n1}, t)

	idoms := FindIdoms(ordered)
	checkIntList(idoms, []int{0, 0, 1, 1, 1, 1, 1}, t)
}

//   0
//   |
//   1
//  / \
// 2   3
//  \ /
//   4
//   |
//   5
func TestDiamond(t *testing.T) {
	g := CreateGraph(nil, nil)
	n1 := g.CreateNode("1", 2)
	n2 := g.CreateNode("2", 1)
	n3 := g.CreateNode("3", 1)
	n4 := g.CreateNode("4", 1)

	g.Connect(g.Entry(), 0, n1)

	g.Connect(n1, 0, n2)
	g.Connect(n1, 1, n3)

	g.Connect(n2, 0, n4)

	g.Connect(n3, 0, n4)

	g.Connect(n4, 0, g.Exit())

	r := CreateTestRegion(g)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []NodeID{g.Entry(), n1, n2, n3, n4, g.Exit()}, t)

	idoms := FindIdoms(ordered)
	checkIntList(idoms, []int{0, 0, 1, 1, 1, 4}, t)

	df := FindFrontiers(ordered, idoms)
	checkIntListList(df, [][]int{[]int{}, []int{}, []int{4}, []int{4}, []int{}, []int{}}, t)
}

//   0
//   |
//   1
//   |\
//   2 \
//  / \ \
// 3   4 6
//  \ / /
//   5 /
//   |/
//   7
//   |
//   8
func TestDoubleDiamond(t *testing.T) {
	g := CreateGraph(nil, nil)
	n1 := g.CreateNode("1", 2)
	n2 := g.CreateNode("2", 2)
	n3 := g.CreateNode("3", 1)
	n4 := g.CreateNode("4", 1)
	n5 := g.CreateNode("5", 1)
	n6 := g.CreateNode("6", 1)
	n7 := g.CreateNode("7", 1)

	g.Connect(g.Entry(), 0, n1)

	g.Connect(n1, 0, n2)
	g.Connect(n1, 1, n6)

	g.Connect(n2, 0, n3)
	g.Connect(n2, 1, n4)

	g.Connect(n3, 0, n5)
	g.Connect(n4, 0, n5)
	g.Connect(n5, 0, n7)
	g.Connect(n6, 0, n7)

	g.Connect(n7, 0, g.Exit())

	r := CreateTestRegion(g)

	builder := CreateSSIBuilder(r, ReversePostorder(r), &SimpleLivenessOracle{})

	checkOrder(builder.nodes, []NodeID{g.Entry(), n1, n2, n3, n4, n5, n6, n7, g.Exit()}, t)

	checkIntList(builder.Idoms, []int{0, 0, 1, 2, 2, 2, 1, 1, 7}, t)

	checkIntListList(builder.df, [][]int{
		[]int{},
		[]int{},
		[]int{7},
		[]int{5},
		[]int{5},
		[]int{7},
		[]int{7},
		[]int{},
		[]int{},
	}, t)

	numVars := 3
	defuse := CreateDefUse(len(builder.nodes), numVars)
	// Var 0
	defuse.AddDef(1, 0)
	// Var 1
	defuse.AddDef(0, 1)
	defuse.AddDef(3, 1)
	// Var 2
	defuse.AddDef(6, 2)

	for i := 0; i < numVars; i++ {
		SSI(builder, i, defuse.VarDefAt[i])
	}

	checkIntListList(builder.PhiFuncs, [][]int{
		[]int{},
		[]int{},
		[]int{},
		[]int{},
		[]int{},
		[]int{1},
		[]int{},
		[]int{1, 2},
		[]int{},
	}, t)
}

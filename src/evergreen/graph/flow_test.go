package graph

import (
	"evergreen/assert"
	"testing"
)

func checkEdge(e *edge, src *node, dst *node, t *testing.T) {
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
	numExits := len(node.exits)
	if len(exits) != numExits {
		t.Errorf("Expected %d exits, got %d", len(exits), numExits)
	} else {
		for i, exit := range exits {
			checkEdge(node.exits[i], node, g.nodes[exit], t)
		}
	}
}

func emitEdge(g *Graph, nid NodeID, flow int) EdgeID {
	return g.IndexedExitEdge(nid, flow)
}

func TestSimpleFlow(t *testing.T) {
	g := CreateGraph()
	n := g.CreateNode(1)
	g.Connect(g.Entry(), 0, n)
	g.Connect(n, 0, g.Exit())

	checkTopology(g, g.Entry(), []NodeID{}, []NodeID{n}, t)
	checkTopology(g, n, []NodeID{g.Entry()}, []NodeID{g.Exit()}, t)
	checkTopology(g, g.Exit(), []NodeID{n}, []NodeID{}, t)
}

func TestSliceEmptySplice(t *testing.T) {
	g := CreateGraph()
	gr0 := g.CreateRegion(2)
	gr1 := g.CreateRegion(2)
	assert.IntEquals(t, len(gr0.exits[0]), 1)
	gr1.Swap(0, 1)
	gr0.Splice(0, gr1)
	assert.IntEquals(t, len(gr0.exits[1]), 1)
}

func TestSliceEdgeEmptySplice(t *testing.T) {
	g := CreateGraph()
	gr0 := g.CreateRegion(2)
	gr1 := g.CreateRegion(2)
	n := g.CreateNode(1)
	gr0.AttachFlow(0, n)

	assert.IntEquals(t, len(gr0.exits[0]), 0)
	gr1.Swap(0, 1)
	gr0.SpliceToEdge(emitEdge(g, n, 0), gr1)
	assert.IntEquals(t, len(gr0.exits[1]), 1)
}

func TestRepeatFlow(t *testing.T) {
	g := CreateGraph()
	gr := g.CreateRegion(2)

	n := g.CreateNode(2)
	gr.AttachFlow(0, n)
	gr.RegisterExit(emitEdge(g, n, 0), 0)
	gr.RegisterExit(emitEdge(g, n, 1), 1)

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
	g := CreateGraph()
	gr := g.CreateRegion(2)
	c := g.CreateNode(1)
	d := g.CreateNode(2)
	b := g.CreateNode(1)

	gr.AttachFlow(0, c)
	gr.RegisterExit(emitEdge(g, c, 0), 0)

	gr.AttachFlow(0, d)
	gr.RegisterExit(emitEdge(g, d, 0), 0)
	gr.RegisterExit(emitEdge(g, d, 1), 1)

	gr.AttachFlow(0, b)
	gr.RegisterExit(emitEdge(g, b, 0), 0)

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

func checkNodeList(actualList []NodeID, expectedList []NodeID, t *testing.T) {
	if len(actualList) != len(expectedList) {
		t.Fatalf("%#v != %#v", actualList, expectedList)
	}
	for i, expected := range expectedList {
		if actualList[i] != expected {
			t.Fatalf("%d: %#v != %#v", i, actualList[i], expected)
		}
	}
}

func checkNodeListList(actualList [][]NodeID, expectedList [][]NodeID, t *testing.T) {
	if len(actualList) != len(expectedList) {
		t.Fatalf("%#v != %#v", actualList, expectedList)
	}
	for i, expected := range expectedList {
		checkNodeList(actualList[i], expected, t)
	}
}

func TestSanity(t *testing.T) {
	g := CreateGraph()

	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode(1)
	n2 := g.CreateNode(1)
	n3 := g.CreateNode(1)

	g.Connect(e, 0, n1)
	g.Connect(n1, 0, n2)
	g.Connect(n2, 0, n3)
	g.Connect(n3, 0, x)

	order, index := ReversePostorder(g)
	checkNodeList(order, []NodeID{e, n1, n2, n3, x}, t)
	assert.IntListEquals(t, index, []int{0, 4, 1, 2, 3})

	idoms := FindDominators(g, order, index)
	checkNodeList(idoms, []NodeID{e, n3, e, n1, n2}, t)
}

func TestDead(t *testing.T) {
	g := CreateGraph()

	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode(1)
	n2 := g.CreateNode(1)
	n3 := g.CreateNode(1)
	n4 := g.CreateNode(1)

	g.Connect(e, 0, n1)
	g.Connect(n1, 0, n2)
	g.Connect(n2, 0, n3)
	g.Connect(n3, 0, x)
	g.Connect(n4, 0, n3)

	order, index := ReversePostorder(g)
	checkNodeList(order, []NodeID{e, n1, n2, n3, x}, t)
	assert.IntListEquals(t, index, []int{0, 4, 1, 2, 3, -1})

	idoms := FindDominators(g, order, index)
	checkNodeList(idoms, []NodeID{e, n3, e, n1, n2, NoNode}, t)
}

func TestLoop(t *testing.T) {
	g := CreateGraph()

	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode(1)
	n2 := g.CreateNode(1)
	n3 := g.CreateNode(1)

	g.Connect(e, 0, n1)
	g.Connect(n1, 0, n2)
	g.Connect(n2, 0, n3)
	g.Connect(n3, 0, n1)

	order, index := ReversePostorder(g)
	checkNodeList(order, []NodeID{e, n1, n2, n3, x}, t)
	assert.IntListEquals(t, index, []int{0, 4, 1, 2, 3})

	idoms := FindDominators(g, order, index)
	checkNodeList(idoms, []NodeID{e, NoNode, e, n1, n2}, t)
}

func TestIrreducible(t *testing.T) {
	g := CreateGraph()

	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode(1)
	n2 := g.CreateNode(2)
	n3 := g.CreateNode(1)
	n4 := g.CreateNode(2)
	n5 := g.CreateNode(1)
	n6 := g.CreateNode(2)

	g.Connect(e, 0, n6)

	g.Connect(n6, 0, n5)
	g.Connect(n6, 1, n4)

	g.Connect(n5, 0, n1)

	g.Connect(n4, 0, n2)
	g.Connect(n4, 1, n3)

	g.Connect(n3, 0, n2)

	g.Connect(n2, 0, n1)
	g.Connect(n2, 1, n3)

	g.Connect(n1, 0, n2)

	order, index := ReversePostorder(g)
	checkNodeList(order, []NodeID{e, n6, n5, n4, n3, n2, n1, x}, t)
	assert.IntListEquals(t, index, []int{0, 7, 6, 5, 4, 3, 2, 1})

	idoms := FindDominators(g, order, index)
	checkNodeList(idoms, []NodeID{e, NoNode, n6, n6, n6, n6, n6, e}, t)
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
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode(2)
	n2 := g.CreateNode(1)
	n3 := g.CreateNode(1)
	n4 := g.CreateNode(1)

	g.Connect(e, 0, n1)

	g.Connect(n1, 0, n2)
	g.Connect(n1, 1, n3)

	g.Connect(n2, 0, n4)

	g.Connect(n3, 0, n4)

	g.Connect(n4, 0, x)

	order, index := ReversePostorder(g)
	checkNodeList(order, []NodeID{e, n1, n2, n3, n4, x}, t)
	assert.IntListEquals(t, index, []int{0, 5, 1, 2, 3, 4})

	idoms := FindDominators(g, order, index)
	checkNodeList(idoms, []NodeID{e, n4, e, n1, n1, n1}, t)

	df := FindDominanceFrontiers(g, idoms)

	checkNodeListList(df, [][]NodeID{
		[]NodeID{}, []NodeID{}, []NodeID{}, []NodeID{n4}, []NodeID{n4}, []NodeID{},
	}, t)
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
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode(2)
	n2 := g.CreateNode(2)
	n3 := g.CreateNode(1)
	n4 := g.CreateNode(1)
	n5 := g.CreateNode(1)
	n6 := g.CreateNode(1)
	n7 := g.CreateNode(1)

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

	builder := CreateSSIBuilder(g, &SimpleLivenessOracle{})

	checkNodeList(builder.order, []NodeID{e, n1, n2, n3, n4, n5, n6, n7, x}, t)
	checkNodeList(builder.Idoms, []NodeID{e, n7, e, n1, n2, n2, n2, n1, n1}, t)

	checkNodeListList(builder.df, [][]NodeID{
		[]NodeID{},
		[]NodeID{},
		[]NodeID{},
		[]NodeID{n7},
		[]NodeID{n5},
		[]NodeID{n5},
		[]NodeID{n7},
		[]NodeID{n7},
		[]NodeID{},
	}, t)

	numVars := 3
	defuse := CreateDefUse(len(g.nodes), numVars)
	// Var 0
	defuse.AddDef(n1, 0)
	// Var 1
	defuse.AddDef(e, 1)
	defuse.AddDef(n3, 1)
	// Var 2
	defuse.AddDef(n6, 2)

	for i := 0; i < numVars; i++ {
		SSI(builder, i, defuse.VarDefAt[i])
	}

	assert.IntListListEquals(

		t, builder.PhiFuncs, [][]int{
			[]int{},
			[]int{},
			[]int{},
			[]int{},
			[]int{},
			[]int{},
			[]int{1},
			[]int{},
			[]int{1, 2},
		})

}

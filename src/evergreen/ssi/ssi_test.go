package ssi

import (
	"evergreen/assert"
	"evergreen/graph"
	"testing"
)

func checkNodeList(actualList []graph.NodeID, expectedList []graph.NodeID, t *testing.T) {
	if len(actualList) != len(expectedList) {
		t.Fatalf("%#v != %#v", actualList, expectedList)
	}
	for i, expected := range expectedList {
		if actualList[i] != expected {
			t.Fatalf("%d: %#v != %#v", i, actualList[i], expected)
		}
	}
}

func checkNodeListList(actualList [][]graph.NodeID, expectedList [][]graph.NodeID, t *testing.T) {
	if len(actualList) != len(expectedList) {
		t.Fatalf("%#v != %#v", actualList, expectedList)
	}
	for i, expected := range expectedList {
		checkNodeList(actualList[i], expected, t)
	}
}

func emitFullEdge(g *graph.Graph, src graph.NodeID, dst graph.NodeID) graph.EdgeID {
	e := g.CreateEdge()
	g.ConnectEdge(src, e, dst)
	return e
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
	g := graph.CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()
	n5 := g.CreateNode()
	n6 := g.CreateNode()
	n7 := g.CreateNode()

	emitFullEdge(g, g.Entry(), n1)

	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, n6)

	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n2, n4)

	emitFullEdge(g, n3, n5)
	emitFullEdge(g, n4, n5)
	emitFullEdge(g, n5, n7)
	emitFullEdge(g, n6, n7)

	emitFullEdge(g, n7, g.Exit())

	builder := CreateSSIBuilder(g, &SimpleLivenessOracle{})

	checkNodeList(builder.order, []graph.NodeID{e, n1, n2, n3, n4, n5, n6, n7, x}, t)
	checkNodeList(builder.Idoms, []graph.NodeID{e, n7, e, n1, n2, n2, n2, n1, n1}, t)

	checkNodeListList(builder.df, [][]graph.NodeID{
		[]graph.NodeID{},
		[]graph.NodeID{},
		[]graph.NodeID{},
		[]graph.NodeID{n7},
		[]graph.NodeID{n5},
		[]graph.NodeID{n5},
		[]graph.NodeID{n7},
		[]graph.NodeID{n7},
		[]graph.NodeID{},
	}, t)

	numVars := 3
	defuse := CreateDefUse(g.NumNodes(), numVars)
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

package graph

import (
	"evergreen/assert"
	"testing"
)

func TestLinear(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n3, n4)
	emitFullEdge(g, n4, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "(0 2 3 4 5 1)")
}

func TestLinearSkip(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, n4) // skip
	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n3, n4)
	emitFullEdge(g, n4, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0 2) (3 4) (5 1)]")
}

func TestClusterDiamond(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)

	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, n3)

	emitFullEdge(g, n2, n4)

	emitFullEdge(g, n3, n4)

	emitFullEdge(g, n4, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0 2) <(3) (4)> (5 1)]")
}

func TestClusterCross(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, e, n2)

	emitFullEdge(g, n1, n3)
	emitFullEdge(g, n1, n4)

	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n2, n4)

	emitFullEdge(g, n3, x)
	emitFullEdge(g, n4, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) <(2) (3)> <(4) (5)> (1)]")
}

func TestLadder2(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n2, x)

	emitFullEdge(g, e, n3)
	emitFullEdge(g, n1, n3)
	emitFullEdge(g, n3, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) (2) <(3) (4)> (1)]")
}

func TestLadder3(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n3, x)

	emitFullEdge(g, e, n4)
	emitFullEdge(g, n1, n4)
	emitFullEdge(g, n2, n4)
	emitFullEdge(g, n4, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) (2) (3) <(4) (5)> (1)]")

	info := findLoops(g)
	makeCluster2(g, info)
}

func TestLadderSkip(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, e, n2)

	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, x)

	emitFullEdge(g, n2, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) (2) (3) (1)]")
}

func TestLadderSkipComplex(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()
	n5 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, e, n2)

	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, n3)
	emitFullEdge(g, n1, n4)

	emitFullEdge(g, n2, n3)

	emitFullEdge(g, n3, x)

	emitFullEdge(g, n4, n5)

	emitFullEdge(g, n5, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) (2) <(3) (5 6)> (4) (1)]")
}

func TestCrossEdgeToLoop(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()

	emitFullEdge(g, e, n1)

	emitFullEdge(g, e, n2)
	emitFullEdge(g, n1, n2)

	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n3, n2)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) (2) {(3 4)}]")

	info := findLoops(g)
	makeCluster2(g, info)
}

func TestInnerOuter(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()
	n5 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, e, n5)

	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, n3)

	emitFullEdge(g, n2, x)

	emitFullEdge(g, n3, n4)

	emitFullEdge(g, n4, n5)

	emitFullEdge(g, n5, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0 2) <(3) (4 5)> (6) (1)]")
}

func Test2Levelif(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()
	n5 := g.CreateNode()
	n6 := g.CreateNode()
	n7 := g.CreateNode()
	n8 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, e, n2)

	emitFullEdge(g, n1, n3)
	emitFullEdge(g, n1, n4)

	emitFullEdge(g, n2, n5)
	emitFullEdge(g, n2, n6)

	emitFullEdge(g, n3, n7)

	emitFullEdge(g, n4, n7)

	emitFullEdge(g, n5, n8)

	emitFullEdge(g, n6, n8)

	emitFullEdge(g, n7, x)

	emitFullEdge(g, n8, x)

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) <[(2) <(4) (5)> (8)] [(3) <(6) (7)> (9)]> (1)]")

	info := findLoops(g)
	makeCluster2(g, info)
}

func assertNodeInfos(actual []lfNodeInfo, expected []lfNodeInfo, t *testing.T) {
	assert.IntEquals(t, len(actual), len(expected))

	for i := 0; i < len(expected); i++ {
		a := actual[i]
		e := expected[i]
		if a.containingHead != e.containingHead {
			t.Errorf("%d loop head: %v vs %v", i, a.containingHead, e.containingHead)
		}
		if a.isHead != e.isHead {
			t.Errorf("%d is head: %v vs %v", i, a.isHead, e.isHead)
		}
		if a.isIrreducible != e.isIrreducible {
			t.Errorf("%d is irreducible: %v vs %v", i, a.isIrreducible, e.isIrreducible)
		}
	}
}

func TestForwardEdgeLoop(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, n3)

	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n2, x)

	emitFullEdge(g, n3, n4)

	emitFullEdge(g, n4, n2)

	info := findLoops(g)

	assertNodeInfos(info, []lfNodeInfo{
		// e
		lfNodeInfo{
			containingHead: NoNode,
		},
		// x
		lfNodeInfo{
			containingHead: NoNode,
		},
		// n1
		lfNodeInfo{
			containingHead: NoNode,
		},
		// n2
		lfNodeInfo{
			containingHead: NoNode,
			isHead:         true,
			isIrreducible:  true,
		},
		// n3
		lfNodeInfo{
			containingHead: n2,
		},
		// n4
		lfNodeInfo{
			containingHead: n2,
		},
	}, t)
}

func TestClassicIrreducible(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, e, n4)

	emitFullEdge(g, n1, n2)

	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n2, n1)
	emitFullEdge(g, n2, x)

	emitFullEdge(g, n3, n4)
	emitFullEdge(g, n3, n2)
	emitFullEdge(g, n3, x)

	emitFullEdge(g, n4, n3)

	info := findLoops(g)

	assertNodeInfos(info, []lfNodeInfo{
		// e
		lfNodeInfo{
			containingHead: NoNode,
		},
		// x
		lfNodeInfo{
			containingHead: NoNode,
		},
		// n1
		lfNodeInfo{
			containingHead: NoNode,
			isHead:         true,
		},
		// n2
		lfNodeInfo{
			containingHead: n1,
			isHead:         true,
		},
		// n3
		lfNodeInfo{
			containingHead: n2,
			isHead:         true,
			isIrreducible:  true,
		},
		// n4
		lfNodeInfo{
			containingHead: n3,
		},
	}, t)
}

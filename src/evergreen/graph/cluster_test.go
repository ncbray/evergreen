package graph

import (
	"evergreen/assert"
	"testing"
)

func assertNodeInfos(actual []NodeInfo, expected []NodeInfo, t *testing.T) {
	assert.IntEquals(t, len(actual), len(expected))

	for i := 0; i < len(expected); i++ {
		a := actual[i]
		e := expected[i]
		if a.IDom != e.IDom {
			t.Errorf("%d idom: %v vs %v", i, a.IDom, e.IDom)
		}
		if a.LoopHead != e.LoopHead {
			t.Errorf("%d loop head: %v vs %v", i, a.LoopHead, e.LoopHead)
		}
		if a.IsHead != e.IsHead {
			t.Errorf("%d is head: %v vs %v", i, a.IsHead, e.IsHead)
		}
		if a.IsIrreducible != e.IsIrreducible {
			t.Errorf("%d is irreducible: %v vs %v", i, a.IsIrreducible, e.IsIrreducible)
		}
	}
}

func assertNodeList(actual []NodeID, expected []NodeID, t *testing.T) {
	assert.IntEquals(t, len(actual), len(expected))

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("%d: %v vs %v", i, actual[i], expected[i])
		}
	}
}

func assertEdgeTypeList(actual []EdgeType, expected []EdgeType, t *testing.T) {
	assert.IntEquals(t, len(actual), len(expected))

	for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("%d: %v vs %v", i, actual[i], expected[i])
		}
	}
}

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

	cluster := MakeCluster(g)
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

	cluster := MakeCluster(g)
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

	cluster := MakeCluster(g)
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

	cluster := MakeCluster(g)
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

	cluster := MakeCluster(g)
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

	cluster := MakeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) (2) (3) <(4) (5)> (1)]")
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

	cluster := MakeCluster(g)
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

	cluster := MakeCluster(g)
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

	cluster := MakeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) (2) {(3 4)}]")
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

	cluster := MakeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) [(2) <(3) (4 5)>] (6) (1)]")
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

	cluster := MakeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) <[(2) <(4) (5)> (8)] [(3) <(6) (7)> (9)]> (1)]")

	info, edges, postorder := AnalyzeStructure(g)

	assertNodeInfos(info, []NodeInfo{
		// e
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// x
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// n1
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// n2
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// n3
		NodeInfo{
			IDom:     n1,
			LoopHead: NoNode,
		},
		// n4
		NodeInfo{
			IDom:     n1,
			LoopHead: NoNode,
		},
		// n5
		NodeInfo{
			IDom:     n2,
			LoopHead: NoNode,
		},
		// n6
		NodeInfo{
			IDom:     n2,
			LoopHead: NoNode,
		},
		// n7
		NodeInfo{
			IDom:     n1,
			LoopHead: NoNode,
		},
		// n8
		NodeInfo{
			IDom:     n2,
			LoopHead: NoNode,
		},
	}, t)

	assertEdgeTypeList(edges, []EdgeType{FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, CROSS, FORWARD, CROSS, FORWARD, CROSS}, t)

	assertNodeList(postorder, []NodeID{
		x,
		n7,
		n3,
		n4,
		n1,
		n8,
		n5,
		n6,
		n2,
		e,
	}, t)
}

func TestDualLoop(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)

	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n1, x)

	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n2, n4)

	emitFullEdge(g, n3, n2)

	emitFullEdge(g, n4, n1)

	info, edges, postorder := AnalyzeStructure(g)

	assertNodeInfos(info, []NodeInfo{
		// e
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// x
		NodeInfo{
			IDom:     n1,
			LoopHead: NoNode,
		},
		// n1
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
			IsHead:   true,
		},
		// n2
		NodeInfo{
			IDom:     n1,
			LoopHead: n1,
			IsHead:   true,
		},
		// n3
		NodeInfo{
			IDom:     n2,
			LoopHead: n2,
		},
		// n4
		NodeInfo{
			IDom:     n2,
			LoopHead: n1,
		},
	}, t)

	assertEdgeTypeList(edges, []EdgeType{FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, BACKWARD, BACKWARD}, t)

	assertNodeList(postorder, []NodeID{
		n3,
		n4,
		n2,
		x,
		n1,
		e,
	}, t)
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

	info, edges, postorder := AnalyzeStructure(g)

	assertNodeInfos(info, []NodeInfo{
		// e
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// x
		NodeInfo{
			IDom:     n2,
			LoopHead: NoNode,
		},
		// n1
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// n2
		NodeInfo{
			IDom:          n1,
			LoopHead:      NoNode,
			IsHead:        true,
			IsIrreducible: true,
		},
		// n3
		NodeInfo{
			IDom:     n1,
			LoopHead: n2,
		},
		// n4
		NodeInfo{
			IDom:     n3,
			LoopHead: n2,
		},
	}, t)

	assertEdgeTypeList(edges, []EdgeType{FORWARD, FORWARD, REENTRY, FORWARD, FORWARD, FORWARD, BACKWARD}, t)

	assertNodeList(postorder, []NodeID{
		n4,
		n3,
		x,
		n2,
		n1,
		e,
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

	info, edges, postorder := AnalyzeStructure(g)

	assertNodeInfos(info, []NodeInfo{
		// e
		NodeInfo{
			LoopHead: NoNode,
			IDom:     e,
		},
		// x
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
		},
		// n1
		NodeInfo{
			IDom:     e,
			LoopHead: NoNode,
			IsHead:   true,
		},
		// n2
		NodeInfo{
			IDom:     e,
			LoopHead: n1,
			IsHead:   true,
		},
		// n3
		NodeInfo{
			IDom:          e,
			LoopHead:      n2,
			IsHead:        true,
			IsIrreducible: true,
		},
		// n4
		NodeInfo{
			IDom:     e,
			LoopHead: n3,
		},
	}, t)

	assertEdgeTypeList(edges, []EdgeType{FORWARD, REENTRY, FORWARD, FORWARD, BACKWARD, CROSS, FORWARD, BACKWARD, FORWARD, BACKWARD}, t)

	assertNodeList(postorder, []NodeID{
		n4,
		x,
		n3,
		n2,
		n1,
		e,
	}, t)
}

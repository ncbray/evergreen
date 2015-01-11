package graph

import (
	"evergreen/assert"
	"testing"
)

func assertNodeInfos(actual []nodeInfo, expected []nodeInfo, t *testing.T) {
	assert.IntEquals(t, len(actual), len(expected))

	for i := 0; i < len(expected); i++ {
		a := actual[i]
		e := expected[i]
		if a.idom != e.idom {
			t.Errorf("%d idom: %v vs %v", i, a.idom, e.idom)
		}
		if a.loopHead != e.loopHead {
			t.Errorf("%d loop head: %v vs %v", i, a.loopHead, e.loopHead)
		}
		if a.isHead != e.isHead {
			t.Errorf("%d is head: %v vs %v", i, a.isHead, e.isHead)
		}
		if a.isIrreducible != e.isIrreducible {
			t.Errorf("%d is irreducible: %v vs %v", i, a.isIrreducible, e.isIrreducible)
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

func assertEdgeTypeList(actual []edgeType, expected []edgeType, t *testing.T) {
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

	info, edges, postorder := analyzeStructure(g)
	makeCluster2(g, info, edges, postorder)
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

	info, edges, postorder := analyzeStructure(g)
	makeCluster2(g, info, edges, postorder)
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

	info, edges, postorder := analyzeStructure(g)
	makeCluster2(g, info, edges, postorder)
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

	cluster := makeCluster(g)
	assert.StringEquals(t, cluster.DumpShort(), "[(0) <[(2) <(4) (5)> (8)] [(3) <(6) (7)> (9)]> (1)]")

	info, edges, postorder := analyzeStructure(g)

	assertNodeInfos(info, []nodeInfo{
		// e
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// x
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// n1
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// n2
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// n3
		nodeInfo{
			idom:     n1,
			loopHead: NoNode,
		},
		// n4
		nodeInfo{
			idom:     n1,
			loopHead: NoNode,
		},
		// n5
		nodeInfo{
			idom:     n2,
			loopHead: NoNode,
		},
		// n6
		nodeInfo{
			idom:     n2,
			loopHead: NoNode,
		},
		// n7
		nodeInfo{
			idom:     n1,
			loopHead: NoNode,
		},
		// n8
		nodeInfo{
			idom:     n2,
			loopHead: NoNode,
		},
	}, t)

	assertEdgeTypeList(edges, []edgeType{FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, CROSS, FORWARD, CROSS, FORWARD, CROSS}, t)

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

	makeCluster2(g, info, edges, postorder)
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

	info, edges, postorder := analyzeStructure(g)

	assertNodeInfos(info, []nodeInfo{
		// e
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// x
		nodeInfo{
			idom:     n1,
			loopHead: NoNode,
		},
		// n1
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
			isHead:   true,
		},
		// n2
		nodeInfo{
			idom:     n1,
			loopHead: n1,
			isHead:   true,
		},
		// n3
		nodeInfo{
			idom:     n2,
			loopHead: n2,
		},
		// n4
		nodeInfo{
			idom:     n2,
			loopHead: n1,
		},
	}, t)

	assertEdgeTypeList(edges, []edgeType{FORWARD, FORWARD, FORWARD, FORWARD, FORWARD, BACKWARD, BACKWARD}, t)

	assertNodeList(postorder, []NodeID{
		n3,
		n4,
		n2,
		x,
		n1,
		e,
	}, t)

	makeCluster2(g, info, edges, postorder)
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

	info, edges, postorder := analyzeStructure(g)

	assertNodeInfos(info, []nodeInfo{
		// e
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// x
		nodeInfo{
			idom:     n2,
			loopHead: NoNode,
		},
		// n1
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// n2
		nodeInfo{
			idom:          n1,
			loopHead:      NoNode,
			isHead:        true,
			isIrreducible: true,
		},
		// n3
		nodeInfo{
			idom:     n1,
			loopHead: n2,
		},
		// n4
		nodeInfo{
			idom:     n3,
			loopHead: n2,
		},
	}, t)

	assertEdgeTypeList(edges, []edgeType{FORWARD, FORWARD, REENTRY, FORWARD, FORWARD, FORWARD, BACKWARD}, t)

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

	info, edges, postorder := analyzeStructure(g)

	assertNodeInfos(info, []nodeInfo{
		// e
		nodeInfo{
			loopHead: NoNode,
			idom:     e,
		},
		// x
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
		},
		// n1
		nodeInfo{
			idom:     e,
			loopHead: NoNode,
			isHead:   true,
		},
		// n2
		nodeInfo{
			idom:     e,
			loopHead: n1,
			isHead:   true,
		},
		// n3
		nodeInfo{
			idom:          e,
			loopHead:      n2,
			isHead:        true,
			isIrreducible: true,
		},
		// n4
		nodeInfo{
			idom:     e,
			loopHead: n3,
		},
	}, t)

	assertEdgeTypeList(edges, []edgeType{FORWARD, REENTRY, FORWARD, FORWARD, BACKWARD, CROSS, FORWARD, BACKWARD, FORWARD, BACKWARD}, t)

	assertNodeList(postorder, []NodeID{
		n4,
		x,
		n3,
		n2,
		n1,
		e,
	}, t)
}

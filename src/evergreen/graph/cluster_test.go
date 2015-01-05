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

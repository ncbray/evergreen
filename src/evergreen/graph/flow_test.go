package graph

import (
	"evergreen/assert"
	"testing"
)

func checkEdgeConsistency(e *edge, t *testing.T) {
	if e.prevEntry != nil && (e.prevEntry.nextEntry != e || e.prevEntry.dst != e.dst) {
		t.Errorf("Bad prevEntry for %#v / %#v", e, e.prevEntry)
	}
	if e.nextEntry != nil && (e.nextEntry.prevEntry != e || e.nextEntry.dst != e.dst) {
		t.Errorf("Bad nextEntry for %#v / %#v", e, e.nextEntry)
	}
	if e.prevExit != nil && (e.prevExit.nextExit != e || e.prevExit.src != e.src) {
		t.Errorf("Bad prevExit for %#v / %#v", e, e.prevExit)
	}
	if e.nextExit != nil && (e.nextExit.prevExit != e || e.nextExit.src != e.src) {
		t.Errorf("Bad nextExit for %#v / %#v", e, e.nextExit)
	}
}

func checkEntryEdges(actual entryEdges, expected []*edge, t *testing.T) {
	current := actual.head
	count := 0
	for current != nil {
		if count < len(expected) {
			other := expected[count]
			if current != other {
				t.Errorf("Expected %#v at position %d, got %#v", other, count, current)
			}
		}
		checkEdgeConsistency(current, t)
		current = current.nextEntry
		count += 1
	}
	if len(expected) != count {
		t.Errorf("Expected %d elements, got %d", len(expected), count)
	}
}

func simpleEntry(count int) (entryEdges, []*edge) {
	l := emptyEntry()
	edges := make([]*edge, count)
	for i := 0; i < count; i++ {
		e := &edge{id: EdgeID(i)}
		l.Append(e)
		edges[i] = e
	}
	return l, edges
}

func TestEntryEdgeEmpty(t *testing.T) {
	l := emptyEntry()
	if l.HasMultipleEdges() {
		t.Errorf("Should not have multiple edges.")
	}
	checkEntryEdges(l, []*edge{}, t)
}

func TestEntryEdgeSimple1Sanity(t *testing.T) {
	l, e := simpleEntry(1)
	if l.HasMultipleEdges() {
		t.Errorf("Should not have multiple edges.")
	}
	checkEntryEdges(l, []*edge{e[0]}, t)
}

func TestEntryEdgeSimple2Sanity(t *testing.T) {
	l, e := simpleEntry(2)
	if !l.HasMultipleEdges() {
		t.Errorf("Should have multiple edges.")
	}
	checkEntryEdges(l, []*edge{e[0], e[1]}, t)
}

func TestEntryEdgeSimple3Sanity(t *testing.T) {
	l, e := simpleEntry(3)
	if !l.HasMultipleEdges() {
		t.Errorf("Should have multiple edges.")
	}
	checkEntryEdges(l, []*edge{e[0], e[1], e[2]}, t)
}

func TestEntryEdgeSimple3x1Beginning(t *testing.T) {
	l, e := simpleEntry(3)
	r := &edge{id: 3}
	l.ReplaceEdge(e[0], r)
	checkEntryEdges(l, []*edge{r, e[1], e[2]}, t)
}

func TestEntryEdgeSimple3x1Middle(t *testing.T) {
	l, e := simpleEntry(3)
	r := &edge{id: 3}
	l.ReplaceEdge(e[1], r)
	checkEntryEdges(l, []*edge{e[0], r, e[2]}, t)
}

func TestEntryEdgeSimple3x1End(t *testing.T) {
	l, e := simpleEntry(3)
	r := &edge{id: 3}
	l.ReplaceEdge(e[2], r)
	checkEntryEdges(l, []*edge{e[0], e[1], r}, t)
}

func TestEntryEdgeSimple3x0Beginning(t *testing.T) {
	l, e := simpleEntry(3)
	other, _ := simpleEntry(0)
	l.ReplaceEdgeWithMultiple(e[0], other)
	checkEntryEdges(l, []*edge{e[1], e[2]}, t)
}

func TestEntryEdgeSimple3x0Middle(t *testing.T) {
	l, e := simpleEntry(3)
	other, _ := simpleEntry(0)
	l.ReplaceEdgeWithMultiple(e[1], other)
	checkEntryEdges(l, []*edge{e[0], e[2]}, t)
}

func TestEntryEdgeSimple3x0End(t *testing.T) {
	l, e := simpleEntry(3)
	other, _ := simpleEntry(0)
	l.ReplaceEdgeWithMultiple(e[2], other)
	checkEntryEdges(l, []*edge{e[0], e[1]}, t)
}

func TestEntryEdgeSimple3x3Beginning(t *testing.T) {
	l, e := simpleEntry(3)
	other, f := simpleEntry(3)
	l.ReplaceEdgeWithMultiple(e[0], other)
	checkEntryEdges(l, []*edge{f[0], f[1], f[2], e[1], e[2]}, t)
}

func TestEntryEdgeSimple3x3Middle(t *testing.T) {
	l, e := simpleEntry(3)
	other, f := simpleEntry(3)
	l.ReplaceEdgeWithMultiple(e[1], other)
	checkEntryEdges(l, []*edge{e[0], f[0], f[1], f[2], e[2]}, t)
}

func TestEntryEdgeSimple3x3End(t *testing.T) {
	l, e := simpleEntry(3)
	other, f := simpleEntry(3)
	l.ReplaceEdgeWithMultiple(e[2], other)
	checkEntryEdges(l, []*edge{e[0], e[1], f[0], f[1], f[2]}, t)
}

func TestEntryEdgeSimple1x3(t *testing.T) {
	l, e := simpleEntry(1)
	other, f := simpleEntry(3)
	l.ReplaceEdgeWithMultiple(e[0], other)
	checkEntryEdges(l, []*edge{f[0], f[1], f[2]}, t)
}

func TestEntryEdgeSimple1x0(t *testing.T) {
	l, e := simpleEntry(1)
	other, _ := simpleEntry(0)
	l.ReplaceEdgeWithMultiple(e[0], other)
	checkEntryEdges(l, []*edge{}, t)
}

func checkEdge(e *edge, src *node, dst *node, t *testing.T) {
	if e.src != src {
		t.Errorf("Got src of %#v, expected %#v", e.src, src)
	}
	if e.dst != dst {
		t.Errorf("Got dst of %#v, expected %#v", e.dst, dst)
	}
	checkEdgeConsistency(e, t)
}

func checkTopology(g *Graph, id NodeID, entries []NodeID, exits []NodeID, t *testing.T) {
	node := g.nodes[id]
	if node == nil {
		t.Error("Node should not be nil")
		return
	}
	{
		it := g.EntryIterator(id)
		count := 0
		for it.HasNext() {
			_, e := it.GetNext()
			if count < len(entries) {
				checkEdge(g.edges[e], g.nodes[entries[count]], node, t)
			}
			count += 1
		}
		if len(entries) != count {
			t.Errorf("Expected %d entries, got %d", len(entries), count)
		}
	}
	{
		it := g.ExitIterator(id)
		count := 0
		for it.HasNext() {
			e, _ := it.GetNext()
			if count < len(exits) {
				checkEdge(g.edges[e], node, g.nodes[exits[count]], t)
			}
			count += 1
		}
		if len(exits) != count {
			t.Errorf("Expected %d entries, got %d", len(exits), count)
		}
	}
}

func emitDanglingEdge(g *Graph, src NodeID) EdgeID {
	e := g.CreateEdge()
	g.ConnectEdgeEntry(src, e)
	return e
}

func emitFullEdge(g *Graph, src NodeID, dst NodeID) EdgeID {
	e := g.CreateEdge()
	g.ConnectEdge(src, e, dst)
	return e
}

func TestSimpleFlow(t *testing.T) {
	g := CreateGraph()
	n := g.CreateNode()
	emitFullEdge(g, g.Entry(), n)
	emitFullEdge(g, n, g.Exit())

	checkTopology(g, g.Entry(), []NodeID{}, []NodeID{n}, t)
	checkTopology(g, n, []NodeID{g.Entry()}, []NodeID{g.Exit()}, t)
	checkTopology(g, g.Exit(), []NodeID{n}, []NodeID{}, t)
}

func numFlowEdges(fb *FlowBuilder, flow int) int {
	count := 0
	current := fb.flows[flow].head
	for current != nil {
		current = current.nextEntry
		count += 1
	}
	return count
}

func TestSliceEmptySplice(t *testing.T) {
	g := CreateGraph()
	fb0 := CreateFlowBuilder(g, g.CreateEdge(), 2)
	assert.IntEquals(t, numFlowEdges(fb0, 0), 1)
	fb1 := fb0.SplitOffFlow(0)
	assert.IntEquals(t, numFlowEdges(fb0, 0), 0)
	fb0.AbsorbExits(fb1)
	assert.IntEquals(t, numFlowEdges(fb0, 0), 1)
}

func TestSliceEdgeEmptySplice(t *testing.T) {
	g := CreateGraph()
	fb0 := CreateFlowBuilder(g, g.CreateEdge(), 2)
	n := g.CreateNode()
	fb0.AttachFlow(0, n)

	assert.IntEquals(t, numFlowEdges(fb0, 0), 0)
	fb1 := fb0.SplitOffEdge(emitDanglingEdge(g, n))
	fb0.AbsorbExits(fb1)
	assert.IntEquals(t, numFlowEdges(fb0, 0), 1)
}

func TestRepeatFlow(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()

	entry := emitDanglingEdge(g, e)

	fb := CreateFlowBuilder(g, entry, 2)

	n := g.CreateNode()
	fb.AttachFlow(0, n)
	fb.RegisterExit(emitDanglingEdge(g, n), 0)
	fb.RegisterExit(emitDanglingEdge(g, n), 1)

	// Normal flow iterates
	fb.AttachFlow(0, n)

	// Stop iterating on failure
	fb.AttachFlow(1, x)

	checkTopology(g, e, []NodeID{}, []NodeID{n}, t)
	checkTopology(g, n, []NodeID{e, n}, []NodeID{n, x}, t)
	checkTopology(g, x, []NodeID{n}, []NodeID{}, t)
}

func TestWhileFlow(t *testing.T) {
	g := CreateGraph()
	e := g.Entry()
	x := g.Exit()

	entry := emitDanglingEdge(g, e)

	fb := CreateFlowBuilder(g, entry, 2)

	c := g.CreateNode()
	d := g.CreateNode()
	b := g.CreateNode()

	fb.AttachFlow(0, c)
	fb.RegisterExit(emitDanglingEdge(g, c), 0)

	fb.AttachFlow(0, d)
	fb.RegisterExit(emitDanglingEdge(g, d), 0)
	fb.RegisterExit(emitDanglingEdge(g, d), 1)

	fb.AttachFlow(0, b)
	fb.RegisterExit(emitDanglingEdge(g, b), 0)

	fb.AttachFlow(0, c)
	fb.AttachFlow(1, x)

	checkTopology(g, e, []NodeID{}, []NodeID{c}, t)
	checkTopology(g, c, []NodeID{e, b}, []NodeID{d}, t)
	checkTopology(g, d, []NodeID{c}, []NodeID{b, x}, t)
	checkTopology(g, b, []NodeID{d}, []NodeID{c}, t)
	checkTopology(g, x, []NodeID{d}, []NodeID{}, t)
}

func TestInsertInEdge1(t *testing.T) {
	g := CreateGraph()

	e := g.Entry()
	x := g.Exit()
	r := g.CreateNode()

	tgt := emitFullEdge(g, e, x)
	re := emitDanglingEdge(g, r)

	g.InsertInEdge(re, tgt)

	checkTopology(g, e, []NodeID{}, []NodeID{r}, t)
	checkTopology(g, r, []NodeID{e}, []NodeID{x}, t)
	checkTopology(g, x, []NodeID{r}, []NodeID{}, t)
}

func TestInsertInEdge3(t *testing.T) {
	g := CreateGraph()

	e := g.Entry()
	x := g.Exit()
	n := g.CreateNode()
	r := g.CreateNode()

	emitFullEdge(g, e, n)
	emitFullEdge(g, n, x)
	tgt := emitFullEdge(g, n, x)
	emitFullEdge(g, n, x)

	re := emitDanglingEdge(g, r)

	g.InsertInEdge(re, tgt)

	checkTopology(g, e, []NodeID{}, []NodeID{n}, t)
	checkTopology(g, n, []NodeID{e}, []NodeID{x, r, x}, t)
	checkTopology(g, r, []NodeID{n}, []NodeID{x}, t)
	checkTopology(g, x, []NodeID{n, r, n}, []NodeID{}, t)
}

func TestKillNode1(t *testing.T) {
	g := CreateGraph()

	e := g.Entry()
	x := g.Exit()
	n := g.CreateNode()

	emitFullEdge(g, e, n)
	emitFullEdge(g, n, x)

	g.KillNode(n)

	checkTopology(g, e, []NodeID{}, []NodeID{x}, t)
	checkTopology(g, n, []NodeID{}, []NodeID{}, t)
	checkTopology(g, x, []NodeID{e}, []NodeID{}, t)
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
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n3, x)

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
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n3, x)
	emitFullEdge(g, n4, n3)

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
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()

	emitFullEdge(g, e, n1)
	emitFullEdge(g, n1, n2)
	emitFullEdge(g, n2, n3)
	emitFullEdge(g, n3, n1)

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
	n1 := g.CreateNode()
	n2 := g.CreateNode()
	n3 := g.CreateNode()
	n4 := g.CreateNode()
	n5 := g.CreateNode()
	n6 := g.CreateNode()

	emitFullEdge(g, e, n6)

	emitFullEdge(g, n6, n5)
	emitFullEdge(g, n6, n4)

	emitFullEdge(g, n5, n1)

	emitFullEdge(g, n4, n2)
	emitFullEdge(g, n4, n3)

	emitFullEdge(g, n3, n2)

	emitFullEdge(g, n2, n1)
	emitFullEdge(g, n2, n3)

	emitFullEdge(g, n1, n2)

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

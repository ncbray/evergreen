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

type TestExit struct {
	flow int
}

type TestNode struct {
	name string
}

func CreateTestEntry() *Node {
	return CreateNode(&TestEntry{}, 1)
}

func CreateTestNode(name string, numExits int) *Node {
	return CreateNode(&TestNode{name: name}, numExits)
}

func CreateTestExit(flow int) *Node {
	return CreateNode(&TestExit{flow: flow}, 0)
}

func CreateTestRegion() *Region {
	r := &Region{
		Entry: CreateTestEntry(),
		Exits: []*Node{
			CreateTestExit(0),
			CreateTestExit(1),
		},
	}
	r.Entry.SetExit(0, r.Exits[0])
	return r
}

func TestSimpleFlow(t *testing.T) {
	r := CreateTestRegion()
	n := CreateTestNode("n", 2)
	r.Connect(0, n)
	r.AttachDefaultExits(n)

	checkTopology(r.GetEntry(), []*Node{}, []*Node{n}, t)
	checkTopology(n, []*Node{r.GetEntry()}, []*Node{r.GetExit(0), r.GetExit(1)}, t)
	checkTopology(r.GetExit(0), []*Node{n}, []*Node{}, t)
	checkTopology(r.GetExit(1), []*Node{n}, []*Node{}, t)
}

func TestRepeatFlow(t *testing.T) {
	l := CreateTestRegion()
	n := CreateTestNode("n", 2)
	l.Connect(0, n)
	l.AttachDefaultExits(n)
	head := l.Head()
	// Normal flow iterates
	l.GetExit(0).TransferEntries(head)
	// Stop iterating on failure
	l.GetExit(1).TransferEntries(l.GetExit(0))

	checkTopology(l.GetEntry(), []*Node{}, []*Node{n}, t)
	checkTopology(n, []*Node{l.GetEntry(), n}, []*Node{n, l.GetExit(0)}, t)
	checkTopology(l.GetExit(0), []*Node{n}, []*Node{}, t)
	checkTopology(l.GetExit(1), []*Node{}, []*Node{}, t)

	r := CreateTestRegion()

	r.Splice(0, l)

	checkTopology(l.GetEntry(), []*Node{}, []*Node{nil}, t)
	checkTopology(l.GetExit(0), []*Node{}, []*Node{}, t)
	checkTopology(l.GetExit(1), []*Node{}, []*Node{}, t)

	checkTopology(r.GetEntry(), []*Node{}, []*Node{n}, t)
	checkTopology(n, []*Node{r.GetEntry(), n}, []*Node{n, r.GetExit(0)}, t)
	checkTopology(r.GetExit(0), []*Node{n}, []*Node{}, t)
	checkTopology(r.GetExit(1), []*Node{}, []*Node{}, t)
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
	decide.SetExit(1, l.GetExit(0))

	r := CreateTestRegion()
	r.Splice(0, l)

	checkTopology(r.GetEntry(), []*Node{}, []*Node{cond}, t)
	checkTopology(cond, []*Node{r.GetEntry(), body}, []*Node{decide, r.GetExit(1)}, t)
	checkTopology(decide, []*Node{cond}, []*Node{body, r.GetExit(0)}, t)
	checkTopology(body, []*Node{decide}, []*Node{cond, r.GetExit(1)}, t)
	checkTopology(r.GetExit(0), []*Node{decide}, []*Node{}, t)
	checkTopology(r.GetExit(1), []*Node{cond, body}, []*Node{}, t)
}

func checkInt(name string, actual int, expected int, t *testing.T) {
	if actual != expected {
		t.Fatalf("%s: %d != %d", name, actual, expected)
	}
}

func checkOrder(actualOrder []*Node, expectedOrder []*Node, t *testing.T) {
	checkInt("len", len(actualOrder), len(expectedOrder), t)
	for i, expected := range expectedOrder {
		if actualOrder[i] != expected {
			t.Fatalf("%d: %#v != %#v", i, actualOrder[i], expected)
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

func TestSanity(t *testing.T) {
	r := CreateTestRegion()
	n1 := CreateTestNode("1", 1)
	n2 := CreateTestNode("2", 1)
	n3 := CreateTestNode("3", 1)

	r.Connect(0, n1)
	r.AttachDefaultExits(n1)

	r.Connect(0, n2)
	r.AttachDefaultExits(n2)

	r.Connect(0, n3)
	r.AttachDefaultExits(n3)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []*Node{r.Entry, n1, n2, n3, r.Exits[0]}, t)

	idoms := FindIdoms(ordered)
	checkIntList(idoms, []int{0, 0, 1, 2, 3}, t)
}

func TestLoop(t *testing.T) {
	r := CreateTestRegion()
	n1 := CreateTestNode("1", 1)
	n2 := CreateTestNode("2", 1)
	n3 := CreateTestNode("3", 1)

	r.Connect(0, n1)
	n1.SetExit(0, n2)
	n2.SetExit(0, n3)
	n3.SetExit(0, n1)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []*Node{r.Entry, n1, n2, n3}, t)

	idoms := FindIdoms(ordered)
	checkIntList(idoms, []int{0, 0, 1, 2}, t)
}

func TestIrreducible(t *testing.T) {
	r := CreateTestRegion()
	n1 := CreateTestNode("1", 1)
	n2 := CreateTestNode("2", 2)
	n3 := CreateTestNode("3", 1)
	n4 := CreateTestNode("4", 2)
	n5 := CreateTestNode("5", 1)
	n6 := CreateTestNode("6", 2)

	r.Connect(0, n6)

	n6.SetExit(0, n5)
	n6.SetExit(1, n4)

	n5.SetExit(0, n1)

	n4.SetExit(0, n2)
	n4.SetExit(1, n3)

	n3.SetExit(0, n2)

	n2.SetExit(0, n1)
	n2.SetExit(1, n3)

	n1.SetExit(0, n2)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []*Node{r.Entry, n6, n5, n4, n3, n2, n1}, t)

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
	r := CreateTestRegion()
	n1 := CreateTestNode("1", 2)
	n2 := CreateTestNode("2", 1)
	n3 := CreateTestNode("3", 1)
	n4 := CreateTestNode("4", 1)

	r.Connect(0, n1)

	n1.SetExit(0, n2)
	n1.SetExit(1, n3)

	n2.SetExit(0, n4)

	n3.SetExit(0, n4)

	r.AttachDefaultExits(n4)

	ordered := ReversePostorder(r)
	checkOrder(ordered, []*Node{r.Entry, n1, n2, n3, n4, r.GetExit(0)}, t)

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
	r := CreateTestRegion()
	n1 := CreateTestNode("1", 2)
	n2 := CreateTestNode("2", 2)
	n3 := CreateTestNode("3", 1)
	n4 := CreateTestNode("4", 1)
	n5 := CreateTestNode("5", 1)
	n6 := CreateTestNode("6", 1)
	n7 := CreateTestNode("7", 1)

	r.Connect(0, n1)

	n1.SetExit(0, n2)
	n1.SetExit(1, n6)

	n2.SetExit(0, n3)
	n2.SetExit(1, n4)

	n3.SetExit(0, n5)
	n4.SetExit(0, n5)
	n5.SetExit(0, n7)
	n6.SetExit(0, n7)
	r.AttachDefaultExits(n7)

	builder := CreateSSIBuilder(r, ReversePostorder(r), &SimpleLivenessOracle{})

	checkOrder(builder.nodes, []*Node{r.Entry, n1, n2, n3, n4, n5, n6, n7, r.GetExit(0)}, t)

	checkIntList(builder.idoms, []int{0, 0, 1, 2, 2, 2, 1, 1, 7}, t)

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

package main

import (
	"evergreen/dub"
	"fmt"
)

type DubEntry struct {
}

func (n *DubEntry) NumExits() int {
	return 1
}

type DubExit struct {
	flow int
}

func (n *DubExit) NumExits() int {
	return 0
}

type DubNode struct {
	name string
}

func (n *DubNode) NumExits() int {
	return 2
}

func CreateDubEntry() *dub.Node {
	return dub.CreateNode(&DubEntry{})
}

func CreateDubNode(name string, numExits int) *dub.Node {
	return dub.CreateNode(&DubNode{name: name})
}

func CreateDubExit(flow int) *dub.Node {
	return dub.CreateNode(&DubExit{flow: flow})
}

func CreateDubRegion() *dub.Region {
	return dub.CreateRegion(
		CreateDubEntry(),
		[]*dub.Node{
			CreateDubExit(0),
			CreateDubExit(1),
		},
	)
}

func main() {

	l := CreateDubRegion()
	cond := CreateDubNode("cond", 2)
	decide := CreateDubNode("decide", 2)
	body := CreateDubNode("body", 2)

	l.Connect(0, cond)
	l.AttachDefaultExits(cond)

	l.Connect(0, decide)
	decide.SetExit(0, body)

	l.AttachDefaultExits(body)
	l.Connect(0, cond)
	decide.SetExit(1, l.GetExit(0))

	nodes := dub.ReversePostorder(l)
	fmt.Println(nodes)

	fmt.Println(dub.GenerateGo())
}

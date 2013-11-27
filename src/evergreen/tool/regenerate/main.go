package main

import (
	"evergreen/dub"
	"evergreen/io"
	"fmt"
	"path/filepath"
)

type DubEntry struct {
}

func (n *DubEntry) NumExits() int {
	return 1
}

func (n *DubEntry) DotNodeStyle() string {
	return `shape=point,label="entry"`
}

func (n *DubEntry) DotEdgeStyle(flow int) string {
	return `color="green"`
}

type DubExit struct {
	flow int
}

func (n *DubExit) NumExits() int {
	return 0
}

func (n *DubExit) DotNodeStyle() string {
	switch n.flow {
	case 0:
		return `shape=invtriangle,label="n"`
	case 1:
		return `shape=invtriangle,label="f"`
	default:
		return `shape=invtriangle,label="?"`
	}
}

func (n *DubExit) DotEdgeStyle(flow int) string {
	panic("Exit has no edges.")
}

type DubBlock struct {
	name string
}

func (n *DubBlock) NumExits() int {
	return 2
}

func (n *DubBlock) DotNodeStyle() string {
	return fmt.Sprintf("shape=box,label=%#v", n.name)
}

func (n *DubBlock) DotEdgeStyle(flow int) string {
	switch flow {
	case 0:
		return `color="green"`
	case 1:
		return `color="goldenrod"`
	default:
		return `label="?"`
	}
}

type DubSwitch struct {
	name string
}

func (n *DubSwitch) NumExits() int {
	return 2
}

func (n *DubSwitch) DotNodeStyle() string {
	return fmt.Sprintf("shape=diamond,label=%#v", n.name)
}

func (n *DubSwitch) DotEdgeStyle(flow int) string {
	switch flow {
	case 0:
		return `color="limegreen"`
	case 1:
		return `color="yellow"`
	default:
		return `label="?"`
	}
}

func CreateDubEntry() *dub.Node {
	return dub.CreateNode(&DubEntry{})
}

func CreateDubBlock(name string, numExits int) *dub.Node {
	return dub.CreateNode(&DubBlock{name: name})
}

func CreateDubSwitch(name string) *dub.Node {
	return dub.CreateNode(&DubSwitch{name: name})
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
	cond := CreateDubBlock("cond", 2)
	decide := CreateDubSwitch("decide")
	body := CreateDubBlock("body", 2)

	l.Connect(0, cond)
	l.AttachDefaultExits(cond)

	l.Connect(0, decide)
	decide.SetExit(0, body)

	l.AttachDefaultExits(body)
	l.Connect(0, cond)
	decide.SetExit(1, l.GetExit(0))

	dot := dub.RegionToDot(l)
	outfile := filepath.Join("output", "test.svg")

	result := make(chan error, 2)
	go func() {
		err := io.WriteDot(dot, outfile)
		result <- err
	}()

	fmt.Println(dub.GenerateGo())

	err := <-result
	if err != nil {
		fmt.Println(err)
	}

}

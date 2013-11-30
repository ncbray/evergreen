package main

import (
	"evergreen/dub"
	"evergreen/io"
	"fmt"
	"path/filepath"
)

type DubRegister uint32

func RegisterName(reg DubRegister) string {
	return fmt.Sprintf("r%d", reg)
}

func AssignOp(op string, dst DubRegister) string {
	return fmt.Sprintf("%s := %s", RegisterName(dst), op)
}

type DubOp interface {
	OpToString() string
}

type GetLocal struct {
	name string
	dst  DubRegister
}

func (n *GetLocal) OpToString() string {
	return AssignOp(n.name, n.dst)
}

type SetLocal struct {
	src  DubRegister
	name string
}

func (n *SetLocal) OpToString() string {
	return fmt.Sprintf("%s := %s", n.name, RegisterName(n.src))
}

type ConstantInt struct {
	value int64
	dst   DubRegister
}

func (n *ConstantInt) OpToString() string {
	return AssignOp(fmt.Sprintf("%d", n.value), n.dst)
}

type BinaryOp struct {
	left  DubRegister
	op    string
	right DubRegister
	dst   DubRegister
}

func (n *BinaryOp) OpToString() string {
	return AssignOp(fmt.Sprintf("%s %s %s", RegisterName(n.left), n.op, RegisterName(n.right)), n.dst)
}

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
	exprs []DubOp
}

func (n *DubBlock) NumExits() int {
	return 2
}

func (n *DubBlock) DotNodeStyle() string {
	s := ""
	for _, expr := range n.exprs {
		s += expr.OpToString() + "\n"
	}
	return fmt.Sprintf("shape=box,label=%#v", s)
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
	expr DubRegister
}

func (n *DubSwitch) NumExits() int {
	return 2
}

func (n *DubSwitch) DotNodeStyle() string {
	return fmt.Sprintf("shape=diamond,label=%#v", RegisterName(n.expr))
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

func CreateDubBlock(exprs []DubOp) *dub.Node {
	return dub.CreateNode(&DubBlock{exprs: exprs})
}

func CreateDubSwitch(expr DubRegister) *dub.Node {
	return dub.CreateNode(&DubSwitch{expr: expr})
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
	cond := CreateDubBlock([]DubOp{
		&GetLocal{name: "counter", dst: 1},
		&GetLocal{name: "limit", dst: 2},
		&BinaryOp{
			left:  1,
			op:    "<",
			right: 2,
			dst:   3,
		},
	})
	decide := CreateDubSwitch(3)
	body := CreateDubBlock([]DubOp{
		&GetLocal{name: "counter", dst: 4},
		&ConstantInt{value: 1, dst: 5},
		&BinaryOp{
			left:  4,
			op:    "+",
			right: 5,
			dst:   6,
		},
		&SetLocal{src: 6, name: "counter"},
	})

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

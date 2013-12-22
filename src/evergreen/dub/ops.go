package dub

import (
	"evergreen/base"
	"fmt"
)

type DubRegister uint32

var NoRegister DubRegister = ^DubRegister(0)

type RegisterInfo struct {
	T string
}

func RegisterName(reg DubRegister) string {
	return fmt.Sprintf("r%d", reg)
}

func formatAssignment(op string, dst DubRegister) string {
	return fmt.Sprintf("%s := %s", RegisterName(dst), op)
}

type DubOp interface {
	OpToString() string
}

type GetLocalOp struct {
	Name string
	Dst  DubRegister
}

func (n *GetLocalOp) OpToString() string {
	return formatAssignment(n.Name, n.Dst)
}

type SetLocalOp struct {
	Src  DubRegister
	Name string
}

func (n *SetLocalOp) OpToString() string {
	return fmt.Sprintf("%s := %s", n.Name, RegisterName(n.Src))
}

type ConstantIntOp struct {
	Value int64
	Dst   DubRegister
}

func (n *ConstantIntOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%d", n.Value), n.Dst)
}

type ConstantRuneOp struct {
	Value rune
	Dst   DubRegister
}

func (n *ConstantRuneOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%v", n.Value), n.Dst)
}

type BinaryOp struct {
	Left  DubRegister
	Op    string
	Right DubRegister
	Dst   DubRegister
}

func (n *BinaryOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%s %s %s", RegisterName(n.Left), n.Op, RegisterName(n.Right)), n.Dst)
}

type Checkpoint struct {
	Dst DubRegister
}

func (n *Checkpoint) OpToString() string {
	return formatAssignment("<checkpoint>", n.Dst)
}

type Recover struct {
	Src DubRegister
}

func (n *Recover) OpToString() string {
	return fmt.Sprintf("<recover> %s", RegisterName(n.Src))
}

type Fail struct {
}

func (n *Fail) OpToString() string {
	return fmt.Sprintf("<fail>")
}

type Read struct {
	Dst DubRegister
}

func (n *Read) OpToString() string {
	return formatAssignment("<read>", n.Dst)
}

// Flow blocks

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
	Flow int
}

func (n *DubExit) NumExits() int {
	return 0
}

func (n *DubExit) DotNodeStyle() string {
	switch n.Flow {
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
	Ops []DubOp
}

func (n *DubBlock) NumExits() int {
	return 2
}

func (n *DubBlock) DotNodeStyle() string {
	s := ""
	for _, op := range n.Ops {
		s += op.OpToString() + "\n"
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
	Cond DubRegister
}

func (n *DubSwitch) NumExits() int {
	return 2
}

func (n *DubSwitch) DotNodeStyle() string {
	return fmt.Sprintf("shape=diamond,label=%#v", RegisterName(n.Cond))
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

func CreateEntry() *base.Node {
	return base.CreateNode(&DubEntry{})
}

func CreateBlock(ops []DubOp) *base.Node {
	return base.CreateNode(&DubBlock{Ops: ops})
}

func CreateSwitch(cond DubRegister) *base.Node {
	return base.CreateNode(&DubSwitch{Cond: cond})
}

func CreateExit(flow int) *base.Node {
	return base.CreateNode(&DubExit{Flow: flow})
}

func CreateRegion() *base.Region {
	r := &base.Region{
		Entry: CreateEntry(),
		Exits: []*base.Node{
			CreateExit(0),
			CreateExit(1),
		},
	}
	r.Entry.SetExit(0, r.Exits[0])
	return r
}

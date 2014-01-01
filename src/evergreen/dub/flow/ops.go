package flow

import (
	"evergreen/base"
	"fmt"
	"strings"
)

const (
	// Real flows, used at runtime
	NORMAL = iota
	FAIL
	EXCEPTION
	// Virtual flows, only for graph construction
	RETURN
)

type DubType interface {
	isDubType()
}

type BoolType struct {
}

func (t *BoolType) isDubType() {
}

type IntType struct {
}

func (t *IntType) isDubType() {
}

type RuneType struct {
}

func (t *RuneType) isDubType() {
}

type StringType struct {
}

func (t *StringType) isDubType() {
}

type ListType struct {
	Type DubType
}

func (t *ListType) isDubType() {
}

type LLField struct {
	Name string
	T    DubType
}

type LLStruct struct {
	Name       string
	Implements *LLStruct
	Abstract   bool
	Fields     []*LLField
}

func (t *LLStruct) isDubType() {
}

type DubRegister uint32

var NoRegister DubRegister = ^DubRegister(0)

type RegisterInfo struct {
	T DubType
}

type LLFunc struct {
	Name        string
	ReturnTypes []DubType
	Registers   []RegisterInfo
	Region      *base.Region
}

func TypeName(t DubType) string {
	switch t := t.(type) {
	case *LLStruct:
		return t.Name
	case *ListType:
		return fmt.Sprintf("[]%s", TypeName(t.Type))
	case *StringType:
		return "string"
	case *RuneType:
		return "rune"
	case *IntType:
		return "int"
	case *BoolType:
		return "bool"
	default:
		panic(t)
	}
}

func RegisterName(reg DubRegister) string {
	return fmt.Sprintf("r%d", reg)
}

func RegisterList(regs []DubRegister) string {
	names := make([]string, len(regs))
	for i, reg := range regs {
		names[i] = RegisterName(reg)
	}
	return strings.Join(names, ", ")
}

func KeyValueList(args []*KeyValue) string {
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = fmt.Sprintf("%s: %s", arg.Key, RegisterName(arg.Value))
	}
	return strings.Join(names, ", ")
}

func formatAssignment(op string, dst DubRegister) string {
	if dst == NoRegister {
		return op
	}
	return fmt.Sprintf("%s := %s", RegisterName(dst), op)
}

type DubOp interface {
	OpToString() string
}

type CoerceOp struct {
	Src DubRegister
	T   DubType
	Dst DubRegister
}

func (n *CoerceOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%s(%s)", TypeName(n.T), RegisterName(n.Src)), n.Dst)
}

type CopyOp struct {
	Src DubRegister
	Dst DubRegister
}

func (n *CopyOp) OpToString() string {
	return fmt.Sprintf("%s := %s", RegisterName(n.Dst), RegisterName(n.Src))
}

type ConstantNilOp struct {
	Dst DubRegister
}

func (n *ConstantNilOp) OpToString() string {
	return formatAssignment("nil", n.Dst)
}

type ConstantIntOp struct {
	Value int64
	Dst   DubRegister
}

func (n *ConstantIntOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%v", n.Value), n.Dst)
}

type ConstantBoolOp struct {
	Value bool
	Dst   DubRegister
}

func (n *ConstantBoolOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%v", n.Value), n.Dst)
}

type ConstantRuneOp struct {
	Value rune
	Dst   DubRegister
}

func (n *ConstantRuneOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%#U", n.Value), n.Dst)
}

type ConstantStringOp struct {
	Value string
	Dst   DubRegister
}

func (n *ConstantStringOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%#v", n.Value), n.Dst)
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

type CallOp struct {
	Name string
	Dst  DubRegister
}

func (n *CallOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%s()", n.Name), n.Dst)
}

type KeyValue struct {
	Key   string
	Value DubRegister
}

type ConstructOp struct {
	Type *LLStruct
	Args []*KeyValue
	Dst  DubRegister
}

func (n *ConstructOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%s{%s}", TypeName(n.Type), KeyValueList(n.Args)), n.Dst)
}

type ConstructListOp struct {
	Type *ListType
	Args []DubRegister
	Dst  DubRegister
}

func (n *ConstructListOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("%s{%s}", TypeName(n.Type), RegisterList(n.Args)), n.Dst)
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

type LookaheadBegin struct {
	Dst DubRegister
}

func (n *LookaheadBegin) OpToString() string {
	return formatAssignment("<lookahead begin>", n.Dst)
}

type LookaheadEnd struct {
	Failed bool
	Src    DubRegister
}

func (n *LookaheadEnd) OpToString() string {
	return fmt.Sprintf("<lookahead end> %v %s", n.Failed, RegisterName(n.Src))
}

type Slice struct {
	Src DubRegister
	Dst DubRegister
}

func (n *Slice) OpToString() string {
	return formatAssignment(fmt.Sprintf("<slice> %s", RegisterName(n.Src)), n.Dst)
}

type AppendOp struct {
	List  DubRegister
	Value DubRegister
	Dst   DubRegister
}

func (n *AppendOp) OpToString() string {
	return formatAssignment(fmt.Sprintf("<append> %s %s", RegisterName(n.List), RegisterName(n.Value)), n.Dst)
}

type ReturnOp struct {
	Exprs []DubRegister
}

func (n *ReturnOp) OpToString() string {
	return fmt.Sprintf("<return> %s", RegisterList(n.Exprs))
}

type Fail struct {
}

func (n *Fail) OpToString() string {
	return fmt.Sprintf("<fail>")
}

type Peek struct {
	Dst DubRegister
}

func (n *Peek) OpToString() string {
	return formatAssignment("<peek>", n.Dst)
}

type Consume struct {
}

func (n *Consume) OpToString() string {
	return "<consume>"
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
	case 2:
		return `shape=invtriangle,label="e"`
	case 3:
		return `shape=invtriangle,label="r"`
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
			CreateExit(2),
			CreateExit(3),
		},
	}
	r.Entry.SetExit(0, r.Exits[0])
	return r
}
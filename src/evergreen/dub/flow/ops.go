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

type IntrinsicType struct {
	Name string
}

func (t *IntrinsicType) isDubType() {
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
	Params      []DubRegister
	ReturnTypes []DubType
	Registers   []RegisterInfo
	CFG         *base.Graph
	Ops         []DubOp
}

func TypeName(t DubType) string {
	switch t := t.(type) {
	case *LLStruct:
		return t.Name
	case *ListType:
		return fmt.Sprintf("[]%s", TypeName(t.Type))
	case *IntrinsicType:
		return t.Name
	default:
		panic(t)
	}
}

func RegisterName(reg DubRegister) string {
	if reg != NoRegister {
		return fmt.Sprintf("r%d", reg)
	} else {
		return "_"
	}
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

func formatMultiAssignment(op string, dsts []DubRegister) string {
	if len(dsts) > 0 {
		return fmt.Sprintf("%s := %s", RegisterList(dsts), op)
	} else {
		return op
	}
}

type DubOp interface {
	isDubOp()
}

type CoerceOp struct {
	Src DubRegister
	T   DubType
	Dst DubRegister
}

func (n *CoerceOp) isDubOp() {
}

type CopyOp struct {
	Src DubRegister
	Dst DubRegister
}

func (n *CopyOp) isDubOp() {
}

type ConstantNilOp struct {
	Dst DubRegister
}

func (n *ConstantNilOp) isDubOp() {
}

type ConstantIntOp struct {
	Value int64
	Dst   DubRegister
}

func (n *ConstantIntOp) isDubOp() {
}

type ConstantBoolOp struct {
	Value bool
	Dst   DubRegister
}

func (n *ConstantBoolOp) isDubOp() {
}

type ConstantRuneOp struct {
	Value rune
	Dst   DubRegister
}

func (n *ConstantRuneOp) isDubOp() {
}

type ConstantStringOp struct {
	Value string
	Dst   DubRegister
}

func (n *ConstantStringOp) isDubOp() {
}

type BinaryOp struct {
	Left  DubRegister
	Op    string
	Right DubRegister
	Dst   DubRegister
}

func (n *BinaryOp) isDubOp() {
}

type CallOp struct {
	Name   string
	Target *LLFunc
	Args   []DubRegister
	Dsts   []DubRegister
}

func (n *CallOp) isDubOp() {
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

func (n *ConstructOp) isDubOp() {
}

type ConstructListOp struct {
	Type *ListType
	Args []DubRegister
	Dst  DubRegister
}

func (n *ConstructListOp) isDubOp() {
}

type Checkpoint struct {
	Dst DubRegister
}

func (n *Checkpoint) isDubOp() {
}

type Recover struct {
	Src DubRegister
}

func (n *Recover) isDubOp() {
}

type LookaheadBegin struct {
	Dst DubRegister
}

func (n *LookaheadBegin) isDubOp() {
}

type LookaheadEnd struct {
	Failed bool
	Src    DubRegister
}

func (n *LookaheadEnd) isDubOp() {
}

type Slice struct {
	Src DubRegister
	Dst DubRegister
}

func (n *Slice) isDubOp() {
}

type AppendOp struct {
	List  DubRegister
	Value DubRegister
	Dst   DubRegister
}

func (n *AppendOp) isDubOp() {
}

type ReturnOp struct {
	Exprs []DubRegister
}

func (n *ReturnOp) isDubOp() {
}

type Fail struct {
}

func (n *Fail) isDubOp() {
}

type Peek struct {
	Dst DubRegister
}

func (n *Peek) isDubOp() {
}

type Consume struct {
}

func (n *Consume) isDubOp() {
}

type TransferOp struct {
	Srcs []DubRegister
	Dsts []DubRegister
}

func (n *TransferOp) isDubOp() {
}

// Flow blocks

type EntryOp struct {
}

func (n *EntryOp) isDubOp() {
}

type SwitchOp struct {
	Cond DubRegister
}

func (n *SwitchOp) isDubOp() {
}

type FlowExitOp struct {
	Flow int
}

func (n *FlowExitOp) isDubOp() {
}

type ExitOp struct {
}

func (n *ExitOp) isDubOp() {
}

type DotStyler struct {
	Decl *LLFunc
}

func opToString(op DubOp) string {
	switch n := op.(type) {
	case *CoerceOp:
		return formatAssignment(fmt.Sprintf("%s(%s)", TypeName(n.T), RegisterName(n.Src)), n.Dst)
	case *CopyOp:
		return fmt.Sprintf("%s := %s", RegisterName(n.Dst), RegisterName(n.Src))
	case *ConstantNilOp:
		return formatAssignment("nil", n.Dst)
	case *ConstantIntOp:
		return formatAssignment(fmt.Sprintf("%v", n.Value), n.Dst)
	case *ConstantBoolOp:
		return formatAssignment(fmt.Sprintf("%v", n.Value), n.Dst)
	case *ConstantRuneOp:
		return formatAssignment(fmt.Sprintf("%#U", n.Value), n.Dst)
	case *ConstantStringOp:
		return formatAssignment(fmt.Sprintf("%#v", n.Value), n.Dst)
	case *BinaryOp:
		return formatAssignment(fmt.Sprintf("%s %s %s", RegisterName(n.Left), n.Op, RegisterName(n.Right)), n.Dst)
	case *CallOp:
		return formatMultiAssignment(fmt.Sprintf("%s(%s)", n.Name, RegisterList(n.Args)), n.Dsts)
	case *ConstructOp:
		return formatAssignment(fmt.Sprintf("%s{%s}", TypeName(n.Type), KeyValueList(n.Args)), n.Dst)
	case *ConstructListOp:
		return formatAssignment(fmt.Sprintf("%s{%s}", TypeName(n.Type), RegisterList(n.Args)), n.Dst)
	case *Checkpoint:
		return formatAssignment("<checkpoint>", n.Dst)
	case *Recover:
		return fmt.Sprintf("<recover> %s", RegisterName(n.Src))
	case *LookaheadBegin:
		return formatAssignment("<lookahead begin>", n.Dst)
	case *LookaheadEnd:
		return fmt.Sprintf("<lookahead end> %v %s", n.Failed, RegisterName(n.Src))
	case *Slice:
		return formatAssignment(fmt.Sprintf("<slice> %s", RegisterName(n.Src)), n.Dst)
	case *AppendOp:
		return formatAssignment(fmt.Sprintf("<append> %s %s", RegisterName(n.List), RegisterName(n.Value)), n.Dst)
	case *ReturnOp:
		return fmt.Sprintf("<return> %s", RegisterList(n.Exprs))
	case *Fail:
		return "<fail>"
	case *Peek:
		return formatAssignment("<peek>", n.Dst)
	case *Consume:
		return "<consume>"
	case *TransferOp:
		return fmt.Sprintf("%s << %s", RegisterList(n.Dsts), RegisterList(n.Srcs))

	default:
		panic(op)
	}
}

func (styler *DotStyler) NodeStyle(node base.NodeID) string {
	op := styler.Decl.Ops[node]
	switch op := op.(type) {
	case *EntryOp:
		return `shape=point,label="entry"`
	case *ExitOp:
		return `shape=point,label="exit"`
	case *FlowExitOp:
		switch op.Flow {
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
	case *SwitchOp:
		return fmt.Sprintf("shape=diamond,label=%#v", RegisterName(op.Cond))
	case DubOp:
		return fmt.Sprintf("shape=box,label=%#v", opToString(op))
	default:
		panic(op)
	}
}

func (styler *DotStyler) EdgeStyle(node base.NodeID, flow int) string {
	op := styler.Decl.Ops[node]
	color := "red"
	switch op.(type) {
	case *SwitchOp:
		switch flow {
		case 0:
			color = "limegreen"
		case 1:
			color = "yellow"
		}
	case *FlowExitOp:
		color = "gray"
	default:
		switch flow {
		case 0:
			color = "green"
		case 1:
			color = "goldenrod"
		}
	}
	return fmt.Sprintf("color=%s", color)
}

func IsNop(op DubOp) bool {
	switch op := op.(type) {
	case *Consume:
		return false
	case *Fail:
		return false
	case *Checkpoint:
		return op.Dst == NoRegister
	case *Peek:
		return false
	case *LookaheadBegin:
		return false
	case *ConstantRuneOp:
		return op.Dst == NoRegister
	case *ConstantStringOp:
		return op.Dst == NoRegister
	case *ConstantIntOp:
		return op.Dst == NoRegister
	case *ConstantBoolOp:
		return op.Dst == NoRegister
	case *ConstantNilOp:
		return op.Dst == NoRegister
	case *CallOp:
		return false
	case *Slice:
		return op.Dst == NoRegister
	case *BinaryOp:
		return op.Dst == NoRegister
	case *AppendOp:
		return op.Dst == NoRegister
	case *CopyOp:
		return op.Dst == NoRegister || op.Dst == op.Src
	case *CoerceOp:
		return op.Dst == NoRegister
	case *Recover:
		return false
	case *LookaheadEnd:
		return false
	case *ReturnOp:
		return false
	case *ConstructOp:
		return op.Dst == NoRegister
	case *ConstructListOp:
		return op.Dst == NoRegister
	case *TransferOp:
		return len(op.Dsts) == 0
	case *EntryOp:
		return false
	case *SwitchOp:
		return false
	case *FlowExitOp:
		return false
	case *ExitOp:
		return false
	default:
		panic(op)
	}
}

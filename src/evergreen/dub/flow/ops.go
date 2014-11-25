package flow

import (
	"evergreen/base"
	core "evergreen/dub/tree"
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

func (scope *RegisterInfo_Scope) Get(ref RegisterInfo_Ref) *RegisterInfo {
	return scope.objects[ref]
}

func (scope *RegisterInfo_Scope) Register(info *RegisterInfo) RegisterInfo_Ref {
	index := RegisterInfo_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return index
}

func (scope *RegisterInfo_Scope) Len() int {
	return len(scope.objects)
}

func (scope *RegisterInfo_Scope) Remap(remap []int, count int) {
	objects := make([]*RegisterInfo, count)
	for i, info := range scope.objects {
		idx := remap[i]
		if idx >= 0 {
			objects[idx] = info
		}
	}
	scope.objects = objects
}

func (scope *RegisterInfo_Scope) Replace(replacement []*RegisterInfo) {
	scope.objects = replacement
}

func TypeName(t core.DubType) string {
	switch t := t.(type) {
	case *core.StructType:
		return t.Name.Text
	case *core.ListType:
		return fmt.Sprintf("[]%s", TypeName(t.Type))
	case *core.BuiltinType:
		return t.Name
	default:
		panic(t)
	}
}

func RegisterName(reg RegisterInfo_Ref) string {
	if reg != NoRegisterInfo {
		return fmt.Sprintf("r%d", reg)
	} else {
		return "_"
	}
}

func RegisterList(regs []RegisterInfo_Ref) string {
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

func formatAssignment(op string, dst RegisterInfo_Ref) string {
	if dst == NoRegisterInfo {
		return op
	}
	return fmt.Sprintf("%s := %s", RegisterName(dst), op)
}

func formatMultiAssignment(op string, dsts []RegisterInfo_Ref) string {
	if len(dsts) > 0 {
		return fmt.Sprintf("%s := %s", RegisterList(dsts), op)
	} else {
		return op
	}
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
		return op.Dst == NoRegisterInfo
	case *Peek:
		return false
	case *LookaheadBegin:
		return false
	case *ConstantRuneOp:
		return op.Dst == NoRegisterInfo
	case *ConstantStringOp:
		return op.Dst == NoRegisterInfo
	case *ConstantIntOp:
		return op.Dst == NoRegisterInfo
	case *ConstantBoolOp:
		return op.Dst == NoRegisterInfo
	case *ConstantNilOp:
		return op.Dst == NoRegisterInfo
	case *CallOp:
		return false
	case *Slice:
		return op.Dst == NoRegisterInfo
	case *BinaryOp:
		return op.Dst == NoRegisterInfo
	case *AppendOp:
		return op.Dst == NoRegisterInfo
	case *CopyOp:
		return op.Dst == NoRegisterInfo || op.Dst == op.Src
	case *CoerceOp:
		return op.Dst == NoRegisterInfo
	case *Recover:
		return false
	case *LookaheadEnd:
		return false
	case *ReturnOp:
		return false
	case *ConstructOp:
		return op.Dst == NoRegisterInfo
	case *ConstructListOp:
		return op.Dst == NoRegisterInfo
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

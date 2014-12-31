package flow

import (
	"evergreen/dub/core"
	"evergreen/graph"
	"fmt"
	"strings"
)

func registerName(reg RegisterInfo_Ref) string {
	if reg != NoRegisterInfo {
		return fmt.Sprintf("r%d", reg)
	} else {
		return "_"
	}
}

func registerList(regs []RegisterInfo_Ref) string {
	names := make([]string, len(regs))
	for i, reg := range regs {
		names[i] = registerName(reg)
	}
	return strings.Join(names, ", ")
}

func keyValueList(args []*KeyValue) string {
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = fmt.Sprintf("%s: %s", arg.Key, registerName(arg.Value))
	}
	return strings.Join(names, ", ")
}

func formatAssignment(op string, dst RegisterInfo_Ref) string {
	if dst == NoRegisterInfo {
		return op
	}
	return fmt.Sprintf("%s := %s", registerName(dst), op)
}

func formatMultiAssignment(op string, dsts []RegisterInfo_Ref) string {
	if len(dsts) > 0 {
		return fmt.Sprintf("%s := %s", registerList(dsts), op)
	} else {
		return op
	}
}

type DotStyler struct {
	Decl *LLFunc
	Core *core.CoreProgram
}

func opToString(coreProg *core.CoreProgram, op DubOp) string {
	switch n := op.(type) {
	case *CoerceOp:
		return formatAssignment(fmt.Sprintf("%s(%s)", core.TypeName(n.T), registerName(n.Src)), n.Dst)
	case *CopyOp:
		return fmt.Sprintf("%s := %s", registerName(n.Dst), registerName(n.Src))
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
		return formatAssignment(fmt.Sprintf("%s %s %s", registerName(n.Left), n.Op, registerName(n.Right)), n.Dst)
	case *CallOp:
		name := coreProg.Function_Scope.Get(n.Target).Name
		return formatMultiAssignment(fmt.Sprintf("%s(%s)", name, registerList(n.Args)), n.Dsts)
	case *ConstructOp:
		return formatAssignment(fmt.Sprintf("%s{%s}", core.TypeName(n.Type), keyValueList(n.Args)), n.Dst)
	case *ConstructListOp:
		return formatAssignment(fmt.Sprintf("%s{%s}", core.TypeName(n.Type), registerList(n.Args)), n.Dst)
	case *Checkpoint:
		return formatAssignment("<checkpoint>", n.Dst)
	case *Recover:
		return fmt.Sprintf("<recover> %s", registerName(n.Src))
	case *LookaheadBegin:
		return formatAssignment("<lookahead begin>", n.Dst)
	case *LookaheadEnd:
		return fmt.Sprintf("<lookahead end> %v %s", n.Failed, registerName(n.Src))
	case *Slice:
		return formatAssignment(fmt.Sprintf("<slice> %s", registerName(n.Src)), n.Dst)
	case *AppendOp:
		return formatAssignment(fmt.Sprintf("<append> %s %s", registerName(n.List), registerName(n.Value)), n.Dst)
	case *ReturnOp:
		return fmt.Sprintf("<return> %s", registerList(n.Exprs))
	case *Fail:
		return "<fail>"
	case *Peek:
		return formatAssignment("<peek>", n.Dst)
	case *Consume:
		return "<consume>"
	case *TransferOp:
		return fmt.Sprintf("%s << %s", registerList(n.Dsts), registerList(n.Srcs))

	default:
		panic(op)
	}
}

func nodeLabel(node graph.NodeID, label string) string {
	return fmt.Sprintf("[%d] %s", node, label)
}

func flowExit(label string) string {
	return fmt.Sprintf("shape=invtriangle,label=%#v", label)
}

func (styler *DotStyler) NodeStyle(node graph.NodeID) string {
	op := styler.Decl.Ops[node]
	switch op := op.(type) {
	case *EntryOp:
		return `shape=point,label="entry"`
	case *ExitOp:
		return `shape=point,label="exit"`
	case *FlowExitOp:
		var label string
		switch op.Flow {
		case 0:
			label = "n"
		case 1:
			label = "f"
		case 2:
			label = "e"
		case 3:
			label = "r"
		default:
			label = "?"
		}
		return flowExit(nodeLabel(node, label))
	case *SwitchOp:
		return fmt.Sprintf("shape=diamond,label=%#v", nodeLabel(node, "?"+registerName(op.Cond)))
	case DubOp:
		return fmt.Sprintf("shape=box,label=%#v", nodeLabel(node, opToString(styler.Core, op)))
	default:
		panic(op)
	}
}

func (styler *DotStyler) EdgeStyle(node graph.NodeID, flow int) string {
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
package flow

import (
	"evergreen/dub/core"
	"evergreen/graph"
	"fmt"
	"strings"
)

func dotString(message string) string {
	return fmt.Sprintf("\"%s\"", graph.EscapeDotString(message))
}

func registerName(reg *RegisterInfo) string {
	if reg != nil {
		return fmt.Sprintf("r%d", reg.Index)
	} else {
		return "_"
	}
}

func registerList(regs []*RegisterInfo) string {
	names := make([]string, len(regs))
	for i, reg := range regs {
		names[i] = registerName(reg)
	}
	return strings.Join(names, ", ")
}

func keyValueList(args []*KeyValue) string {
	if len(args) == 0 {
		return ""
	}
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = fmt.Sprintf("        %s: %s,\n", arg.Key, registerName(arg.Value))
	}
	return "\n" + strings.Join(names, "")
}

func formatAssignment(op string, dst *RegisterInfo) string {
	if dst == nil {
		return op
	}
	return fmt.Sprintf("%s := %s", registerName(dst), op)
}

func formatMultiAssignment(op string, dsts []*RegisterInfo) string {
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

func callableName(coreProg *core.CoreProgram, c core.Callable) string {
	switch c := c.(type) {
	case *core.Function:
		return c.Name
	case *core.IntrinsicFunction:
		return "!" + c.Name
	default:
		panic(c)
	}
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
	case *ConstantFloat32Op:
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
		name := callableName(coreProg, n.Target)
		return formatMultiAssignment(fmt.Sprintf("%s(%s)", name, registerList(n.Args)), n.Dsts)
	case *ConstructOp:
		return formatAssignment(fmt.Sprintf("%s{%s}\n", core.TypeName(n.Type), keyValueList(n.Args)), n.Dst)
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

func (styler *DotStyler) NodeStyle(node graph.NodeID) string {
	op := styler.Decl.Ops[node]
	switch op := op.(type) {
	case *EntryOp:
		return `shape=point,label="entry"`
	case *ExitOp:
		return `shape=point,label="exit"`
	case *SwitchOp:
		return fmt.Sprintf("shape=diamond,label=%s", dotString(nodeLabel(node, "?"+registerName(op.Cond))))
	case DubOp:
		return fmt.Sprintf("shape=box,label=%s", dotString(nodeLabel(node, opToString(styler.Core, op))))
	default:
		panic(op)
	}
}

func (styler *DotStyler) BlockLabel(node graph.NodeID) (string, bool) {
	op := styler.Decl.Ops[node]
	switch op := op.(type) {
	case *EntryOp, *ExitOp, *SwitchOp:
		return "", false
	case DubOp:
		return nodeLabel(node, opToString(styler.Core, op)), true
	default:
		panic(op)
	}
}

func (styler *DotStyler) EdgeStyle(src graph.NodeID, e graph.EdgeID, dst graph.NodeID) string {
	flow := styler.Decl.Edges[e]
	color := "red"
	switch flow {
	case NORMAL:
		color = graph.NORMAL_EDGE_COLOR
	case COND_TRUE:
		color = graph.TRUE_EDGE_COLOR
	case COND_FALSE:
		color = graph.FALSE_EDGE_COLOR
	case FAIL:
		color = graph.FAIL_EDGE_COLOR
	case RETURN:
		color = graph.RETURN_EDGE_COLOR
	}
	return fmt.Sprintf("color=%s", color)
}

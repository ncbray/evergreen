package flow

import (
	"evergreen/go/core"
	"evergreen/graph"
	"fmt"
	"strings"
)

func RegisterName(reg *Register) string {
	if reg != nil {
		return fmt.Sprintf("r%d", reg.Index)
	} else {
		return "_"
	}
}

func addDst(op string, dst *Register) string {
	if dst == nil {
		return op
	}
	return fmt.Sprintf("%s := %s", RegisterName(dst), op)
}

func addDsts(op string, dsts []*Register) string {
	if len(dsts) == 0 {
		return op
	}
	return fmt.Sprintf("%s := %s", registerList(dsts), op)
}

func registerList(args []*Register) string {
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = RegisterName(arg)
	}
	return strings.Join(names, ", ")
}

func namedArgList(args []*NamedArg) string {
	if len(args) == 0 {
		return ""
	}
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = fmt.Sprintf("        %s: %s,\n", arg.Name, RegisterName(arg.Arg))
	}
	return "\n" + strings.Join(names, "")
}

func typeName(t core.GoType) string {
	switch t := t.(type) {
	case *core.StructType:
		return t.Name
	case *core.InterfaceType:
		return t.Name
	case *core.SliceType:
		return "[]" + typeName(t.Element)
	case *core.PointerType:
		return "*" + typeName(t.Element)
	case *core.ExternalType:
		return t.Name
	default:
		panic(t)
	}
}

func callableName(c core.Callable) string {
	switch c := c.(type) {
	case *core.Function:
		return c.Name
	case *core.IntrinsicFunction:
		return "!" + c.Name
	default:
		panic(c)
	}
}

func opToString(coreProg *core.CoreProgram, op GoOp) string {
	switch op := op.(type) {
	case *Entry:
		return "entry"
	case *Exit:
		return "exit"
	case *Switch:
		return fmt.Sprintf("switch %s", RegisterName(op.Cond))
	case *Transfer:
		return fmt.Sprintf("%s << %s", registerList(op.Dsts), registerList(op.Srcs))
	case *ConstantNil:
		return addDst("nil", op.Dst)
	case *ConstantInt:
		return addDst(fmt.Sprintf("%v", op.Value), op.Dst)
	case *ConstantFloat32:
		return addDst(fmt.Sprintf("%v", op.Value), op.Dst)
	case *ConstantBool:
		return addDst(fmt.Sprintf("%v", op.Value), op.Dst)
	case *ConstantRune:
		return addDst(fmt.Sprintf("%#U", op.Value), op.Dst)
	case *ConstantString:
		return addDst(fmt.Sprintf("%#v", op.Value), op.Dst)
	case *Attr:
		return addDst(fmt.Sprintf("%s.%s", RegisterName(op.Expr), op.Name), op.Dst)
	case *BinaryOp:
		return addDst(fmt.Sprintf("%s %s %s", RegisterName(op.Left), op.Op, RegisterName(op.Right)), op.Dst)
	case *Call:
		return addDsts(fmt.Sprintf("%s(%s)", callableName(op.Target), registerList(op.Args)), op.Dsts)
	case *MethodCall:
		return addDsts(fmt.Sprintf("%s.%s(%s)", RegisterName(op.Expr), op.Name, registerList(op.Args)), op.Dsts)
	case *ConstructStruct:
		prefix := ""
		if op.AddrTaken {
			prefix = "&"
		}
		return addDst(fmt.Sprintf("%s%s{%s}\n", prefix, typeName(op.Type), namedArgList(op.Args)), op.Dst)
	case *ConstructSlice:
		return addDst(fmt.Sprintf("%s{%s}", typeName(op.Type), registerList(op.Args)), op.Dst)
	case *Coerce:
		return addDst(fmt.Sprintf("%s(%s)", typeName(op.Type), RegisterName(op.Src)), op.Dst)
	case *Return:
		return fmt.Sprintf("return %s", registerList(op.Args))
	case *Nop:
		return "nop"
	default:
		panic(op)
	}
}

type DotStyler struct {
	Func *FlowFunc
	Core *core.CoreProgram
}

func nodeLabel(node graph.NodeID, label string) string {
	return fmt.Sprintf("[%d] %s", node, label)
}

func dotString(message string) string {
	return fmt.Sprintf("\"%s\"", graph.EscapeDotString(message))
}

func (styler *DotStyler) NodeStyle(node graph.NodeID) string {
	op := styler.Func.Ops[node]
	switch op := op.(type) {
	case *Entry:
		return `shape=point,label="entry"`
	case *Exit:
		return `shape=point,label="exit"`
	case *Switch:
		return fmt.Sprintf("shape=diamond,label=%s", dotString(nodeLabel(node, "?"+RegisterName(op.Cond))))
	case GoOp:
		return fmt.Sprintf("shape=box,label=%s", dotString(nodeLabel(node, opToString(styler.Core, op))))
	default:
		panic(op)
	}
}

func (styler *DotStyler) EdgeStyle(src graph.NodeID, e graph.EdgeID, dst graph.NodeID) string {
	flow := styler.Func.Edges[e]
	color := "red"
	switch flow {
	case NORMAL:
		color = graph.NORMAL_EDGE_COLOR
	case COND_TRUE:
		color = graph.TRUE_EDGE_COLOR
	case COND_FALSE:
		color = graph.FALSE_EDGE_COLOR
	case RETURN:
		color = graph.RETURN_EDGE_COLOR
	}
	return fmt.Sprintf("color=%s", color)
}

func (styler *DotStyler) BlockLabel(node graph.NodeID) (string, bool) {
	op := styler.Func.Ops[node]
	switch op := op.(type) {
	case *Entry, *Exit, *Switch:
		return "", false
	case GoOp:
		return nodeLabel(node, opToString(styler.Core, op)), true
	default:
		panic(op)
	}
}

package flow

import (
	"evergreen/base"
	core "evergreen/go/tree"
	"fmt"
	"strings"
)

func RegisterName(reg Register_Ref) string {
	if reg != NoRegister {
		return fmt.Sprintf("r%d", reg)
	} else {
		return "_"
	}
}

func addDst(op string, dst Register_Ref) string {
	if dst == NoRegister {
		return op
	}
	return fmt.Sprintf("%s := %s", RegisterName(dst), op)
}

func addDsts(op string, dsts []Register_Ref) string {
	if len(dsts) == 0 {
		return op
	}
	return fmt.Sprintf("%s := %s", registerList(dsts), op)
}

func registerList(args []Register_Ref) string {
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = RegisterName(arg)
	}
	return strings.Join(names, ", ")
}

func namedArgList(args []*NamedArg) string {
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = fmt.Sprintf("%s: %s", arg.Name, RegisterName(arg.Arg))
	}
	return strings.Join(names, ", ")
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

func OpToString(op GoOp) string {
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
		return addDsts(fmt.Sprintf("%s(%s)", op.Name, registerList(op.Args)), op.Dsts)
	case *MethodCall:
		return addDsts(fmt.Sprintf("%s.%s(%s)", RegisterName(op.Expr), op.Name, registerList(op.Args)), op.Dsts)
	case *ConstructStruct:
		prefix := ""
		if op.AddrTaken {
			prefix = "&"
		}
		return addDst(fmt.Sprintf("%s%s{%s}", prefix, typeName(op.Type), namedArgList(op.Args)), op.Dst)
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
	Ops []GoOp
}

func (styler *DotStyler) NodeStyle(node base.NodeID) string {
	op := styler.Ops[node]
	switch op := op.(type) {
	case *Entry:
		return `shape=point,label="entry"`
	case *Exit:
		return `shape=point,label="exit"`
	case *Switch:
		return fmt.Sprintf("shape=diamond,label=%#v", RegisterName(op.Cond))
	case GoOp:
		return fmt.Sprintf("shape=box,label=%#v", OpToString(op))
	default:
		panic(op)
	}
}

func (styler *DotStyler) EdgeStyle(node base.NodeID, flow int) string {
	op := styler.Ops[node]
	color := "red"
	switch op.(type) {
	case *Switch:
		switch flow {
		case 0:
			color = "limegreen"
		case 1:
			color = "yellow"
		}
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
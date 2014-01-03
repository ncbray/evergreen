package flow

import (
	"evergreen/base"
	"fmt"
)

func AddDef(reg DubRegister, node int, defuse *base.DefUseCollector) {
	if reg != NoRegister {
		defuse.AddDef(node, int(reg))
	}
}

func AddUse(reg DubRegister, node int, defuse *base.DefUseCollector) {
	defuse.AddUse(node, int(reg))
}

func collectDefUse(node *base.Node, defuse *base.DefUseCollector) {
	switch op := node.Data.(type) {
	case *DubEntry, *DubExit:
	case *Consume, *Fail:
	case *Checkpoint:
		AddDef(op.Dst, node.Name, defuse)
	case *Peek:
		AddDef(op.Dst, node.Name, defuse)
	case *LookaheadBegin:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantRuneOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantStringOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantIntOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantBoolOp:
		AddDef(op.Dst, node.Name, defuse)
	case *ConstantNilOp:
		AddDef(op.Dst, node.Name, defuse)
	case *CallOp:
		AddDef(op.Dst, node.Name, defuse)
	case *Slice:
		AddDef(op.Dst, node.Name, defuse)
	case *BinaryOp:
		AddUse(op.Left, node.Name, defuse)
		AddUse(op.Right, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *AppendOp:
		AddUse(op.List, node.Name, defuse)
		AddUse(op.Value, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *CopyOp:
		AddUse(op.Src, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *CoerceOp:
		AddUse(op.Src, node.Name, defuse)
		AddDef(op.Dst, node.Name, defuse)
	case *Recover:
		AddUse(op.Src, node.Name, defuse)
	case *LookaheadEnd:
		AddUse(op.Src, node.Name, defuse)
	case *DubSwitch:
		AddUse(op.Cond, node.Name, defuse)
	case *ReturnOp:
		for _, arg := range op.Exprs {
			AddUse(arg, node.Name, defuse)
		}
	case *ConstructOp:
		for _, arg := range op.Args {
			AddUse(arg.Value, node.Name, defuse)
		}
		AddDef(op.Dst, node.Name, defuse)
	case *ConstructListOp:
		for _, arg := range op.Args {
			AddUse(arg, node.Name, defuse)
		}
		AddDef(op.Dst, node.Name, defuse)
	default:
		panic(op)
	}
}

func SSI(decl *LLFunc) {
	order := base.ReversePostorder(decl.Region)
	defuse := base.FindDefUse(order, len(decl.Registers), collectDefUse)
	live := base.FindLiveVars(order, defuse)

	builder := base.CreateSSIBuilder(decl.Region, order, live)
	for i := 0; i < len(decl.Registers); i++ {
		base.SSI(builder, i, defuse.VarDefAt[i])
	}

	fmt.Println(decl.Name)
	for i := 0; i < len(order); i++ {
		fmt.Println(live.LiveSet(i), builder.PhiFuncs[i])
	}
	fmt.Println()

}

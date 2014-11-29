package flow

import (
	"evergreen/base"
	core "evergreen/go/tree"
)

func (scope *Register_Scope) Get(ref Register_Ref) *Register {
	return scope.objects[ref]
}

func (scope *Register_Scope) Register(info *Register) Register_Ref {
	index := Register_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return index
}

func (scope *Register_Scope) Len() int {
	return len(scope.objects)
}

type GoFlowBuilder struct {
	decl *LLFunc
	CFG  *base.Graph
}

func (builder *GoFlowBuilder) MakeRegister(t core.GoType) Register_Ref {
	return builder.decl.Register_Scope.Register(&Register{T: t})
}

func (builder *GoFlowBuilder) EmitOp(op GoOp, exit_count int) base.NodeID {
	id := builder.decl.CFG.CreateNode(exit_count)
	if int(id) != len(builder.decl.Ops) {
		panic(op)
	}
	builder.decl.Ops = append(builder.decl.Ops, op)
	return id
}

func MakeGoFlowBuilder(decl *LLFunc) *GoFlowBuilder {
	decl.CFG = base.CreateGraph()
	decl.Ops = []GoOp{
		&Entry{},
		&Exit{},
	}
	return &GoFlowBuilder{
		decl: decl,
	}
}

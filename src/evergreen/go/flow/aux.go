// Package flow implements a graph IR for the Go language.
package flow

import (
	"evergreen/go/core"
	"evergreen/graph"
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

func (scope *FlowFunc_Scope) Get(ref FlowFunc_Ref) *FlowFunc {
	return scope.objects[ref]
}

func (scope *FlowFunc_Scope) Register(info *FlowFunc) FlowFunc_Ref {
	index := FlowFunc_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return index
}

func (scope *FlowFunc_Scope) Len() int {
	return len(scope.objects)
}

func (scope *FlowFunc_Scope) Iter() *funcIterator {
	return &funcIterator{scope: scope, current: -1}
}

type funcIterator struct {
	scope   *FlowFunc_Scope
	current int
}

func (iter *funcIterator) Next() bool {
	iter.current += 1
	return iter.current < len(iter.scope.objects)
}

func (iter *funcIterator) Value() (FlowFunc_Ref, *FlowFunc) {
	return FlowFunc_Ref(iter.current), iter.scope.objects[iter.current]
}

type GoFlowBuilder struct {
	decl *FlowFunc
	CFG  *graph.Graph
}

func (builder *GoFlowBuilder) MakeRegister(name string, t core.GoType) Register_Ref {
	return builder.decl.Register_Scope.Register(&Register{Name: name, T: t})
}

func (builder *GoFlowBuilder) EmitOp(op GoOp, exit_count int) graph.NodeID {
	id := builder.decl.CFG.CreateNode(exit_count)
	if int(id) != len(builder.decl.Ops) {
		panic(op)
	}
	builder.decl.Ops = append(builder.decl.Ops, op)
	return id
}

func (builder *GoFlowBuilder) EmitEdge(nid graph.NodeID, flow int) graph.EdgeID {
	return builder.decl.CFG.IndexedExitEdge(nid, flow)
}

func (builder *GoFlowBuilder) EmitConnection(src graph.NodeID, flow int, dst graph.NodeID) graph.EdgeID {
	g := builder.decl.CFG
	edge := g.IndexedExitEdge(src, flow)
	g.Connect(edge, dst)
	return edge
}

func MakeGoFlowBuilder(decl *FlowFunc) *GoFlowBuilder {
	decl.CFG = graph.CreateGraph()
	decl.Ops = []GoOp{
		&Entry{},
		&Exit{},
	}
	return &GoFlowBuilder{
		decl: decl,
	}
}

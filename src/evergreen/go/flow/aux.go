// Package flow implements a graph IR for the Go language.
package flow

import (
	"evergreen/go/core"
	"evergreen/graph"
)

const (
	NORMAL = iota
	COND_TRUE
	COND_FALSE
	RETURN
	NUM_FLOWS
)

func (scope *Register_Scope) Get(ref Register_Ref) *Register {
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *Register_Scope) Register(info *Register) Register_Ref {
	info.Index = Register_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info.Index
}

func (scope *Register_Scope) Len() int {
	return len(scope.objects)
}

func (scope *FlowFunc_Scope) Get(ref FlowFunc_Ref) *FlowFunc {
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *FlowFunc_Scope) Register(info *FlowFunc) FlowFunc_Ref {
	info.Index = FlowFunc_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info.Index
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

func AllocNode(decl *FlowFunc, op GoOp) graph.NodeID {
	n := decl.CFG.CreateNode()
	if int(n) != len(decl.Ops) {
		panic(op)
	}
	decl.Ops = append(decl.Ops, op)
	return n

}

func AllocEdge(decl *FlowFunc, flow int) graph.EdgeID {
	e := decl.CFG.CreateEdge()
	if int(e) != len(decl.Edges) {
		panic(flow)
	}
	decl.Edges = append(decl.Edges, flow)
	return e
}

type GoFlowBuilder struct {
	decl *FlowFunc
}

func (builder *GoFlowBuilder) MakeRegister(name string, t core.GoType) Register_Ref {
	return builder.decl.Register_Scope.Register(&Register{Name: name, T: t})
}

func (builder *GoFlowBuilder) EmitOp(op GoOp) graph.NodeID {
	return AllocNode(builder.decl, op)
}

func (builder *GoFlowBuilder) EmitEdge(nid graph.NodeID, flow int) graph.EdgeID {
	e := AllocEdge(builder.decl, flow)
	builder.decl.CFG.ConnectEdgeEntry(nid, e)
	return e
}

func (builder *GoFlowBuilder) EmitConnection(src graph.NodeID, flow int, dst graph.NodeID) graph.EdgeID {
	e := AllocEdge(builder.decl, flow)
	builder.decl.CFG.ConnectEdge(src, e, dst)
	return e
}

func (builder *GoFlowBuilder) ConnectEdgeExit(e graph.EdgeID, dst graph.NodeID) {
	builder.decl.CFG.ConnectEdgeExit(e, dst)
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

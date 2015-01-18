// Package flow implements a graph IR for the Dub language.
package flow

import (
	"evergreen/dub/core"
	"evergreen/graph"
)

const (
	NORMAL = iota
	COND_TRUE
	COND_FALSE
	FAIL
	EXCEPTION
	RETURN
	NUM_FLOWS
)

type edgeTypeInfo struct {
	IsLocalFlow   bool
	AsInlinedFlow int
}

var EdgeTypeInfo = []edgeTypeInfo{
	edgeTypeInfo{
		IsLocalFlow:   true,
		AsInlinedFlow: NORMAL,
	},
	edgeTypeInfo{
		IsLocalFlow:   true,
		AsInlinedFlow: COND_TRUE,
	},
	edgeTypeInfo{
		IsLocalFlow:   true,
		AsInlinedFlow: COND_FALSE,
	},
	edgeTypeInfo{
		IsLocalFlow:   false,
		AsInlinedFlow: FAIL,
	},
	edgeTypeInfo{
		IsLocalFlow:   false,
		AsInlinedFlow: EXCEPTION,
	},
	edgeTypeInfo{
		IsLocalFlow:   false,
		AsInlinedFlow: NORMAL,
	},
}

func (scope *RegisterInfo_Scope) Get(ref RegisterInfo_Ref) *RegisterInfo {
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *RegisterInfo_Scope) Register(info *RegisterInfo) *RegisterInfo {
	info.Index = RegisterInfo_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info
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
			info.Index = RegisterInfo_Ref(idx)
		} else {
			info.Index = ^RegisterInfo_Ref(0)
		}
	}
	scope.objects = objects
	scope.check()
}

func (scope *RegisterInfo_Scope) Replace(replacement []*RegisterInfo) {
	scope.objects = replacement
	scope.check()
}

func (scope *RegisterInfo_Scope) check() {
	for i, reg := range scope.objects {
		if reg.Index != RegisterInfo_Ref(i) {
			panic(i)
		}
	}
}

func AllocNode(decl *LLFunc, op DubOp) graph.NodeID {
	n := decl.CFG.CreateNode()
	if int(n) != len(decl.Ops) {
		panic(op)
	}
	decl.Ops = append(decl.Ops, op)
	return n

}

func AllocEdge(decl *LLFunc, flow int) graph.EdgeID {
	e := decl.CFG.CreateEdge()
	if int(e) != len(decl.Edges) {
		panic(flow)
	}
	decl.Edges = append(decl.Edges, flow)
	return e
}

// TODO plumb through type index.
func MayHaveSideEffects(c core.Callable) bool {
	switch c := c.(type) {
	case *core.IntrinsicFunction:
		// HACK should be comparaing against builtins, not names.
		switch c.Name {
		case "append":
			// Not entirely true, but close enough for now.
			return false
		case "position":
			return false
		default:
			panic(c)
		}
	}
	return true
}

func IsNop(op DubOp) bool {
	switch op := op.(type) {
	case *Consume:
		return false
	case *Fail:
		return false
	case *Checkpoint:
		return op.Dst == nil
	case *Peek:
		return false
	case *LookaheadBegin:
		return false
	case *ConstantRuneOp:
		return op.Dst == nil
	case *ConstantStringOp:
		return op.Dst == nil
	case *ConstantIntOp:
		return op.Dst == nil
	case *ConstantFloat32Op:
		return op.Dst == nil
	case *ConstantBoolOp:
		return op.Dst == nil
	case *ConstantNilOp:
		return op.Dst == nil
	case *CallOp:
		return len(op.Dsts) == 0 && !MayHaveSideEffects(op.Target)
	case *Slice:
		return op.Dst == nil
	case *BinaryOp:
		return op.Dst == nil
	case *CopyOp:
		return op.Dst == nil || op.Dst == op.Src
	case *CoerceOp:
		return op.Dst == nil
	case *Recover:
		return false
	case *LookaheadEnd:
		return false
	case *ReturnOp:
		return false
	case *ConstructOp:
		return op.Dst == nil
	case *ConstructListOp:
		return op.Dst == nil
	case *TransferOp:
		return len(op.Dsts) == 0
	case *EntryOp:
		return false
	case *SwitchOp:
		return false
	case *ExitOp:
		return false
	default:
		panic(op)
	}
}

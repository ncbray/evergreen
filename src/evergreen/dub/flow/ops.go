// Package flow implements a graph IR for the Dub language.
package flow

import (
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
	case *ConstantFloat32Op:
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
	case *ExitOp:
		return false
	default:
		panic(op)
	}
}

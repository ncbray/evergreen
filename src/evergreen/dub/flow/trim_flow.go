package flow

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"evergreen/graph"
)

func TrimFlow(status compiler.PassStatus, program *DubProgram) {
	status.Begin()
	defer status.End()

	// TODO use whole-program analysis to agressively find dead flows.
	lut := map[core.Function_Ref]int{}

	// HACK assumes two flows.
	numFlows := 2

	flows := make([][]bool, numFlows)
	for i := 0; i < numFlows; i++ {
		flows[i] = make([]bool, len(program.LLFuncs))
	}

	// Find the exit flows of every function.
	for i, f := range program.LLFuncs {
		lut[f.F] = i
		it := f.CFG.EntryIterator(f.CFG.Exit())
		for it.HasNext() {
			e, _ := it.GetNext()
			op := f.Ops[e]
			flowExit, ok := op.(*FlowExitOp)
			if !ok {
				panic(op)
			}
			flows[flowExit.Flow][i] = true
		}
	}

	for _, f := range program.LLFuncs {
		g := f.CFG
		for node, op := range f.Ops {
			switch op := op.(type) {
			case *CallOp:
				tgt, ok := lut[op.Target]
				if !ok {
					panic(op.Target)
				}
				n := graph.NodeID(node)
				iter := g.ExitIterator(n)
				for iter.HasNext() {
					e, _ := iter.GetNext()
					possible := flows[f.Edges[e]][tgt]
					if !possible {
						g.KillEdge(e)
					}
				}
			default:
			}
		}
	}
}

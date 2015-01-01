package flow

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"evergreen/graph"
)

func TrimFlow(status compiler.PassStatus, program *DubProgram) {
	status.Begin()
	defer status.End()

	// HACK assumes a particular order of flow enums.
	exitFlowToLocal := []int{NORMAL, FAIL, EXCEPTION, NORMAL}
	// HACK assumes two flows.
	numFlows := 2

	// TODO use whole-program analysis to agressively find dead flows.
	lut := map[core.Function_Ref]int{}

	flows := make([][]bool, numFlows)
	for i := 0; i < numFlows; i++ {
		flows[i] = make([]bool, len(program.LLFuncs))
	}

	// Find the exit flows of every function.
	for i, f := range program.LLFuncs {
		lut[f.F] = i
		it := f.CFG.EntryIterator(f.CFG.Exit())
		for it.HasNext() {
			_, e := it.GetNext()
			exitFlow := f.Edges[e]
			localFlow := exitFlowToLocal[exitFlow]
			flows[localFlow][i] = true
		}
	}

	// For each call site, kill edges that will not be taken in practice.
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

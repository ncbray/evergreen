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

	flows := make([][]bool, NUM_FLOWS)
	for i := 0; i < NUM_FLOWS; i++ {
		flows[i] = make([]bool, len(program.LLFuncs))
	}

	// Find the exit flows of every function.
	for i, f := range program.LLFuncs {
		lut[f.F.Index] = i
		it := f.CFG.EntryIterator(f.CFG.Exit())
		for it.HasNext() {
			_, e := it.GetNext()
			exitFlow := f.Edges[e]
			localFlow := EdgeTypeInfo[exitFlow].AsInlinedFlow
			flows[localFlow][i] = true
		}
	}

	// For each call site, kill edges that will not be taken in practice.
	for _, f := range program.LLFuncs {
		g := f.CFG
		for node, op := range f.Ops {
			switch op := op.(type) {
			case *CallOp:
				var tgt int
				switch c := op.Target.(type) {
				case *core.CallableFunction:
					var ok bool
					tgt, ok = lut[c.Func.Index]
					if !ok {
						panic(c.Func)
					}
				default:
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

package flow

import (
	"evergreen/base"
	"evergreen/dub/core"
	"evergreen/framework"
)

func TrimFlow(status framework.PassStatus, program []*DubPackage) {
	status.Begin()
	defer status.End()

	// TODO use whole-program analysis to agressively find dead flows.
	// HACK assumes no cross-package function references.
	for _, pkg := range program {
		lut := map[*core.Function]int{}

		// HACK assumes two flows.
		numFlows := 2

		flows := make([][]bool, numFlows)
		for i := 0; i < numFlows; i++ {
			flows[i] = make([]bool, len(pkg.Funcs))
		}

		// Find the exit flows of every function.
		for i, f := range pkg.Funcs {
			lut[f.F] = i
			it := base.EntryIterator(f.CFG, f.CFG.Exit())
			for it.Next() {
				op := f.Ops[it.Value()]
				flowExit, ok := op.(*FlowExitOp)
				if !ok {
					panic(op)
				}
				flows[flowExit.Flow][i] = true
			}
		}

		for _, f := range pkg.Funcs {
			for node, op := range f.Ops {
				switch op := op.(type) {
				case *CallOp:
					tgt, ok := lut[op.Target]
					if !ok {
						panic(op.Target)
					}
					for i := 0; i < numFlows; i++ {
						possible := flows[i][tgt]
						if !possible {
							f.CFG.Disconnect(base.NodeID(node), i)
						}
					}
				default:
				}
			}
		}
	}
}

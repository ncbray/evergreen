package tree

import (
	"evergreen/compiler"
	"evergreen/go/core"
)

func GoProgramBackend(status compiler.PassStatus, prog *ProgramAST, coreProg *core.CoreProgram, srcRoot string, runner *compiler.TaskRunner) {
	status.Begin()
	defer status.End()

	// Compact simple expressions back into tree form.
	Consolidate(status.Pass("consolidate"), prog)

	// Give everything names: variables, etc.
	Nameify(status.Pass("nameify"), prog, coreProg)

	// Generate the sources.
	OutputProgram(status.Pass("output"), prog, coreProg, srcRoot, runner)
}

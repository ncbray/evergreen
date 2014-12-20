package tree

import (
	"evergreen/compiler"
)

func GoProgramBackend(status compiler.PassStatus, prog *ProgramAST, srcRoot string, runner *compiler.TaskRunner) {
	status.Begin()
	defer status.End()

	// Compact simple expressions back into tree form.
	Consolidate(status.Pass("consolidate"), prog)

	// Give everything names: variables, etc.
	Nameify(status.Pass("nameify"), prog)

	// Generate the sources.
	OutputProgram(status.Pass("output"), prog, srcRoot, runner)
}

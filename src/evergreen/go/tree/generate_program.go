// Package tree implements a tree IR for the Go language.
package tree

import (
	"evergreen/compiler"
	"evergreen/go/core"
	"evergreen/io"
	"evergreen/text"
	"path/filepath"
)

func OutputFile(file *FileAST, dirname string, runner *compiler.TaskRunner) {
	if file.Name == "" {
		panic(file)
	}

	runner.Run(func() {
		filename := filepath.Join(dirname, file.Name)
		b, w := text.BufferedCodeWriter()
		GenerateFile(file, w)
		io.WriteFile(filename, []byte(b.String()))
	})
}

func OutputPackage(pkg *PackageAST, coreProg *core.CoreProgram, dirname string, runner *compiler.TaskRunner) {
	p := pkg.P

	path := []string{dirname}
	path = append(path, p.Path...)
	pkgdir := filepath.Join(path...)
	for _, file := range pkg.Files {
		OutputFile(file, pkgdir, runner)
	}
}

func OutputProgram(status compiler.PassStatus, prog *ProgramAST, coreProg *core.CoreProgram, dirname string, runner *compiler.TaskRunner) {
	status.Begin()
	defer status.End()

	for _, pkg := range prog.Packages {
		OutputPackage(pkg, coreProg, dirname, runner)
	}
}

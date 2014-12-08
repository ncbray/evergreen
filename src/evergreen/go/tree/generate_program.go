// Package tree implements a tree IR for the Go language.
package tree

import (
	"evergreen/compiler"
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

func OutputPackage(pkg *PackageAST, dirname string, runner *compiler.TaskRunner) {
	path := []string{dirname}
	path = append(path, pkg.P.Path...)
	pkgdir := filepath.Join(path...)
	for _, file := range pkg.Files {
		OutputFile(file, pkgdir, runner)
	}
}

func OutputProgram(status compiler.PassStatus, prog *ProgramAST, dirname string, runner *compiler.TaskRunner) {
	status.Begin()
	defer status.End()

	for _, pkg := range prog.Packages {
		OutputPackage(pkg, dirname, runner)
	}
}

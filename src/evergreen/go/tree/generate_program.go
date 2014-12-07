// Package tree implements a tree IR for the Go language.
package tree

import (
	"bytes"
	"evergreen/base"
	"evergreen/framework"
	"evergreen/io"
	"path/filepath"
)

func OutputFile(file *FileAST, dirname string, runner *framework.TaskRunner) {
	if file.Name == "" {
		panic(file)
	}

	runner.Run(func() {
		filename := filepath.Join(dirname, file.Name)
		b := &bytes.Buffer{}
		w := &base.CodeWriter{Out: b}
		GenerateFile(file, w)
		io.WriteFile(filename, []byte(b.String()))
	})
}

func OutputPackage(pkg *PackageAST, dirname string, runner *framework.TaskRunner) {
	path := []string{dirname}
	path = append(path, pkg.P.Path...)
	pkgdir := filepath.Join(path...)
	for _, file := range pkg.Files {
		OutputFile(file, pkgdir, runner)
	}
}

func OutputProgram(status framework.PassStatus, prog *ProgramAST, dirname string, runner *framework.TaskRunner) {
	status.Begin()
	defer status.End()

	for _, pkg := range prog.Packages {
		OutputPackage(pkg, dirname, runner)
	}
}

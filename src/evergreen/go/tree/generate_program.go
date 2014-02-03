package tree

import (
	"bytes"
	"evergreen/base"
	"evergreen/io"
	"path/filepath"
)

func OutputFile(file *File, dirname string) {
	if file.Name == "" {
		panic(file)
	}
	filename := filepath.Join(dirname, file.Name)
	b := &bytes.Buffer{}
	w := &base.CodeWriter{Out: b}
	GenerateFile(file, w)
	io.WriteFile(filename, []byte(b.String()))
}

func OutputPackage(pkg *Package, dirname string) {
	path := []string{dirname}
	path = append(path, pkg.Path...)
	pkgdir := filepath.Join(path...)
	for _, file := range pkg.Files {
		OutputFile(file, pkgdir)
	}
}

func OutputProgram(prog *Program, dirname string) {
	for _, pkg := range prog.Packages {
		if !pkg.Extern {
			OutputPackage(pkg, dirname)
		}
	}
}

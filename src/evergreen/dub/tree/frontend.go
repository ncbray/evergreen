// Package tree implements a tree IR for the Dub language.
package tree

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func parsePackage(status compiler.PassStatus, p compiler.LocationProvider, path []string, filenames []string) *Package {
	files := make([]*File, len(filenames))
	for i, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			status.GlobalError(err.Error())
			return nil
		}
		stream := []rune(string(data))
		offset := p.AddFile(filename, stream)
		files[i] = ParseDub(data, offset, status.Task(filename))
		if status.ShouldHalt() {
			return nil
		}
		files[i].Name = filepath.Base(filename)
	}

	pkg := &Package{
		Path:  path,
		Files: files,
	}

	return pkg
}

func extendPath(path []string, next string) []string {
	newPath := make([]string, len(path)+1)
	copy(newPath, path)
	newPath[len(path)] = next
	return newPath
}

func parsePackageTree(status compiler.PassStatus, p compiler.LocationProvider, root string, path []string, packages []*Package) []*Package {
	dir := filepath.Join(root, strings.Join(path, string(filepath.Separator)))
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		status.GlobalError(err.Error())
		return nil
	}
	dubfiles := []string{}
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if file.IsDir() {
			newPath := extendPath(path, name)
			packages = parsePackageTree(status, p, root, newPath, packages)
		} else {
			if strings.HasSuffix(name, ".dub") {
				filename := file.Name()
				fullpath := filepath.Join(dir, filename)
				dubfiles = append(dubfiles, fullpath)
			}
		}
	}
	if len(dubfiles) > 0 {
		pkg := parsePackage(status, p, path, dubfiles)
		if pkg != nil {
			packages = append(packages, pkg)
		}
	}
	return packages
}

func parseProgram(status compiler.PassStatus, p compiler.LocationProvider, root string) *Program {
	status.Begin()
	defer status.End()

	packages := parsePackageTree(status, p, root, []string{}, []*Package{})
	if status.ShouldHalt() {
		return nil
	}
	return &Program{
		Packages: packages,
	}
}

func DubProgramFrontend(status compiler.PassStatus, p compiler.LocationProvider, root string) (*Program, *core.CoreProgram) {
	status.Begin()
	defer status.End()
	program := parseProgram(status.Pass("parse"), p, root)
	if status.ShouldHalt() {
		return nil, nil
	}
	coreProg := SemanticPass(program, status.Pass("semantic"))
	if status.ShouldHalt() {
		return nil, nil
	}
	return program, coreProg
}

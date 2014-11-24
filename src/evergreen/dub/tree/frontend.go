package tree

import (
	"evergreen/framework"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func parsePackage(status framework.Status, p framework.LocationProvider, path []string, filenames []string) *Package {
	files := make([]*File, len(filenames))
	for i, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			status.Error("%s", err.Error())
			return nil
		}
		offset := p.AddFile(filename, []rune(string(data)))
		files[i] = ParseDub(data, offset, status)

		if status.ShouldHalt() {
			return nil
		}
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

func parsePackageTree(status framework.Status, p framework.LocationProvider, root string, path []string, packages []*Package) []*Package {
	dir := filepath.Join(root, strings.Join(path, string(filepath.Separator)))
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		status.Error(err.Error())
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
		pkg := parsePackage(status.CreateChild(), p, path, dubfiles)
		if pkg != nil {
			packages = append(packages, pkg)
		}
	}
	return packages
}

func ParseProgram(status framework.Status, p framework.LocationProvider, root string) *Program {
	packages := parsePackageTree(status.CreateChild(), p, root, []string{}, []*Package{})
	if status.ShouldHalt() {
		return nil
	}
	program := &Program{
		Packages: packages,
		Builtins: MakeBuiltinTypeIndex(),
	}
	SemanticPass(program, status.CreateChild())
	return program
}

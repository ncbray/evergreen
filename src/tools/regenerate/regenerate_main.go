package main

import (
	"evergreen/base"
	"evergreen/dub"
	"evergreen/dub/flow"
	"evergreen/dub/tree"
	"evergreen/framework"
	gotree "evergreen/go/tree"
	"evergreen/io"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type IOManager struct {
	limit chan bool
	group sync.WaitGroup
}

func (m *IOManager) Create() {
	m.group.Add(1)
}

func (m *IOManager) Aquire() {
	<-m.limit
}

func (m *IOManager) Release() {
	m.limit <- true
	m.group.Done()
}

func (m *IOManager) Flush() {
	m.group.Wait()
}

func (m *IOManager) WriteFile(filename string, data []byte) {
	m.Create()
	go func() {
		m.Aquire()
		defer m.Release()
		io.WriteFile(filename, data)
	}()
}

func CreateIOManager() *IOManager {
	limit := 8
	manager := &IOManager{limit: make(chan bool, limit)}
	for i := 0; i < limit; i++ {
		manager.limit <- true
	}
	return manager
}

func parsePackage(status framework.Status, p framework.LocationProvider, path []string, filenames []string) (*tree.Package, *tree.ModuleScope) {
	fmt.Printf("Processing %s...\n", strings.Join(path, "."))

	files := make([]*tree.File, len(filenames))
	for i, filename := range filenames {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			status.Error("%s", err)
			return nil, nil
		}
		p.AddFile(filename, []rune(string(data)))
		files[i] = tree.ParseDub(data, status)

		if status.ShouldHalt() {
			return nil, nil
		}
	}

	pkg := &tree.Package{
		Path:  path,
		Files: files,
	}

	glbls := tree.MakeDubGlobals()
	tree.SemanticPass(pkg, glbls, status)
	if status.ShouldHalt() {
		return nil, nil
	}

	return pkg, glbls
}

func processDub(status framework.Status, p framework.LocationProvider, manager *IOManager, path []string, filenames []string) {

	pkg, glbls := parsePackage(status, p, path, filenames)
	if status.ShouldHalt() {
		return
	}

	gbuilder := &dub.GlobalDubBuilder{Types: map[tree.ASTType]flow.DubType{}}

	gbuilder.String = &flow.IntrinsicType{Name: "string"}
	gbuilder.Types[glbls.String] = gbuilder.String

	gbuilder.Rune = &flow.IntrinsicType{Name: "rune"}
	gbuilder.Types[glbls.Rune] = gbuilder.Rune

	gbuilder.Int = &flow.IntrinsicType{Name: "int"}
	gbuilder.Types[glbls.Int] = gbuilder.Int

	gbuilder.Int64 = &flow.IntrinsicType{Name: "int64"}
	gbuilder.Types[glbls.Int64] = gbuilder.Int64

	gbuilder.Bool = &flow.IntrinsicType{Name: "bool"}
	gbuilder.Types[glbls.Bool] = gbuilder.Bool

	gbuilder.Graph = &flow.IntrinsicType{Name: "graph"}
	gbuilder.Types[glbls.Graph] = gbuilder.Graph

	// Preallocate the translated structures.
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *tree.FuncDecl:
			case *tree.StructDecl:
				gbuilder.Types[decl] = &flow.LLStruct{}
			default:
				panic(decl)
			}
		}
	}

	structs := []*flow.LLStruct{}
	funcs := []*flow.LLFunc{}
	tests := []*tree.Test{}

	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *tree.FuncDecl:
				f := dub.LowerAST(decl, gbuilder)
				flow.SSI(f)
				funcs = append(funcs, f)

				if dump {
					styler := &flow.DotStyler{Decl: f}
					dot := base.GraphToDot(f.CFG, styler)
					parts := []string{"output"}
					parts = append(parts, path...)
					parts = append(parts, fmt.Sprintf("%s.svg", f.Name))
					outfile := filepath.Join(parts...)
					manager.Create()
					go func(dot string, outfile string) {
						manager.Aquire()
						defer manager.Release()

						// Dump flowgraph
						io.WriteDot(dot, outfile)
					}(dot, outfile)
				}
			case *tree.StructDecl:
				t, _ := gbuilder.Types[decl]
				s, _ := t.(*flow.LLStruct)
				structs = append(structs, dub.LowerStruct(decl, s, gbuilder))
			default:
				panic(decl)
			}
		}
		tests = append(tests, file.Tests...)
	}

	// Analysis
	for _, s := range structs {
		if s.Implements != nil {
			s.Implements.Abstract = true
		}
	}

	GenerateGo(path, structs, funcs, tests, gbuilder)
}

func GenerateGo(original_path []string, structs []*flow.LLStruct, funcs []*flow.LLFunc, tests []*tree.Test, gbuilder *dub.GlobalDubBuilder) {
	root := "generated"
	if replace {
		root = "evergreen"
	}

	path := []string{root}
	path = append(path, original_path...)
	leaf := path[len(path)-1]

	link := flow.MakeLinker()

	packages := []*gotree.Package{}

	pkg, index := flow.ExternBuiltinRuntime()
	packages = append(packages, pkg)

	pkg, state := flow.ExternParserRuntime()
	packages = append(packages, pkg)

	pkg, graph := flow.ExternGraph()
	packages = append(packages, pkg)

	files := []*gotree.File{}
	files = append(files, flow.GenerateGo(leaf, structs, funcs, index, state, graph, link))

	if !replace && len(tests) != 0 {
		pkg, t := dub.ExternTestingPackage()
		packages = append(packages, pkg)
		pkg, stateT := dub.ExternRuntimePackage()
		packages = append(packages, pkg)

		files = append(files, dub.GenerateTests(leaf, tests, gbuilder, t, stateT, link))
	}

	packages = append(packages, &gotree.Package{
		Path:  path,
		Files: files,
	})

	prog := &gotree.Program{
		Packages: packages,
	}

	link.Finish()

	// Compact simple expressions back into tree form.
	gotree.Retree(prog)

	// Give everything names: variables, etc.
	gotree.Nameify(prog)

	// Generate the sources.
	gotree.OutputProgram(prog, "src")
}

var dump bool
var replace bool

func extendPath(path []string, next string) []string {
	newPath := make([]string, len(path)+1)
	copy(newPath, path)
	newPath[len(path)] = next
	return newPath
}

func processPackage(status framework.Status, p framework.LocationProvider, manager *IOManager, root string, path []string) {
	dir := filepath.Join(root, strings.Join(path, string(filepath.Separator)))
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dubfiles := []string{}
	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if file.IsDir() {
			newPath := extendPath(path, name)
			processPackage(status, p, manager, root, newPath)
		} else {
			if strings.HasSuffix(name, ".dub") {
				filename := file.Name()
				fullpath := filepath.Join(dir, filename)
				dubfiles = append(dubfiles, fullpath)
			}
		}
	}
	if len(dubfiles) > 0 {
		processDub(status.CreateChild(), p, manager, path, dubfiles)
	}
}

func main() {
	flag.BoolVar(&dump, "dump", false, "Dump flowgraphs to disk.")
	flag.BoolVar(&replace, "replace", false, "Replace the existing implementation.")
	flag.Parse()
	p := framework.MakeProvider()
	status := framework.MakeStatus(p)
	manager := CreateIOManager()

	root_dir := "dub"
	processPackage(status.CreateChild(), p, manager, root_dir, []string{})
	manager.Flush()
	if status.ShouldHalt() {
		fmt.Printf("%d errors\n", status.ErrorCount())
		os.Exit(1)
	}
}

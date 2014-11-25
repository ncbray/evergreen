package main

import (
	"evergreen/base"
	"evergreen/dub"
	"evergreen/dub/flow"
	"evergreen/dub/tree"
	core "evergreen/dub/tree"
	"evergreen/framework"
	gotree "evergreen/go/tree"
	"evergreen/io"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
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

type DubPackage struct {
	Path    []string
	Structs []*core.StructType
	Funcs   []*flow.LLFunc
	Tests   []*tree.Test
}

func lowerPackage(program *tree.Program, pkg *tree.Package) *DubPackage {
	dubPkg := &DubPackage{
		Path:    pkg.Path,
		Structs: []*core.StructType{},
		Funcs:   []*flow.LLFunc{},
		Tests:   []*tree.Test{},
	}

	// Lower to flow IR
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *tree.FuncDecl:
				f := dub.LowerAST(program, decl)
				flow.SSI(f)
				dubPkg.Funcs = append(dubPkg.Funcs, f)
			case *tree.StructDecl:
				dubPkg.Structs = append(dubPkg.Structs, decl.T)
			default:
				panic(decl)
			}
		}
		dubPkg.Tests = append(dubPkg.Tests, file.Tests...)
	}

	return dubPkg
}

func lowerProgram(program *tree.Program) []*DubPackage {
	dubPackages := []*DubPackage{}
	for _, pkg := range program.Packages {
		dubPackages = append(dubPackages, lowerPackage(program, pkg))
	}
	return dubPackages
}

func dumpProgram(manager *IOManager, program []*DubPackage) {
	for _, dubPkg := range program {
		for _, f := range dubPkg.Funcs {
			styler := &flow.DotStyler{Decl: f}
			dot := base.GraphToDot(f.CFG, styler)
			parts := []string{"output"}
			parts = append(parts, dubPkg.Path...)
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
	}
}

func analyizeProgram(program []*DubPackage) {
	for _, dubPkg := range program {
		for _, s := range dubPkg.Structs {
			if s.Implements != nil {
				s.Implements.IsParent = true
			}
		}
	}
}

func GenerateGo(program []*DubPackage) {
	root := "generated"
	if replace {
		root = "evergreen"
	}

	link := flow.MakeLinker()

	packages := []*gotree.Package{}

	pkg, index := flow.ExternBuiltinRuntime()
	packages = append(packages, pkg)

	pkg, state := flow.ExternParserRuntime()
	packages = append(packages, pkg)

	pkg, graph := flow.ExternGraph()
	packages = append(packages, pkg)

	pkg, t := flow.ExternTestingPackage()
	packages = append(packages, pkg)

	for _, dubPkg := range program {
		path := []string{root}
		path = append(path, dubPkg.Path...)
		leaf := path[len(path)-1]

		files := []*gotree.File{}
		files = append(files, flow.GenerateGo(leaf, dubPkg.Structs, dubPkg.Funcs, index, state, graph, link))

		if !replace && len(dubPkg.Tests) != 0 {
			files = append(files, dub.GenerateTests(leaf, dubPkg.Tests, t, state, link))
		}
		packages = append(packages, &gotree.Package{
			Path:  path,
			Files: files,
		})
	}

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

func processProgram(status framework.Status, p framework.LocationProvider, manager *IOManager, root string) {
	program := tree.ParseProgram(status.CreateChild(), p, root)
	if status.ShouldHalt() {
		return
	}
	flowProgram := lowerProgram(program)

	if dump {
		dumpProgram(manager, flowProgram)
	}

	analyizeProgram(flowProgram)
	GenerateGo(flowProgram)
}

func entryPoint(p framework.LocationProvider, status framework.Status) {
	manager := CreateIOManager()

	root_dir := "dub"
	processProgram(status, p, manager, root_dir)
	manager.Flush()
}

var dump bool
var replace bool
var cpuprofile string
var memprofile string

func main() {
	flag.BoolVar(&dump, "dump", false, "Dump flowgraphs to disk.")
	flag.BoolVar(&replace, "replace", false, "Replace the existing implementation.")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to this file")
	flag.Parse()

	p := framework.MakeProvider()
	status := framework.MakeStatus(p)

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			status.Error(err.Error())
			return
		} else {
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	}

	entryPoint(p, status)

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			status.Error(err.Error())
			return
		} else {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}

	if status.ShouldHalt() {
		fmt.Printf("%d errors\n", status.ErrorCount())
		os.Exit(1)
	}
}

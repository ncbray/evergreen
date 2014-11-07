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

func processDub(status framework.Status, p framework.LocationProvider, manager *IOManager, language_name string, ir_name string, filename string) {
	fmt.Printf("Processing %s.%s...\n", language_name, ir_name)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		status.Error("%s", err)
		return
	}
	p.AddFile(filename, []rune(string(data)))
	file := tree.ParseDub(data, status)
	if status.ShouldHalt() {
		return
	}
	glbls := tree.SemanticPass(file, status)
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

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *tree.FuncDecl:
		case *tree.StructDecl:
			gbuilder.Types[decl] = &flow.LLStruct{}
		default:
			panic(decl)
		}
	}

	structs := []*flow.LLStruct{}
	funcs := []*flow.LLFunc{}
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *tree.FuncDecl:
			f := dub.LowerAST(decl, gbuilder)
			flow.SSI(f)
			funcs = append(funcs, f)

			if dump {
				styler := &flow.DotStyler{Decl: f}
				dot := base.GraphToDot(f.CFG, styler)
				outfile := filepath.Join("output", language_name, ir_name, fmt.Sprintf("%s.svg", f.Name))
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

	// Analysis
	for _, s := range structs {
		if s.Implements != nil {
			s.Implements.Abstract = true
		}
	}

	GenerateGo(language_name, ir_name, file, structs, funcs, gbuilder)
}

func GenerateGo(language_name string, ir_name string, file *tree.File, structs []*flow.LLStruct, funcs []*flow.LLFunc, gbuilder *dub.GlobalDubBuilder) {
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

	files := []*gotree.File{}
	files = append(files, flow.GenerateGo(ir_name, structs, funcs, index, state, link))

	if !replace && len(file.Tests) != 0 {
		pkg, t := dub.ExternTestingPackage()
		packages = append(packages, pkg)
		pkg, stateT := dub.ExternRuntimePackage()
		packages = append(packages, pkg)

		files = append(files, dub.GenerateTests(language_name, file.Tests, gbuilder, t, stateT, link))
	}

	packages = append(packages, &gotree.Package{
		Path:  []string{root, language_name, ir_name},
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

func processLanguage(status framework.Status, p framework.LocationProvider, manager *IOManager, language_name string, dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".dub") {
			filename := file.Name()
			fullpath := filepath.Join(dir, filename)
			ir_name := filename[:len(filename)-4]
			processDub(status.CreateChild(), p, manager, language_name, ir_name, fullpath)
		}
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
	files, err := ioutil.ReadDir(root_dir)
	if err != nil {
		status.Error(err.Error())
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			processLanguage(status.CreateChild(), p, manager, file.Name(), filepath.Join(root_dir, file.Name()))
		}
	}
	manager.Flush()
	if status.ShouldHalt() {
		os.Exit(1)
	}
}

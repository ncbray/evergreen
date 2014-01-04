package main

import (
	"evergreen/base"
	"evergreen/dub"
	"evergreen/dub/flow"
	"evergreen/dub/tree"
	"evergreen/framework"
	"evergreen/io"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func processDub(status framework.Status, manager *IOManager, name string) {
	fmt.Printf("Processing %s...\n", name)
	filename := fmt.Sprintf("dub/%s.dub", name)
	data, _ := ioutil.ReadFile(filename)
	file := tree.ParseDub(filename, data, status)
	if status.ShouldHalt() {
		return
	}
	glbls := tree.SemanticPass(file, status)
	if status.ShouldHalt() {
		return
	}

	gbuilder := &dub.GlobalDubBuilder{Types: map[tree.ASTType]flow.DubType{}}

	gbuilder.String = &flow.StringType{}
	gbuilder.Types[glbls.String] = gbuilder.String

	gbuilder.Rune = &flow.RuneType{}
	gbuilder.Types[glbls.Rune] = gbuilder.Rune

	gbuilder.Int = &flow.IntType{}
	gbuilder.Types[glbls.Int] = gbuilder.Int

	gbuilder.Bool = &flow.BoolType{}
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
				styler := &flow.DotStyler{}
				dot := base.RegionToDot(f.Region, styler)
				outfile := filepath.Join("output", name, fmt.Sprintf("%s.svg", f.Name))
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

	code := flow.GenerateGo(name, structs, funcs)
	manager.WriteFile(fmt.Sprintf("src/generated/%s/tree/generated_parser.go", name), []byte(code))

	if len(file.Tests) != 0 {
		tests := dub.GenerateTests(name, file.Tests, gbuilder)
		manager.WriteFile(fmt.Sprintf("src/generated/%s/tree/generated_parser_test.go", name), []byte(tests))
	}
}

var dump bool

func main() {
	flag.BoolVar(&dump, "dump", false, "Dump flowgraphs to disk.")
	flag.Parse()
	status := framework.MakeStatus()
	manager := CreateIOManager()
	processDub(status.CreateChild(), manager, "dub")
	manager.Flush()
	if status.ShouldHalt() {
		os.Exit(1)
	}
}

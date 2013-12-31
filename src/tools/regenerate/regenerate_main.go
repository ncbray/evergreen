package main

import (
	"evergreen/base"
	"evergreen/dasm"
	"evergreen/dub"
	"evergreen/dubx"
	"evergreen/io"
	"flag"
	"fmt"
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

func processDASM(manager *IOManager, name string) {
	fmt.Printf("Processing %s...\n", name)
	file := dasm.ParseDASM(fmt.Sprintf("dasm/%s.dasm", name))
	glbls := dasm.SemanticPass(file)
	gbuilder := &dasm.GlobalDubBuilder{Types: map[dubx.ASTType]dub.DubType{}}

	gbuilder.String = &dub.StringType{}
	gbuilder.Types[glbls.String] = gbuilder.String

	gbuilder.Rune = &dub.RuneType{}
	gbuilder.Types[glbls.Rune] = gbuilder.Rune

	gbuilder.Int = &dub.IntType{}
	gbuilder.Types[glbls.Int] = gbuilder.Int

	gbuilder.Bool = &dub.BoolType{}
	gbuilder.Types[glbls.Bool] = gbuilder.Bool

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *dubx.FuncDecl:
		case *dubx.StructDecl:
			gbuilder.Types[decl] = &dub.LLStruct{}
		default:
			panic(decl)
		}
	}

	structs := []*dub.LLStruct{}
	funcs := []*dub.LLFunc{}
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *dubx.FuncDecl:
			f := dasm.LowerAST(decl, gbuilder)
			funcs = append(funcs, f)

			if dump {
				dot := base.RegionToDot(f.Region)
				outfile := filepath.Join("output", name, fmt.Sprintf("%s.svg", f.Name))
				manager.Create()
				go func(dot string, outfile string) {
					manager.Aquire()
					defer manager.Release()

					// Dump flowgraph
					io.WriteDot(dot, outfile)
				}(dot, outfile)
			}
		case *dubx.StructDecl:
			t, _ := gbuilder.Types[decl]
			s, _ := t.(*dub.LLStruct)
			structs = append(structs, dasm.LowerStruct(decl, s, gbuilder))
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

	code := dub.GenerateGo(name, structs, funcs)
	manager.WriteFile(fmt.Sprintf("src/generated/%s/parser.go", name), []byte(code))

	if len(file.Tests) != 0 {
		tests := dasm.GenerateTests(name, file.Tests, gbuilder)
		manager.WriteFile(fmt.Sprintf("src/generated/%s/parser_test.go", name), []byte(tests))
	}
}

var dump bool

func main() {
	flag.BoolVar(&dump, "dump", false, "Dump flowgraphs to disk.")
	flag.Parse()
	manager := CreateIOManager()
	processDASM(manager, "math")
	processDASM(manager, "dubx")
	manager.Flush()
}

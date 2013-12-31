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
)

func processDASM(name string) {
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
				// Dump flowgraph
				dot := base.RegionToDot(f.Region)
				outfile := filepath.Join("output", name, fmt.Sprintf("%s.svg", f.Name))
				io.WriteDot(dot, outfile)
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
	//mt.Println(code)
	io.WriteFile(fmt.Sprintf("src/generated/%s/parser.go", name), []byte(code))

	if len(file.Tests) != 0 {
		tests := dasm.GenerateTests(name, file.Tests, gbuilder)
		//fmt.Println(tests)
		io.WriteFile(fmt.Sprintf("src/generated/%s/parser_test.go", name), []byte(tests))
	}
}

var dump bool

func main() {
	flag.BoolVar(&dump, "dump", true, "Dump flowgraphs to disk.")
	flag.Parse()
	processDASM("math")
	processDASM("dubx")
}

package main

import (
	"evergreen/base"
	"evergreen/dasm"
	"evergreen/dub"
	"evergreen/io"
	"fmt"
	"path/filepath"
)

func processDASM(name string) {
	file := dasm.ParseDASM(fmt.Sprintf("dasm/%s.dasm", name))
	glbls := dasm.SemanticPass(file)
	gbuilder := &dasm.GlobalDubBuilder{Types: map[dasm.ASTType]dub.DubType{}}

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
		case *dasm.FuncDecl:
		case *dasm.StructDecl:
			gbuilder.Types[decl] = &dub.LLStruct{}
		default:
			panic(decl)
		}
	}

	structs := []*dub.LLStruct{}
	funcs := []*dub.LLFunc{}
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *dasm.FuncDecl:
			f := dasm.LowerAST(decl, gbuilder)
			funcs = append(funcs, f)

			// Dump flowgraph
			dot := base.RegionToDot(f.Region)
			outfile := filepath.Join("output", name, fmt.Sprintf("%s.svg", f.Name))
			io.WriteDot(dot, outfile)
		case *dasm.StructDecl:
			t, _ := gbuilder.Types[decl]
			s, _ := t.(*dub.LLStruct)
			structs = append(structs, dasm.LowerStruct(decl, s, gbuilder))
		default:
			panic(decl)
		}
	}
	for _, s := range structs {
		if s.Implements != nil {
			s.Implements.Abstract = true
		}
	}

	code := dub.GenerateGo(name, structs, funcs)
	fmt.Println(code)
	io.WriteFile(fmt.Sprintf("src/generated/%s/parser.go", name), []byte(code))
}

func main() {
	processDASM("math")
	processDASM("dubx")
}

package golang

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"evergreen/dub/flow"
	dstcore "evergreen/go/core"
	dstflow "evergreen/go/flow"
	"evergreen/go/transform"
	ast "evergreen/go/tree"
	"evergreen/graph"
	"evergreen/io"
	"fmt"
	"path/filepath"
)

type DubToGoContext struct {
	index *dstcore.BuiltinTypeIndex
	state *dstcore.StructType
	graph *dstcore.StructType
	t     *dstcore.StructType
	link  DubToGoLinker
}

func generateGoFunc(f *flow.LLFunc, ctx *DubToGoContext) ast.Decl {
	flowDecl := translateFlow(f, ctx)

	if false {
		dot := graph.GraphToDot(flowDecl.CFG, &dstflow.DotStyler{Ops: flowDecl.Ops})
		parts := []string{"output", "translate"}
		parts = append(parts, fmt.Sprintf("%s.svg", flowDecl.Name))
		outfile := filepath.Join(parts...)
		io.WriteDot(dot, outfile)
	}

	return transform.RetreeFunc(flowDecl)
}

func addTags(base *core.StructType, parent *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if parent != nil {
		decls = addTags(base, parent.Implements, ctx, decls)
		decl := &ast.FuncDecl{
			Name:            tagName(parent),
			Type:            &ast.FuncTypeRef{},
			Body:            []ast.Stmt{},
			LocalInfo_Scope: &ast.LocalInfo_Scope{},
		}
		recv := decl.CreateLocalInfo("node", ast.RefForType(goType(base, ctx)))
		decl.Recv = decl.MakeParam(recv)
		decls = append(decls, decl)
	}
	return decls
}

func declForType(t dstcore.GoType, ctx *DubToGoContext) ast.Decl {
	switch t := t.(type) {
	case *dstcore.TypeDefType:
		return &ast.TypeDefDecl{
			Name: t.Name,
			Type: ast.RefForType(t.Type),
			T:    t,
		}
	case *dstcore.StructType:
		fields := []*ast.FieldDecl{}
		for _, f := range t.Fields {
			fields = append(fields, &ast.FieldDecl{
				Name: f.Name,
				Type: ast.RefForType(f.Type),
			})
		}

		return &ast.StructDecl{
			Name:   t.Name,
			Fields: fields,
			T:      t,
		}
	case *dstcore.InterfaceType:
		fields := []*ast.FieldDecl{}
		for _, f := range t.Fields {
			fields = append(fields, &ast.FieldDecl{
				Name: f.Name,
				Type: ast.RefForType(f.Type),
			})
		}
		return &ast.InterfaceDecl{
			Name:   t.Name,
			Fields: fields,
			T:      t,
		}
	default:
		panic(t)
	}
}

func generateGoStruct(s *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if s.IsParent {
		decls = append(decls, declForType(ctx.link.GetType(s, STRUCT), ctx))
	} else {
		if s.Scoped {
			ref := declForType(ctx.link.GetType(s, REF), ctx)
			scope := declForType(ctx.link.GetType(s, SCOPE), ctx)
			decls = append(decls, ref, scope)
		}
		decls = append(decls, declForType(ctx.link.GetType(s, STRUCT), ctx))
	}
	return decls
}

func generateTreeForStruct(s *core.StructType, ctx *DubToGoContext) []ast.Decl {
	decls := []ast.Decl{}
	if !s.IsParent {
		if s.Scoped {
			decls = append(decls, &ast.VarDecl{
				Name: "No" + s.Name,
				Type: ctx.link.TypeRef(s, REF),
				Expr: &ast.UnaryExpr{
					Op: "^",
					Expr: &ast.TypeCoerce{
						Type: ctx.link.TypeRef(s, REF),
						Expr: &ast.IntLiteral{Value: 0},
					},
				},
				Const: true,
			})
		}
		decls = addTags(s, s.Implements, ctx, decls)
	}
	return decls
}

func externParserRuntime() *dstcore.StructType {
	p := &dstcore.Package{
		Extern: true,
		Path:   []string{"evergreen", "dub", "runtime"},
	}
	stateT := &dstcore.StructType{
		Name:    "State",
		Package: p,
	}
	return stateT
}

func externTesting() *dstcore.StructType {
	p := &dstcore.Package{
		Extern: true,
		Path:   []string{"testing"},
	}
	tT := &dstcore.StructType{
		Name:    "T",
		Package: p,
	}
	return tT
}

func externGraph() *dstcore.StructType {
	p := &dstcore.Package{
		Extern: true,
		Path:   []string{"evergreen", "graph"},
	}
	graphT := &dstcore.StructType{
		Name:    "Graph",
		Package: p,
	}
	return graphT
}

func generateGoFile(package_name string, dubPkg *flow.DubPackage, structToDecls map[*core.StructType][]ast.Decl, ctx *DubToGoContext) *ast.FileAST {
	imports := []*ast.Import{}

	decls := []ast.Decl{}
	for _, t := range dubPkg.Structs {
		decls = generateGoStruct(t, ctx, decls)
		more, _ := structToDecls[t]
		decls = append(decls, more...)
	}
	for _, f := range dubPkg.Funcs {
		decls = append(decls, generateGoFunc(f, ctx))
	}

	file := &ast.FileAST{
		Name:    "generated_dub.go",
		Package: package_name,
		Imports: imports,
		Decls:   decls,
	}
	return file
}

func GenerateGo(status compiler.PassStatus, program *flow.DubProgram, coreProg *core.CoreProgram, root string, generate_tests bool) *ast.ProgramAST {
	status.Begin()
	defer status.End()

	ctx := &DubToGoContext{
		index: makeBuiltinTypes(),
		state: externParserRuntime(),
		graph: externGraph(),
		t:     externTesting(),
		link:  makeLinker(),
	}

	// Translate package identities.
	packages := make([]*dstcore.Package, len(program.Packages))
	for i, dubPkg := range program.Packages {
		path := []string{root}
		path = append(path, dubPkg.Path...)
		packages[i] = &dstcore.Package{
			Path: path,
		}
	}

	// Translate types.
	createTypeMapping(program, coreProg, ctx.link)
	createTypes(program, coreProg, ctx)

	// For each type, generate declarations that cannot be derived from the flow IR.
	structToDecls := map[*core.StructType][]ast.Decl{}
	for _, s := range coreProg.Structures {
		structToDecls[s] = generateTreeForStruct(s, ctx)
	}

	packageDecls := make([]*ast.PackageAST, len(program.Packages))
	for i, dubPkg := range program.Packages {
		p := packages[i]
		leaf := p.Path[len(p.Path)-1]

		files := []*ast.FileAST{
			generateGoFile(leaf, dubPkg, structToDecls, ctx),
		}
		if generate_tests && len(dubPkg.Tests) != 0 {
			files = append(files, GenerateTests(leaf, dubPkg.Tests, ctx))
		}
		packageDecls[i] = &ast.PackageAST{
			Files: files,
			P:     p,
		}
	}

	return &ast.ProgramAST{
		Builtins: ctx.index,
		Packages: packageDecls,
	}
}

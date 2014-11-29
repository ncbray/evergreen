package golang

import (
	"evergreen/base"
	"evergreen/dub/flow"
	core "evergreen/dub/tree"
	dstflow "evergreen/go/flow"
	"evergreen/go/transform"
	ast "evergreen/go/tree"
	"evergreen/io"
	"fmt"
	"path/filepath"
)

type DubToGoContext struct {
	index *ast.BuiltinTypeIndex
	state *ast.StructType
	graph *ast.StructType
	t     *ast.StructType
	link  DubToGoLinker
}

func GenerateGoFunc(f *flow.LLFunc, ctx *DubToGoContext) ast.Decl {
	flowDecl := translateFlow(f, ctx)

	if false {
		dot := base.GraphToDot(flowDecl.CFG, &dstflow.DotStyler{Ops: flowDecl.Ops})
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

func DeclForType(t ast.GoType, ctx *DubToGoContext) ast.Decl {
	switch t := t.(type) {
	case *ast.TypeDefType:
		return &ast.TypeDefDecl{
			Name: t.Name,
			Type: ast.RefForType(t.Type),
			T:    t,
		}
	case *ast.StructType:
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
	case *ast.InterfaceType:
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

func GenerateScopeHelpers(s *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	ref := DeclForType(ctx.link.GetType(s, REF), ctx)

	noRef := &ast.VarDecl{
		Name: "No" + s.Name.Text,
		Type: ctx.link.TypeRef(s, REF),
		Expr: &ast.UnaryExpr{
			Op: "^",
			Expr: &ast.TypeCoerce{
				Type: ctx.link.TypeRef(s, REF),
				Expr: &ast.IntLiteral{Value: 0},
			},
		},
		Const: true,
	}

	scope := DeclForType(ctx.link.GetType(s, SCOPE), ctx)

	decls = append(decls, ref, noRef, scope)
	return decls
}

func GenerateGoStruct(s *core.StructType, ctx *DubToGoContext, decls []ast.Decl) []ast.Decl {
	if s.IsParent {
		if s.Scoped {
			panic(s.Name)
		}
		if len(s.Fields) != 0 {
			panic(s.Name)
		}
		decls = append(decls, DeclForType(ctx.link.GetType(s, STRUCT), ctx))
	} else {
		if s.Scoped {
			decls = GenerateScopeHelpers(s, ctx, decls)
		}
		decls = append(decls, DeclForType(ctx.link.GetType(s, STRUCT), ctx))
		decls = addTags(s, s.Implements, ctx, decls)
	}
	return decls
}

func externParserRuntime() (*ast.PackageAST, *ast.StructType) {
	stateT := &ast.StructType{
		Name: "State",
	}
	state := &ast.StructDecl{
		Name: "State",
		T:    stateT,
	}
	pkg := &ast.PackageAST{
		Extern: true,
		Path:   []string{"evergreen", "dub", "runtime"},
		Files: []*ast.FileAST{
			&ast.FileAST{
				Decls: []ast.Decl{
					state,
				},
			},
		},
	}
	return pkg, stateT
}

func externTestingPackage() (*ast.PackageAST, *ast.StructType) {
	tT := &ast.StructType{
		Name: "T",
	}
	t := &ast.StructDecl{
		Name: "T",
		T:    tT,
	}
	pkg := &ast.PackageAST{
		Extern: true,
		Path:   []string{"testing"},
		Files: []*ast.FileAST{
			&ast.FileAST{
				Decls: []ast.Decl{
					t,
				},
			},
		},
	}
	return pkg, tT
}

func externGraph() (*ast.PackageAST, *ast.StructType) {
	graphT := &ast.StructType{
		Name: "Graph",
	}
	graph := &ast.StructDecl{
		Name: "Graph",
		T:    graphT,
	}
	pkg := &ast.PackageAST{
		Extern: true,
		Path:   []string{"evergreen", "base"},
		Files: []*ast.FileAST{
			&ast.FileAST{
				Decls: []ast.Decl{
					graph,
				},
			},
		},
	}
	return pkg, graphT
}

func generateGoFile(package_name string, dubPkg *flow.DubPackage, ctx *DubToGoContext) *ast.FileAST {
	imports := []*ast.Import{}

	decls := []ast.Decl{}
	for _, t := range dubPkg.Structs {
		decls = GenerateGoStruct(t, ctx, decls)
	}
	for _, f := range dubPkg.Funcs {
		decls = append(decls, GenerateGoFunc(f, ctx))
	}

	file := &ast.FileAST{
		Name:    "generated_dub.go",
		Package: package_name,
		Imports: imports,
		Decls:   decls,
	}
	return file
}

func GenerateGo(program []*flow.DubPackage, root string, generate_tests bool) *ast.ProgramAST {
	link := makeLinker()

	packages := []*ast.PackageAST{}

	index := makeBuiltinTypes()

	pkg, state := externParserRuntime()
	packages = append(packages, pkg)

	pkg, graph := externGraph()
	packages = append(packages, pkg)

	pkg, t := externTestingPackage()
	packages = append(packages, pkg)

	ctx := &DubToGoContext{
		index: index,
		state: state,
		graph: graph,
		t:     t,
		link:  link,
	}

	createTypeMapping(program, ctx.link)
	createTypes(program, ctx)

	for _, dubPkg := range program {
		path := []string{root}
		path = append(path, dubPkg.Path...)
		leaf := path[len(path)-1]

		files := []*ast.FileAST{}
		files = append(files, generateGoFile(leaf, dubPkg, ctx))

		if generate_tests && len(dubPkg.Tests) != 0 {
			files = append(files, GenerateTests(leaf, dubPkg.Tests, ctx))
		}
		packages = append(packages, &ast.PackageAST{
			Path:  path,
			Files: files,
		})
	}

	return &ast.ProgramAST{
		Builtins: index,
		Packages: packages,
	}
}

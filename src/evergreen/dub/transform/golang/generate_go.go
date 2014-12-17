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
	core  *core.CoreProgram
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

func declForType(t dstcore.GoType) ast.Decl {
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
		decls = append(decls, declForType(ctx.link.GetType(s, STRUCT)))
	} else {
		if s.Scoped {
			ref := declForType(ctx.link.GetType(s, REF))
			scope := declForType(ctx.link.GetType(s, SCOPE))
			decls = append(decls, ref, scope)
		}
		decls = append(decls, declForType(ctx.link.GetType(s, STRUCT)))
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

func generateGoFile(dubPkg *flow.DubPackage, auxDeclsForStruct map[*core.StructType][]ast.Decl, flowFuncs []*dstflow.LLFunc, file *ast.FileAST, ctx *DubToGoContext) {
	file.Name = "generated_dub.go"

	for _, t := range dubPkg.Structs {
		file.Decls = generateGoStruct(t, ctx, file.Decls)
		more, _ := auxDeclsForStruct[t]
		file.Decls = append(file.Decls, more...)
	}
	for _, f := range dubPkg.Funcs {
		file.Decls = append(file.Decls, transform.RetreeFunc(flowFuncs[f.F]))
	}
}

func createFuncs(program *flow.DubProgram, coreProg *core.CoreProgram, ctx *DubToGoContext) []*dstflow.LLFunc {
	flowFuncs := make([]*dstflow.LLFunc, coreProg.Function_Scope.Len())

	// TODO iterate over Dub funcs directly.
	for _, p := range program.Packages {
		for _, f := range p.Funcs {
			flowFuncs[f.F] = translateFlow(f, ctx)
		}
	}
	return flowFuncs
}

func dumpFuncs(flowFuncs []*dstflow.LLFunc) {
	for _, f := range flowFuncs {
		dot := graph.GraphToDot(f.CFG, &dstflow.DotStyler{Ops: f.Ops})
		parts := []string{"output", "translate"}
		parts = append(parts, fmt.Sprintf("%s.svg", f.Name))
		outfile := filepath.Join(parts...)
		io.WriteDot(dot, outfile)
	}
}

type TreeBypass struct {
	DeclsForStruct map[*core.StructType][]ast.Decl
	Tests          []*ast.FileAST
}

func pathLeaf(path []string) string {
	return path[len(path)-1]
}

func GenerateGo(status compiler.PassStatus, program *flow.DubProgram, coreProg *core.CoreProgram, root string, generate_tests bool, dump bool) *ast.ProgramAST {
	status.Begin()
	defer status.End()

	ctx := &DubToGoContext{
		index: makeBuiltinTypes(),
		state: externParserRuntime(),
		graph: externGraph(),
		t:     externTesting(),
		link:  makeLinker(),
		core:  coreProg,
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

	// Translate functions.
	flowFuncs := createFuncs(program, coreProg, ctx)
	if dump {
		dumpFuncs(flowFuncs)
	}

	bypass := generateTreeBypass(program, coreProg, generate_tests, ctx)
	return generateTree(packages, flowFuncs, bypass, program, ctx)
}

func generateTreeBypass(program *flow.DubProgram, coreProg *core.CoreProgram, generate_tests bool, ctx *DubToGoContext) *TreeBypass {
	bypass := &TreeBypass{
		DeclsForStruct: map[*core.StructType][]ast.Decl{},
		Tests:          make([]*ast.FileAST, len(program.Packages)),
	}

	// For each type, generate declarations that cannot be derived from the flow IR.
	for _, s := range coreProg.Structures {
		bypass.DeclsForStruct[s] = generateTreeForStruct(s, ctx)
	}

	// For each package, generate tests that cannot be derived from the flow IR
	if generate_tests {
		for i, dubPkg := range program.Packages {
			if len(dubPkg.Tests) != 0 {
				bypass.Tests[i] = GenerateTests(pathLeaf(dubPkg.Path), dubPkg.Tests, ctx)
			}
		}
	}
	return bypass
}

func generateTree(packages []*dstcore.Package, flowFuncs []*dstflow.LLFunc, bypass *TreeBypass, program *flow.DubProgram, ctx *DubToGoContext) *ast.ProgramAST {
	packageDecls := make([]*ast.PackageAST, len(packages))
	fileDecls := make([]*ast.FileAST, len(packages))
	for i, p := range packages {
		leaf := pathLeaf(p.Path)

		file := &ast.FileAST{
			Package: leaf,
			Imports: []*ast.Import{},
		}
		fileDecls[i] = file
		pkg := &ast.PackageAST{
			Files: []*ast.FileAST{file},
			P:     p,
		}
		packageDecls[i] = pkg
	}

	for i, dubPkg := range program.Packages {
		generateGoFile(dubPkg, bypass.DeclsForStruct, flowFuncs, fileDecls[i], ctx)
		if bypass.Tests[i] != nil {
			packageDecls[i].Files = append(packageDecls[i].Files, bypass.Tests[i])
		}
	}

	return &ast.ProgramAST{
		Builtins: ctx.index,
		Packages: packageDecls,
	}
}

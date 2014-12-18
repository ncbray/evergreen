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

func generateTreeForStruct(s *core.StructType, bypass *TreeBypass, ctx *DubToGoContext) {
	if !s.IsParent {
		if s.Scoped {
			bypass.DeclsForStruct[ctx.link.GetType(s, REF)] = []ast.Decl{
				&ast.VarDecl{
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
				},
			}
		}
		bypass.DeclsForStruct[ctx.link.GetType(s, STRUCT)] = addTags(s, s.Implements, ctx, []ast.Decl{})
	}
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

func generateGoFile(dubPkg *flow.DubPackage, auxDeclsForStruct map[dstcore.GoType][]ast.Decl, flowFuncs []*dstflow.LLFunc, types []dstcore.GoType, file *ast.FileAST) {
	file.Name = "generated_dub.go"

	for _, t := range types {
		file.Decls = append(file.Decls, declForType(t))
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
	DeclsForStruct map[dstcore.GoType][]ast.Decl
	Tests          []*ast.FileAST
}

func pathLeaf(path []string) string {
	return path[len(path)-1]
}

func getPackage(t dstcore.GoType) *dstcore.Package {
	switch t := t.(type) {
	case *dstcore.StructType:
		return t.Package
	case *dstcore.InterfaceType:
		return t.Package
	case *dstcore.TypeDefType:
		return t.Package
	default:
		panic(t)
	}

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
		path := append([]string{root}, dubPkg.Path...)
		packages[i] = &dstcore.Package{
			Path:  path,
			Index: i,
		}
	}

	// Translate types.
	types := createTypeMapping(program, coreProg, packages, ctx.link)
	createTypes(program, coreProg, ctx)

	// Bucket types for each package.
	packageTypes := make([][]dstcore.GoType, len(program.Packages))
	for _, t := range types {
		pIndex := getPackage(t).Index
		packageTypes[pIndex] = append(packageTypes[pIndex], t)
	}

	// Translate functions.
	flowFuncs := createFuncs(program, coreProg, ctx)
	if dump {
		dumpFuncs(flowFuncs)
	}

	bypass := generateTreeBypass(program, coreProg, generate_tests, ctx)
	return generateTree(packages, flowFuncs, packageTypes, bypass, program, ctx)
}

func generateTreeBypass(program *flow.DubProgram, coreProg *core.CoreProgram, generate_tests bool, ctx *DubToGoContext) *TreeBypass {
	bypass := &TreeBypass{
		DeclsForStruct: map[dstcore.GoType][]ast.Decl{},
		Tests:          make([]*ast.FileAST, len(program.Packages)),
	}

	// For each type, generate declarations that cannot be derived from the flow IR.
	for _, s := range coreProg.Structures {
		generateTreeForStruct(s, bypass, ctx)
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

func generateTree(packages []*dstcore.Package, flowFuncs []*dstflow.LLFunc, packageTypes [][]dstcore.GoType, bypass *TreeBypass, program *flow.DubProgram, ctx *DubToGoContext) *ast.ProgramAST {
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
		generateGoFile(dubPkg, bypass.DeclsForStruct, flowFuncs, packageTypes[i], fileDecls[i])
		if bypass.Tests[i] != nil {
			packageDecls[i].Files = append(packageDecls[i].Files, bypass.Tests[i])
		}
	}

	return &ast.ProgramAST{
		Builtins: ctx.index,
		Packages: packageDecls,
	}
}

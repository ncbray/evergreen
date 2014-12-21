package golang

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"evergreen/dub/flow"
	dstcore "evergreen/go/core"
	dstflow "evergreen/go/flow"
	"evergreen/go/transform"
	ast "evergreen/go/tree"
)

type DubToGoContext struct {
	index *dstcore.BuiltinTypeIndex
	state *dstcore.StructType
	graph *dstcore.StructType
	t     *dstcore.StructType
	link  DubToGoLinker
	core  *core.CoreProgram
}

func generateTreeForStruct(s *core.StructType, bypass *transform.TreeBypass, ctx *DubToGoContext) {
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

func createFuncs(program *flow.DubProgram, coreProg *core.CoreProgram, packages []*dstcore.Package, ctx *DubToGoContext) []*dstflow.LLFunc {
	flowFuncs := make([]*dstflow.LLFunc, coreProg.Function_Scope.Len())

	// TODO iterate over Dub funcs directly.
	for i, p := range program.Packages {
		dstPkg := packages[i]
		for _, f := range p.Funcs {
			dstF := translateFlow(f, ctx)
			dstF.Package = dstPkg
			flowFuncs[f.F] = dstF
		}
	}
	return flowFuncs
}

func pathLeaf(path []string) string {
	return path[len(path)-1]
}

func GenerateGo(status compiler.PassStatus, program *flow.DubProgram, coreProg *core.CoreProgram, root string, generate_tests bool) (*dstflow.FlowProgram, *transform.TreeBypass) {
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

	// Translate functions.
	flowFuncs := createFuncs(program, coreProg, packages, ctx)

	flowFuncs = append(flowFuncs, createTags(program, coreProg, packages, ctx)...)

	bypass := generateTreeBypass(program, coreProg, generate_tests, ctx)

	flowProg := &dstflow.FlowProgram{
		Packages:  packages,
		Types:     types,
		Functions: flowFuncs,
		Builtins:  ctx.index,
	}

	return flowProg, bypass
}

func generateTreeBypass(program *flow.DubProgram, coreProg *core.CoreProgram, generate_tests bool, ctx *DubToGoContext) *transform.TreeBypass {
	bypass := &transform.TreeBypass{
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

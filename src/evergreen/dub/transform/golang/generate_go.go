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
	index   *dstcore.BuiltinTypeIndex
	state   *dstcore.StructType
	graph   *dstcore.StructType
	t       *dstcore.StructType
	link    DubToGoLinker
	core    *core.CoreProgram
	dstCore *dstcore.CoreProgram
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

func externParserRuntime(coreProg *dstcore.CoreProgram) *dstcore.StructType {
	p := coreProg.Package_Scope.Register(&dstcore.Package{
		Extern: true,
		Path:   []string{"evergreen", "dub", "runtime"},
	})
	stateT := &dstcore.StructType{
		Name:    "State",
		Package: p,
	}
	return stateT
}

func externTesting(coreProg *dstcore.CoreProgram) *dstcore.StructType {
	p := coreProg.Package_Scope.Register(&dstcore.Package{
		Extern: true,
		Path:   []string{"testing"},
	})
	tT := &dstcore.StructType{
		Name:    "T",
		Package: p,
	}
	return tT
}

func externGraph(coreProg *dstcore.CoreProgram) *dstcore.StructType {
	p := coreProg.Package_Scope.Register(&dstcore.Package{
		Extern: true,
		Path:   []string{"evergreen", "graph"},
	})
	graphT := &dstcore.StructType{
		Name:    "Graph",
		Package: p,
	}
	return graphT
}

func createFuncs(dubCoreProg *core.CoreProgram, dubFlowProg *flow.DubProgram, goCoreProg *dstcore.CoreProgram, goFlowProg *dstflow.FlowProgram, packages []dstcore.Package_Ref, ctx *DubToGoContext) []dstflow.FlowFunc_Ref {
	flowFuncs := make([]dstflow.FlowFunc_Ref, dubCoreProg.Function_Scope.Len())

	// TODO iterate over Dub funcs directly.
	for i, p := range dubFlowProg.Packages {
		dstPkg := packages[i]
		for _, f := range p.Funcs {
			dstCoreFunc, dstFlowFunc := translateFlow(f, ctx)
			dstCoreFunc.Package = dstPkg
			goCoreProg.Function_Scope.Register(dstCoreFunc)
			flowFuncs[f.F] = goFlowProg.FlowFunc_Scope.Register(dstFlowFunc)
		}
	}
	return flowFuncs
}

func pathLeaf(path []string) string {
	return path[len(path)-1]
}

func GenerateGo(status compiler.PassStatus, program *flow.DubProgram, coreProg *core.CoreProgram, root string, generate_tests bool) (*dstflow.FlowProgram, *dstcore.CoreProgram, *transform.TreeBypass) {
	status.Begin()
	defer status.End()

	dstCoreProg := &dstcore.CoreProgram{
		Package_Scope:  &dstcore.Package_Scope{},
		Function_Scope: &dstcore.Function_Scope{},
	}

	// Translate package identities.
	packages := make([]dstcore.Package_Ref, len(program.Packages))
	for i, dubPkg := range program.Packages {
		path := append([]string{root}, dubPkg.Path...)
		packages[i] = dstCoreProg.Package_Scope.Register(&dstcore.Package{
			Path: path,
		})
	}

	ctx := &DubToGoContext{
		index:   dstcore.MakeBuiltinTypeIndex(),
		state:   externParserRuntime(dstCoreProg),
		graph:   externGraph(dstCoreProg),
		t:       externTesting(dstCoreProg),
		link:    makeLinker(),
		core:    coreProg,
		dstCore: dstCoreProg,
	}

	// Translate types.
	types := createTypeMapping(program, coreProg, packages, ctx.link)
	createTypes(program, coreProg, ctx)

	flowProg := &dstflow.FlowProgram{
		Packages:       packages,
		Types:          types,
		Builtins:       ctx.index,
		FlowFunc_Scope: &dstflow.FlowFunc_Scope{},
	}

	// Translate functions.
	createFuncs(coreProg, program, dstCoreProg, flowProg, packages, ctx)
	createTags(coreProg, program, dstCoreProg, flowProg, packages, ctx)

	bypass := generateTreeBypass(program, coreProg, generate_tests, ctx)

	return flowProg, dstCoreProg, bypass
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

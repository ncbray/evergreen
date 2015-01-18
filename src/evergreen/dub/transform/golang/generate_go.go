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
	index       *dstcore.BuiltinTypeIndex
	state       *dstcore.StructType
	graph       *dstcore.StructType
	t           *dstcore.StructType
	link        DubToGoLinker
	core        *core.CoreProgram
	dstCore     *dstcore.CoreProgram
	functionMap []*dstcore.Function
}

func makeExterns(goCoreProg *dstcore.CoreProgram, ctx *DubToGoContext) {
	runtimePkg := goCoreProg.Package_Scope.Register(&dstcore.Package{
		Extern: true,
		Path:   []string{"evergreen", "dub", "runtime"},
	})
	ctx.state = &dstcore.StructType{
		Name:    "State",
		Package: runtimePkg,
	}

	graphPkg := goCoreProg.Package_Scope.Register(&dstcore.Package{
		Extern: true,
		Path:   []string{"evergreen", "graph"},
	})
	ctx.graph = &dstcore.StructType{
		Name:    "Graph",
		Package: graphPkg,
	}

	testingPkg := goCoreProg.Package_Scope.Register(&dstcore.Package{
		Extern: true,
		Path:   []string{"testing"},
	})
	ctx.t = &dstcore.StructType{
		Name:    "T",
		Package: testingPkg,
	}
}

func createFuncs(dubCoreProg *core.CoreProgram, dubFlowProg *flow.DubProgram, goCoreProg *dstcore.CoreProgram, goFlowProg *dstflow.FlowProgram, packages []*dstcore.Package, ctx *DubToGoContext) []*dstflow.FlowFunc {
	n := dubCoreProg.Function_Scope.Len()
	coreMap := make([]*dstcore.Function, n)
	for i := 0; i < n; i++ {
		f := dubCoreProg.Function_Scope.Get(core.Function_Ref(i))
		coreMap[i] = goCoreProg.Function_Scope.Register(&dstcore.Function{
			Name:    f.Name,
			Package: nil,
		})
	}
	ctx.functionMap = coreMap

	flowFuncs := make([]*dstflow.FlowFunc, n)

	// TODO iterate over Dub funcs directly.
	for i, p := range dubFlowProg.Packages {
		dstPkg := packages[i]
		for _, f := range p.Funcs {
			dstFlowFunc := translateFlow(f, ctx)
			flowFuncs[f.F.Index] = goFlowProg.FlowFunc_Scope.Register(dstFlowFunc)
			dstcore.InsertFunctionIntoPackage(goCoreProg, dstPkg, dstFlowFunc.Function)
		}
	}
	return flowFuncs
}

func pathLeaf(path []string) string {
	return path[len(path)-1]
}

func GenerateGo(status compiler.PassStatus, program *flow.DubProgram, coreProg *core.CoreProgram, rootPackage []string, generate_tests bool) (*dstflow.FlowProgram, *dstcore.CoreProgram, *transform.TreeBypass) {
	status.Begin()
	defer status.End()

	dstCoreProg := &dstcore.CoreProgram{
		Package_Scope:  &dstcore.Package_Scope{},
		Function_Scope: &dstcore.Function_Scope{},
	}

	// Translate package identities.
	packages := make([]*dstcore.Package, len(program.Packages))
	for i, dubPkg := range program.Packages {
		path := append(rootPackage, dubPkg.Path...)
		packages[i] = dstCoreProg.Package_Scope.Register(&dstcore.Package{
			Path: path,
		})
	}

	ctx := &DubToGoContext{
		index:   dstcore.MakeBuiltinTypeIndex(),
		link:    makeLinker(),
		core:    coreProg,
		dstCore: dstCoreProg,
	}
	makeExterns(dstCoreProg, ctx)

	// Translate types.
	types := createTypeMapping(program, coreProg, packages, ctx.link)
	createTypes(program, coreProg, ctx)

	flowProg := &dstflow.FlowProgram{
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
		Tests: make([]*ast.FileAST, len(program.Packages)),
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

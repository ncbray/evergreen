package tree

import (
	"evergreen/compiler"
	"evergreen/dub/core"
	"fmt"
	"strings"
)

type semanticScope struct {
	parent *semanticScope
	locals map[string]*LocalInfo
}

func (scope *semanticScope) localInfo(name string) (*LocalInfo, bool) {
	for scope != nil {
		info, ok := scope.locals[name]
		if ok {
			return info, true
		}
		scope = scope.parent
	}
	return nil, false
}

func childScope(scope *semanticScope) *semanticScope {
	return &semanticScope{
		parent: scope,
		locals: map[string]*LocalInfo{},
	}
}

var unresolvedType core.DubType = nil

var unresolvedScalar = []core.DubType{unresolvedType}

func TypeMatches(actual core.DubType, expected core.DubType, exact bool) bool {
	if actual == unresolvedType || expected == unresolvedType {
		return true
	}

	switch actual := actual.(type) {
	case *core.StructType:
		other, ok := expected.(*core.StructType)
		if !ok {
			return false
		}
		if exact {
			return actual == other
		}
		current := actual
		for current != nil {
			if current == other {
				return true
			}
			current = current.Implements
		}
		return false
	case *core.NilType:
		_, ok := expected.(*core.StructType)
		return ok
	case *core.BuiltinType:
		other, ok := expected.(*core.BuiltinType)
		if !ok {
			return false
		}
		return actual.Name == other.Name
	case *core.ListType:
		other, ok := expected.(*core.ListType)
		if !ok {
			return false
		}
		return TypeMatches(actual.Type, other.Type, true)
	default:
		panic(actual)
	}
}

func IsDiscard(name string) bool {
	return name == "_"
}

func semanticTargetPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, t core.DubType, define bool, scope *semanticScope) {
	switch expr := expr.(type) {
	case *NameRef:
		name := expr.Name.Text
		if IsDiscard(name) {
			expr.Local = nil
			return
		}
		var info *LocalInfo
		var exists bool
		if define {
			_, exists = scope.localInfo(expr.Name.Text)
			if exists {
				ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				return
			}
			info = decl.LocalInfo_Scope.Register(&LocalInfo{Name: name, T: t})
			scope.locals[expr.Name.Text] = info
		} else {
			info, exists = scope.localInfo(name)
			if !exists {
				ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Tried to assign to unknown variable %#v", name))
				return
			}
			// TODO type check
		}
		expr.Local = info
	default:
		panic(expr)
	}
}

func scalarSemanticExprPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, scope *semanticScope) core.DubType {
	types := semanticExprPass(ctx, decl, expr, scope)
	if len(types) != 1 {
		ctx.Status.GlobalError(fmt.Sprintf("expected a single value, got %d instead", len(types)))
		return unresolvedType
	}
	return types[0]
}

func scalarReturn(t core.DubType) []core.DubType {
	return []core.DubType{t}
}

// TODO rewrite binary op to explicitly promote types.
func binaryOpType(ctx *semanticPassContext, lt *core.BuiltinType, op string, rt *core.BuiltinType) (*core.BuiltinType, bool) {
	builtins := ctx.Program.Index
	switch op {
	case "+", "-", "*", "/":
		switch lt {
		case builtins.Int:
			switch rt {
			case builtins.Int:
				return builtins.Int, true
			}
		case builtins.Float32:
			switch rt {
			case builtins.Float32:
				return builtins.Float32, true
			}
		}
	case "<", "<=", ">", ">=", "==", "!=":
		switch lt {
		case builtins.Int:
			switch rt {
			case builtins.Int:
				return builtins.Bool, true
			}
		case builtins.Float32:
			switch rt {
			case builtins.Float32:
				return builtins.Bool, true
			}
		}
	}
	return nil, false
}

func resolve(ctx *semanticPassContext, name string) (namedElement, bool) {
	result, ok := ctx.Module.Namespace[name]
	if ok {
		return result, true
	}
	result, ok = ctx.Program.Namespace[name]
	return result, ok
}

func semanticExprPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, scope *semanticScope) []core.DubType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(ctx, decl, expr.Block, scope)
		return nil
	case *Choice:
		for _, block := range expr.Blocks {
			semanticBlockPass(ctx, decl, block, childScope(scope))
		}
		return nil
	case *Optional:
		semanticBlockPass(ctx, decl, expr.Block, scope)
		return nil
	case *If:
		semanticExprPass(ctx, decl, expr.Expr, scope)
		// TODO check condition type
		semanticBlockPass(ctx, decl, expr.Block, childScope(scope))
		semanticBlockPass(ctx, decl, expr.Else, childScope(scope))
		return nil
	case *BinaryOp:
		l := scalarSemanticExprPass(ctx, decl, expr.Left, scope)
		r := scalarSemanticExprPass(ctx, decl, expr.Right, scope)
		if l == nil || r == nil {
			return unresolvedScalar
		}
		lt, ok := l.(*core.BuiltinType)
		if !ok {
			panic(l)
		}
		rt, ok := r.(*core.BuiltinType)
		if !ok {
			panic(r)
		}
		t, ok := binaryOpType(ctx, lt, expr.Op, rt)
		if ok {
			expr.T = t
		} else {
			expr.T = unresolvedType
			sig := fmt.Sprintf("%s%s%s", lt.Name, expr.Op, rt.Name)
			ctx.Status.GlobalError(fmt.Sprintf("unsupported binary op %s", sig))
		}
		return scalarReturn(expr.T)
	case *NameRef:
		name := expr.Name.Text
		info, found := scope.localInfo(name)
		if !found {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return scalarReturn(unresolvedType)
		}
		expr.Local = info
		return scalarReturn(info.T)
	case *Assign:
		var t []core.DubType
		if expr.Expr != nil {
			t = semanticExprPass(ctx, decl, expr.Expr, scope)
		}
		if expr.Type != nil {
			t = scalarReturn(semanticTypePass(ctx, expr.Type))
		}
		if len(expr.Targets) != len(t) {
			ctx.Status.GlobalError(fmt.Sprintf("Expected %d values but got %d", len(expr.Targets), len(t)))
			t = make([]core.DubType, len(expr.Targets))
			for i, _ := range expr.Targets {
				t[i] = unresolvedType
			}
		}
		for i, target := range expr.Targets {
			semanticTargetPass(ctx, decl, target, t[i], expr.Define, scope)
		}
		return t
	case *Slice:
		semanticBlockPass(ctx, decl, expr.Block, scope)
		return scalarReturn(ctx.Program.Index.String)
	case *StringMatch:
		return scalarReturn(ctx.Program.Index.String)
	case *RuneMatch:
		return scalarReturn(ctx.Program.Index.Rune)
	case *RuneLiteral:
		return scalarReturn(ctx.Program.Index.Rune)
	case *StringLiteral:
		return scalarReturn(ctx.Program.Index.String)
	case *IntLiteral:
		return scalarReturn(ctx.Program.Index.Int)
	case *Float32Literal:
		return scalarReturn(ctx.Program.Index.Float32)
	case *BoolLiteral:
		return scalarReturn(ctx.Program.Index.Bool)
	case *NilLiteral:
		return scalarReturn(ctx.Program.Index.Nil)
	case *Return:
		if len(decl.ReturnTypes) != len(expr.Exprs) {
			ctx.Status.GlobalError(fmt.Sprintf("wrong number of return types: %d vs. %d", len(expr.Exprs), len(decl.ReturnTypes)))
		}
		for i, e := range expr.Exprs {
			at := scalarSemanticExprPass(ctx, decl, e, scope)
			if i < len(decl.ReturnTypes) {
				et := ResolveType(decl.ReturnTypes[i])
				if !TypeMatches(at, et, false) {
					ctx.Status.GlobalError(fmt.Sprintf("return: %s vs. %s", core.TypeName(at), core.TypeName(et)))
				}

			}
		}
		return nil
	case *Fail:
		return nil
	case *Call:
		name := expr.Name.Text
		fd, ok := resolve(ctx, name)
		if !ok {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return unresolvedScalar
		}
		f, ok := AsFunc(fd)
		if !ok {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("%#v is not callable", name))
			return unresolvedScalar
		}
		args := make([]core.DubType, len(expr.Args))
		for i, e := range expr.Args {
			args[i] = scalarSemanticExprPass(ctx, decl, e, scope)
		}
		types := ReturnTypes(ctx, f, args)
		expr.Target = f
		expr.T = types
		return types
	case *Construct:
		t := semanticTypePass(ctx, expr.Type)
		st, ok := t.(*core.StructType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			aft := scalarSemanticExprPass(ctx, decl, arg.Expr, scope)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := f.Type
					if !TypeMatches(aft, eft, false) {
						ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("Expected type %s, but got %s", core.TypeName(eft), core.TypeName(aft)))
					}
				} else {
					ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", core.TypeName(t), fn))
				}
			}
		}
		return scalarReturn(t)
	case *ConstructList:
		t := semanticTypePass(ctx, expr.Type)
		lt, ok := t.(*core.ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			at := scalarSemanticExprPass(ctx, decl, arg, scope)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					ctx.Status.GlobalError(fmt.Sprintf("%s vs. %s", core.TypeName(at), core.TypeName(lt.Type)))
				}
			}
		}
		return scalarReturn(t)
	case *Coerce:
		t := semanticTypePass(ctx, expr.Type)
		scalarSemanticExprPass(ctx, decl, expr.Expr, scope)
		// TODO type check
		return scalarReturn(t)
	default:
		panic(expr)
	}
}

func semanticTypePass(ctx *semanticPassContext, node ASTTypeRef) core.DubType {
	switch node := node.(type) {
	case *TypeRef:
		name := node.Name.Text
		d, ok := resolve(ctx, name)
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			node.T = unresolvedType
			return unresolvedType
		}
		t, ok := AsType(d)
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("%#v is not a type", name))
			node.T = unresolvedType
			return unresolvedType
		}
		node.T = t
		return t
	case *QualifiedTypeRef:
		node.T = unresolvedType

		mname := node.Package.Text
		pkg, ok := ctx.Module.Namespace[mname]
		if !ok {
			ctx.Status.LocationError(node.Package.Pos, fmt.Sprintf("Could not resolve name %#v", mname))
			return unresolvedType
		}
		scope, ok := AsPackage(pkg)
		if !ok {
			ctx.Status.LocationError(node.Package.Pos, fmt.Sprintf("%#v is not a package", mname))
			return unresolvedType
		}
		name := node.Name.Text
		d, ok := scope.Namespace[name]
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return unresolvedType
		}
		t, ok := AsType(d)
		if !ok {
			ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("%#v is not a type", name))
			return unresolvedType
		}
		node.T = t
		return t
	case *ListTypeRef:
		t := semanticTypePass(ctx, node.Type)
		// TODO memoize list types
		node.T = &core.ListType{Type: t}
		return node.T
	default:
		panic(node)
	}
}

func refLocation(node ASTTypeRef) int {
	switch node := node.(type) {
	case *TypeRef:
		return node.Name.Pos
	case *ListTypeRef:
		// TODO more precise location
		return refLocation(node.Type)
	default:
		panic(node)
	}
}

func semanticStructTypePass(ctx *semanticPassContext, node ASTTypeRef) *core.StructType {
	t := semanticTypePass(ctx, node)
	if t == unresolvedType {
		return nil
	}
	st, ok := t.(*core.StructType)
	if !ok {
		ctx.Status.LocationError(refLocation(node), "Not a structure type.")
		return nil
	}
	return st
}

func semanticBlockPass(ctx *semanticPassContext, decl *FuncDecl, block []ASTExpr, scope *semanticScope) {
	for _, expr := range block {
		semanticExprPass(ctx, decl, expr, scope)
	}
}

func semanticFuncSignaturePass(ctx *semanticPassContext, decl *FuncDecl) {
	for _, p := range decl.Params {
		semanticTypePass(ctx, p.Type)
	}
	for _, t := range decl.ReturnTypes {
		semanticTypePass(ctx, t)
	}
}

func semanticFuncBodyPass(ctx *semanticPassContext, decl *FuncDecl) {
	scope := childScope(nil)
	for _, p := range decl.Params {
		semanticTargetPass(ctx, decl, p.Name, ResolveType(p.Type), true, scope)
	}
	semanticBlockPass(ctx, decl, decl.Block, scope)
}

func semanticStructPass(ctx *semanticPassContext, decl *StructDecl) {
	t := decl.T
	t.Name = decl.Name.Text
	t.Scoped = decl.Scoped
	t.Contains = make([]*core.StructType, len(decl.Contains))
	for i, c := range decl.Contains {
		t.Contains[i] = semanticStructTypePass(ctx, c)
	}
	if decl.Implements != nil {
		t.Implements = semanticStructTypePass(ctx, decl.Implements)
	}
	t.Fields = make([]*core.FieldType, len(decl.Fields))
	for i, f := range decl.Fields {
		t.Fields[i] = &core.FieldType{
			Name: f.Name.Text,
			Type: semanticTypePass(ctx, f.Type),
		}
	}
}

func semanticDestructurePass(ctx *semanticPassContext, decl *FuncDecl, d Destructure, scope *semanticScope) core.DubType {
	switch d := d.(type) {
	case *DestructureStruct:
		t := semanticTypePass(ctx, d.Type)
		st, ok := t.(*core.StructType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			aft := semanticDestructurePass(ctx, decl, arg.Destructure, scope)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := f.Type
					if !TypeMatches(aft, eft, false) {
						ctx.Status.GlobalError(fmt.Sprintf("%s.%s: %s vs. %s", core.TypeName(t), fn, core.TypeName(aft), core.TypeName(eft)))
					}
				} else {
					ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", core.TypeName(t), fn))
				}
			}
		}
		return t
	case *DestructureList:
		t := semanticTypePass(ctx, d.Type)
		lt, ok := t.(*core.ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			at := semanticDestructurePass(ctx, decl, arg, scope)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					ctx.Status.GlobalError(fmt.Sprintf("%s vs. %s", core.TypeName(at), core.TypeName(lt.Type)))
				}
			}
		}
		return t
	case *DestructureValue:
		return scalarSemanticExprPass(ctx, decl, d.Expr, scope)
	default:
		panic(d)
	}
}

func semanticTestPass(ctx *semanticPassContext, tst *Test) {
	// HACK no real context
	scope := childScope(nil)
	types := semanticExprPass(ctx, nil, tst.Rule, scope)
	if len(types) != 1 {
		panic(types)
	}
	tst.Type = types[0]

	at := semanticDestructurePass(ctx, nil, tst.Destructure, scope)
	if !TypeMatches(at, tst.Type, false) {
		ctx.Status.GlobalError(fmt.Sprintf("destructure %s vs. %s", core.TypeName(at), core.TypeName(tst.Type)))
	}
}

type ProgramScope struct {
	Index     *core.BuiltinTypeIndex
	Namespace map[string]namedElement
}

type namedElement interface {
	isNamedElement()
}

type NamedType struct {
	T core.DubType
}

func (element *NamedType) isNamedElement() {
}

func AsType(node namedElement) (core.DubType, bool) {
	switch node := node.(type) {
	case *NamedType:
		return node.T, true
	default:
		return nil, false
	}
}

type NamedCallable struct {
	Func core.Callable
}

func (element *NamedCallable) isNamedElement() {
}

func AsFunc(node namedElement) (core.Callable, bool) {
	switch node := node.(type) {
	case *NamedCallable:
		return node.Func, true
	default:
		return nil, false
	}
}

type NamedPackage struct {
	Scope *ModuleScope
}

func (element *NamedPackage) isNamedElement() {
}

func AsPackage(node namedElement) (*ModuleScope, bool) {
	switch node := node.(type) {
	case *NamedPackage:
		return node.Scope, true
	default:
		return nil, false
	}
}

type ModuleScope struct {
	Path      []string
	Namespace map[string]namedElement
}

func GetField(node *core.StructType, name string) *core.FieldType {
	for _, decl := range node.Fields {
		if decl.Name == name {
			return decl
		}
	}
	return nil
}

func ResolveType(ref ASTTypeRef) core.DubType {
	switch ref := ref.(type) {
	case *TypeRef:
		return ref.T
	case *QualifiedTypeRef:
		return ref.T
	case *ListTypeRef:
		return ref.T
	default:
		panic(ref)
	}
}

func ReturnTypes(ctx *semanticPassContext, node core.Callable, args []core.DubType) []core.DubType {
	builtins := ctx.Program.Index

	// TODO check argument types
	switch node := node.(type) {
	case *core.Function:
		f := ctx.Functions[node.Index]
		types := make([]core.DubType, len(f.ReturnTypes))
		for i, t := range f.ReturnTypes {
			types[i] = ResolveType(t)
		}
		return types
	case *core.IntrinsicFunction:
		switch node {
		case ctx.Program.Index.Append:
			if len(args) != 2 {
				panic(args)
			}
			return []core.DubType{args[0]}
		case ctx.Program.Index.Position:
			if len(args) != 0 {
				panic(args)
			}
			return []core.DubType{builtins.Int}
		default:
			panic(node)
		}
	default:
		panic(node)
	}
}

func MakeBuiltinTypeIndex() *core.BuiltinTypeIndex {
	return &core.BuiltinTypeIndex{
		String:   &core.BuiltinType{Name: "string"},
		Rune:     &core.BuiltinType{Name: "rune"},
		Int:      &core.BuiltinType{Name: "int"},
		Int64:    &core.BuiltinType{Name: "int64"},
		Float32:  &core.BuiltinType{Name: "float32"},
		Bool:     &core.BuiltinType{Name: "bool"},
		Graph:    &core.BuiltinType{Name: "graph"},
		Nil:      &core.NilType{},
		Append:   &core.IntrinsicFunction{Name: "append"},
		Position: &core.IntrinsicFunction{Name: "position"},
	}
}

func addIntrinsticFunction(f *core.IntrinsicFunction, namespace map[string]namedElement) {
	namespace[f.Name] = &NamedCallable{
		Func: f,
	}
}

func addBuiltinType(t *core.BuiltinType, namespace map[string]namedElement) {
	namespace[t.Name] = &NamedType{
		T: t,
	}
}

func MakeProgramScope(program *Program) *ProgramScope {
	ns := map[string]namedElement{}
	programScope := &ProgramScope{
		Index:     program.Builtins,
		Namespace: ns,
	}
	builtins := program.Builtins

	addBuiltinType(builtins.String, ns)
	addBuiltinType(builtins.Rune, ns)
	addBuiltinType(builtins.Int, ns)
	addBuiltinType(builtins.Int64, ns)
	addBuiltinType(builtins.Float32, ns)
	addBuiltinType(builtins.Bool, ns)
	addBuiltinType(builtins.Graph, ns)

	addIntrinsticFunction(builtins.Append, ns)
	addIntrinsticFunction(builtins.Position, ns)

	return programScope
}

type semanticPassContext struct {
	Program        *ProgramScope
	Module         *ModuleScope
	ModuleContexts []*semanticPassContext
	Status         compiler.PassStatus
	Core           *core.CoreProgram
	Functions      []*FuncDecl
}

func resolveImport(ctx *semanticPassContext, imp *ImportDecl) {
	pos := imp.Path.Pos
	path := imp.Path.Value
	parts := strings.Split(path, "/")

	// HACK O(n^2)
	for _, other := range ctx.ModuleContexts {
		otherPath := other.Module.Path
		found := strings.Join(otherPath, "/") == path
		if found {
			name := parts[len(parts)-1]
			// HACK should use file-local namespace.
			_, exists := ctx.Module.Namespace[name]
			if exists {
				ctx.Status.LocationError(pos, fmt.Sprintf("Tried to redefine %#v", name))
			} else {
				ctx.Module.Namespace[name] = &NamedPackage{Scope: other.Module}
			}
			return
		}
	}
	ctx.Status.LocationError(pos, fmt.Sprintf("cannot find module %#v", path))
}

func indexModule(ctx *semanticPassContext, pkg *Package) {
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			resolveImport(ctx, imp)
		}

		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				name := decl.Name.Text
				_, exists := ctx.Module.Namespace[name]
				if exists {
					ctx.Status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				} else {
					f := &core.Function{
						Name: name,
						File: file.F,
					}
					decl.F = ctx.Core.Function_Scope.Register(f)
					ctx.Functions = append(ctx.Functions, decl)

					ctx.Module.Namespace[name] = &NamedCallable{
						Func: decl.F,
					}
				}
			case *StructDecl:
				name := decl.Name.Text
				_, exists := ctx.Module.Namespace[name]
				if exists {
					ctx.Status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				} else {
					st := &core.StructType{
						File: file.F,
					}
					decl.T = st
					ctx.Core.Structures = append(ctx.Core.Structures, st)
					ctx.Module.Namespace[name] = &NamedType{T: st}
				}
			default:
				panic(decl)
			}
		}
	}
}

func resolveSignatures(ctx *semanticPassContext, pkg *Package) {
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				// Needed for resolving calls in the next step.
				semanticFuncSignaturePass(ctx, decl)
			case *StructDecl:
				// Needed for resolving field reference types.
				semanticStructPass(ctx, decl)
			default:
				panic(decl)
			}
		}
	}
}

func semanticModulePass(ctx *semanticPassContext, pkg *Package) {
	// Resolve the declaration contents.
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				semanticFuncBodyPass(ctx, decl)
			case *StructDecl:
			default:
				panic(decl)
			}
		}
		for _, tst := range file.Tests {
			semanticTestPass(ctx, tst)
		}
	}
}

func SemanticPass(program *Program, status compiler.PassStatus) *core.CoreProgram {
	status.Begin()
	defer status.End()

	programScope := MakeProgramScope(program)
	coreProg := &core.CoreProgram{
		Builtins:       program.Builtins,
		Package_Scope:  &core.Package_Scope{},
		File_Scope:     &core.File_Scope{},
		Function_Scope: &core.Function_Scope{},
	}

	for _, pkg := range program.Packages {
		corePkg := &core.Package{Path: pkg.Path}
		packageRef := coreProg.Package_Scope.Register(corePkg)
		pkg.P = packageRef

		for _, file := range pkg.Files {
			coreFile := &core.File{Name: file.Name, Package: pkg.P}
			fileRef := coreProg.File_Scope.Register(coreFile)
			file.F = fileRef
			corePkg.Files = append(corePkg.Files, fileRef)
		}
	}

	ctxs := make([]*semanticPassContext, len(program.Packages))
	for i, pkg := range program.Packages {
		moduleScope := &ModuleScope{
			Path:      pkg.Path,
			Namespace: map[string]namedElement{},
		}
		ctxs[i] = &semanticPassContext{
			Program:        programScope,
			Module:         moduleScope,
			ModuleContexts: ctxs,
			Status:         status,
			Core:           coreProg,
		}
	}

	for i, pkg := range program.Packages {
		indexModule(ctxs[i], pkg)
	}
	if status.ShouldHalt() {
		return nil
	}
	for i, pkg := range program.Packages {
		resolveSignatures(ctxs[i], pkg)
	}
	if status.ShouldHalt() {
		return nil
	}
	for i, pkg := range program.Packages {
		semanticModulePass(ctxs[i], pkg)
	}
	return coreProg
}

func (scope *LocalInfo_Scope) Get(ref LocalInfo_Ref) *LocalInfo {
	if scope.objects[ref].Index != ref {
		panic(scope.objects[ref].Index)
	}
	return scope.objects[ref]
}

func (scope *LocalInfo_Scope) Register(info *LocalInfo) *LocalInfo {
	info.Index = LocalInfo_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return info
}

func (scope *LocalInfo_Scope) Len() int {
	return len(scope.objects)
}

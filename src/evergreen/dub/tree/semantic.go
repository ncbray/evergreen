package tree

import (
	"evergreen/framework"
	"fmt"
	"strings"
)

type semanticScope struct {
	parent *semanticScope
	locals map[string]LocalInfo_Ref
}

func (scope *semanticScope) localInfo(name string) (LocalInfo_Ref, bool) {
	for scope != nil {
		info, ok := scope.locals[name]
		if ok {
			return info, true
		}
		scope = scope.parent
	}
	return NoLocalInfo, false
}

func childScope(scope *semanticScope) *semanticScope {
	return &semanticScope{parent: scope, locals: map[string]LocalInfo_Ref{}}
}

var unresolvedType DubType = nil

func TypeMatches(actual DubType, expected DubType, exact bool) bool {
	if actual == unresolvedType || expected == unresolvedType {
		return true
	}

	switch actual := actual.(type) {
	case *StructType:
		other, ok := expected.(*StructType)
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
	case *NilType:
		_, ok := expected.(*StructType)
		return ok
	case *BuiltinType:
		other, ok := expected.(*BuiltinType)
		if !ok {
			return false
		}
		return actual.Name == other.Name
	case *ListType:
		other, ok := expected.(*ListType)
		if !ok {
			return false
		}
		return TypeMatches(actual.Type, other.Type, true)
	default:
		panic(actual)
	}
}

func TypeName(t DubType) string {
	switch t := t.(type) {
	case *StructType:
		return t.Name.Text
	case *BuiltinType:
		return t.Name
	case *ListType:
		return fmt.Sprintf("[]%s", TypeName(t.Type))
	default:
		panic(t)
	}
}

func IsDiscard(name string) bool {
	return name == "_"
}

func semanticTargetPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, t DubType, define bool, scope *semanticScope) {
	switch expr := expr.(type) {
	case *NameRef:
		name := expr.Name.Text
		if IsDiscard(name) {
			expr.Local = NoLocalInfo
			return
		}
		var info LocalInfo_Ref
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

func scalarSemanticExprPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, scope *semanticScope) DubType {
	types := semanticExprPass(ctx, decl, expr, scope)
	if len(types) != 1 {
		ctx.Status.Error("expected a single value, got %d instead", len(types))
		return unresolvedType
	}
	return types[0]
}

func scalarReturn(t DubType) []DubType {
	return []DubType{t}
}

func semanticExprPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, scope *semanticScope) []DubType {
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
		return nil
	case *BinaryOp:
		l := scalarSemanticExprPass(ctx, decl, expr.Left, scope)
		r := scalarSemanticExprPass(ctx, decl, expr.Right, scope)
		if l == nil || r == nil {
			return nil
		}
		lt, ok := l.(*BuiltinType)
		if !ok {
			panic(l)
		}
		rt, ok := r.(*BuiltinType)
		if !ok {
			panic(r)
		}
		sig := fmt.Sprintf("%s%s%s", lt.Name, expr.Op, rt.Name)
		t, ok := ctx.Program.BinaryOps[sig].(DubType)
		if !ok {
			panic(sig)
		}
		expr.T = t
		return scalarReturn(t)
	case *NameRef:
		name := expr.Name.Text
		info, found := scope.localInfo(name)
		if !found {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return scalarReturn(unresolvedType)
		}
		expr.Local = info
		return scalarReturn(decl.LocalInfo_Scope.Get(info).T)
	case *Assign:
		var t []DubType
		if expr.Expr != nil {
			t = semanticExprPass(ctx, decl, expr.Expr, scope)
		}
		if expr.Type != nil {
			t = scalarReturn(semanticTypePass(ctx, expr.Type))
		}
		if len(expr.Targets) != len(t) {
			ctx.Status.Error("Expected %d values but got %d", len(expr.Targets), len(t))
			t = make([]DubType, len(expr.Targets))
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
	case *BoolLiteral:
		return scalarReturn(ctx.Program.Index.Bool)
	case *NilLiteral:
		return scalarReturn(ctx.Program.Index.Nil)
	case *Return:
		if len(decl.ReturnTypes) != len(expr.Exprs) {
			ctx.Status.Error("wrong number of return types: %d vs. %d", len(expr.Exprs), len(decl.ReturnTypes))
		}
		for i, e := range expr.Exprs {
			at := scalarSemanticExprPass(ctx, decl, e, scope)
			if i < len(decl.ReturnTypes) {
				et := ResolveType(decl.ReturnTypes[i])
				if !TypeMatches(at, et, false) {
					ctx.Status.Error("return: %s vs. %s", TypeName(at), TypeName(et))
				}

			}
		}
		return nil
	case *Position:
		return scalarReturn(ctx.Program.Index.Int)
	case *Fail:
		return nil
	case *Call:

		name := expr.Name.Text
		// HACK resolve other scopes?
		fd, ok := ctx.Module.Namespace[name]
		if !ok {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return scalarReturn(unresolvedType)
		}
		f, ok := AsFunc(fd)
		if !ok {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("%#v is not callable", name))
			return scalarReturn(unresolvedType)
		}
		for _, e := range expr.Args {
			// TODO check argument types
			scalarSemanticExprPass(ctx, decl, e, scope)
		}
		types := ReturnTypes(f)
		expr.Target = f
		expr.T = types
		return types
	case *Append:
		t := scalarSemanticExprPass(ctx, decl, expr.List, scope)
		scalarSemanticExprPass(ctx, decl, expr.Expr, scope)
		// TODO type check arguments
		expr.T = t
		return scalarReturn(t)
	case *Construct:
		t := semanticTypePass(ctx, expr.Type)
		st, ok := t.(*StructType)
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
						ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("Expected type %s, but got %s", TypeName(eft), TypeName(aft)))
					}
				} else {
					ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", TypeName(t), fn))
				}
			}
		}
		return scalarReturn(t)
	case *ConstructList:
		t := semanticTypePass(ctx, expr.Type)
		lt, ok := t.(*ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			at := scalarSemanticExprPass(ctx, decl, arg, scope)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					ctx.Status.Error("%s vs. %s", TypeName(at), TypeName(lt.Type))
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

func semanticTypePass(ctx *semanticPassContext, node ASTTypeRef) DubType {
	switch node := node.(type) {
	case *TypeRef:
		var t DubType
		name := node.Name.Text
		d, ok := ctx.Module.Namespace[name]
		if ok {
			t, ok = AsType(d)
			if !ok {
				ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("%#v is not a type", name))
				node.T = unresolvedType
				return unresolvedType
			}
		} else {
			t, ok = GetBuiltinType(ctx.Program.Index, name)
			if !ok {
				ctx.Status.LocationError(node.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
				node.T = unresolvedType
				return unresolvedType
			}
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
		node.T = &ListType{Type: t}
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

func semanticStructTypePass(ctx *semanticPassContext, node ASTTypeRef) *StructType {
	t := semanticTypePass(ctx, node)
	if t == unresolvedType {
		return nil
	}
	st, ok := t.(*StructType)
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
	t.Name = decl.Name
	t.Scoped = decl.Scoped
	t.Contains = make([]*StructType, len(decl.Contains))
	for i, c := range decl.Contains {
		t.Contains[i] = semanticStructTypePass(ctx, c)
	}
	if decl.Implements != nil {
		t.Implements = semanticStructTypePass(ctx, decl.Implements)
	}
	t.Fields = make([]*FieldType, len(decl.Fields))
	for i, f := range decl.Fields {
		t.Fields[i] = &FieldType{
			Name: f.Name,
			Type: semanticTypePass(ctx, f.Type),
		}
	}
}

func semanticDestructurePass(ctx *semanticPassContext, decl *FuncDecl, d Destructure, scope *semanticScope) DubType {
	switch d := d.(type) {
	case *DestructureStruct:
		t := semanticTypePass(ctx, d.Type)
		st, ok := t.(*StructType)
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
						ctx.Status.Error("%s.%s: %s vs. %s", TypeName(t), fn, TypeName(aft), TypeName(eft))
					}
				} else {
					ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", TypeName(t), fn))
				}
			}
		}
		return t
	case *DestructureList:
		t := semanticTypePass(ctx, d.Type)
		lt, ok := t.(*ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			at := semanticDestructurePass(ctx, decl, arg, scope)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					ctx.Status.Error("%s vs. %s", TypeName(at), TypeName(lt.Type))
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
		ctx.Status.Error(fmt.Sprintf("destructure %s vs. %s", TypeName(at), TypeName(tst.Type)))
	}
}

type ProgramScope struct {
	Index     *BuiltinTypeIndex
	BinaryOps map[string]DubType
}

type namedElement interface {
	isNamedElement()
}

type NamedType struct {
	T DubType
}

func (element *NamedType) isNamedElement() {
}

func AsType(node namedElement) (DubType, bool) {
	switch node := node.(type) {
	case *NamedType:
		return node.T, true
	default:
		return nil, false
	}
}

type NamedCallable struct {
	Func ASTCallable
}

func (element *NamedCallable) isNamedElement() {
}

func AsFunc(node namedElement) (ASTCallable, bool) {
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

func GetField(node *StructType, name string) *FieldType {
	for _, decl := range node.Fields {
		if decl.Name.Text == name {
			return decl
		}
	}
	return nil
}

func ResolveType(ref ASTTypeRef) DubType {
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

func ReturnTypes(node ASTCallable) []DubType {
	switch node := node.(type) {
	case *FuncDecl:
		types := make([]DubType, len(node.ReturnTypes))
		for i, t := range node.ReturnTypes {
			types[i] = ResolveType(t)
		}
		return types
	default:
		panic(node)
	}
}

func MakeBuiltinTypeIndex() *BuiltinTypeIndex {
	return &BuiltinTypeIndex{
		String: &BuiltinType{"string"},
		Rune:   &BuiltinType{"rune"},
		Int:    &BuiltinType{"int"},
		Int64:  &BuiltinType{"int64"},
		Bool:   &BuiltinType{"bool"},
		Graph:  &BuiltinType{"graph"},
		Nil:    &NilType{},
	}
}

var BuiltinTypeNames = []string{
	"string",
	"rune",
	"int",
	"int64",
	"bool",
	"graph",
}

func GetBuiltinType(index *BuiltinTypeIndex, name string) (DubType, bool) {
	switch name {
	case "string":
		return index.String, true
	case "rune":
		return index.Rune, true
	case "int":
		return index.Int, true
	case "int64":
		return index.Int64, true
	case "bool":
		return index.Bool, true
	case "graph":
		return index.Graph, true
	default:
		return nil, false
	}
}

var binaryOps = []string{
	"int+int:int",
	"int-int:int",
	"int*int:int",
	"int/int:int",
	"int<int:bool",
	"int<=int:bool",
	"int>int:bool",
	"int>=int:bool",
	"int==int:bool",
	"int!=int:bool",
}

func MakeProgramScope(program *Program) *ProgramScope {
	programScope := &ProgramScope{
		BinaryOps: map[string]DubType{},
		Index:     program.Builtins,
	}

	for _, desc := range binaryOps {
		for i := 0; i < len(desc); i++ {
			if desc[i] == ':' {
				expr := desc[:i]
				out := desc[i+1:]
				outT, ok := GetBuiltinType(program.Builtins, out)
				if !ok {
					panic(desc)
				}
				programScope.BinaryOps[expr] = outT
			}
		}
	}

	return programScope
}

type semanticPassContext struct {
	Program        *ProgramScope
	Module         *ModuleScope
	ModuleContexts []*semanticPassContext
	Status         framework.Status
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
					ctx.Module.Namespace[name] = &NamedCallable{Func: decl}
				}
			case *StructDecl:
				name := decl.Name.Text
				_, exists := ctx.Module.Namespace[name]
				if exists {
					ctx.Status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				} else {
					decl.T = &StructType{}
					ctx.Module.Namespace[name] = &NamedType{T: decl.T}
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

func SemanticPass(program *Program, status framework.Status) {
	programScope := MakeProgramScope(program)
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
			Status:         status.CreateChild(),
		}
	}

	for i, pkg := range program.Packages {
		indexModule(ctxs[i], pkg)
	}
	if status.ShouldHalt() {
		return
	}
	for i, pkg := range program.Packages {
		resolveSignatures(ctxs[i], pkg)
	}
	if status.ShouldHalt() {
		return
	}
	for i, pkg := range program.Packages {
		semanticModulePass(ctxs[i], pkg)
	}
}

func (scope *LocalInfo_Scope) Get(ref LocalInfo_Ref) *LocalInfo {
	return scope.objects[ref]
}

func (scope *LocalInfo_Scope) Register(info *LocalInfo) LocalInfo_Ref {
	index := LocalInfo_Ref(len(scope.objects))
	scope.objects = append(scope.objects, info)
	return index
}

func (scope *LocalInfo_Scope) Len() int {
	return len(scope.objects)
}

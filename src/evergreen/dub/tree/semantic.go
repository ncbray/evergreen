package tree

import (
	"evergreen/framework"
	"fmt"
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

var unresolvedType ASTType = nil

func TypeMatches(actual ASTType, expected ASTType, exact bool) bool {
	if actual == unresolvedType || expected == unresolvedType {
		return true
	}

	switch actual := actual.(type) {
	case *StructDecl:
		other, ok := expected.(*StructDecl)
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
			if current.Implements != nil {
				t := ResolveType(current.Implements)
				var ok bool
				current, ok = t.(*StructDecl)
				if !ok {
					panic(t)
				}
			} else {
				current = nil
			}
		}
		return false
	case *NilType:
		_, ok := expected.(*StructDecl)
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

func TypeName(t ASTType) string {
	switch t := t.(type) {
	case *StructDecl:
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

func semanticTargetPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, t ASTType, define bool, scope *semanticScope) {
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

func scalarSemanticExprPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, scope *semanticScope) ASTType {
	types := semanticExprPass(ctx, decl, expr, scope)
	if len(types) != 1 {
		ctx.Status.Error("expected a single value, got %d instead", len(types))
		return unresolvedType
	}
	return types[0]
}

func semanticExprPass(ctx *semanticPassContext, decl *FuncDecl, expr ASTExpr, scope *semanticScope) []ASTType {
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
		t, ok := ctx.Program.BinaryOps[sig].(ASTType)
		if !ok {
			panic(sig)
		}
		expr.T = t
		return []ASTType{t}
	case *NameRef:
		name := expr.Name.Text
		info, found := scope.localInfo(name)
		if !found {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return []ASTType{unresolvedType}
		}
		expr.Local = info
		return []ASTType{decl.LocalInfo_Scope.Get(info).T}
	case *Assign:
		var t []ASTType
		if expr.Expr != nil {
			t = semanticExprPass(ctx, decl, expr.Expr, scope)
		}
		if expr.Type != nil {
			t = []ASTType{semanticTypePass(ctx, expr.Type)}
		}
		if len(expr.Targets) != len(t) {
			ctx.Status.Error("Expected %d values but got %d", len(expr.Targets), len(t))
			t = make([]ASTType, len(expr.Targets))
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
		return []ASTType{ctx.Program.Index.String}
	case *StringMatch:
		return []ASTType{ctx.Program.Index.String}
	case *RuneMatch:
		return []ASTType{ctx.Program.Index.Rune}
	case *RuneLiteral:
		return []ASTType{ctx.Program.Index.Rune}
	case *StringLiteral:
		return []ASTType{ctx.Program.Index.String}
	case *IntLiteral:
		return []ASTType{ctx.Program.Index.Int}
	case *BoolLiteral:
		return []ASTType{ctx.Program.Index.Bool}
	case *NilLiteral:
		return []ASTType{ctx.Program.Index.Nil}
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
		return []ASTType{ctx.Program.Index.Int}
	case *Fail:
		return nil
	case *Call:

		name := expr.Name.Text
		// HACK resolve other scopes?
		fd, ok := ctx.Module.Module[name]
		if !ok {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return []ASTType{unresolvedType}
		}
		f, ok := AsFunc(fd)
		if !ok {
			ctx.Status.LocationError(expr.Name.Pos, fmt.Sprintf("%#v is not callable", name))
			return []ASTType{unresolvedType}
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
		return []ASTType{t}
	case *Construct:
		t := semanticTypePass(ctx, expr.Type)
		st, ok := t.(*StructDecl)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			aft := scalarSemanticExprPass(ctx, decl, arg.Expr, scope)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := ResolveType(f.Type)
					if !TypeMatches(aft, eft, false) {
						ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("Expected type %s, but got %s", TypeName(eft), TypeName(aft)))
					}
				} else {
					ctx.Status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", TypeName(t), fn))
				}
			}
		}
		return []ASTType{t}
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
		return []ASTType{t}
	case *Coerce:
		t := semanticTypePass(ctx, expr.Type)
		scalarSemanticExprPass(ctx, decl, expr.Expr, scope)
		// TODO type check
		return []ASTType{t}
	default:
		panic(expr)
	}
}

func semanticTypePass(ctx *semanticPassContext, node ASTTypeRef) ASTType {
	switch node := node.(type) {
	case *TypeRef:
		name := node.Name.Text
		d, ok := ctx.Module.Module[name]
		if !ok {
			d, ok = GetBuiltinType(ctx.Program.Index, name)
		}
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
	case *ListTypeRef:
		t := semanticTypePass(ctx, node.Type)
		// TODO memoize list types
		node.T = &ListType{Type: t}
		return node.T
	default:
		panic(node)
	}
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
	for _, t := range decl.Contains {
		semanticTypePass(ctx, t)
	}
	if decl.Implements != nil {
		semanticTypePass(ctx, decl.Implements)
	}
	for _, f := range decl.Fields {
		semanticTypePass(ctx, f.Type)
	}
}

func semanticDestructurePass(ctx *semanticPassContext, decl *FuncDecl, d Destructure, scope *semanticScope) ASTType {
	switch d := d.(type) {
	case *DestructureStruct:
		t := semanticTypePass(ctx, d.Type)
		st, ok := t.(*StructDecl)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			aft := semanticDestructurePass(ctx, decl, arg.Destructure, scope)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := ResolveType(f.Type)
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
	BinaryOps map[string]ASTDecl
}

type ModuleScope struct {
	Module map[string]ASTDecl
}

func AsType(node ASTDecl) (ASTType, bool) {
	switch node := node.(type) {
	case *StructDecl:
		return node, true
	case *BuiltinType:
		return node, true
	default:
		return nil, false
	}
}

func AsFunc(node ASTDecl) (ASTCallable, bool) {
	switch node := node.(type) {
	case *FuncDecl:
		return node, true
	default:
		return nil, false
	}
}

func GetField(node *StructDecl, name string) *FieldDecl {
	for _, decl := range node.Fields {
		if decl.Name.Text == name {
			return decl
		}
	}
	return nil
}

func ResolveType(ref ASTTypeRef) ASTType {
	switch ref := ref.(type) {
	case *TypeRef:
		return ref.T
	case *ListTypeRef:
		return ref.T
	default:
		panic(ref)
	}
}

func ReturnTypes(node ASTCallable) []ASTType {
	switch node := node.(type) {
	case *FuncDecl:
		types := make([]ASTType, len(node.ReturnTypes))
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

func GetBuiltinType(index *BuiltinTypeIndex, name string) (ASTDecl, bool) {
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

func MakeProgramScope(index *BuiltinTypeIndex) *ProgramScope {
	glbls := &ProgramScope{
		BinaryOps: map[string]ASTDecl{},
		Index:     index,
	}

	for _, desc := range binaryOps {
		for i := 0; i < len(desc); i++ {
			if desc[i] == ':' {
				expr := desc[:i]
				out := desc[i+1:]
				outT, ok := GetBuiltinType(index, out)
				if !ok {
					panic(desc)
				}
				glbls.BinaryOps[expr] = outT
			}
		}
	}

	return glbls
}

type semanticPassContext struct {
	Program *ProgramScope
	Module  *ModuleScope
	Status  framework.Status
}

func SemanticPass(program *ProgramScope, pkg *Package, status framework.Status) {
	module := &ModuleScope{
		Module: map[string]ASTDecl{},
	}
	ctx := &semanticPassContext{
		Program: program,
		Module:  module,
		Status:  status,
	}

	// Index the package namespace.
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				name := decl.Name.Text
				_, exists := ctx.Module.Module[name]
				if exists {
					ctx.Status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				} else {
					ctx.Module.Module[name] = decl
				}
			case *StructDecl:
				name := decl.Name.Text
				_, exists := ctx.Module.Module[name]
				if exists {
					ctx.Status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				} else {
					ctx.Module.Module[name] = decl
				}
			default:
				panic(decl)
			}
		}
	}
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *FuncDecl:
				// Needed for resolving calls in the next step.
				semanticFuncSignaturePass(ctx, decl)
			case *StructDecl:
				// Needed for resolving field types.
				semanticStructPass(ctx, decl)
			default:
				panic(decl)
			}
		}
	}

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

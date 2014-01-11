package tree

import (
	"evergreen/framework"
	"fmt"
)

type semanticScope struct {
	parent *semanticScope
	locals map[string]int
}

func (scope *semanticScope) localInfo(name string) (int, bool) {
	for scope != nil {
		info, ok := scope.locals[name]
		if ok {
			return info, true
		}
		scope = scope.parent
	}
	return -1, false
}

func childScope(scope *semanticScope) *semanticScope {
	return &semanticScope{parent: scope, locals: map[string]int{}}
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
		} else {
			current := actual
			for current != nil {
				if current == other {
					return true
				}
				t := ResolveType(current.Implements)
				var ok bool
				current, ok = t.(*StructDecl)
				if !ok {
					panic(t)
				}
			}
			return false
		}
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

func semanticTargetPass(decl *FuncDecl, expr ASTExpr, t ASTType, define bool, scope *semanticScope, glbls *ModuleScope, status framework.Status) {
	switch expr := expr.(type) {
	case *NameRef:
		name := expr.Name.Text
		if IsDiscard(name) {
			expr.Info = -1
			return
		}
		var info int
		var exists bool
		if define {
			_, exists = scope.localInfo(expr.Name.Text)
			if exists {
				status.LocationError(expr.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
				return
			}

			info = len(decl.Locals)
			decl.Locals = append(decl.Locals, &LocalInfo{Name: name, T: t})
			scope.locals[expr.Name.Text] = info
		} else {
			info, exists = scope.localInfo(name)
			if !exists {
				status.LocationError(expr.Name.Pos, fmt.Sprintf("Tried to assign to unknown variable %#v", name))
				return
			}
			// TODO type check
		}
		expr.Info = info
	default:
		panic(expr)
	}
}

func scalarSemanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope, glbls *ModuleScope, status framework.Status) ASTType {
	types := semanticExprPass(decl, expr, scope, glbls, status)
	if len(types) != 1 {
		status.Error("expected a single value, got %d instead", len(types))
		return unresolvedType
	}
	return types[0]
}

func semanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope, glbls *ModuleScope, status framework.Status) []ASTType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(decl, expr.Block, scope, glbls, status)
		return nil
	case *Choice:
		for _, block := range expr.Blocks {
			semanticBlockPass(decl, block, childScope(scope), glbls, status)
		}
		return nil
	case *Optional:
		semanticBlockPass(decl, expr.Block, scope, glbls, status)
		return nil
	case *If:
		semanticExprPass(decl, expr.Expr, scope, glbls, status)
		// TODO check condition type
		semanticBlockPass(decl, expr.Block, childScope(scope), glbls, status)
		return nil
	case *BinaryOp:
		l := scalarSemanticExprPass(decl, expr.Left, scope, glbls, status)
		r := scalarSemanticExprPass(decl, expr.Right, scope, glbls, status)
		lt, ok := l.(*BuiltinType)
		if !ok {
			panic(l)
		}
		rt, ok := r.(*BuiltinType)
		if !ok {
			panic(r)
		}
		sig := fmt.Sprintf("%s%s%s", lt.Name, expr.Op, rt.Name)
		t, ok := glbls.BinaryOps[sig]
		if !ok {
			panic(sig)
		}
		expr.T = t
		return []ASTType{t}
	case *NameRef:
		name := expr.Name.Text
		info, found := scope.localInfo(name)
		if !found {
			status.LocationError(expr.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			return []ASTType{unresolvedType}
		}
		expr.Info = info
		return []ASTType{decl.Locals[info].T}
	case *Assign:
		var t []ASTType
		if expr.Expr != nil {
			t = semanticExprPass(decl, expr.Expr, scope, glbls, status)
		}
		if expr.Type != nil {
			t = []ASTType{semanticTypePass(expr.Type, glbls, status)}
		}
		if len(expr.Targets) != len(t) {
			status.Error("Expected %d values but got %d", len(expr.Targets), len(t))
			t = make([]ASTType, len(expr.Targets))
			for i, _ := range expr.Targets {
				t[i] = unresolvedType
			}
		}
		for i, target := range expr.Targets {
			semanticTargetPass(decl, target, t[i], expr.Define, scope, glbls, status)
		}
		return t
	case *Slice:
		semanticBlockPass(decl, expr.Block, scope, glbls, status)
		return []ASTType{glbls.String}
	case *StringMatch:
		return []ASTType{glbls.String}
	case *RuneMatch:
		return []ASTType{glbls.Rune}
	case *RuneLiteral:
		return []ASTType{glbls.Rune}
	case *StringLiteral:
		return []ASTType{glbls.String}
	case *IntLiteral:
		return []ASTType{glbls.Int}
	case *BoolLiteral:
		return []ASTType{glbls.Bool}
	case *Return:
		if len(decl.ReturnTypes) != len(expr.Exprs) {
			status.Error("wrong number of return types: %d vs. %d", len(expr.Exprs), len(decl.ReturnTypes))
		}
		for i, e := range expr.Exprs {
			at := scalarSemanticExprPass(decl, e, scope, glbls, status)
			if i < len(decl.ReturnTypes) {
				et := ResolveType(decl.ReturnTypes[i])
				if !TypeMatches(at, et, false) {
					status.Error("return: %s vs. %s", TypeName(at), TypeName(et))
				}

			}
		}
		return nil
	case *Position:
		return []ASTType{glbls.Int}
	case *Fail:
		return nil
	case *Call:
		types := glbls.ReturnTypes(expr.Name.Text)
		expr.T = types
		return types
	case *Append:
		t := scalarSemanticExprPass(decl, expr.List, scope, glbls, status)
		scalarSemanticExprPass(decl, expr.Expr, scope, glbls, status)
		// TODO type check arguments
		expr.T = t
		return []ASTType{t}
	case *Construct:
		t := semanticTypePass(expr.Type, glbls, status)
		st, ok := t.(*StructDecl)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			aft := scalarSemanticExprPass(decl, arg.Expr, scope, glbls, status)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := ResolveType(f.Type)
					if !TypeMatches(aft, eft, false) {
						status.Error("%s.%s: %s vs. %s", TypeName(t), fn, TypeName(aft), TypeName(eft))
					}
				} else {
					status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", TypeName(t), fn))
				}
			}
		}
		return []ASTType{t}
	case *ConstructList:
		t := semanticTypePass(expr.Type, glbls, status)
		lt, ok := t.(*ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			at := scalarSemanticExprPass(decl, arg, scope, glbls, status)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					status.Error("%s vs. %s", TypeName(at), TypeName(lt.Type))
				}
			}
		}
		return []ASTType{t}
	case *Coerce:
		t := semanticTypePass(expr.Type, glbls, status)
		scalarSemanticExprPass(decl, expr.Expr, scope, glbls, status)
		// TODO type check
		return []ASTType{t}
	default:
		panic(expr)
	}
}

func semanticTypePass(node ASTTypeRef, glbls *ModuleScope, status framework.Status) ASTType {
	switch node := node.(type) {
	case *TypeRef:
		name := node.Name.Text
		d, ok := glbls.Module[name]
		if !ok {
			d, ok = glbls.Builtin[name]
		}
		if !ok {
			status.LocationError(node.Name.Pos, fmt.Sprintf("Could not resolve name %#v", name))
			node.T = unresolvedType
			return unresolvedType
		}
		t, ok := AsType(d)
		if !ok {
			status.LocationError(node.Name.Pos, fmt.Sprintf("%#v is not a type", name))
			node.T = unresolvedType
			return unresolvedType
		}
		node.T = t
		return t
	case *ListTypeRef:
		t := semanticTypePass(node.Type, glbls, status)
		// TODO memoize list types
		node.T = &ListType{Type: t}
		return node.T
	default:
		panic(node)
	}
}

func semanticBlockPass(decl *FuncDecl, block []ASTExpr, scope *semanticScope, glbls *ModuleScope, status framework.Status) {
	for _, expr := range block {
		semanticExprPass(decl, expr, scope, glbls, status)
	}
}

func semanticFuncSignaturePass(decl *FuncDecl, glbls *ModuleScope, status framework.Status) {
	for _, t := range decl.ReturnTypes {
		semanticTypePass(t, glbls, status)
	}
}

func semanticFuncBodyPass(decl *FuncDecl, glbls *ModuleScope, status framework.Status) {
	semanticBlockPass(decl, decl.Block, childScope(nil), glbls, status)
}

func semanticStructPass(decl *StructDecl, glbls *ModuleScope, status framework.Status) {
	if decl.Implements != nil {
		semanticTypePass(decl.Implements, glbls, status)
	}
	for _, f := range decl.Fields {
		semanticTypePass(f.Type, glbls, status)
	}
}

func semanticDestructurePass(decl *FuncDecl, d Destructure, scope *semanticScope, glbls *ModuleScope, status framework.Status) ASTType {
	switch d := d.(type) {
	case *DestructureStruct:
		t := semanticTypePass(d.Type, glbls, status)
		st, ok := t.(*StructDecl)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			aft := semanticDestructurePass(decl, arg.Destructure, scope, glbls, status)
			if st != nil {
				fn := arg.Name.Text
				f := GetField(st, fn)
				if f != nil {
					eft := ResolveType(f.Type)
					if !TypeMatches(aft, eft, false) {
						status.Error("%s.%s: %s vs. %s", TypeName(t), fn, TypeName(aft), TypeName(eft))
					}
				} else {
					status.LocationError(arg.Name.Pos, fmt.Sprintf("%s does not have field %s", TypeName(t), fn))
				}
			}
		}
		return t
	case *DestructureList:
		t := semanticTypePass(d.Type, glbls, status)
		lt, ok := t.(*ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			at := semanticDestructurePass(decl, arg, scope, glbls, status)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					status.Error("%s vs. %s", TypeName(at), TypeName(lt.Type))
				}
			}
		}
		return t
	case *DestructureValue:
		return scalarSemanticExprPass(decl, d.Expr, scope, glbls, status)
	default:
		panic(d)
	}
}

func semanticTestPass(tst *Test, glbls *ModuleScope, status framework.Status) {
	types := glbls.ReturnTypes(tst.Rule.Text)
	if len(types) != 1 {
		panic(types)
	}
	tst.Type = types[0]
	// HACK no real context
	at := semanticDestructurePass(nil, tst.Destructure, nil, glbls, status)
	if !TypeMatches(at, tst.Type, false) {
		status.Error("%s vs. %s", TypeName(at), TypeName(tst.Type))
	}
}

type ModuleScope struct {
	Builtin map[string]ASTDecl
	Module  map[string]ASTDecl

	BinaryOps map[string]*BuiltinType

	String *BuiltinType
	Rune   *BuiltinType
	Int    *BuiltinType
	Bool   *BuiltinType
	Void   *BuiltinType
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

func AsFunc(node ASTDecl) (ASTFunc, bool) {
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

func ReturnTypes(node ASTFunc) []ASTType {
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

func (glbls *ModuleScope) ReturnTypes(name string) []ASTType {
	// HACK resolve other scopes?
	decl, ok := glbls.Module[name]
	if !ok {
		panic(name)
	}
	f, ok := AsFunc(decl)
	if !ok {
		panic(name)
	}
	return ReturnTypes(f)
}

func SemanticPass(file *File, status framework.Status) *ModuleScope {
	glbls := &ModuleScope{
		Builtin:   map[string]ASTDecl{},
		Module:    map[string]ASTDecl{},
		BinaryOps: map[string]*BuiltinType{},
	}
	glbls.String = &BuiltinType{"string"}
	glbls.Builtin["string"] = glbls.String

	glbls.Rune = &BuiltinType{"rune"}
	glbls.Builtin["rune"] = glbls.Rune

	glbls.Int = &BuiltinType{"int"}
	glbls.Builtin["int"] = glbls.Int

	glbls.Bool = &BuiltinType{"bool"}
	glbls.Builtin["bool"] = glbls.Bool

	glbls.BinaryOps["int+int"] = glbls.Int
	glbls.BinaryOps["int-int"] = glbls.Int
	glbls.BinaryOps["int*int"] = glbls.Int
	glbls.BinaryOps["int/int"] = glbls.Int

	// Index the module namespace.
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			name := decl.Name.Text
			_, exists := glbls.Module[name]
			if exists {
				status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
			} else {
				glbls.Module[name] = decl
			}
		case *StructDecl:
			name := decl.Name.Text
			_, exists := glbls.Module[name]
			if exists {
				status.LocationError(decl.Name.Pos, fmt.Sprintf("Tried to redefine %#v", name))
			} else {
				glbls.Module[name] = decl
			}
		default:
			panic(decl)
		}
	}
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			// Needed for resolving calls in the next step.
			semanticFuncSignaturePass(decl, glbls, status)
		case *StructDecl:
			// Needed for resolving field types.
			semanticStructPass(decl, glbls, status)
		default:
			panic(decl)
		}
	}

	// Resolve the declaration contents.
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			semanticFuncBodyPass(decl, glbls, status)
		case *StructDecl:
		default:
			panic(decl)
		}
	}
	for _, tst := range file.Tests {
		semanticTestPass(tst, glbls, status)
	}
	return glbls
}

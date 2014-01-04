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

func TypeMatches(actual ASTType, expected ASTType, exact bool) bool {
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

func semanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope, glbls *ModuleScope, status framework.Status) ASTType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(decl, expr.Block, scope, glbls, status)
		return glbls.Void
	case *Choice:
		for _, block := range expr.Blocks {
			semanticBlockPass(decl, block, childScope(scope), glbls, status)
		}
		return glbls.Void
	case *Optional:
		semanticBlockPass(decl, expr.Block, scope, glbls, status)
		return glbls.Void
	case *If:
		semanticExprPass(decl, expr.Expr, scope, glbls, status)
		// TODO check condition type
		semanticBlockPass(decl, expr.Block, childScope(scope), glbls, status)
		return glbls.Void
	case *BinaryOp:
		l := semanticExprPass(decl, expr.Left, scope, glbls, status)
		r := semanticExprPass(decl, expr.Right, scope, glbls, status)
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
		return t
	case *GetName:
		name := expr.Name.Text
		info, found := scope.localInfo(name)
		if !found {
			panic(fmt.Sprintf("Could not resolve name %#v", name))
		}
		expr.Info = info
		return decl.Locals[info].T
	case *Assign:
		var t ASTType
		if expr.Expr != nil {
			t = semanticExprPass(decl, expr.Expr, scope, glbls, status)
		}
		if expr.Type != nil {
			t = semanticTypePass(expr.Type, glbls, status)
		}
		name := expr.Name.Text
		if t == nil {
			panic(fmt.Sprintf("%s: Cannot infer the type of %#v", decl.Name.Text, name))
		}
		var info int
		var exists bool
		if expr.Define {
			_, exists = scope.localInfo(expr.Name.Text)
			if exists {
				panic(fmt.Sprintf("Tried to redefine %#v", name))
			}

			info = len(decl.Locals)
			decl.Locals = append(decl.Locals, &LocalInfo{Name: name, T: t})
			scope.locals[expr.Name.Text] = info
		} else {
			info, exists = scope.localInfo(name)
			if !exists {
				panic(fmt.Sprintf("%s: Tried to assign to unknown variable %#v", decl.Name.Text, name))
			}
		}
		expr.Info = info
		return t
	case *Slice:
		semanticBlockPass(decl, expr.Block, scope, glbls, status)
		return glbls.String
	case *StringMatch:
		return glbls.String
	case *RuneMatch:
		return glbls.Rune
	case *RuneLiteral:
		return glbls.Rune
	case *StringLiteral:
		return glbls.String
	case *IntLiteral:
		return glbls.Int
	case *BoolLiteral:
		return glbls.Bool
	case *Return:
		if len(decl.ReturnTypes) != len(expr.Exprs) {
			status.Error("wrong number of return types: %d vs. %d", len(expr.Exprs), len(decl.ReturnTypes))
		}
		for i, e := range expr.Exprs {
			at := semanticExprPass(decl, e, scope, glbls, status)
			if i < len(decl.ReturnTypes) {
				et := ResolveType(decl.ReturnTypes[i])
				if !TypeMatches(at, et, false) {
					status.Error("return: %s vs. %s", TypeName(at), TypeName(et))
				}

			}
		}
		return glbls.Void
	case *Fail:
		return glbls.Void
	case *Call:
		t := glbls.ReturnType(expr.Name.Text)
		expr.T = t
		return t
	case *Append:
		t := semanticExprPass(decl, expr.List, scope, glbls, status)
		semanticExprPass(decl, expr.Expr, scope, glbls, status)
		expr.T = t
		return t
	case *Construct:
		t := semanticTypePass(expr.Type, glbls, status)
		st, ok := t.(*StructDecl)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			aft := semanticExprPass(decl, arg.Expr, scope, glbls, status)
			if st != nil {
				fn := arg.Name.Text
				eft := FieldType(st, fn)
				if !TypeMatches(aft, eft, false) {
					status.Error("%s.%s: %s vs. %s", TypeName(t), fn, TypeName(aft), TypeName(eft))
				}
			}
		}
		return t
	case *ConstructList:
		t := semanticTypePass(expr.Type, glbls, status)
		lt, ok := t.(*ListType)
		if t != nil && !ok {
			panic(t)
		}
		for _, arg := range expr.Args {
			at := semanticExprPass(decl, arg, scope, glbls, status)
			if lt != nil {
				if !TypeMatches(at, lt.Type, false) {
					status.Error("%s vs. %s", TypeName(at), TypeName(lt.Type))
				}
			}
		}
		return t
	case *Coerce:
		t := semanticTypePass(expr.Type, glbls, status)
		semanticExprPass(decl, expr.Expr, scope, glbls, status)
		return t
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
			status.Error("Unknown name %#v", name)
			panic(name)
		}
		t, ok := AsType(d)
		if !ok {
			status.Error("%#v is not a type", name)
			panic(name)
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
				eft := FieldType(st, fn)
				if !TypeMatches(aft, eft, false) {
					status.Error("%s.%s: %s vs. %s", TypeName(t), fn, TypeName(aft), TypeName(eft))
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
		return semanticExprPass(decl, d.Expr, scope, glbls, status)
	default:
		panic(d)
	}
}

func semanticTestPass(tst *Test, glbls *ModuleScope, status framework.Status) {
	tst.Type = glbls.ReturnType(tst.Rule.Text)
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

func FieldType(node *StructDecl, name string) ASTType {
	for _, decl := range node.Fields {
		if decl.Name.Text == name {
			return ResolveType(decl.Type)
		}
	}
	panic(name)
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

func ReturnType(node ASTFunc) ASTType {
	switch node := node.(type) {
	case *FuncDecl:
		// HACK assume single return value
		if len(node.ReturnTypes) == 0 {
			return nil
		}
		if len(node.ReturnTypes) != 1 {
			panic(node.Name.Text)
		}
		return ResolveType(node.ReturnTypes[0])
	default:
		panic(node)
	}
}

func (glbls *ModuleScope) ReturnType(name string) ASTType {
	// HACK resolve other scopes?
	decl, ok := glbls.Module[name]
	if !ok {
		panic(name)
	}
	f, ok := AsFunc(decl)
	if !ok {
		panic(name)
	}
	return ReturnType(f)
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

	glbls.Void = &BuiltinType{"void"}
	glbls.Builtin["void"] = glbls.Void

	glbls.BinaryOps["int+int"] = glbls.Int
	glbls.BinaryOps["int-int"] = glbls.Int
	glbls.BinaryOps["int*int"] = glbls.Int
	glbls.BinaryOps["int/int"] = glbls.Int

	// Index the module namespace.
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			glbls.Module[decl.Name.Text] = decl
		case *StructDecl:
			glbls.Module[decl.Name.Text] = decl
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

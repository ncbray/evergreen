package dasm

import (
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

func semanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope, glbls *ModuleScope) ASTType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(decl, expr.Block, scope, glbls)
		return glbls.Void
	case *Optional:
		semanticBlockPass(decl, expr.Block, scope, glbls)
		return glbls.Void
	case *If:
		semanticExprPass(decl, expr.Expr, scope, glbls)
		// TODO check condition type
		semanticBlockPass(decl, expr.Block, childScope(scope), glbls)
		return glbls.Void
	case *BinaryOp:
		semanticExprPass(decl, expr.Left, scope, glbls)
		semanticExprPass(decl, expr.Right, scope, glbls)
		// HACK assume compare
		t := glbls.Bool
		expr.T = t
		return t
	case *GetName:
		info, found := scope.localInfo(expr.Name)
		if !found {
			panic(fmt.Sprintf("Could not resolve name %#v", expr.Name))
		}
		expr.Info = info
		return decl.Locals[info].T
	case *Assign:
		var t ASTType
		if expr.Expr != nil {
			t = semanticExprPass(decl, expr.Expr, scope, glbls)
		}
		if expr.Type != nil {
			t = semanticTypePass(expr.Type, glbls)
		}
		if t == nil {
			panic(fmt.Sprintf("%s: Cannot infer the type of %#v", decl.Name, expr.Name))
		}
		var info int
		var exists bool
		if expr.Define {
			_, exists = scope.localInfo(expr.Name)
			if exists {
				panic(fmt.Sprintf("Tried to redefine %#v", expr.Name))
			}

			info = len(decl.Locals)
			decl.Locals = append(decl.Locals, &LocalInfo{Name: expr.Name, T: t})
			scope.locals[expr.Name] = info
		} else {
			info, exists = scope.localInfo(expr.Name)
			if !exists {
				panic(fmt.Sprintf("%s: Tried to assign to unknown variable %#v", decl.Name, expr.Name))
			}
		}
		expr.Info = info
		return t
	case *Slice:
		semanticBlockPass(decl, expr.Block, scope, glbls)
		return glbls.String
	case *Match:
		return glbls.String
	case *Read:
		return glbls.Rune
	case *RuneLiteral:
		return glbls.Rune
	case *StringLiteral:
		return glbls.String
	case *IntLiteral:
		return glbls.Int
	case *Return:
		for _, e := range expr.Exprs {
			semanticExprPass(decl, e, scope, glbls)
		}
		return glbls.Void
	case *Fail:
		return glbls.Void
	case *Call:
		t := glbls.ReturnType(expr.Name)
		expr.T = t
		return t
	case *Append:
		t := semanticExprPass(decl, expr.List, scope, glbls)
		semanticExprPass(decl, expr.Value, scope, glbls)
		expr.T = t
		return t
	case *Construct:
		t := semanticTypePass(expr.Type, glbls)
		for _, arg := range expr.Args {
			semanticExprPass(decl, arg.Value, scope, glbls)
		}
		return t
	case *ConstructList:
		t := semanticTypePass(expr.Type, glbls)
		for _, arg := range expr.Args {
			semanticExprPass(decl, arg, scope, glbls)
		}
		return t
	default:
		panic(expr)
	}
}

func semanticTypePass(node ASTTypeRef, glbls *ModuleScope) ASTType {
	switch node := node.(type) {
	case *TypeRef:
		d, ok := glbls.Module[node.Name]
		if !ok {
			d, ok = glbls.Builtin[node.Name]
		}
		if !ok {
			panic(node.Name)
		}
		t, ok := d.AsType()
		if !ok {
			panic(node.Name)
		}
		node.T = t
		return t
	case *ListTypeRef:
		t := semanticTypePass(node.Type, glbls)
		// TODO memoize list types
		node.T = &ListType{Type: t}
		return node.T
	default:
		panic(node)
	}
}

func semanticBlockPass(decl *FuncDecl, block []ASTExpr, scope *semanticScope, glbls *ModuleScope) {
	for _, expr := range block {
		semanticExprPass(decl, expr, scope, glbls)
	}
}

func semanticFuncSignaturePass(decl *FuncDecl, glbls *ModuleScope) {
	for _, t := range decl.ReturnTypes {
		semanticTypePass(t, glbls)
	}
}

func semanticFuncBodyPass(decl *FuncDecl, glbls *ModuleScope) {
	semanticBlockPass(decl, decl.Block, childScope(nil), glbls)
}

func semanticStructPass(decl *StructDecl, glbls *ModuleScope) {
	if decl.Implements != nil {
		semanticTypePass(decl.Implements, glbls)
	}
	for _, f := range decl.Fields {
		semanticTypePass(f.Type, glbls)
	}
}

func semanticDestructurePass(d Destructure, general ASTType, glbls *ModuleScope) {
	switch d := d.(type) {
	case *DestructureStruct:
		semanticTypePass(d.Type, glbls)
		switch t := d.Type.Resolve().(type) {
		case *StructDecl:
			d.Actual = t
		default:
			panic(t)
		}
		switch t := general.(type) {
		case *StructDecl:
			d.General = t
		default:
			panic(t)
		}
		for _, arg := range d.Args {
			semanticDestructurePass(arg.Destructure, d.Actual.FieldType(arg.Name), glbls)
		}
	case *DestructureList:
		semanticTypePass(d.Type, glbls)
		t := d.Type.Resolve()
		dt, ok := t.(*ListType)
		if !ok {
			panic(t)
		}
		for _, arg := range d.Args {
			semanticDestructurePass(arg, dt.Type, glbls)
		}

	case *DestructureString, *DestructureRune, *DestructureInt:
		// Leaf
	default:
		panic(d)
	}
}

func semanticTestPass(tst *Test, glbls *ModuleScope) {
	general := glbls.ReturnType(tst.Rule)
	semanticDestructurePass(tst.Destructure, general, glbls)
}

type ModuleScope struct {
	Builtin map[string]Decl
	Module  map[string]Decl

	String *BuiltinType
	Rune   *BuiltinType
	Int    *BuiltinType
	Bool   *BuiltinType
	Void   *BuiltinType
}

func (glbls *ModuleScope) ReturnType(name string) ASTType {
	// HACK resolve other scopes?
	decl, ok := glbls.Module[name]
	if !ok {
		panic(name)
	}
	f, ok := decl.AsFunc()
	if !ok {
		panic(name)
	}
	return f.ReturnType()
}

func SemanticPass(file *File) *ModuleScope {
	glbls := &ModuleScope{
		Builtin: map[string]Decl{},
		Module:  map[string]Decl{},
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

	// Index the module namespace.
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			glbls.Module[decl.Name] = decl
		case *StructDecl:
			glbls.Module[decl.Name] = decl
		default:
			panic(decl)
		}
	}
	// Resolve function signatures.
	// Needed for resolving calls in the next step.
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			semanticFuncSignaturePass(decl, glbls)
		case *StructDecl:
		default:
			panic(decl)
		}
	}

	// Resolve the declaration contents.
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *FuncDecl:
			semanticFuncBodyPass(decl, glbls)
		case *StructDecl:
			semanticStructPass(decl, glbls)
		default:
			panic(decl)
		}
	}
	for _, tst := range file.Tests {
		semanticTestPass(tst, glbls)
	}
	return glbls
}

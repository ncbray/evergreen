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

func semanticExprPass(decl *FuncDecl, expr ASTExpr, scope *semanticScope, glbls *ModuleScope) ASTType {
	switch expr := expr.(type) {
	case *Repeat:
		semanticBlockPass(decl, expr.Block, scope, glbls)
		return glbls.Void
	case *Choice:
		for _, block := range expr.Blocks {
			semanticBlockPass(decl, block, childScope(scope), glbls)
		}
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
		l := semanticExprPass(decl, expr.Left, scope, glbls)
		r := semanticExprPass(decl, expr.Right, scope, glbls)
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
		semanticExprPass(decl, expr.Expr, scope, glbls)
		expr.T = t
		return t
	case *Construct:
		t := semanticTypePass(expr.Type, glbls)
		for _, arg := range expr.Args {
			semanticExprPass(decl, arg.Expr, scope, glbls)
		}
		return t
	case *ConstructList:
		t := semanticTypePass(expr.Type, glbls)
		for _, arg := range expr.Args {
			semanticExprPass(decl, arg, scope, glbls)
		}
		return t
	case *Coerce:
		t := semanticTypePass(expr.Type, glbls)
		semanticExprPass(decl, expr.Expr, scope, glbls)
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
		t, ok := AsType(d)
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

func semanticDestructurePass(decl *FuncDecl, d Destructure, scope *semanticScope, glbls *ModuleScope) {
	switch d := d.(type) {
	case *DestructureStruct:
		semanticTypePass(d.Type, glbls)
		for _, arg := range d.Args {
			semanticDestructurePass(decl, arg.Destructure, scope, glbls)
		}
	case *DestructureList:
		semanticTypePass(d.Type, glbls)
		for _, arg := range d.Args {
			semanticDestructurePass(decl, arg, scope, glbls)
		}

	case *DestructureValue:
		semanticExprPass(decl, d.Expr, scope, glbls)
	default:
		panic(d)
	}
}

func semanticTestPass(tst *Test, glbls *ModuleScope) {
	tst.Type = glbls.ReturnType(tst.Rule)
	// HACK no real context
	semanticDestructurePass(nil, tst.Destructure, nil, glbls)
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
		if decl.Name == name {
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
			panic(node.Name)
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

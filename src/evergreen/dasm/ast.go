package dasm

import (
	"evergreen/dubx"
)

type ASTExpr interface {
	IsASTExpr()
}

type ASTType interface {
	IsASTType()
}

type ASTTypeRef interface {
	IsASTTypeRef()
}

func ResolveType(ref ASTTypeRef) ASTType {
	switch ref := ref.(type) {
	case *dubx.TypeRef:
		return ref.T
	case *dubx.ListTypeRef:
		return ref.T
	default:
		panic(ref)
	}
}

type LocalInfo struct {
	Name string
	T    ASTType
}

type Decl interface {
	AsType() (ASTType, bool)
	AsFunc() (ASTFunc, bool)
	IsASTDecl()
}

type ASTFunc interface {
	ReturnType() ASTType
}

type FuncDecl struct {
	Name        string
	ReturnTypes []dubx.ASTTypeRef
	Block       []dubx.ASTExpr
	Locals      []*LocalInfo
}

func (node *FuncDecl) AsType() (ASTType, bool) {
	return nil, false
}

func (node *FuncDecl) AsFunc() (ASTFunc, bool) {
	return node, true
}

func (node *FuncDecl) ReturnType() ASTType {
	// HACK assume single return value
	if len(node.ReturnTypes) == 0 {
		return nil
	}
	if len(node.ReturnTypes) != 1 {
		panic(node.Name)
	}
	return ResolveType(node.ReturnTypes[0])
}

func (node *FuncDecl) IsASTDecl() {
}

type FieldDecl struct {
	Name string
	Type ASTTypeRef
}

type StructDecl struct {
	Name       string
	Implements ASTTypeRef
	Fields     []*FieldDecl
}

func (node *StructDecl) FieldType(name string) ASTType {
	for _, decl := range node.Fields {
		if decl.Name == name {
			return ResolveType(decl.Type)
		}
	}

	panic(name)
}

func (node *StructDecl) AsType() (ASTType, bool) {
	return node, true
}

func (node *StructDecl) AsFunc() (ASTFunc, bool) {
	return nil, false
}

func (node *StructDecl) IsASTDecl() {
}

func (node *StructDecl) IsASTType() {
}

type BuiltinType struct {
	Name string
}

func (node *BuiltinType) AsType() (ASTType, bool) {
	return node, true
}

func (node *BuiltinType) AsFunc() (ASTFunc, bool) {
	return nil, false
}

func (node *BuiltinType) IsASTDecl() {
}

func (node *BuiltinType) IsASTType() {
}

type ListType struct {
	Type ASTType
}

func (node *ListType) IsASTType() {
}

type File struct {
	Decls []Decl
	Tests []*dubx.Test
}

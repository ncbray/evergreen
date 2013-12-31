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

type LocalInfo struct {
	Name string
	T    ASTType
}

type ASTDecl interface {
	IsASTDecl()
}

type ASTFunc interface {
	IsASTFunc()
}

type FuncDecl struct {
	Name        string
	ReturnTypes []dubx.ASTTypeRef
	Block       []dubx.ASTExpr
	Locals      []*LocalInfo
}

func (node *FuncDecl) IsASTDecl() {
}

func (node *FuncDecl) IsASTFunc() {
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

func (node *StructDecl) IsASTDecl() {
}

func (node *StructDecl) IsASTType() {
}

type BuiltinType struct {
	Name string
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
	Decls []ASTDecl
	Tests []*dubx.Test
}

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

type If struct {
	Expr  ASTExpr
	Block []ASTExpr
}

func (node *If) IsASTExpr() {
}

type Repeat struct {
	Block []ASTExpr
	Min   int
}

func (node *Repeat) IsASTExpr() {
}

type Choice struct {
	// TODO wrap blocks
	Blocks [][]ASTExpr
}

func (node *Choice) IsASTExpr() {
}

type Optional struct {
	Block []ASTExpr
}

func (node *Optional) IsASTExpr() {
}

type Slice struct {
	Block []ASTExpr
}

func (node *Slice) IsASTExpr() {
}

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	Right ASTExpr
	T     ASTType
}

func (node *BinaryOp) IsASTExpr() {
}

type KeyValue struct {
	Key   string
	Value ASTExpr
}

type Construct struct {
	Type ASTTypeRef
	Args []*KeyValue
}

func (node *Construct) IsASTExpr() {
}

type ConstructList struct {
	Type ASTTypeRef
	Args []ASTExpr
}

func (node *ConstructList) IsASTExpr() {
}

type Coerce struct {
	Type ASTTypeRef
	Expr ASTExpr
}

func (node *Coerce) IsASTExpr() {
}

type GetName struct {
	Name string
	Info int
}

func (node *GetName) IsASTExpr() {
}

type Assign struct {
	Expr   ASTExpr
	Name   string
	Info   int
	Type   ASTTypeRef
	Define bool
}

func (node *Assign) IsASTExpr() {
}

type Append struct {
	List  ASTExpr
	Value ASTExpr
	T     ASTType
}

func (node *Append) IsASTExpr() {
}

type Return struct {
	Exprs []ASTExpr
}

func (node *Return) IsASTExpr() {
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
	ReturnTypes []ASTTypeRef
	Block       []ASTExpr
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

type Test struct {
	Name        string
	Rule        string
	Type        ASTType
	Input       string
	Destructure dubx.Destructure
}

type File struct {
	Decls []Decl
	Tests []*Test
}

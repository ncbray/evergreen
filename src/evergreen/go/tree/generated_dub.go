package tree

import (
	"evergreen/go/core"
)

type TypeRef interface {
	isTypeRef()
}

type NameRef struct {
	Name string
	T    core.GoType
}

func (node *NameRef) isTypeRef() {
}

type PointerRef struct {
	Element TypeRef
	T       core.GoType
}

func (node *PointerRef) isTypeRef() {
}

type SliceRef struct {
	Element TypeRef
	T       core.GoType
}

func (node *SliceRef) isTypeRef() {
}

type Param struct {
	Name string
	Type TypeRef
	Info *LocalInfo
}

type FuncTypeRef struct {
	Params  []*Param
	Results []*Param
}

func (node *FuncTypeRef) isTypeRef() {
}

type Stmt interface {
	isStmt()
}

type Expr interface {
	isExpr()
	isStmt()
}

type Target interface {
	isTarget()
}

type IntLiteral struct {
	Value int
}

func (node *IntLiteral) isStmt() {
}

func (node *IntLiteral) isExpr() {
}

type Float32Literal struct {
	Value float32
}

func (node *Float32Literal) isStmt() {
}

func (node *Float32Literal) isExpr() {
}

type BoolLiteral struct {
	Value bool
}

func (node *BoolLiteral) isStmt() {
}

func (node *BoolLiteral) isExpr() {
}

type StringLiteral struct {
	Value string
}

func (node *StringLiteral) isStmt() {
}

func (node *StringLiteral) isExpr() {
}

type RuneLiteral struct {
	Value rune
}

func (node *RuneLiteral) isStmt() {
}

func (node *RuneLiteral) isExpr() {
}

type NilLiteral struct {
}

func (node *NilLiteral) isStmt() {
}

func (node *NilLiteral) isExpr() {
}

type KeywordExpr struct {
	Name string
	Expr Expr
}

type StructLiteral struct {
	Type TypeRef
	Args []*KeywordExpr
}

func (node *StructLiteral) isStmt() {
}

func (node *StructLiteral) isExpr() {
}

type ListLiteral struct {
	Type TypeRef
	Args []Expr
}

func (node *ListLiteral) isStmt() {
}

func (node *ListLiteral) isExpr() {
}

type LocalInfo_Ref uint32

type LocalInfo_Scope struct {
	objects []*LocalInfo
}

type LocalInfo struct {
	Name  string
	T     TypeRef
	Index LocalInfo_Ref
}

type GetName struct {
	Text string
}

func (node *GetName) isStmt() {
}

func (node *GetName) isExpr() {
}

type SetName struct {
	Text string
}

func (node *SetName) isTarget() {
}

type GetLocal struct {
	Info *LocalInfo
}

func (node *GetLocal) isStmt() {
}

func (node *GetLocal) isExpr() {
}

type SetLocal struct {
	Info *LocalInfo
}

func (node *SetLocal) isTarget() {
}

type GetGlobal struct {
	Text string
}

func (node *GetGlobal) isStmt() {
}

func (node *GetGlobal) isExpr() {
}

type GetFunction struct {
	Func core.Callable
}

func (node *GetFunction) isStmt() {
}

func (node *GetFunction) isExpr() {
}

type SetDiscard struct {
}

func (node *SetDiscard) isTarget() {
}

type UnaryExpr struct {
	Op   string
	Expr Expr
}

func (node *UnaryExpr) isStmt() {
}

func (node *UnaryExpr) isExpr() {
}

type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}

func (node *BinaryExpr) isStmt() {
}

func (node *BinaryExpr) isExpr() {
}

type Selector struct {
	Expr Expr
	Text string
}

func (node *Selector) isStmt() {
}

func (node *Selector) isExpr() {
}

type Index struct {
	Expr  Expr
	Index Expr
}

func (node *Index) isStmt() {
}

func (node *Index) isExpr() {
}

type Call struct {
	Expr Expr
	Args []Expr
}

func (node *Call) isStmt() {
}

func (node *Call) isExpr() {
}

type TypeAssert struct {
	Expr Expr
	Type TypeRef
}

func (node *TypeAssert) isStmt() {
}

func (node *TypeAssert) isExpr() {
}

type TypeCoerce struct {
	Type TypeRef
	Expr Expr
}

func (node *TypeCoerce) isStmt() {
}

func (node *TypeCoerce) isExpr() {
}

type Assign struct {
	Sources []Expr
	Op      string
	Targets []Target
}

func (node *Assign) isStmt() {
}

type Var struct {
	Name string
	Type TypeRef
	Expr Expr
	Info *LocalInfo
}

func (node *Var) isStmt() {
}

type Block struct {
	Body []Stmt
}

type BlockStmt struct {
	Block *Block
}

func (node *BlockStmt) isStmt() {
}

type If struct {
	Cond Expr
	T    *Block
	F    *Block
}

func (node *If) isStmt() {
}

type For struct {
	Block *Block
}

func (node *For) isStmt() {
}

type Goto struct {
	Text string
}

func (node *Goto) isStmt() {
}

type Label struct {
	Text string
}

func (node *Label) isStmt() {
}

type Return struct {
	Args []Expr
}

func (node *Return) isStmt() {
}

type Decl interface {
	isDecl()
}

type VarDecl struct {
	Name  string
	Type  TypeRef
	Expr  Expr
	Const bool
}

func (node *VarDecl) isDecl() {
}

type FuncDecl struct {
	Name            string
	Recv            *Param
	Type            *FuncTypeRef
	Block           *Block
	Package         *core.Package
	LocalInfo_Scope *LocalInfo_Scope
}

func (node *FuncDecl) isDecl() {
}

type FieldDecl struct {
	Name string
	Type TypeRef
}

type StructDecl struct {
	Name   string
	Fields []*FieldDecl
	T      *core.StructType
}

func (node *StructDecl) isDecl() {
}

type InterfaceDecl struct {
	Name   string
	Fields []*FieldDecl
	T      *core.InterfaceType
}

func (node *InterfaceDecl) isDecl() {
}

type TypeDefDecl struct {
	Name string
	Type TypeRef
	T    *core.TypeDefType
}

func (node *TypeDefDecl) isDecl() {
}

type OpaqueDecl struct {
	T *core.ExternalType
}

func (node *OpaqueDecl) isDecl() {
}

type Import struct {
	Name string
	Path string
}

type FileAST struct {
	Name    string
	Package string
	Imports []*Import
	Decls   []Decl
}

type PackageAST struct {
	Files []*FileAST
	P     *core.Package
}

type ProgramAST struct {
	Builtins *core.BuiltinTypeIndex
	Packages []*PackageAST
}

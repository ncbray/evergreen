package dasm

type ASTExpr interface {
	isASTExpr()
}

type ASTType interface {
	isASTType()
}

type ASTTypeRef interface {
	Resolve() ASTType
	isASTTypeRef()
}

type TypeRef struct {
	Name string
	T    ASTType
}

func (node *TypeRef) Resolve() ASTType {
	return node.T
}

func (node *TypeRef) isASTTypeRef() {
}

type ListTypeRef struct {
	Type ASTTypeRef
	T    ASTType
}

func (node *ListTypeRef) Resolve() ASTType {
	return node.T
}

func (node *ListTypeRef) isASTTypeRef() {
}

type If struct {
	Expr  ASTExpr
	Block []ASTExpr
}

func (node *If) isASTExpr() {
}

type Repeat struct {
	Block []ASTExpr
	Min   int
}

func (node *Repeat) isASTExpr() {
}

type Optional struct {
	Block []ASTExpr
}

func (node *Optional) isASTExpr() {
}

type Slice struct {
	Block []ASTExpr
}

func (node *Slice) isASTExpr() {
}

type BinaryOp struct {
	Left  ASTExpr
	Op    string
	Right ASTExpr
	T     ASTType
}

func (node *BinaryOp) isASTExpr() {
}

type Call struct {
	Name string
	T    ASTType
}

func (node *Call) isASTExpr() {
}

type KeyValue struct {
	Key   string
	Value ASTExpr
}

type Construct struct {
	Type ASTTypeRef
	Args []*KeyValue
}

func (node *Construct) isASTExpr() {
}

type ConstructList struct {
	Type ASTTypeRef
	Args []ASTExpr
}

func (node *ConstructList) isASTExpr() {
}

type GetName struct {
	Name string
	Info int
}

func (node *GetName) isASTExpr() {
}

type RuneLiteral struct {
	Value rune
}

func (node *RuneLiteral) isASTExpr() {
}

type StringLiteral struct {
	Value string
}

func (node *StringLiteral) isASTExpr() {
}

type Assign struct {
	Expr   ASTExpr
	Name   string
	Info   int
	Type   ASTTypeRef
	Define bool
}

func (node *Assign) isASTExpr() {
}

type Read struct {
}

func (node *Read) isASTExpr() {
}

type Append struct {
	List  ASTExpr
	Value ASTExpr
	T     ASTType
}

func (node *Append) isASTExpr() {
}

type Return struct {
	Exprs []ASTExpr
}

func (node *Return) isASTExpr() {
}

type Fail struct {
}

func (node *Fail) isASTExpr() {
}

type LocalInfo struct {
	Name string
	T    ASTType
}

type Decl interface {
	AsType() (ASTType, bool)
	AsFunc() (ASTFunc, bool)
	isASTDecl()
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
	return node.ReturnTypes[0].Resolve()
}

func (node *FuncDecl) isASTDecl() {
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

func (node *StructDecl) AsType() (ASTType, bool) {
	return node, true
}

func (node *StructDecl) AsFunc() (ASTFunc, bool) {
	return nil, false
}

func (node *StructDecl) isASTDecl() {
}

func (node *StructDecl) isASTType() {
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

func (node *BuiltinType) isASTDecl() {
}

func (node *BuiltinType) isASTType() {
}

type ListType struct {
	Type ASTType
}

func (node *ListType) isASTType() {
}

type Destructure interface {
	isDestructure()
}

type DestructureList struct {
	Type *ListTypeRef
	Args []Destructure
}

func (node *DestructureList) isDestructure() {
}

type DestructureField struct {
	Name        string
	Destructure Destructure
}

type DestructureStruct struct {
	Type *TypeRef
	Args []*DestructureField
}

func (node *DestructureStruct) isDestructure() {
}

type DestructureString struct {
	Value string
}

func (node *DestructureString) isDestructure() {
}

type Test struct {
	Name        string
	Rule        string
	Input       string
	Destructure Destructure
}

type File struct {
	Decls []Decl
	Tests []*Test
}

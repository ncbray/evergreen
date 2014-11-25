package tree

type Type interface {
	isType()
}

type TypeRef struct {
	Name string
	Impl TypeImpl
}

func (node *TypeRef) isType() {
}

type PointerType struct {
	Element Type
}

func (node *PointerType) isType() {
}

type SliceType struct {
	Element Type
}

func (node *SliceType) isType() {
}

type Param struct {
	Name string
	Type Type
	Info int
}

type FuncType struct {
	Params  []*Param
	Results []*Param
}

func (node *FuncType) isType() {
}

type Stmt interface {
	isStmt()
}

type Expr interface {
	isExpr()
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
	Type *TypeRef
	Args []*KeywordExpr
}

func (node *StructLiteral) isStmt() {
}

func (node *StructLiteral) isExpr() {
}

type ListLiteral struct {
	Type *SliceType
	Args []Expr
}

func (node *ListLiteral) isStmt() {
}

func (node *ListLiteral) isExpr() {
}

type LocalInfo_Ref uint32

const NoLocalInfo = ^LocalInfo_Ref(0)

type LocalInfo_Scope struct {
	objects []*LocalInfo
}

type LocalInfo struct {
	Name string
	T    Type
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
	Info int
}

func (node *GetLocal) isStmt() {
}

func (node *GetLocal) isExpr() {
}

type SetLocal struct {
	Info int
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
	Type Type
}

func (node *TypeAssert) isStmt() {
}

func (node *TypeAssert) isExpr() {
}

type TypeCoerce struct {
	Type Type
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
	Type Type
	Expr Expr
	Info int
}

func (node *Var) isStmt() {
}

type BlockStmt struct {
	Body []Stmt
}

func (node *BlockStmt) isStmt() {
}

type If struct {
	Cond Expr
	Body []Stmt
	Else Stmt
}

func (node *If) isStmt() {
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
	Type  Type
	Expr  Expr
	Const bool
}

func (node *VarDecl) isDecl() {
}

type FuncDecl struct {
	Name            string
	Recv            *Param
	Type            *FuncType
	Body            []Stmt
	Package         *Package
	LocalInfo_Scope *LocalInfo_Scope
}

func (node *FuncDecl) isDecl() {
}

type Field struct {
	Name string
	Type Type
}

type TypeImpl interface {
	isTypeImpl()
}

type StructDecl struct {
	Name    string
	Fields  []*Field
	Package *Package
}

func (node *StructDecl) isDecl() {
}

func (node *StructDecl) isTypeImpl() {
}

type InterfaceDecl struct {
	Name    string
	Fields  []*Field
	Package *Package
}

func (node *InterfaceDecl) isDecl() {
}

func (node *InterfaceDecl) isTypeImpl() {
}

type ExternalType struct {
	Name    string
	Package *Package
}

func (node *ExternalType) isDecl() {
}

func (node *ExternalType) isTypeImpl() {
}

type TypeDef struct {
	Name    string
	Type    Type
	Package *Package
}

func (node *TypeDef) isDecl() {
}

func (node *TypeDef) isTypeImpl() {
}

type Import struct {
	Name string
	Path string
}

type File struct {
	Name    string
	Package string
	Imports []*Import
	Decls   []Decl
}

type Package struct {
	Path   []string
	Files  []*File
	Extern bool
}

type Program struct {
	Packages []*Package
}
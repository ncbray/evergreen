package tree

type Type interface {
	isType()
}
type TypeRef struct {
	Name string
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

type Stmt interface {
	isStmt()
}
type Expr interface {
	isExpr()
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

type NameRef struct {
	Text string
}

func (node *NameRef) isStmt() {
}
func (node *NameRef) isExpr() {
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

type Assign struct {
	Sources []Expr
	Op      string
	Targets []Expr
}

func (node *Assign) isStmt() {
}

type BlockStmt struct {
	Body []Stmt
}

func (node *BlockStmt) isStmt() {
}

type If struct {
	Cond Expr
	Body []Stmt
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
type Param struct {
	Name string
	Type Type
}
type FuncDecl struct {
	Name    string
	Params  []*Param
	Returns []*Param
	Body    []Stmt
}

func (node *FuncDecl) isDecl() {
}

type Field struct {
	Name string
	Type Type
}
type StructDecl struct {
	Name   string
	Fields []*Field
}

func (node *StructDecl) isDecl() {
}

type Import struct {
	Name string
	Path string
}
type File struct {
	Package string
	Imports []*Import
	Decls   []Decl
}

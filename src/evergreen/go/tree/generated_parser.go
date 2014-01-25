package tree

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

type Assign struct {
	Sources []Expr
	Op      string
	Targets []Expr
}

func (node *Assign) isStmt() {
}

type If struct {
	Cond Expr
	Body []Stmt
}

func (node *If) isStmt() {
}

type FuncDecl struct {
	Name string
	Body []Stmt
}

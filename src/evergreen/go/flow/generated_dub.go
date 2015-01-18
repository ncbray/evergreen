package flow

import (
	"evergreen/go/core"
	"evergreen/graph"
)

type Register_Ref uint32

type Register_Scope struct {
	objects []*Register
}

type Register struct {
	Name  string
	T     core.GoType
	Index Register_Ref
}

type FlowFunc_Ref uint32

type FlowFunc_Scope struct {
	objects []*FlowFunc
}

type FlowFunc struct {
	Function       *core.Function
	Recv           *Register
	Params         []*Register
	Results        []*Register
	CFG            *graph.Graph
	Ops            []GoOp
	Edges          []int
	Register_Scope *Register_Scope
	Index          FlowFunc_Ref
}

type GoOp interface {
	isGoOp()
}

type ConstantNil struct {
	Dst *Register
}

func (node *ConstantNil) isGoOp() {
}

type ConstantInt struct {
	Value int64
	Dst   *Register
}

func (node *ConstantInt) isGoOp() {
}

type ConstantFloat32 struct {
	Value float32
	Dst   *Register
}

func (node *ConstantFloat32) isGoOp() {
}

type ConstantBool struct {
	Value bool
	Dst   *Register
}

func (node *ConstantBool) isGoOp() {
}

type ConstantRune struct {
	Value rune
	Dst   *Register
}

func (node *ConstantRune) isGoOp() {
}

type ConstantString struct {
	Value string
	Dst   *Register
}

func (node *ConstantString) isGoOp() {
}

type BinaryOp struct {
	Left  *Register
	Op    string
	Right *Register
	Dst   *Register
}

func (node *BinaryOp) isGoOp() {
}

type Attr struct {
	Expr *Register
	Name string
	Dst  *Register
}

func (node *Attr) isGoOp() {
}

type Call struct {
	Target core.Callable
	Args   []*Register
	Dsts   []*Register
}

func (node *Call) isGoOp() {
}

type MethodCall struct {
	Expr *Register
	Name string
	Args []*Register
	Dsts []*Register
}

func (node *MethodCall) isGoOp() {
}

type NamedArg struct {
	Name string
	Arg  *Register
}

type ConstructStruct struct {
	Type      *core.StructType
	AddrTaken bool
	Args      []*NamedArg
	Dst       *Register
}

func (node *ConstructStruct) isGoOp() {
}

type ConstructSlice struct {
	Type *core.SliceType
	Args []*Register
	Dst  *Register
}

func (node *ConstructSlice) isGoOp() {
}

type Coerce struct {
	Src  *Register
	Type core.GoType
	Dst  *Register
}

func (node *Coerce) isGoOp() {
}

type Transfer struct {
	Srcs []*Register
	Dsts []*Register
}

func (node *Transfer) isGoOp() {
}

type Return struct {
	Args []*Register
}

func (node *Return) isGoOp() {
}

type Nop struct {
}

func (node *Nop) isGoOp() {
}

type Entry struct {
}

func (node *Entry) isGoOp() {
}

type Switch struct {
	Cond *Register
}

func (node *Switch) isGoOp() {
}

type Exit struct {
}

func (node *Exit) isGoOp() {
}

type FlowProgram struct {
	Types          []core.GoType
	Builtins       *core.BuiltinTypeIndex
	FlowFunc_Scope *FlowFunc_Scope
}

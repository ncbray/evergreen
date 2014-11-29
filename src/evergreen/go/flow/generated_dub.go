package flow

import (
	"evergreen/base"
	"evergreen/go/tree"
)

type Register_Ref uint32

const NoRegister = ^Register_Ref(0)

type Register_Scope struct {
	objects []*Register
}

type Register struct {
	Name string
	T    tree.GoType
}

type LLFunc struct {
	Name           string
	Params         []Register_Ref
	Results        []Register_Ref
	CFG            *base.Graph
	Ops            []GoOp
	Register_Scope *Register_Scope
}

type GoOp interface {
	isGoOp()
}

type ConstantNil struct {
	Dst Register_Ref
}

func (node *ConstantNil) isGoOp() {
}

type ConstantInt struct {
	Value int64
	Dst   Register_Ref
}

func (node *ConstantInt) isGoOp() {
}

type ConstantBool struct {
	Value bool
	Dst   Register_Ref
}

func (node *ConstantBool) isGoOp() {
}

type ConstantRune struct {
	Value rune
	Dst   Register_Ref
}

func (node *ConstantRune) isGoOp() {
}

type ConstantString struct {
	Value string
	Dst   Register_Ref
}

func (node *ConstantString) isGoOp() {
}

type BinaryOp struct {
	Left  Register_Ref
	Op    string
	Right Register_Ref
	Dst   Register_Ref
}

func (node *BinaryOp) isGoOp() {
}

type Attr struct {
	Expr Register_Ref
	Name string
	Dst  Register_Ref
}

func (node *Attr) isGoOp() {
}

type Call struct {
	Name string
	Args []Register_Ref
	Dsts []Register_Ref
}

func (node *Call) isGoOp() {
}

type MethodCall struct {
	Expr Register_Ref
	Name string
	Args []Register_Ref
	Dsts []Register_Ref
}

func (node *MethodCall) isGoOp() {
}

type NamedArg struct {
	Name string
	Arg  Register_Ref
}

type ConstructStruct struct {
	Type      *tree.StructType
	AddrTaken bool
	Args      []*NamedArg
	Dst       Register_Ref
}

func (node *ConstructStruct) isGoOp() {
}

type ConstructSlice struct {
	Type *tree.SliceType
	Args []Register_Ref
	Dst  Register_Ref
}

func (node *ConstructSlice) isGoOp() {
}

type Coerce struct {
	Src  Register_Ref
	Type tree.GoType
	Dst  Register_Ref
}

func (node *Coerce) isGoOp() {
}

type Transfer struct {
	Srcs []Register_Ref
	Dsts []Register_Ref
}

func (node *Transfer) isGoOp() {
}

type Return struct {
	Args []Register_Ref
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
	Cond Register_Ref
}

func (node *Switch) isGoOp() {
}

type Exit struct {
}

func (node *Exit) isGoOp() {
}
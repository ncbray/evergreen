package interpreter

type Op interface {
}

type Locals []int

type StoreConst struct {
	Const  int
	Target int
}

type Return struct {
	Args Locals
}

type ConditionalJump struct {
	Arg      int
	Location int
}

type Jump struct {
	Location int
}

type Call struct {
	Func    int
	Args    Locals
	Targets Locals
}

type BinOp int

const (
	ADD BinOp = iota
	SUB
	MUL
	DIV
	REM
)

type BinaryOp struct {
	Op     BinOp
	Left   int
	Right  int
	Target int
}

type Function struct {
	Name      string
	NumParams int
	NumLocals int
	Constants []Object
	Body      []Op
}

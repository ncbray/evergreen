package interpreter

const (
	NORMAL = iota
	FAIL
	NUM_FLOWS
)

type Type interface {
}

type Object interface {
	Type() Type
}

type I32Type struct {
}

var i32Type Type = &I32Type{}

type I32 struct {
	Value int32
}

func (o *I32) Type() Type {
	return i32Type
}

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

type Function struct {
	Name      string
	NumParams int
	NumLocals int
	Constants []Object
	Body      []Op
}

type StackFrame struct {
	F        *Function
	Location int
	Locals   []Object
	Parent   *StackFrame
}

func toBool(o Object) bool {
	switch o := o.(type) {
	case *I32:
		return o.Value != 0
	default:
		panic(o)
	}
}

type Interpreter struct {
	Frame   *StackFrame
	Flow    int
	Temp    []Object
	TempLen int
}

func (i *Interpreter) GatherTemp(args Locals) {
	for idx, lcl := range args {
		i.Temp[idx] = i.Frame.Locals[lcl]
	}
	i.TempLen = len(args)
}

func (i *Interpreter) SetTemp(args []Object) {
	for idx, o := range args {
		i.Temp[idx] = o
	}
	i.TempLen = len(args)
}

func (i *Interpreter) Invoke(f *Function) {
	if i.TempLen != f.NumParams {
		panic(i.TempLen)
	}

	// Create the new stack frame.
	i.Frame = &StackFrame{
		F:        f,
		Location: 0,
		Locals:   make([]Object, f.NumLocals),
		Parent:   i.Frame,
	}

	// Set the function arguments.
	for idx := 0; idx < i.TempLen; idx++ {
		i.Frame.Locals[idx] = i.Temp[idx]
	}
}

func (i *Interpreter) Run() {
	for {
		switch op := i.Frame.F.Body[i.Frame.Location].(type) {
		case *ConditionalJump:
			if toBool(i.Frame.Locals[op.Arg]) {
				i.Frame.Location = op.Location
				continue
			}
		case *Jump:
			i.Frame.Location = op.Location
			continue
		case *StoreConst:
			i.Frame.Locals[op.Target] = i.Frame.F.Constants[op.Const]
		case *Return:
			i.GatherTemp(op.Args)
			i.Frame = i.Frame.Parent
			if i.Frame == nil {
				// Returned off end of stack.
				return
			}
			// TODO return value assignment.
			continue // Skip PC increment
		default:
			panic(op)
		}
		i.Frame.Location += 1
	}
}

func CreateInterpreter() *Interpreter {
	i := &Interpreter{
		Flow:    NORMAL,
		Temp:    make([]Object, 10),
		TempLen: 0,
	}
	return i
}

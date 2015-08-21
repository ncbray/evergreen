package interpreter

const (
	NORMAL = iota
	FAIL
	NUM_FLOWS
)

type StackFrame struct {
	F        *Function
	Location int
	Locals   []Object
	Targets  Locals
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

func binop(op BinOp, l Object, r Object) Object {
	switch l := l.(type) {
	case *I32:
		switch r := r.(type) {
		case *I32:
			switch op {
			case ADD:
				return &I32{Value: l.Value + r.Value}
			case SUB:
				return &I32{Value: l.Value - r.Value}
			case MUL:
				return &I32{Value: l.Value * r.Value}
			case DIV:
				return &I32{Value: l.Value / r.Value}
			case REM:
				return &I32{Value: l.Value % r.Value}
			default:
				panic(op)
			}
		default:
			panic(r)
		}
	default:
		panic(l)
	}
}

type Interpreter struct {
	Frame   *StackFrame
	Flow    int
	Temp    []Object
	TempLen int
	Funcs   []*Function
}

func (i *Interpreter) GatherTemp(args Locals) {
	for idx, lcl := range args {
		i.Temp[idx] = i.Frame.Locals[lcl]
	}
	i.TempLen = len(args)
}

func (i *Interpreter) ScatterTemp(targets Locals) {
	for idx, lcl := range targets {
		i.Frame.Locals[lcl] = i.Temp[idx]
	}
}

func (i *Interpreter) SetTemp(args []Object) {
	for idx, o := range args {
		i.Temp[idx] = o
	}
	i.TempLen = len(args)
}

func (i *Interpreter) Invoke(uid int) {
	f := i.Funcs[uid]
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
		case *BinaryOp:
			i.Frame.Locals[op.Target] = binop(op.Op, i.Frame.Locals[op.Left], i.Frame.Locals[op.Right])
		case *GetAttr:
			expr := i.Frame.Locals[op.Expr]
			switch expr := expr.(type) {
			case *Struct:
				i.Frame.Locals[op.Target] = expr.Slots[op.Slot]
			default:
				panic(expr)
			}
		case *Call:
			i.GatherTemp(op.Args)
			i.Frame.Targets = op.Targets
			i.Invoke(op.Func)
			continue
		case *Return:
			i.GatherTemp(op.Args)
			i.Frame = i.Frame.Parent
			if i.Frame == nil {
				// Returned off end of stack.
				return
			}
			i.ScatterTemp(i.Frame.Targets)
		default:
			panic(op)
		}
		i.Frame.Location += 1
	}
}

func CreateInterpreter(funcs []*Function) *Interpreter {
	i := &Interpreter{
		Flow:    NORMAL,
		Temp:    make([]Object, 10),
		TempLen: 0,
		Funcs:   funcs,
	}
	return i
}

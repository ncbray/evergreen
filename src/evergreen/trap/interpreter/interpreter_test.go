package interpreter

import (
	"testing"
)

func callAndReturnInt(i *Interpreter, f int, args []Object, value int32, t *testing.T) {
	i.SetTemp(args)
	i.Invoke(f)
	i.Run()
	if i.Flow != NORMAL {
		t.Errorf("Expected normal flow, got %d", i.Flow)
	}
	if i.TempLen != 1 {
		t.Errorf("Expected 1 return value, got %d", i.TempLen)
	}
	ret := i.Temp[0]
	switch ret := ret.(type) {
	case *I32:
		if ret.Value != value {
			t.Errorf("Expected return value of %d, got %d", value, ret.Value)
		}
	default:
		t.Errorf("Return value is not int32: %#v", ret)
	}
}

func TestSanity(t *testing.T) {
	i := CreateInterpreter([]*Function{})
	if i.Flow != NORMAL {
		t.Errorf("Expected normal flow, got %d", i.Flow)
	}
}

func TestReturnInt(t *testing.T) {
	b := CreateProgramBuilder()
	funcs := []*Function{
		&Function{
			Name:      "Foo",
			NumParams: 1,
			NumLocals: 1,
			Constants: []Object{
				b.i32(7),
				b.i32(11),
			},
			Body: []Op{
				&ConditionalJump{Arg: 0, Location: 3},
				&StoreConst{Const: 0, Target: 0},
				&Jump{Location: 4},
				&StoreConst{Const: 1, Target: 0},
				&Return{Args: Locals{0}},
			},
		},
	}

	i := CreateInterpreter(funcs)

	callAndReturnInt(i, 0, []Object{b.i32(0)}, 7, t)
	callAndReturnInt(i, 0, []Object{b.i32(1)}, 11, t)
}

func TestSwap(t *testing.T) {
	b := CreateProgramBuilder()
	funcs := []*Function{
		&Function{
			Name:      "Swap",
			NumParams: 2,
			NumLocals: 2,
			Constants: []Object{},
			Body: []Op{
				&Return{Args: Locals{1, 0}},
			},
		},
		&Function{
			Name:      "Main",
			NumParams: 2,
			NumLocals: 2,
			Constants: []Object{},
			Body: []Op{
				&Call{Func: 0, Args: Locals{0, 1}, Targets: Locals{0, 1}},
				&BinaryOp{Op: SUB, Left: 0, Right: 1, Target: 0},
				&Return{Args: Locals{0}},
			},
		},
	}

	i := CreateInterpreter(funcs)

	callAndReturnInt(i, 1, []Object{b.i32(5), b.i32(11)}, 6, t)
	callAndReturnInt(i, 1, []Object{b.i32(7), b.i32(4)}, -3, t)
}

func TestGetAttr(t *testing.T) {
	b := CreateProgramBuilder()
	funcs := []*Function{
		&Function{
			Name:      "Get0",
			NumParams: 1,
			NumLocals: 2,
			Constants: []Object{},
			Body: []Op{
				&GetAttr{Expr: 0, Slot: 0, Target: 1},
				&Return{Args: Locals{1}},
			},
		},
		&Function{
			Name:      "Get1",
			NumParams: 1,
			NumLocals: 2,
			Constants: []Object{},
			Body: []Op{
				&GetAttr{Expr: 0, Slot: 1, Target: 1},
				&Return{Args: Locals{1}},
			},
		},
	}

	i := CreateInterpreter(funcs)

	o := &Struct{Slots: []Object{b.i32(13), b.i32(17)}}

	callAndReturnInt(i, 0, []Object{o}, 13, t)
	callAndReturnInt(i, 1, []Object{o}, 17, t)
}

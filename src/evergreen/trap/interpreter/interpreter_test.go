package interpreter

import (
	"testing"
)

func callAndReturnInt(i *Interpreter, f *Function, args []Object, value int32, t *testing.T) {
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
	i := CreateInterpreter()
	if i.Flow != NORMAL {
		t.Errorf("Expected normal flow, got %d", i.Flow)
	}
}

func TestReturnInt(t *testing.T) {
	i := CreateInterpreter()

	f := &Function{
		Name:      "Foo",
		NumParams: 1,
		NumLocals: 1,
		Constants: []Object{
			&I32{Value: 7},
			&I32{Value: 11},
		},
		Body: []Op{
			&ConditionalJump{Arg: 0, Location: 3},
			&StoreConst{Const: 0, Target: 0},
			&Jump{Location: 4},
			&StoreConst{Const: 1, Target: 0},
			&Return{Args: Locals{0}},
		},
	}

	callAndReturnInt(i, f, []Object{&I32{Value: 0}}, 7, t)
	callAndReturnInt(i, f, []Object{&I32{Value: 1}}, 11, t)
}

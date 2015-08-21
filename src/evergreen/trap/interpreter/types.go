package interpreter

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

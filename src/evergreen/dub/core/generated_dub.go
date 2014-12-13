package core

type DubType interface {
	isDubType()
}

type BuiltinType struct {
	Name string
}

func (node *BuiltinType) isDubType() {
}

type NilType struct {
}

func (node *NilType) isDubType() {
}

type ListType struct {
	Type DubType
}

func (node *ListType) isDubType() {
}

type FieldType struct {
	Name string
	Type DubType
}

type StructType struct {
	Name       string
	Implements *StructType
	Fields     []*FieldType
	Scoped     bool
	Contains   []*StructType
	IsParent   bool
}

func (node *StructType) isDubType() {
}

type Function_Ref uint32

type Function_Scope struct {
	objects []*Function
}

type Function struct {
	Name string
}

const NoFunction = ^Function_Ref(0)

type CoreProgram struct {
	Structures     []*StructType
	Function_Scope *Function_Scope
}

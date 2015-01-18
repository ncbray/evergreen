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
	File       File_Ref
}

func (node *StructType) isDubType() {
}

type Package_Ref uint32

const NoPackage = ^Package_Ref(0)

type Package_Scope struct {
	objects []*Package
}

type Package struct {
	Path  []string
	Files []File_Ref
}

type File_Ref uint32

const NoFile = ^File_Ref(0)

type File_Scope struct {
	objects []*File
}

type File struct {
	Name    string
	Package Package_Ref
}

type Function_Ref uint32

const NoFunction = ^Function_Ref(0)

type Function_Scope struct {
	objects []*Function
}

type Function struct {
	Name string
	File File_Ref
}

type IntrinsticFunction struct {
	Name string
}

type Callable interface {
	isCallable()
}

type CallableFunction struct {
	Func Function_Ref
}

func (node *CallableFunction) isCallable() {
}

type CallableIntrinstic struct {
	Func *IntrinsticFunction
}

func (node *CallableIntrinstic) isCallable() {
}

type BuiltinTypeIndex struct {
	String  *BuiltinType
	Rune    *BuiltinType
	Int     *BuiltinType
	Int64   *BuiltinType
	Float32 *BuiltinType
	Bool    *BuiltinType
	Graph   *BuiltinType
	Nil     *NilType
}

type CoreProgram struct {
	Structures     []*StructType
	Package_Scope  *Package_Scope
	File_Scope     *File_Scope
	Function_Scope *Function_Scope
}
